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

package mgw

import (
	"context"
	"encoding/json"
	"github.com/SENERGY-Platform/device-command/pkg/command/dependencies/impl/mgw/mqtt"
	"github.com/SENERGY-Platform/device-command/pkg/command/dependencies/interfaces"
	"github.com/SENERGY-Platform/device-command/pkg/configuration"
	"github.com/SENERGY-Platform/external-task-worker/lib/messages"
	"github.com/google/uuid"
	"log"
	"strings"
)

func ComFactory(ctx context.Context, config configuration.Config, responseListener func(msg messages.ProtocolMsg) error, errorListener func(msg messages.ProtocolMsg) error) (producer interfaces.Producer, err error) {
	correlationService := &CorrelationImpl{}

	client, err := mqtt.New(ctx, config.MgwMqttBroker, config.MgwMqttClientId, config.MgwMqttUser, config.MgwMqttPw)
	if err != nil {
		return producer, err
	}

	impl := &ComImpl{config: config, client: client, correlation: correlationService}

	err = impl.mgwSubscriptions(client, responseListener, errorListener)
	if err != nil {
		return producer, err
	}
	return impl, nil
}

func (this *ComImpl) mgwSubscriptions(client *mqtt.Mqtt, listener func(msg messages.ProtocolMsg) error, errorListener func(msg messages.ProtocolMsg) error) error {
	respTopic := "response/#"
	err := client.Subscribe(respTopic, 2, func(topic string, message []byte) {
		msg := Command{}
		err := json.Unmarshal(message, &msg)
		if err != nil {
			log.Println("ERROR: unable to unmarshal response to mgw command wrapper", err)
			return
		}
		if strings.HasPrefix(msg.CommandId, this.config.MgwCorrelationIdPrefix) {
			convertedMsg, err := this.convertResponseMessage(msg)
			if err != nil {
				log.Println("ERROR: unable to convert response", err)
				return
			}
			go func() {
				err = listener(convertedMsg)
				if err != nil {
					log.Println("ERROR: unable to handle response", err)
					return
				}
			}()
		}
	})
	if err != nil {
		return err
	}

	errTopicPrefix := "error/command/"
	errTopic := errTopicPrefix + "#"
	err = client.Subscribe(errTopic, 2, func(topic string, message []byte) {
		correlationId := strings.Replace(topic, errTopicPrefix, "", 1)
		if strings.HasPrefix(correlationId, this.config.MgwCorrelationIdPrefix) {
			convertedMsg, err := this.convertErrorMessage(correlationId, string(message))
			if err != nil {
				log.Println("ERROR: unable to convert error response", err)
				return
			}
			go func() {
				err = errorListener(convertedMsg)
				if err != nil {
					log.Println("ERROR: unable to handle response", err)
					return
				}
			}()
		}
	})
	if err != nil {
		return err
	}

	return nil
}

type ComImpl struct {
	config      configuration.Config
	client      *mqtt.Mqtt
	correlation Correlation
}

func (this *ComImpl) SendCommand(msg messages.ProtocolMsg) (err error) {
	topic, payload, err := this.convertCommandMessage(msg)
	if err != nil {
		return err
	}
	return this.client.Publish(topic, 2, false, payload)
}

func (this *ComImpl) convertResponseMessage(msg Command) (messages.ProtocolMsg, error) {
	target, err := this.getCorrelatedTask(msg.CommandId)
	if err != nil {
		return messages.ProtocolMsg{}, err
	}
	if target.Response.Output == nil {
		target.Response.Output = map[string]string{}
	}
	target.Response.Output[this.config.MgwProtocolSegment] = msg.Data
	return target, err
}

func (this *ComImpl) convertErrorMessage(commandId string, payload string) (messages.ProtocolMsg, error) {
	target, err := this.getCorrelatedTask(commandId)
	if err != nil {
		return messages.ProtocolMsg{}, err
	}
	if target.Response.Output == nil {
		target.Response.Output = map[string]string{}
	}
	target.Response.Output[this.config.MgwProtocolSegment] = payload
	return target, err
}

func (this *ComImpl) getCorrelatedTask(id string) (messages.ProtocolMsg, error) {
	return this.correlation.Get(id)
}

func (this *ComImpl) convertCommandMessage(source messages.ProtocolMsg) (mqttTopic string, mqttMessage []byte, err error) {
	if source.Request.Input == nil {
		source.Request.Input = map[string]string{}
	}
	idSuffix := IdProvider()
	correlationId := this.config.MgwCorrelationIdPrefix + idSuffix
	err = this.correlation.Set(correlationId, source)
	if err != nil {
		return
	}

	mqttTopic = "command/" + source.Metadata.Device.LocalId + "/" + source.Metadata.Service.LocalId

	data := source.Request.Input[this.config.MgwProtocolSegment]
	target := Command{
		CommandId: correlationId,
		Data:      data,
	}
	mqttMessage, err = json.Marshal(target)

	return
}

type Command struct {
	CommandId string `json:"command_id"`
	Data      string `json:"data"`
}

var IdProvider = DefaultIdProviderImpl

func DefaultIdProviderImpl() string {
	return uuid.NewString()
}
