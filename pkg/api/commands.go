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

package api

import (
	"encoding/json"
	"github.com/SENERGY-Platform/device-command/pkg/auth"
	"github.com/SENERGY-Platform/device-command/pkg/configuration"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"sync"
)

func init() {
	endpoints = append(endpoints, CommandEndpoints)
}

func CommandEndpoints(config configuration.Config, router *httprouter.Router, command Command) {
	router.POST("/commands", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		cmd := CommandMessage{}
		err = json.NewDecoder(request.Body).Decode(&cmd)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		if cmd.FunctionId == "" {
			http.Error(writer, "expect function_id in body", http.StatusBadRequest)
			return
		}

		if cmd.DeviceId != "" && cmd.ServiceId != "" {
			code, result := command.DeviceCommand(token, cmd.DeviceId, cmd.ServiceId, cmd.FunctionId, cmd.Input)
			writer.Header().Set("Content-Type", "application/json; charset=utf-8")
			writer.WriteHeader(code)
			json.NewEncoder(writer).Encode(result)
			return
		}

		if cmd.GroupId != "" {
			code, result := command.GroupCommand(token, cmd.GroupId, cmd.FunctionId, cmd.AspectId, cmd.DeviceClassId, cmd.Input)
			writer.Header().Set("Content-Type", "application/json; charset=utf-8")
			writer.WriteHeader(code)
			json.NewEncoder(writer).Encode(result)
			return
		}

		http.Error(writer, "missing device_id, service_id or group_id", http.StatusBadRequest)

	})

	router.POST("/commands/batch", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		batch := BatchRequest{}
		err = json.NewDecoder(request.Body).Decode(&batch)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		err = batch.Validate()
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		result := runBatch(token, command, batch)
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(writer).Encode(result)
		return
	})
}

func runBatch(token auth.Token, command Command, batch BatchRequest) []BatchResultElement {
	if len(batch) == 0 {
		return []BatchResultElement{}
	}

	result := make([]BatchResultElement, len(batch))
	wg := sync.WaitGroup{}
	mux := sync.Mutex{}

	for i, cmd := range batch {
		wg.Add(1)
		go func(i int, cmd CommandMessage) {
			defer wg.Done()
			if cmd.DeviceId != "" && cmd.ServiceId != "" {
				code, temp := command.DeviceCommand(token, cmd.DeviceId, cmd.ServiceId, cmd.FunctionId, cmd.Input)
				mux.Lock()
				defer mux.Unlock()
				result[i] = BatchResultElement{
					StatusCode: code,
					Message:    temp,
				}
				return
			}

			if cmd.GroupId != "" {
				code, temp := command.GroupCommand(token, cmd.GroupId, cmd.FunctionId, cmd.AspectId, cmd.DeviceClassId, cmd.Input)
				mux.Lock()
				defer mux.Unlock()
				result[i] = BatchResultElement{
					StatusCode: code,
					Message:    temp,
				}
				return
			}
		}(i, cmd)
	}
	wg.Wait()
	return result
}
