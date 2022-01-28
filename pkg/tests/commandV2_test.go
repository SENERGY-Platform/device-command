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

package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/SENERGY-Platform/device-command/pkg/api"
	"github.com/SENERGY-Platform/device-command/pkg/commandV2"
	"github.com/SENERGY-Platform/device-command/pkg/configuration"
	"github.com/SENERGY-Platform/external-task-worker/lib/com/comswitch"
	"github.com/SENERGY-Platform/external-task-worker/lib/com/kafka"
	"github.com/SENERGY-Platform/external-task-worker/lib/devicerepository/model"
	"github.com/SENERGY-Platform/external-task-worker/lib/messages"
	"github.com/SENERGY-Platform/external-task-worker/lib/test/mock"
	"github.com/Shopify/sarama"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestCommandV2(t *testing.T) {
	wg := &sync.WaitGroup{}
	defer wg.Wait()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	config, err := configuration.Load("../../config.json")
	if err != nil {
		t.Error(err)
		return
	}

	config.ServerPort, err = GetFreePort()
	if err != nil {
		t.Error(err)
		return
	}
	config.HttpCommandConsumerPort, err = GetFreePort()
	if err != nil {
		t.Error(err)
		return
	}

	config, err = timescaleEnv(config, ctx, wg, map[string]map[string]map[string]interface{}{
		"color_event": {
			"urn:infai:ses:service:color_event": {
				"struct.hue":        176,
				"struct.saturation": 70,
				"struct.brightness": 65,
				"struct.on":         true,
				"struct.status":     200,
			},
		},
	})
	if err != nil {
		t.Error(err)
		return
	}

	devices := map[string]map[string]interface{}{
		"testOwner": {
			"/devices/urn:infai:ses:device:a486084b-3323-4cbc-9f6b-d797373ae866": model.Device{
				Id:           "urn:infai:ses:device:a486084b-3323-4cbc-9f6b-d797373ae866",
				LocalId:      "d1",
				Name:         "d1Name",
				DeviceTypeId: "urn:infai:ses:device-type:755d892f-ec47-40ce-926a-59201328c138",
			},
			"/devices/lamp": model.Device{
				Id:           "lamp",
				LocalId:      "lamp",
				Name:         "lamp",
				DeviceTypeId: "urn:infai:ses:device-type:eb4a3337-01a1-4434-9dcc-064b3955eeef",
			},
			"/devices/color_event": model.Device{
				Id:           "color_event",
				LocalId:      "color_event",
				Name:         "color_event",
				DeviceTypeId: "urn:infai:ses:device-type:color_event",
			},
		},
	}

	config, err = iotEnv(config, ctx, wg, devices)
	if err != nil {
		t.Error(err)
		return
	}

	config, err = kafkaEnv(config, ctx, wg)
	if err != nil {
		t.Error(err)
		return
	}
	time.Sleep(1 * time.Second)

	//connector mock
	maxWait, err := time.ParseDuration(config.KafkaConsumerMaxWait)
	if err != nil {
		t.Error(err)
		return
	}
	producer, err := kafka.PrepareProducerWithConfig(ctx, config.KafkaUrl, kafka.ProducerConfig{
		AsyncFlushFrequency: 100 * time.Millisecond,
		AsyncCompression:    sarama.CompressionNone,
		SyncCompression:     sarama.CompressionNone,
		Sync:                config.Sync,
		SyncIdempotent:      config.SyncIdempotent,
		PartitionNum:        int(config.PartitionNum),
		ReplicationFactor:   int(config.ReplicationFactor),
		AsyncFlushMessages:  int(config.AsyncFlushMessages),
		TopicConfigMap:      config.KafkaTopicConfigs,
	})
	if err != nil {
		t.Log(config.KafkaUrl)
		t.Error(err)
		return
	}
	err = kafka.NewConsumer(ctx, kafka.ConsumerConfig{
		KafkaUrl:       config.KafkaUrl,
		GroupId:        "test-connector-mock",
		Topic:          "connector",
		MinBytes:       int(config.KafkaConsumerMinBytes),
		MaxBytes:       int(config.KafkaConsumerMaxBytes),
		MaxWait:        maxWait,
		TopicConfigMap: config.KafkaTopicConfigs,
	}, func(_ string, msg []byte, time time.Time) error {
		t.Log(string(msg))
		message := messages.ProtocolMsg{}
		err := json.Unmarshal(msg, &message)
		if err != nil {
			t.Error(err)
			return nil
		}
		switch message.Metadata.Service.Id {
		case "urn:infai:ses:service:4932d451-3300-4a22-a508-ec740e5789b3":
			if message.Request.Input["data"] != "21" {
				t.Error(message.Request.Input)
			}
		case "urn:infai:ses:service:6d6067a3-ed4e-45ec-a7eb-b1695340d2f1":
			message.Response.Output = map[string]string{"data": `{"value": 13, "lastUpdate": 42}`}
		case "urn:infai:ses:service:36fd778e-b04d-4d72-bed5-1b77ed1164b9":
			//create timeout
			return nil
		case "urn:infai:ses:service:1b0ef253-16f7-4b65-8a15-fe79fccf7e70":
			input := map[string]float64{}
			err = json.Unmarshal([]byte(message.Request.Input["data"]), &input)
			if err != nil {
				t.Error(err)
				return nil
			}
			//values are truncated to integers, not rounded
			if input["brightness"] != float64(65) {
				t.Error(message.Request.Input)
			}
			if input["hue"] != float64(176) {
				t.Error(message.Request.Input)
			}
			if input["saturation"] != float64(70) {
				t.Error(message.Request.Input)
			}
			if input["duration"] != float64(1) {
				t.Error(message.Request.Input)
			}
		default:
			t.Error("unknown service-id", message.Metadata.Service.Id)
		}
		resp, err := json.Marshal(message)
		if err != nil {
			t.Error(err)
			return nil
		}
		err = producer.ProduceWithKey(message.Metadata.ResponseTo, "key", string(resp))
		if err != nil {
			t.Error(err)
			return nil
		}
		return nil
	}, func(err error) {
		t.Error(err)
	})

	cmd, err := command.NewWithFactories(ctx, config, comswitch.Factory, mock.Marshaller)
	if err != nil {
		t.Error(err)
		return
	}
	err = api.Start(ctx, config, cmd)
	if err != nil {
		t.Error(err)
		return
	}

	time.Sleep(1 * time.Second)

	t.Run("device setTemperature", sendCommand(config, api.CommandMessage{
		FunctionId: "urn:infai:ses:controlling-function:99240d90-02dd-4d4f-a47c-069cfe77629c",
		Input:      21,
		DeviceId:   "urn:infai:ses:device:a486084b-3323-4cbc-9f6b-d797373ae866",
		ServiceId:  "urn:infai:ses:service:4932d451-3300-4a22-a508-ec740e5789b3",
	}, 200, "[null]"))

	t.Run("device getTemperature", sendCommand(config, api.CommandMessage{
		FunctionId: "urn:infai:ses:measuring-function:f2769eb9-b6ad-4f7e-bd28-e4ea043d2f8b",
		DeviceId:   "urn:infai:ses:device:a486084b-3323-4cbc-9f6b-d797373ae866",
		ServiceId:  "urn:infai:ses:service:6d6067a3-ed4e-45ec-a7eb-b1695340d2f1",
	}, 200, "[13]"))

	t.Run("device timeout", sendCommand(config, api.CommandMessage{
		FunctionId: "urn:infai:ses:measuring-function:00549f18-88b5-44c7-adb1-f558e8d53d1d",
		DeviceId:   "urn:infai:ses:device:a486084b-3323-4cbc-9f6b-d797373ae866",
		ServiceId:  "urn:infai:ses:service:36fd778e-b04d-4d72-bed5-1b77ed1164b9",
	}, http.StatusRequestTimeout, `"timeout"`))

	t.Run("invalid command", sendCommand(config, api.CommandMessage{
		FunctionId: "foobar",
		DeviceId:   "urn:infai:ses:device:a486084b-3323-4cbc-9f6b-d797373ae866",
		ServiceId:  "urn:infai:ses:service:6d6067a3-ed4e-45ec-a7eb-b1695340d2f1",
	}, 500, `"unable to load function: not found"`))

	t.Run("device color", sendCommand(config, api.CommandMessage{
		FunctionId: "urn:infai:ses:controlling-function:c54e2a89-1fb8-4ecb-8993-a7b40b355599",
		Input: map[string]interface{}{
			"r": 50,
			"g": 168,
			"b": 162,
		},
		DeviceId:  "lamp",
		ServiceId: "urn:infai:ses:service:1b0ef253-16f7-4b65-8a15-fe79fccf7e70",
	}, 200, "[null]"))

	t.Run("device event color", sendCommand(config, api.CommandMessage{
		FunctionId: "urn:infai:ses:measuring-function:bdb6a7c8-4a3d-4fe0-bab3-ce02e09b5869",
		DeviceId:   "color_event",
		ServiceId:  "urn:infai:ses:service:color_event",
	}, 200, `[{"b":158,"g":166,"r":50}]`))

	t.Run("device event on", sendCommand(config, api.CommandMessage{
		FunctionId: "urn:infai:ses:measuring-function:20d3c1d3-77d7-4181-a9f3-b487add58cd0",
		DeviceId:   "color_event",
		ServiceId:  "urn:infai:ses:service:color_event",
	}, 200, `["on"]`))

	t.Run("device batch", sendCommandBatch(config, api.BatchRequest{
		{
			FunctionId: "urn:infai:ses:controlling-function:99240d90-02dd-4d4f-a47c-069cfe77629c",
			Input:      21,
			DeviceId:   "urn:infai:ses:device:a486084b-3323-4cbc-9f6b-d797373ae866",
			ServiceId:  "urn:infai:ses:service:4932d451-3300-4a22-a508-ec740e5789b3",
		},
		{
			FunctionId: "urn:infai:ses:measuring-function:f2769eb9-b6ad-4f7e-bd28-e4ea043d2f8b",
			DeviceId:   "urn:infai:ses:device:a486084b-3323-4cbc-9f6b-d797373ae866",
			ServiceId:  "urn:infai:ses:service:6d6067a3-ed4e-45ec-a7eb-b1695340d2f1",
		},
		{
			FunctionId: "urn:infai:ses:measuring-function:00549f18-88b5-44c7-adb1-f558e8d53d1d",
			DeviceId:   "urn:infai:ses:device:a486084b-3323-4cbc-9f6b-d797373ae866",
			ServiceId:  "urn:infai:ses:service:36fd778e-b04d-4d72-bed5-1b77ed1164b9",
		},
		{
			FunctionId: "foobar",
			DeviceId:   "urn:infai:ses:device:a486084b-3323-4cbc-9f6b-d797373ae866",
			ServiceId:  "urn:infai:ses:service:6d6067a3-ed4e-45ec-a7eb-b1695340d2f1",
		},
		{
			FunctionId: "urn:infai:ses:controlling-function:c54e2a89-1fb8-4ecb-8993-a7b40b355599",
			Input: map[string]interface{}{
				"r": 50,
				"g": 168,
				"b": 162,
			},
			DeviceId:  "lamp",
			ServiceId: "urn:infai:ses:service:1b0ef253-16f7-4b65-8a15-fe79fccf7e70",
		},
		{
			FunctionId: "urn:infai:ses:measuring-function:bdb6a7c8-4a3d-4fe0-bab3-ce02e09b5869",
			DeviceId:   "color_event",
			ServiceId:  "urn:infai:ses:service:color_event",
		},
		{
			FunctionId: "urn:infai:ses:measuring-function:20d3c1d3-77d7-4181-a9f3-b487add58cd0",
			DeviceId:   "color_event",
			ServiceId:  "urn:infai:ses:service:color_event",
		},
	}, 200, `[{"status_code":200,"message":[null]},{"status_code":200,"message":[13]},{"status_code":408,"message":"timeout"},{"status_code":500,"message":"unable to load function: not found"},{"status_code":200,"message":[null]},{"status_code":200,"message":[{"b":158,"g":166,"r":50}]},{"status_code":200,"message":["on"]}]`))
}

func sendCommandBatch(config configuration.Config, commandMessage api.BatchRequest, expectedCode int, expectedContent string) func(t *testing.T) {
	return func(t *testing.T) {
		buff := &bytes.Buffer{}
		err := json.NewEncoder(buff).Encode(commandMessage)
		if err != nil {
			t.Error(err)
			return
		}
		req, err := http.NewRequest("POST", "http://localhost:"+config.ServerPort+"/commands/batch", buff)
		if err != nil {
			t.Error(err)
			return
		}
		req.Header.Set("Authorization", testToken)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Error(err)
			return
		}
		defer resp.Body.Close()
		actualContent, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
			return
		}
		if resp.StatusCode != expectedCode {
			t.Error(resp.StatusCode, string(actualContent))
			return
		}
		if strings.TrimSpace(string(actualContent)) != expectedContent {
			t.Error("\n" + expectedContent + "\n" + string(actualContent))
		}
	}
}
