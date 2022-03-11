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
		timeout := request.URL.Query().Get("timeout")
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
			code, result := command.DeviceCommand(token, cmd.DeviceId, cmd.ServiceId, cmd.FunctionId, cmd.AspectId, cmd.Input, timeout)
			writer.Header().Set("Content-Type", "application/json; charset=utf-8")
			writer.WriteHeader(code)
			json.NewEncoder(writer).Encode(result)
			return
		}

		if cmd.GroupId != "" {
			code, result := command.GroupCommand(token, cmd.GroupId, cmd.FunctionId, cmd.AspectId, cmd.DeviceClassId, cmd.Input, timeout)
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
		timeout := request.URL.Query().Get("timeout")
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

		result := runBatch(token, command, batch, timeout)
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(writer).Encode(result)
		return
	})
}

func runBatch(token auth.Token, command Command, batch BatchRequest, timeout string) []BatchResultElement {
	if len(batch) == 0 {
		return []BatchResultElement{}
	}

	result := make([]BatchResultElement, len(batch))
	wg := sync.WaitGroup{}
	mux := sync.Mutex{}

	resultIndexMap := map[string][]int{}
	for i, cmd := range batch {
		hash := cmd.Hash()
		resultIndexMap[hash] = append(resultIndexMap[hash], i)
	}
	isAlreadySend := map[string]bool{}

	for _, cmd := range batch {
		hash := cmd.Hash()
		if !isAlreadySend[hash] {
			isAlreadySend[hash] = true
			wg.Add(1)
			go func(resultIndexes []int, cmd CommandMessage) {
				defer wg.Done()
				var code int
				var temp interface{}
				if cmd.DeviceId != "" && cmd.ServiceId != "" {
					code, temp = command.DeviceCommand(token, cmd.DeviceId, cmd.ServiceId, cmd.FunctionId, cmd.AspectId, cmd.Input, timeout)
				} else if cmd.GroupId != "" {
					code, temp = command.GroupCommand(token, cmd.GroupId, cmd.FunctionId, cmd.AspectId, cmd.DeviceClassId, cmd.Input, timeout)
				}
				mux.Lock()
				defer mux.Unlock()
				for _, index := range resultIndexes {
					result[index] = BatchResultElement{
						StatusCode: code,
						Message:    temp,
					}
				}
				return
			}(resultIndexMap[hash], cmd)
		}
	}
	wg.Wait()
	return result
}
