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
	"github.com/SENERGY-Platform/device-command/pkg/command"
	"github.com/SENERGY-Platform/device-command/pkg/command/dependencies/impl/cloud"
	"github.com/SENERGY-Platform/device-command/pkg/command/dependencies/impl/mgw"
	"github.com/SENERGY-Platform/device-command/pkg/configuration"
	"github.com/SENERGY-Platform/external-task-worker/lib/com/kafka"
	"github.com/SENERGY-Platform/external-task-worker/lib/devicerepository/model"
	"github.com/SENERGY-Platform/external-task-worker/lib/messages"
	"github.com/Shopify/sarama"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestCommandUnscaled(t *testing.T) {
	testCommand("")(t)
}

func TestCommandScaled(t *testing.T) {
	testCommand("scalingSuffix")(t)
}

func testCommand(scalingSuffix string) func(t *testing.T) {
	return func(t *testing.T) {
		wg := &sync.WaitGroup{}
		defer wg.Wait()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		config, err := configuration.Load("../../config.json")
		if err != nil {
			t.Error(err)
			return
		}
		config.TopicSuffixForScaling = scalingSuffix
		config.Debug = true

		config.ServerPort, err = GetFreePort()
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
				"/devices/urn:infai:ses:device:timestamp-test": model.Device{
					Id:           "urn:infai:ses:device:timestamp-test",
					LocalId:      "d1-timestamp",
					Name:         "d1Name-timestamp",
					DeviceTypeId: "urn:infai:ses:device-type:24b294e8-4676-4782-8dc9-a008c0d94770",
				},
				"/devices/urn:infai:ses:device:a486084b-3323-4cbc-9f6b-d797373ae866": model.Device{
					Id:           "urn:infai:ses:device:a486084b-3323-4cbc-9f6b-d797373ae866",
					LocalId:      "d1",
					Name:         "d1Name",
					DeviceTypeId: "urn:infai:ses:device-type:755d892f-ec47-40ce-926a-59201328c138",
				},
				"/devices/temperature2": model.Device{
					Id:           "temperature2",
					LocalId:      "d1",
					Name:         "d1Name",
					DeviceTypeId: "urn:infai:ses:device-type:755d892f-ec47-40ce-926a-59201328c138",
				},
				"/devices/temperature3": model.Device{
					Id:           "temperature3",
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
				"/devices/lamp2": model.Device{
					Id:           "lamp2",
					LocalId:      "lamp2",
					Name:         "lamp2",
					DeviceTypeId: "urn:infai:ses:device-type:eb4a3337-01a1-4434-9dcc-064b3955eeef",
				},
				"/devices/color_event": model.Device{
					Id:           "color_event",
					LocalId:      "color_event",
					Name:         "color_event",
					DeviceTypeId: "urn:infai:ses:device-type:color_event",
				},
				"/device-groups/group_temperature": model.DeviceGroup{
					Id:        "group_temperature",
					DeviceIds: []string{"urn:infai:ses:device:a486084b-3323-4cbc-9f6b-d797373ae866", "temperature2", "temperature3"},
				},
				"/device-groups/group_color": model.DeviceGroup{
					Id:        "group_color",
					DeviceIds: []string{"color_event", "lamp", "lamp2"},
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
		producer.Log(log.New(os.Stdout, "[KAFKA-PRODUCER] ", 0))

		serviceCallCount := map[string]int{}

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
			case "urn:infai:ses:service:ec456e2a-81ed-4466-a119-daecfbb2d033":
				message.Response.Output = map[string]string{"data": `{"value": "clear", "lastUpdate": 0, "lastUpdate_unit": "unit"}`}
			case "urn:infai:ses:service:4932d451-3300-4a22-a508-ec740e5789b3":
				if message.Request.Input["data"] != "21" {
					t.Error(message.Request.Input)
				}
			case "urn:infai:ses:service:6d6067a3-ed4e-45ec-a7eb-b1695340d2f1":
				serviceCallCount[message.Metadata.Service.Id] = serviceCallCount[message.Metadata.Service.Id] + 1
				message.Response.Output = map[string]string{"data": `{"value": 13, "lastUpdate": 42}`}
			case "urn:infai:ses:service:4b6c4567-f256-4dbd-a562-d13442ad4530":
				//create timeout
				return nil
			case "urn:infai:ses:service:36fd778e-b04d-4d72-bed5-1b77ed1164b9":
				//create timeout
				return nil
			case "urn:infai:ses:service:1199edee-fbf7-44fb-9228-9a9db69bbdd4":
				message.Response.Output = map[string]string{"data": `{"brightness": 65, "hue": 176, "saturation": 70, "kelvin": 0, "on": true, "status": 200}`}
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
				return nil
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

		cmd, err := command.NewWithFactories(ctx, config, cloud.ComFactory, mgw.MarshallerFactory, cloud.IotFactory, cloud.TimescaleFactory)
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

		//some services are called as event (timescale call), some as request
		t.Run("device group color", sendCommand(config, api.CommandMessage{
			FunctionId: "urn:infai:ses:measuring-function:bdb6a7c8-4a3d-4fe0-bab3-ce02e09b5869",
			GroupId:    "group_color",
		}, 200, `[{"b":158,"g":166,"r":50},{"b":158,"g":166,"r":50},{"b":158,"g":166,"r":50}]`))

		//some services return a timeout, some return 13
		t.Run("device group getTemperature", sendCommand(config, api.CommandMessage{
			FunctionId: "urn:infai:ses:measuring-function:f2769eb9-b6ad-4f7e-bd28-e4ea043d2f8b",
			GroupId:    "group_temperature",
		}, 200, "[13,13,13]"))

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
				FunctionId: "urn:infai:ses:measuring-function:20d3c1d3-77d7-4181-a9f3-b487add58cd0",
				DeviceId:   "color_event",
				ServiceId:  "urn:infai:ses:service:color_event",
			},
		}, 200, `[{"status_code":200,"message":[null]},{"status_code":200,"message":[13]},{"status_code":408,"message":"timeout"},{"status_code":500,"message":"unable to load function: not found"},{"status_code":200,"message":[null]},{"status_code":200,"message":[{"b":158,"g":166,"r":50}]},{"status_code":200,"message":["on"]},{"status_code":200,"message":[null]},{"status_code":200,"message":[13]},{"status_code":200,"message":["on"]}]`))

		t.Run("new timestamp", sendCommandBatch(config, api.BatchRequest{
			{
				FunctionId: "urn:infai:ses:measuring-function:3b4e0766-0d67-4658-b249-295902cd3290",
				DeviceId:   "urn:infai:ses:device:timestamp-test",
				ServiceId:  "urn:infai:ses:service:ec456e2a-81ed-4466-a119-daecfbb2d033",
			},
		}, 200, `[{"status_code":200,"message":["1970-01-01T01:00:00+01:00"]}]`))

		t.Run("check callcount", func(t *testing.T) {
			if serviceCallCount["urn:infai:ses:service:6d6067a3-ed4e-45ec-a7eb-b1695340d2f1"] != 5 {
				t.Error(serviceCallCount)
			}
		})
	}
}

