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
)

func init() {
	endpoints = append(endpoints, CommandEndpoints)
}

type CommandMessage struct {
	FunctionId string      `json:"function_id"` //mandatory
	Input      interface{} `json:"input"`

	//device command
	DeviceId  string `json:"device_id,omitempty"`
	ServiceId string `json:"service_id,omitempty"`

	//group command
	GroupId       string `json:"group_id,omitempty"`
	AspectId      string `json:"aspect_id,omitempty"`
	DeviceClassId string `json:"device_class_id,omitempty"`
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
}
