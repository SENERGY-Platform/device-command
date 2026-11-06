/*
 * Copyright 2019 InfAI (CC SES)
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
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/SENERGY-Platform/device-command/pkg/api/util"
	"github.com/SENERGY-Platform/device-command/pkg/auth"
	"github.com/SENERGY-Platform/device-command/pkg/command"
	"github.com/SENERGY-Platform/device-command/pkg/configuration"
	"github.com/SENERGY-Platform/device-command/pkg/metrics"
	"github.com/SENERGY-Platform/service-commons/pkg/accesslog"
)

type Command interface {
	Command(ctx context.Context, token auth.Token, cmd command.CommandMessage, timeout string, preferEventValue bool) (code int, resp interface{})
	Batch(ctx context.Context, token auth.Token, batch command.BatchRequest, timeout string, preferEventValue bool) []command.BatchResultElement
}

func Start(ctx context.Context, config configuration.Config, cmd Command, metrics *metrics.Metrics) error {
	handler, err := NewHandler(config, cmd, metrics)
	if err != nil {
		return err
	}
	server := &http.Server{Addr: ":" + config.ServerPort, Handler: handler}
	go func() {
		config.GetLogger().Info("listening", "address", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			config.GetLogger().Error("FATAL: unable to start http server", "error", err)
			log.Fatal("FATAL:", err)
		}
	}()
	go func() {
		<-ctx.Done()
		config.GetLogger().Info("api shutdown", "result", server.Shutdown(context.Background()))
	}()
	return nil
}

func NewHandler(config configuration.Config, cmd Command, metrics *metrics.Metrics) (handler http.Handler, err error) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})
	mux.Handle("GET /metrics", metrics)
	mux.Handle("POST /commands", commandHandler(config, cmd, metrics))
	mux.Handle("POST /commands/batch", batchHandler(config, cmd, metrics))

	handler = mux
	switch {
	case config.RequestUserIdp == "jwt":
	case config.RequestUserIdp == "auth":
		handler, err = util.NewAuthIdp(config, handler)
	case strings.HasPrefix(config.RequestUserIdp, "user:"):
		handler, err = util.NewPresetUserIdp(config, handler)
	case strings.HasPrefix(config.RequestUserIdp, "mgw:"):
		handler, err = util.NewMgwRequestUserIdp(config, handler)
	default:
		return nil, errors.New("unknown request_user_idp configured")
	}
	if err != nil {
		return nil, err
	}
	handler = util.NewVersionHeaderMiddleware(handler)
	handler = util.NewCors(handler)
	handler = accesslog.New(handler)
	return handler, nil
}
