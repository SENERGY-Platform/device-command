/*
 * Copyright 2022 InfAI (CC SES)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cloud

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/SENERGY-Platform/device-command/pkg/command/dependencies/interfaces"
	"github.com/SENERGY-Platform/device-command/pkg/configuration"
	"github.com/SENERGY-Platform/device-command/pkg/messages"
	"github.com/SENERGY-Platform/service-commons/pkg/kafka"
)

// a response that is older than this has no request waiting for it any more; it is skipped
// unread, so that a consumer that starts behind does not unmarshal a backlog for nobody
const maxResponseAge = time.Hour

func ComFactory(ctx context.Context, config configuration.Config, responseListener func(msg messages.ProtocolMsg) error, errorListener func(msg messages.ProtocolMsg) error) (producer interfaces.Producer, err error) {
	if !isSet(config.KafkaUrl) {
		return nil, errors.New("missing kafka_url: the cloud com_impl receives responses only over kafka")
	}
	resp := getListener(responseListener)
	if config.ResponseWorkerCount > 1 {
		resp = getQueuedResponseHandler(ctx, config.ResponseWorkerCount, config.ResponseWorkerCount, resp)
	}
	if isSet(config.ResponseTopic) {
		err = startConsumer(ctx, config, config.ResponseTopic, resp)
		if err != nil {
			return nil, err
		}
	}
	if isSet(config.ErrorTopic) && !isHttpHandler(config.ErrorTopic) {
		err = startConsumer(ctx, config, config.ErrorTopic, getListener(errorListener))
		if err != nil {
			return nil, err
		}
	}
	return &Producer{
		ctx:         ctx,
		config:      config,
		writers:     map[string]*kafka.Producer{},
		knownTopics: map[string]bool{},
	}, nil
}

func startConsumer(ctx context.Context, config configuration.Config, topic string, listener func(msg []byte) error) error {
	slog.Info("consume topic", "topic", topic)
	if config.InitTopics {
		err := initTopics(config.KafkaUrl, config.KafkaTopicConfigs, 1, 1, topic)
		if err != nil {
			slog.Error("unable to create topic", "topic", topic, "error", err)
			return err
		}
	}
	return kafka.NewMultiConsumer(ctx, kafka.Config{
		KafkaUrl:      config.KafkaUrl,
		ConsumerGroup: config.KafkaConsumerGroup,
		StartOffset:   kafka.FirstOffset,
	}, []string{topic}, func(delivery kafka.Message) error {
		if age := time.Since(delivery.Time); age > maxResponseAge {
			slog.Warn("skip kafka message older than 1h", "topic", delivery.Topic, "age", age.String())
			return nil
		}
		return listener(delivery.Value)
	})
}

// Producer sends a protocol message to the topic named by the handler of its protocol, or by
// http POST if the handler is an http(s) url.
type Producer struct {
	ctx         context.Context
	config      configuration.Config
	mux         sync.Mutex
	writers     map[string]*kafka.Producer
	knownTopics map[string]bool
}

func (this *Producer) SendCommand(msg messages.ProtocolMsg) (err error) {
	message, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	topic := msg.Metadata.Protocol.Handler
	if isHttpHandler(topic) {
		return sendHttp(topic, message)
	}
	writer, err := this.getWriter(topic)
	if err != nil {
		return err
	}
	return writer.Produce(msg.Metadata.Device.Id, message)
}

func (this *Producer) getWriter(topic string) (*kafka.Producer, error) {
	this.mux.Lock()
	defer this.mux.Unlock()
	if this.config.InitTopics && !this.knownTopics[topic] {
		err := initTopics(this.config.KafkaUrl, this.config.KafkaTopicConfigs, int(this.config.PartitionNum), int(this.config.ReplicationFactor), topic)
		if err != nil {
			slog.Warn("unable to ensure topic", "topic", topic, "error", err)
		} else {
			this.knownTopics[topic] = true
		}
	}
	if writer, ok := this.writers[topic]; ok {
		return writer, nil
	}
	writer, err := kafka.NewProducer(this.ctx, kafka.Config{KafkaUrl: this.config.KafkaUrl}, topic)
	if err != nil {
		return nil, err
	}
	this.writers[topic] = writer
	return writer, nil
}

func sendHttp(url string, message []byte) error {
	resp, err := http.Post(url, "application/json", strings.NewReader(string(message)))
	if err != nil {
		slog.Error("unable to send command by http", "url", url, "error", err)
		return err
	}
	defer resp.Body.Close()
	respMsg, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		err = errors.New("http producer: " + resp.Status + " " + string(respMsg))
		slog.Error("unable to send command by http", "url", url, "error", err)
		return err
	}
	return nil
}

// getListener never returns an error: a message that cannot be read or handled is logged and
// skipped. The consumer would otherwise retry it and hold back every response behind it.
func getListener(listener func(msg messages.ProtocolMsg) error) func(msg []byte) error {
	return func(msg []byte) error {
		message := messages.ProtocolMsg{}
		err := json.Unmarshal(msg, &message)
		if err != nil {
			slog.Error("unable to unmarshal protocol message", "error", err, "msg", string(msg))
			return nil
		}
		err = listener(message)
		if err != nil {
			slog.Error("unable to handle protocol message", "error", err, "msg", string(msg))
		}
		return nil
	}
}

func getQueuedResponseHandler(ctx context.Context, workerCount int64, queueSize int64, respHandler func(msg []byte) error) func(msg []byte) (err error) {
	queue := make(chan []byte, queueSize)
	for i := int64(0); i < workerCount; i++ {
		go func() {
			for msg := range queue {
				err := respHandler(msg)
				if err != nil {
					slog.Error("unable to handle response", "error", err, "msg", string(msg))
				}
			}
		}()
	}
	go func() {
		<-ctx.Done()
		close(queue)
	}()
	return func(msg []byte) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = errors.New(fmt.Sprint(r))
			}
		}()
		queue <- msg
		return err
	}
}

func isSet(value string) bool {
	return value != "" && value != "-"
}

func isHttpHandler(handler string) bool {
	return strings.HasPrefix(handler, "http://") || strings.HasPrefix(handler, "https://")
}