func sendCommandBatch(config configuration.Config, commandMessage api.BatchRequest, expectedCode int, expectedContent string) func(t *testing.T) {
	return func(t *testing.T) {
		buff := &bytes.Buffer{}
		err := json.NewEncoder(buff).Encode(commandMessage)
		if err != nil {
			t.Error(err)
			return
		}
		req, err := http.NewRequest("POST", "http://localhost:"+config.ServerPort+"/commands/batch?timeout=5s", buff)
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

func sendCommand(config configuration.Config, commandMessage api.CommandMessage, expectedCode int, expectedContent string) func(t *testing.T) {
	return func(t *testing.T) {
		buff := &bytes.Buffer{}
		err := json.NewEncoder(buff).Encode(commandMessage)
		if err != nil {
			t.Error(err)
			return
		}
		req, err := http.NewRequest("POST", "http://localhost:"+config.ServerPort+"/commands?timeout=5s", buff)
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

const testToken = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJqdGkiOiIwOGM0N2E4OC0yYzc5LTQyMGYtODEwNC02NWJkOWViYmU0MWUiLCJleHAiOjE1NDY1MDcyMzMsIm5iZiI6MCwiaWF0IjoxNTQ2NTA3MTczLCJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjgwMDEvYXV0aC9yZWFsbXMvbWFzdGVyIiwiYXVkIjoiZnJvbnRlbmQiLCJzdWIiOiJ0ZXN0T3duZXIiLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJmcm9udGVuZCIsIm5vbmNlIjoiOTJjNDNjOTUtNzViMC00NmNmLTgwYWUtNDVkZDk3M2I0YjdmIiwiYXV0aF90aW1lIjoxNTQ2NTA3MDA5LCJzZXNzaW9uX3N0YXRlIjoiNWRmOTI4ZjQtMDhmMC00ZWI5LTliNjAtM2EwYWUyMmVmYzczIiwiYWNyIjoiMCIsImFsbG93ZWQtb3JpZ2lucyI6WyIqIl0sInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJ1c2VyIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsibWFzdGVyLXJlYWxtIjp7InJvbGVzIjpbInZpZXctcmVhbG0iLCJ2aWV3LWlkZW50aXR5LXByb3ZpZGVycyIsIm1hbmFnZS1pZGVudGl0eS1wcm92aWRlcnMiLCJpbXBlcnNvbmF0aW9uIiwiY3JlYXRlLWNsaWVudCIsIm1hbmFnZS11c2VycyIsInF1ZXJ5LXJlYWxtcyIsInZpZXctYXV0aG9yaXphdGlvbiIsInF1ZXJ5LWNsaWVudHMiLCJxdWVyeS11c2VycyIsIm1hbmFnZS1ldmVudHMiLCJtYW5hZ2UtcmVhbG0iLCJ2aWV3LWV2ZW50cyIsInZpZXctdXNlcnMiLCJ2aWV3LWNsaWVudHMiLCJtYW5hZ2UtYXV0aG9yaXphdGlvbiIsIm1hbmFnZS1jbGllbnRzIiwicXVlcnktZ3JvdXBzIl19LCJhY2NvdW50Ijp7InJvbGVzIjpbIm1hbmFnZS1hY2NvdW50IiwibWFuYWdlLWFjY291bnQtbGlua3MiLCJ2aWV3LXByb2ZpbGUiXX19LCJyb2xlcyI6WyJ1c2VyIl19.ykpuOmlpzj75ecSI6cHbCATIeY4qpyut2hMc1a67Ycg`

func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port, nil
}

func GetFreePort() (string, error) {
	temp, err := getFreePort()
	return strconv.Itoa(temp), err
}
