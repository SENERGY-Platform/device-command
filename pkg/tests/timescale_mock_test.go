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
	"github.com/SENERGY-Platform/device-command/pkg/configuration"
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
	"time"
)

type TimescaleMockRequest struct {
	DeviceId   string
	ServiceId  string
	ColumnName string
}

type TimescaleMockResponse struct {
	Time  *string     `json:"time"`
	Value interface{} `json:"value"`
}

func timescaleEnv(initialConfig configuration.Config, ctx context.Context, wg *sync.WaitGroup, values map[string]map[string]interface{}) (config configuration.Config, err error) {
	config = initialConfig

	get := func(req TimescaleMockRequest) (value interface{}) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				value = nil
			}
		}()
		log.Printf("timescaleEnv get d=%v s=%v\n", req.DeviceId, req.ServiceId)
		value = values[req.DeviceId][req.ServiceId]
		return value
	}

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		msg := []TimescaleMockRequest{}
		err = json.NewDecoder(request.Body).Decode(&msg)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		result := []TimescaleMockResponse{}

		now := time.Now().String()
		for _, req := range msg {
			if req.ColumnName != "" {
				http.Error(writer, "test expects empty column in mgw last value request", http.StatusBadRequest)
				return
			}
			result = append(result, TimescaleMockResponse{
				Value: get(req),
				Time:  &now,
			})
		}

		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(writer).Encode(result)
		return
	}))
	wg.Add(1)
	go func() {
		<-ctx.Done()
		server.Close()
		wg.Done()
	}()

	config.TimescaleWrapperUrl = server.URL

	return config, nil
}

type LastMessageResponse struct {
	Time  string                 `json:"time"`
	Value map[string]interface{} `json:"value"`
}

func timescaleCloudEnv(initialConfig configuration.Config, ctx context.Context, wg *sync.WaitGroup, values map[string]map[string]map[string]interface{}) (config configuration.Config, err error) {
	config = initialConfig

	getLastMessage := func(deviceId string, serviceId string) (value LastMessageResponse) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				value = LastMessageResponse{}
			}
		}()
		return LastMessageResponse{
			Time:  time.Now().UTC().Format(time.RFC3339),
			Value: values[deviceId][serviceId],
		}
	}

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		deviceId := request.URL.Query().Get("device_id")
		serviceId := request.URL.Query().Get("service_id")
		result := getLastMessage(deviceId, serviceId)
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(writer).Encode(result)
		return
	}))
	wg.Add(1)
	go func() {
		<-ctx.Done()
		server.Close()
		wg.Done()
	}()

	config.TimescaleWrapperUrl = server.URL

	return config, nil
}
