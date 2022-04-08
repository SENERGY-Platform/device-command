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
	"github.com/SENERGY-Platform/device-command/pkg/command/dependencies/interfaces"
	"github.com/SENERGY-Platform/device-command/pkg/configuration"
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
)

type TimescaleMockRequest struct {
	DeviceId   string
	ServiceId  string
	ColumnName string
}

func timescaleEnv(initialConfig configuration.Config, ctx context.Context, wg *sync.WaitGroup, values map[string]map[string]map[string]interface{}) (config configuration.Config, err error) {
	config = initialConfig

	get := func(req TimescaleMockRequest) (value interface{}) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				value = nil
			}
		}()
		value = values[req.DeviceId][req.ServiceId][req.ColumnName]
		return value
	}

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		msg := []TimescaleMockRequest{}
		err = json.NewDecoder(request.Body).Decode(&msg)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		result := []interfaces.TimescaleResponse{}

		for _, req := range msg {
			result = append(result, interfaces.TimescaleResponse{
				Value: get(req),
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

type TimescaleMockQuery = []MockQueriesRequestElement

type MockQueriesRequestElement struct {
	DeviceId  string
	ServiceId string
	Limit     int
	Columns   []MockQueriesRequestElementColumn
}

type MockQueriesRequestElementColumn struct {
	Name string
}

type MockTimescaleQueryResponse = [][][]interface{} //[query-index][0][column-index + 1]

func timescaleCloudEnv(initialConfig configuration.Config, ctx context.Context, wg *sync.WaitGroup, values map[string]map[string]map[string]interface{}) (config configuration.Config, err error) {
	config = initialConfig

	get := func(req MockQueriesRequestElement) (value [][]interface{}) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				value = nil
			}
		}()
		columnValues := []interface{}{""}
		for _, c := range req.Columns {
			columnValues = append(columnValues, values[req.DeviceId][req.ServiceId][c.Name])
		}
		value = [][]interface{}{columnValues}
		return value
	}

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		msg := TimescaleMockQuery{}
		err = json.NewDecoder(request.Body).Decode(&msg)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		result := MockTimescaleQueryResponse{}

		for _, req := range msg {
			result = append(result, get(req))
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
