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
	"github.com/SENERGY-Platform/device-command/pkg/auth"
	"github.com/SENERGY-Platform/device-command/pkg/configuration"
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
)

func iotEnv(initialConfig configuration.Config, ctx context.Context, wg *sync.WaitGroup, userSpecificEndpoints map[string]map[string]interface{}) (config configuration.Config, err error) {
	config = initialConfig

	endpoints := map[string]interface{}{}
	err = json.Unmarshal([]byte(iotExport), &endpoints)
	if err != nil {
		return config, err
	}

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		path := request.URL.Path
		log.Println("IOT: GET", path)
		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		element, ok := endpoints[path]
		log.Println("IOT: endpoint", ok, element)
		if ok {
			writer.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(writer).Encode(element)
			return
		}
		element, ok = userSpecificEndpoints[token.GetUserId()][path]
		log.Println("IOT: user endpoint", ok, element)
		if ok {
			writer.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(writer).Encode(element)
			return
		}
		log.Println("IOT: ERROR: not found: ", path)
		http.Error(writer, "not found", http.StatusNotFound)
	}))
	wg.Add(1)
	go func() {
		<-ctx.Done()
		server.Close()
		wg.Done()
	}()

	config.DeviceManagerUrl = server.URL
	config.DeviceRepositoryUrl = server.URL
	config.PermissionsUrl = server.URL

	return config, nil
}

var iotExport = iotExport1
