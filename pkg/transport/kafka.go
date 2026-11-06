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

package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/SENERGY-Platform/device-command/pkg/tracing"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/SENERGY-Platform/device-command/pkg/configuration"
	"github.com/SENERGY-Platform/device-command/pkg/messages"
	"github.com/SENERGY-Platform/service-commons/pkg/kafka"
)

// a reply that is older than this has no request waiting for it any more; it is skipped unread,
// so that a consumer that starts behind does not work through a backlog for nobody
const maxReplyAge = time.Hour

// Kafka sends a command to the topic named by the handler of its protocol, or by http POST if the
// handler is an http(s) url. Replies arrive on the response and error topics of this instance.
//
// The automation-engine sends to the same connectors. Its replies stay apart as long as response
// and error topics are not shared: on a shared topic a shared consumer group would split the
// replies between both services, and each would lose the ones the other read. Within the own
// topics a reply is matched by its task id only.
type Kafka struct {
	ctx         context.Context
	config      configuration.Config
	mux         sync.Mutex
	writers     map[string]*kafka.Producer
	knownTopics map[string]bool
}

func NewKafka(ctx context.Context, config configuration.Config, receive func(Reply)) (*Kafka, error) {
	if !isSet(config.KafkaUrl) {
		return nil, errors.New("missing kafka_url: the cloud com_impl receives replies only over kafka")
	}
	if isSet(config.ResponseTopic) {
		err := consume(ctx, config, config.ResponseTopic, false, receive)
		if err != nil {
			return nil, err
		}
	}
	if isSet(config.ErrorTopic) && !isHttpHandler(config.ErrorTopic) {
		err := consume(ctx, config, config.ErrorTopic, true, receive)
		if err != nil {
			return nil, err
		}
	}
	return &Kafka{
		ctx:         ctx,
		config:      config,
		writers:     map[string]*kafka.Producer{},
		knownTopics: map[string]bool{},
	}, nil
}

// consume never fails on a message: one that cannot be read is logged and skipped, as the consumer
// would otherwise retry it and hold back every reply behind it.
func consume(ctx context.Context, config configuration.Config, topic string, failed bool, receive func(Reply)) error {
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
		if age := time.Since(delivery.Time); age > maxReplyAge {
			slog.Warn("skip kafka message older than 1h", "topic", delivery.Topic, "age", age.String())
			return nil
		}
		msg := messages.ProtocolMsg{}
		err := json.Unmarshal(delivery.Value, &msg)
		if err != nil {
			slog.Error("unable to unmarshal protocol message", "error", err, "msg", string(delivery.Value))
			return nil
		}
		receive(Reply{TaskId: msg.TaskInfo.TaskId, Output: msg.Response.Output, Failed: failed})
		return nil
	})
}

func (this *Kafka) Send(ctx context.Context, msg messages.ProtocolMsg) error {
	message, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	topic := msg.Metadata.Protocol.Handler
	if isHttpHandler(topic) {
		return sendHttp(ctx, topic, message)
	}
	writer, err := this.writer(topic)
	if err != nil {
		return err
	}
	return writer.Produce(msg.Metadata.Device.Id, message)
}

func (this *Kafka) writer(topic string) (*kafka.Producer, error) {
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

// sendHttp posts to a protocol handler of the platform, which continues the trace
func sendHttp(ctx context.Context, url string, message []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(message))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	err = tracing.Inject(ctx, req)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
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

func isSet(value string) bool {
	return value != "" && value != "-"
}

func isHttpHandler(handler string) bool {
	return strings.HasPrefix(handler, "http://") || strings.HasPrefix(handler, "https://")
}
