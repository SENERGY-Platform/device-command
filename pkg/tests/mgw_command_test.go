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
	"context"
	"encoding/json"
	"fmt"
	"github.com/SENERGY-Platform/device-command/pkg/api"
	"github.com/SENERGY-Platform/device-command/pkg/command"
	"github.com/SENERGY-Platform/device-command/pkg/command/dependencies/impl/mgw"
	"github.com/SENERGY-Platform/device-command/pkg/command/dependencies/impl/mgw/mqtt"
	"github.com/SENERGY-Platform/device-command/pkg/configuration"
	"github.com/SENERGY-Platform/external-task-worker/lib/devicerepository/model"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestMgwCommand(t *testing.T) {
	fallbackFile := filepath.Join(t.TempDir(), "iot_fallback.json")
	t.Run("without auth overwrite to fill fallback file", testMgwCommand(fallbackFile, false))
	t.Run("with auth overwrite", testMgwCommand(fallbackFile, true))
}

func testMgwCommand(fallbackPath string, useAuthOverwriteFallback bool) func(t *testing.T) {
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
		config.Debug = true
		config.ComImpl = "mgw"
		config.MarshallerImpl = "mgw"
		config.UseIotFallback = true
		config.TimescaleImpl = "mgw"
		config.IotFallbackFile = fallbackPath
		config.OverwriteAuthToken = useAuthOverwriteFallback

		config.ServerPort, err = GetFreePort()
		if err != nil {
			t.Error(err)
			return
		}

		config, err = timescaleEnv(config, ctx, wg, map[string]map[string]map[string]interface{}{
			"color_event_lid": {
				"getStatus": {
					"hue":        176,
					"saturation": 70,
					"brightness": 65,
					"on":         true,
					"status":     200,
				},
			},
			"status_event_lid": {
				"getStatus": {
					"": "on",
				},
			},
		})
		if err != nil {
			t.Error(err)
			return
		}

		devices := map[string]map[string]interface{}{
			"testOwner": {
				"/devices/status_event": model.Device{
					Id:           "status_event",
					LocalId:      "status_event_lid",
					Name:         "status_event",
					DeviceTypeId: "urn:infai:ses:device-type:status_event",
				},
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
					LocalId:      "color_event_lid",
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

		config, err = mqttEnv(config, ctx, wg)
		if err != nil {
			t.Error(err)
			return
		}
		time.Sleep(1 * time.Second)

		mqttClient, err := mqtt.New(ctx, config.MgwMqttBroker, "testMgwCommand-client", "", "")
		if err != nil {
			t.Error(err)
			return
		}

		err = mqttClient.Subscribe("#", 2, func(topic string, payload []byte) {
			t.Log(topic, string(payload))
			if strings.HasPrefix(topic, "command/") {
				msg := mgw.Command{}
				err = json.Unmarshal(payload, &msg)
				if err != nil {
					t.Error(err)
					return
				}
				parts := strings.Split(topic, "/")
				deviceLocalId := parts[1]
				serviceLocalId := parts[2]
				t.Log(deviceLocalId, serviceLocalId)
				switch serviceLocalId {
				case "113-1-5:get":
					msg.Data = `{"value": "clear", "lastUpdate": 0, "lastUpdate_unit": "unit"}`
				case "67-1-1":
					if msg.Data != "21" {
						t.Error(msg.Data)
					}
					msg.Data = ""
				case "49-1-1:get":
					msg.Data = `{"value": 13, "lastUpdate": 42}`
				case "67-1-1:get":
					//create timeout
					return
				case "128-1-0:get":
					//create timeout
					return
				case "getStatus":
					msg.Data = `{"brightness": 65, "hue": 176, "saturation": 70, "kelvin": 0, "on": true, "status": 200}`
				case "setColor":
					input := map[string]float64{}
					err = json.Unmarshal([]byte(msg.Data), &input)
					if err != nil {
						t.Error(err)
						return
					}
					//values are truncated to integers, not rounded
					if input["brightness"] != float64(65) {
						t.Error(msg.Data)
					}
					if input["hue"] != float64(176) {
						t.Error(msg.Data)
					}
					if input["saturation"] != float64(70) {
						t.Error(msg.Data)
					}
					if input["duration"] != float64(1) {
						t.Error(msg.Data)
					}
					msg.Data = ""
				default:
					t.Error("unknown service-id", serviceLocalId)
					return
				}
				resp, err := json.Marshal(msg)
				if err != nil {
					t.Error(err)
					return
				}
				err = mqttClient.Publish(strings.Join([]string{"response", deviceLocalId, serviceLocalId}, "/"), 2, false, resp)
				if err != nil {
					t.Error(err)
					return
				}
			}
		})

		cmd, err := command.New(ctx, config)
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

		t.Run("device setTemperature", sendCommand(config, command.CommandMessage{
			FunctionId: "urn:infai:ses:controlling-function:99240d90-02dd-4d4f-a47c-069cfe77629c",
			Input:      21,
			DeviceId:   "urn:infai:ses:device:a486084b-3323-4cbc-9f6b-d797373ae866",
			ServiceId:  "urn:infai:ses:service:4932d451-3300-4a22-a508-ec740e5789b3",
		}, 200, "[null]"))

		t.Run("device getTemperature", sendCommand(config, command.CommandMessage{
			FunctionId: "urn:infai:ses:measuring-function:f2769eb9-b6ad-4f7e-bd28-e4ea043d2f8b",
			DeviceId:   "urn:infai:ses:device:a486084b-3323-4cbc-9f6b-d797373ae866",
			ServiceId:  "urn:infai:ses:service:6d6067a3-ed4e-45ec-a7eb-b1695340d2f1",
		}, 200, "[13]"))

		t.Run("device timeout", sendCommand(config, command.CommandMessage{
			FunctionId: "urn:infai:ses:measuring-function:00549f18-88b5-44c7-adb1-f558e8d53d1d",
			DeviceId:   "urn:infai:ses:device:a486084b-3323-4cbc-9f6b-d797373ae866",
			ServiceId:  "urn:infai:ses:service:36fd778e-b04d-4d72-bed5-1b77ed1164b9",
		}, http.StatusRequestTimeout, `"timeout"`))

		t.Run("invalid command", sendCommand(config, command.CommandMessage{
			FunctionId: "foobar",
			DeviceId:   "urn:infai:ses:device:a486084b-3323-4cbc-9f6b-d797373ae866",
			ServiceId:  "urn:infai:ses:service:6d6067a3-ed4e-45ec-a7eb-b1695340d2f1",
		}, 500, `"unable to load function: value not found in fallback: function.foobar"`))

		t.Run("device color", sendCommand(config, command.CommandMessage{
			FunctionId: "urn:infai:ses:controlling-function:c54e2a89-1fb8-4ecb-8993-a7b40b355599",
			Input: map[string]interface{}{
				"r": 50,
				"g": 168,
				"b": 162,
			},
			DeviceId:  "lamp",
			ServiceId: "urn:infai:ses:service:1b0ef253-16f7-4b65-8a15-fe79fccf7e70",
		}, 200, "[null]"))

		t.Run("device event color", sendCommand(config, command.CommandMessage{
			FunctionId: "urn:infai:ses:measuring-function:bdb6a7c8-4a3d-4fe0-bab3-ce02e09b5869",
			DeviceId:   "color_event",
			ServiceId:  "urn:infai:ses:service:color_event",
		}, 200, `[{"b":158,"g":166,"r":50}]`))

		t.Run("device event on", sendCommand(config, command.CommandMessage{
			FunctionId: "urn:infai:ses:measuring-function:20d3c1d3-77d7-4181-a9f3-b487add58cd0",
			DeviceId:   "color_event",
			ServiceId:  "urn:infai:ses:service:color_event",
		}, 200, `[true]`))

		//some services are called as event (timescale call), some as request
		t.Run("device group color", sendCommand(config, command.CommandMessage{
			FunctionId: "urn:infai:ses:measuring-function:bdb6a7c8-4a3d-4fe0-bab3-ce02e09b5869",
			GroupId:    "group_color",
		}, 200, `[{"b":158,"g":166,"r":50},{"b":158,"g":166,"r":50},{"b":158,"g":166,"r":50}]`))

		//some services return a timeout, some return 13
		t.Run("device group getTemperature", sendCommand(config, command.CommandMessage{
			FunctionId: "urn:infai:ses:measuring-function:f2769eb9-b6ad-4f7e-bd28-e4ea043d2f8b",
			GroupId:    "group_temperature",
		}, 200, "[13,13,13]"))

		t.Run("device batch", sendCommandBatch(config, command.BatchRequest{
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
		}, 200, `[{"status_code":200,"message":[null]},{"status_code":200,"message":[13]},{"status_code":408,"message":"timeout"},{"status_code":500,"message":"unable to load function: value not found in fallback: function.foobar"},{"status_code":200,"message":[null]},{"status_code":200,"message":[{"b":158,"g":166,"r":50}]},{"status_code":200,"message":[true]}]`))

		zeroTimestamp := time.UnixMilli(0).Format(time.RFC3339)

		t.Run("new timestamp", sendCommandBatch(config, command.BatchRequest{
			{
				FunctionId: "urn:infai:ses:measuring-function:3b4e0766-0d67-4658-b249-295902cd3290",
				DeviceId:   "urn:infai:ses:device:timestamp-test",
				ServiceId:  "urn:infai:ses:service:ec456e2a-81ed-4466-a119-daecfbb2d033",
			},
		}, 200, fmt.Sprintf(`[{"status_code":200,"message":["%v"]}]`, zeroTimestamp)))

		t.Run("device group air getTemperature", sendCommand(config, command.CommandMessage{
			FunctionId: "urn:infai:ses:measuring-function:f2769eb9-b6ad-4f7e-bd28-e4ea043d2f8b",
			GroupId:    "group_temperature",
			AspectId:   "urn:infai:ses:aspect:a14c5efb-b0b6-46c3-982e-9fded75b5ab6",
		}, 200, "[13,13,13]"))

		t.Run("device group outside air getTemperature", sendCommand(config, command.CommandMessage{
			FunctionId: "urn:infai:ses:measuring-function:f2769eb9-b6ad-4f7e-bd28-e4ea043d2f8b",
			GroupId:    "group_temperature",
			AspectId:   "urn:infai:ses:aspect:outside_air",
		}, 200, "[13,13,13]"))

		t.Run("device group outside foo-aspect getTemperature", sendCommand(config, command.CommandMessage{
			FunctionId: "urn:infai:ses:measuring-function:f2769eb9-b6ad-4f7e-bd28-e4ea043d2f8b",
			GroupId:    "group_temperature",
			AspectId:   "urn:infai:ses:aspect:foo-aspect",
		}, 200, "[]"))

		t.Run("device getTemperature with aspect", sendCommand(config, command.CommandMessage{
			FunctionId: "urn:infai:ses:measuring-function:f2769eb9-b6ad-4f7e-bd28-e4ea043d2f8b",
			DeviceId:   "urn:infai:ses:device:a486084b-3323-4cbc-9f6b-d797373ae866",
			ServiceId:  "urn:infai:ses:service:6d6067a3-ed4e-45ec-a7eb-b1695340d2f1",
			AspectId:   "urn:infai:ses:aspect:a14c5efb-b0b6-46c3-982e-9fded75b5ab6",
		}, 200, "[13]"))

		t.Run("device event status", sendCommand(config, command.CommandMessage{
			FunctionId: "urn:infai:ses:measuring-function:20d3c1d3-77d7-4181-a9f3-b487add58cd0",
			DeviceId:   "status_event",
			ServiceId:  "urn:infai:ses:service:status_event",
		}, 200, `[true]`))
	}
}

func TestMgwPlainTextCommandWithDockerTimescale(t *testing.T) {
	wg := &sync.WaitGroup{}
	defer wg.Wait()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	config, err := configuration.Load("../../config.json")
	if err != nil {
		t.Error(err)
		return
	}
	config.Debug = true
	config.ComImpl = "mgw"
	config.MarshallerImpl = "mgw"
	config.UseIotFallback = true
	config.TimescaleImpl = "mgw"
	config.IotFallbackFile = filepath.Join(t.TempDir(), "iot_fallback.json")

	config.ServerPort, err = GetFreePort()
	if err != nil {
		t.Error(err)
		return
	}

	devices := map[string]map[string]interface{}{
		"testOwner": {
			"/devices/status_event": model.Device{
				Id:           "status_event",
				LocalId:      "status_event_lid",
				Name:         "status_event",
				DeviceTypeId: "urn:infai:ses:device-type:status_event",
			},
			"/devices/status_event_2": model.Device{
				Id:           "status_event_2",
				LocalId:      "status_event_lid_2",
				Name:         "status_event 2",
				DeviceTypeId: "urn:infai:ses:device-type:status_event",
			},
		},
	}

	config, err = iotEnv(config, ctx, wg, devices)
	if err != nil {
		t.Error(err)
		return
	}

	config, err = mqttEnv(config, ctx, wg)
	if err != nil {
		t.Error(err)
		return
	}
	time.Sleep(1 * time.Second)

	config, err = lastValueEnv(config, ctx, wg)
	if err != nil {
		t.Error(err)
		return
	}

	time.Sleep(2 * time.Second)

	mqttClient, err := mqtt.New(ctx, config.MgwMqttBroker, "testMgwCommand-client", "", "")
	if err != nil {
		t.Error(err)
		return
	}

	err = mqttClient.Publish("event/status_event_lid/getStatus", 2, true, []byte("on"))
	if err != nil {
		t.Error(err)
		return
	}

	err = mqttClient.Subscribe("#", 2, func(topic string, payload []byte) {
		t.Log(topic, string(payload))
		if strings.HasPrefix(topic, "command/") {
			msg := mgw.Command{}
			err = json.Unmarshal(payload, &msg)
			if err != nil {
				t.Error(err)
				return
			}
			parts := strings.Split(topic, "/")
			deviceLocalId := parts[1]
			serviceLocalId := parts[2]
			t.Log(deviceLocalId, serviceLocalId)
			switch serviceLocalId {
			case "":
				t.Error("wtf")
			default:
				t.Error("unknown service-id", serviceLocalId)
				return
			}
			resp, err := json.Marshal(msg)
			if err != nil {
				t.Error(err)
				return
			}
			err = mqttClient.Publish(strings.Join([]string{"response", deviceLocalId, serviceLocalId}, "/"), 2, false, resp)
			if err != nil {
				t.Error(err)
				return
			}
		}
	})

	cmd, err := command.New(ctx, config)
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

	t.Run("known", sendCommand(config, command.CommandMessage{
		FunctionId: "urn:infai:ses:measuring-function:20d3c1d3-77d7-4181-a9f3-b487add58cd0",
		DeviceId:   "status_event",
		ServiceId:  "urn:infai:ses:service:status_event",
	}, 200, `[true]`))

	t.Run("unknown", sendCommand(config, command.CommandMessage{
		FunctionId: "urn:infai:ses:measuring-function:20d3c1d3-77d7-4181-a9f3-b487add58cd0",
		DeviceId:   "status_event_2",
		ServiceId:  "urn:infai:ses:service:status_event",
	}, 513, `"unable to get event value: missing last value in mgw-last-value"`))

	t.Run("batch", sendCommandBatch(config, command.BatchRequest{
		{
			FunctionId: "urn:infai:ses:measuring-function:20d3c1d3-77d7-4181-a9f3-b487add58cd0",
			DeviceId:   "status_event",
			ServiceId:  "urn:infai:ses:service:status_event",
		},
		{
			FunctionId: "urn:infai:ses:measuring-function:20d3c1d3-77d7-4181-a9f3-b487add58cd0",
			DeviceId:   "status_event_2",
			ServiceId:  "urn:infai:ses:service:status_event",
		},
	}, 200, `[{"status_code":200,"message":[true]},{"status_code":513,"message":"unable to get event value: missing last value in mgw-last-value"}]`))

}

func TestShellyError(t *testing.T) {
	wg := &sync.WaitGroup{}
	defer wg.Wait()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	config, err := configuration.Load("../../config.json")
	if err != nil {
		t.Error(err)
		return
	}
	config.Debug = true
	config.ComImpl = "mgw"
	config.MarshallerImpl = "mgw"
	config.UseIotFallback = true
	config.TimescaleImpl = "mgw"
	config.IotFallbackFile = filepath.Join(t.TempDir(), "iot_fallback.json")

	config.ServerPort, err = GetFreePort()
	if err != nil {
		t.Error(err)
		return
	}

	devices := map[string]map[string]interface{}{
		"testOwner": {
			"/devices/status_event": model.Device{
				Id:           "status_event",
				LocalId:      "status_event_lid",
				Name:         "status_event",
				DeviceTypeId: "urn:infai:ses:device-type:1d5375f0-7d7f-46ab-956e-1d7e7ef51826",
			},
			"/devices/status_event_2": model.Device{
				Id:           "status_event_2",
				LocalId:      "status_event_lid_2",
				Name:         "status_event 2",
				DeviceTypeId: "urn:infai:ses:device-type:1d5375f0-7d7f-46ab-956e-1d7e7ef51826",
			},
		},
	}

	config, err = iotEnv(config, ctx, wg, devices)
	if err != nil {
		t.Error(err)
		return
	}

	config, err = mqttEnv(config, ctx, wg)
	if err != nil {
		t.Error(err)
		return
	}
	time.Sleep(1 * time.Second)

	config, err = lastValueEnv(config, ctx, wg)
	if err != nil {
		t.Error(err)
		return
	}

	time.Sleep(2 * time.Second)

	mqttClient, err := mqtt.New(ctx, config.MgwMqttBroker, "testMgwCommand-client", "", "")
	if err != nil {
		t.Error(err)
		return
	}

	err = mqttClient.Publish("event/status_event_lid/relay", 2, true, []byte("on"))
	if err != nil {
		t.Error(err)
		return
	}

	err = mqttClient.Subscribe("#", 2, func(topic string, payload []byte) {
		t.Log(topic, string(payload))
		if strings.HasPrefix(topic, "command/") {
			msg := mgw.Command{}
			err = json.Unmarshal(payload, &msg)
			if err != nil {
				t.Error(err)
				return
			}
			parts := strings.Split(topic, "/")
			deviceLocalId := parts[1]
			serviceLocalId := parts[2]
			t.Log(deviceLocalId, serviceLocalId)
			switch serviceLocalId {
			case "":
				t.Error("wtf")
			default:
				t.Error("unknown service-id", serviceLocalId)
				return
			}
			resp, err := json.Marshal(msg)
			if err != nil {
				t.Error(err)
				return
			}
			err = mqttClient.Publish(strings.Join([]string{"response", deviceLocalId, serviceLocalId}, "/"), 2, false, resp)
			if err != nil {
				t.Error(err)
				return
			}
		}
	})

	cmd, err := command.New(ctx, config)
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

	t.Run("known", sendCommand(config, command.CommandMessage{
		FunctionId: "urn:infai:ses:measuring-function:20d3c1d3-77d7-4181-a9f3-b487add58cd0",
		DeviceId:   "status_event",
		ServiceId:  "urn:infai:ses:service:940fd269-27f7-4f2e-afbf-eddbf0feb4c8",
	}, 200, `[true]`))

	t.Run("unknown", sendCommand(config, command.CommandMessage{
		FunctionId: "urn:infai:ses:measuring-function:20d3c1d3-77d7-4181-a9f3-b487add58cd0",
		DeviceId:   "status_event_2",
		ServiceId:  "urn:infai:ses:service:940fd269-27f7-4f2e-afbf-eddbf0feb4c8",
	}, 513, `"unable to get event value: missing last value in mgw-last-value"`))

	t.Run("batch", sendCommandBatch(config, command.BatchRequest{
		{
			FunctionId: "urn:infai:ses:measuring-function:20d3c1d3-77d7-4181-a9f3-b487add58cd0",
			DeviceId:   "status_event",
			ServiceId:  "urn:infai:ses:service:940fd269-27f7-4f2e-afbf-eddbf0feb4c8",
		},
		{
			FunctionId: "urn:infai:ses:measuring-function:20d3c1d3-77d7-4181-a9f3-b487add58cd0",
			DeviceId:   "status_event_2",
			ServiceId:  "urn:infai:ses:service:940fd269-27f7-4f2e-afbf-eddbf0feb4c8",
		},
	}, 200, `[{"status_code":200,"message":[true]},{"status_code":513,"message":"unable to get event value: missing last value in mgw-last-value"}]`))

}

func TestCharacteristicError(t *testing.T) {
	iotExport = iotExport2
	defer func() {
		iotExport = iotExport1
	}()
	wg := &sync.WaitGroup{}
	defer wg.Wait()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	config, err := configuration.Load("../../config.json")
	if err != nil {
		t.Error(err)
		return
	}
	config.Debug = true
	config.ComImpl = "mgw"
	config.MarshallerImpl = "mgw"
	config.UseIotFallback = true
	config.TimescaleImpl = "mgw"
	config.IotFallbackFile = filepath.Join(t.TempDir(), "iot_fallback.json")

	config.ServerPort, err = GetFreePort()
	if err != nil {
		t.Error(err)
		return
	}

	devices := map[string]map[string]interface{}{
		"testOwner": {
			"/devices/urn:infai:ses:device:79fa231c-45d9-4266-b5fb-2c051bfd8c0d": model.Device{
				Id:           "urn:infai:ses:device:79fa231c-45d9-4266-b5fb-2c051bfd8c0d",
				LocalId:      "d1",
				Name:         "device",
				DeviceTypeId: "urn:infai:ses:device-type:ca30a161-0bd4-49b8-86eb-8c48e29eb34e",
			},
		},
	}

	config, err = iotEnv(config, ctx, wg, devices)
	if err != nil {
		t.Error(err)
		return
	}

	config, err = mqttEnv(config, ctx, wg)
	if err != nil {
		t.Error(err)
		return
	}
	time.Sleep(1 * time.Second)

	config, err = lastValueEnv(config, ctx, wg)
	if err != nil {
		t.Error(err)
		return
	}

	time.Sleep(2 * time.Second)

	mqttClient, err := mqtt.New(ctx, config.MgwMqttBroker, "testMgwCommand-client", "", "")
	if err != nil {
		t.Error(err)
		return
	}

	err = mqttClient.Publish("event/status_event_lid/relay", 2, true, []byte("on"))
	if err != nil {
		t.Error(err)
		return
	}

	err = mqttClient.Subscribe("#", 2, func(topic string, payload []byte) {
		t.Log(topic, string(payload))
		if strings.HasPrefix(topic, "command/") {
			msg := mgw.Command{}
			err = json.Unmarshal(payload, &msg)
			if err != nil {
				t.Error(err)
				return
			}
			parts := strings.Split(topic, "/")
			deviceLocalId := parts[1]
			serviceLocalId := parts[2]
			t.Log(deviceLocalId, serviceLocalId)
			switch serviceLocalId {
			case "31712f59-1f7e-445e-8bf1-523cdc5c3a96":
				msg.Data = ""
			case "":
				t.Error("wtf")
			default:
				t.Error("unknown service-id", serviceLocalId)
				return
			}
			resp, err := json.Marshal(msg)
			if err != nil {
				t.Error(err)
				return
			}
			err = mqttClient.Publish(strings.Join([]string{"response", deviceLocalId, serviceLocalId}, "/"), 2, false, resp)
			if err != nil {
				t.Error(err)
				return
			}
		}
	})

	cmd, err := command.New(ctx, config)
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

	t.Run("batch", sendCommandBatch(config, command.BatchRequest{
		{
			FunctionId: "urn:infai:ses:controlling-function:ced44f01-7328-43e3-8db0-ecd12f448758",
			DeviceId:   "urn:infai:ses:device:79fa231c-45d9-4266-b5fb-2c051bfd8c0d",
			ServiceId:  "urn:infai:ses:service:31712f59-1f7e-445e-8bf1-523cdc5c3a96",
			AspectId:   "urn:infai:ses:aspect:4f6b9747-1d30-4db1-bbf8-c5904db32771",
			Input:      map[string]interface{}{"iterations": 1, "customOrder": true, "segment_ids": []string{"1"}},
		},
	}, 200, `[{"status_code":200,"message":[null]}]`))

}
