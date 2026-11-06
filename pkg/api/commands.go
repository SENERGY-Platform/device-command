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
	"fmt"
	"net/http"
	"strconv"

	"github.com/SENERGY-Platform/device-command/pkg/auth"
	"github.com/SENERGY-Platform/device-command/pkg/command"
	"github.com/SENERGY-Platform/device-command/pkg/configuration"
	"github.com/SENERGY-Platform/device-command/pkg/metrics"
)

func commandHandler(config configuration.Config, cmd Command, metrics *metrics.Metrics) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		token, preferEventValue, ok := readRequest(config, writer, request, metrics)
		if !ok {
			return
		}
		msg := command.CommandMessage{}
		err := json.NewDecoder(request.Body).Decode(&msg)
		if err != nil {
			respondError(config, writer, request, token, http.StatusBadRequest, err.Error())
			return
		}
		if msg.FunctionId == "" {
			respondError(config, writer, request, token, http.StatusBadRequest, "expect function_id in body")
			return
		}
		code, result := cmd.Command(request.Context(), token, msg, request.URL.Query().Get("timeout"), preferEventValue)
		if code != http.StatusOK {
			config.GetLogger().Warn("error response", "request-url", request.URL.String(), "user", token.GetUserId(), "response-status-code", code, "response-body", fmt.Sprintf("%#v", result))
		}
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		writer.WriteHeader(code)
		json.NewEncoder(writer).Encode(result)
	})
}

func batchHandler(config configuration.Config, cmd Command, metrics *metrics.Metrics) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		token, preferEventValue, ok := readRequest(config, writer, request, metrics)
		if !ok {
			return
		}
		batch := command.BatchRequest{}
		err := json.NewDecoder(request.Body).Decode(&batch)
		if err != nil {
			respondError(config, writer, request, token, http.StatusBadRequest, err.Error())
			return
		}
		err = batch.Validate()
		if err != nil {
			respondError(config, writer, request, token, http.StatusBadRequest, err.Error())
			return
		}
		result := cmd.Batch(request.Context(), token, batch, request.URL.Query().Get("timeout"), preferEventValue)
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(writer).Encode(result)
	})
}

// readRequest reads what both endpoints take besides the body: the token and the
// prefer_event_value query parameter. It answers the request itself if one of them is invalid.
func readRequest(config configuration.Config, writer http.ResponseWriter, request *http.Request, metrics *metrics.Metrics) (token auth.Token, preferEventValue bool, ok bool) {
	token, err := auth.GetParsedToken(request)
	if err != nil {
		respondError(config, writer, request, token, http.StatusBadRequest, err.Error())
		return token, false, false
	}
	metrics.LogRequest(token.GetUserId(), request.Pattern)
	if value := request.URL.Query().Get("prefer_event_value"); value != "" {
		preferEventValue, err = strconv.ParseBool(value)
		if err != nil {
			respondError(config, writer, request, token, http.StatusBadRequest, err.Error())
			return token, false, false
		}
	}
	return token, preferEventValue, true
}

func respondError(config configuration.Config, writer http.ResponseWriter, request *http.Request, token auth.Token, code int, message string) {
	config.GetLogger().Warn("error response", "request-url", request.URL.String(), "user", token.GetUserId(), "response-status-code", code, "response-body", message)
	http.Error(writer, message, code)
}
