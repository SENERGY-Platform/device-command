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
	"github.com/SENERGY-Platform/device-command/pkg/command"
	"github.com/SENERGY-Platform/device-command/pkg/configuration"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func init() {
	endpoints = append(endpoints, CommandEndpoints)
}

func CommandEndpoints(config configuration.Config, router *httprouter.Router, cmd Command) {
	router.POST("/commands", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		preferEventValueStr := request.URL.Query().Get("prefer_event_value")
		preferEventValue := false
		if preferEventValueStr != "" {
			preferEventValue, err = strconv.ParseBool(preferEventValueStr)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusBadRequest)
				return
			}
		}
		timeout := request.URL.Query().Get("timeout")
		msg := command.CommandMessage{}
		err = json.NewDecoder(request.Body).Decode(&msg)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		if msg.FunctionId == "" {
			http.Error(writer, "expect function_id in body", http.StatusBadRequest)
			return
		}

		code, result := cmd.Command(token, msg, timeout, preferEventValue)
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		writer.WriteHeader(code)
		json.NewEncoder(writer).Encode(result)
	})

	router.POST("/commands/batch", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		preferEventValueStr := request.URL.Query().Get("prefer_event_value")
		preferEventValue := false
		if preferEventValueStr != "" {
			preferEventValue, err = strconv.ParseBool(preferEventValueStr)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusBadRequest)
				return
			}
		}
		timeout := request.URL.Query().Get("timeout")
		batch := command.BatchRequest{}
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

		result := cmd.Batch(token, batch, timeout, preferEventValue)
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(writer).Encode(result)
		return
	})
}
