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
	"context"
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/SENERGY-Platform/device-command/pkg/configuration"
	"github.com/SENERGY-Platform/device-command/pkg/messages"
	paho "github.com/eclipse/paho.mqtt.golang"
)

const (
	mqttCommandTopicPrefix  = "command/"
	mqttResponseTopicFilter = "response/#"
	mqttErrorTopicPrefix    = "error/command/"
)

// MqttCommand is the payload of a command to, and of a response from, a connector on a gateway.
// The connector echoes the command id; it is mgw_correlation_id_prefix followed by the task id.
type MqttCommand struct {
	CommandId string `json:"command_id"`
	Data      string `json:"data"`
}

// Mqtt sends a command to command/<device local id>/<service local id> on the broker of a gateway.
// The connector there knows one protocol segment, mgw_protocol_segment, which carries the data.
type Mqtt struct {
	client  paho.Client
	prefix  string
	segment string
}

func NewMqtt(ctx context.Context, config configuration.Config, receive func(Reply)) (*Mqtt, error) {
	//every client of the broker gets every response; the prefix tells ours apart
	if config.MgwCorrelationIdPrefix == "" {
		return nil, errors.New("missing mgw_correlation_id_prefix: without it responses to other clients of the broker are read as own")
	}
	result := &Mqtt{prefix: config.MgwCorrelationIdPrefix, segment: config.MgwProtocolSegment}

	//the session is clean, so the subscriptions are renewed on every (re)connect
	subscribed := make(chan error, 1)
	once := sync.Once{}
	options := paho.NewClientOptions().
		SetPassword(config.MgwMqttPw).
		SetUsername(config.MgwMqttUser).
		SetAutoReconnect(true).
		SetCleanSession(true).
		SetClientID(config.MgwMqttClientId).
		AddBroker(config.MgwMqttBroker).
		SetWriteTimeout(10 * time.Second).
		SetOrderMatters(false).
		SetConnectionLostHandler(func(_ paho.Client, err error) {
			slog.Error("connection to mqtt broker lost", "error", err)
		}).
		SetOnConnectHandler(func(client paho.Client) {
			slog.Info("(re)connected to mqtt broker")
			err := result.subscribe(client, receive)
			first := false
			once.Do(func() {
				first = true
				subscribed <- err
			})
			if err != nil && !first {
				slog.Error("FATAL: unable to renew mqtt subscriptions", "error", err)
				log.Fatal("FATAL: ", err)
			}
		})
	result.client = paho.NewClient(options)
	if token := result.client.Connect(); token.Wait() && token.Error() != nil {
		slog.Error("unable to connect to mqtt broker", "error", token.Error())
		return nil, token.Error()
	}
	err := <-subscribed
	if err != nil {
		result.client.Disconnect(0)
		return nil, err
	}
	go func() {
		<-ctx.Done()
		result.client.Disconnect(0)
	}()
	return result, nil
}

func (this *Mqtt) subscribe(client paho.Client, receive func(Reply)) error {
	token := client.Subscribe(mqttResponseTopicFilter, 2, func(_ paho.Client, message paho.Message) {
		if reply, ok := this.readResponse(message.Payload()); ok {
			receive(reply)
		}
	})
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	token = client.Subscribe(mqttErrorTopicPrefix+"#", 2, func(_ paho.Client, message paho.Message) {
		if reply, ok := this.readError(message.Topic(), message.Payload()); ok {
			receive(reply)
		}
	})
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// readResponse returns false for a response to another client of the broker, such as the
// automation-engine, which sends commands to the same connectors with a prefix of its own.
func (this *Mqtt) readResponse(payload []byte) (Reply, bool) {
	response := MqttCommand{}
	err := json.Unmarshal(payload, &response)
	if err != nil {
		slog.Error("unable to unmarshal mqtt response", "error", err)
		return Reply{}, false
	}
	taskId, ok := strings.CutPrefix(response.CommandId, this.prefix)
	if !ok {
		return Reply{}, false
	}
	return Reply{TaskId: taskId, Output: map[string]string{this.segment: response.Data}}, true
}

// readError returns false for an error of a command of another client of the broker.
func (this *Mqtt) readError(topic string, payload []byte) (Reply, bool) {
	taskId, ok := strings.CutPrefix(strings.TrimPrefix(topic, mqttErrorTopicPrefix), this.prefix)
	if !ok {
		return Reply{}, false
	}
	return Reply{TaskId: taskId, Output: map[string]string{this.segment: string(payload)}, Failed: true}, true
}

func (this *Mqtt) Send(msg messages.ProtocolMsg) error {
	payload, err := json.Marshal(MqttCommand{
		CommandId: this.prefix + msg.TaskInfo.TaskId,
		Data:      msg.Request.Input[this.segment],
	})
	if err != nil {
		return err
	}
	topic := mqttCommandTopicPrefix + msg.Metadata.Device.LocalId + "/" + msg.Metadata.Service.LocalId
	token := this.client.Publish(topic, 2, false, payload)
	if token.Wait() && token.Error() != nil {
		slog.Error("unable to publish mqtt command", "topic", topic, "error", token.Error())
		return token.Error()
	}
	return nil
}
