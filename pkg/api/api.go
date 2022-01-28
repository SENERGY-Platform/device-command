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
	"fmt"
	"github.com/SENERGY-Platform/device-command/pkg/api/util"
	"github.com/SENERGY-Platform/device-command/pkg/auth"
	"github.com/SENERGY-Platform/device-command/pkg/configuration"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"reflect"
	"runtime"
	"runtime/debug"
)

type Command interface {
	DeviceCommand(token auth.Token, deviceId string, serviceId string, functionId string, input interface{}, timeout string) (code int, resp interface{})
	GroupCommand(token auth.Token, groupId string, functionId string, aspectId string, deviceClassId string, input interface{}, timeout string) (code int, resp interface{})
}

var endpoints = []func(config configuration.Config, router *httprouter.Router, command Command){}

func Start(ctx context.Context, config configuration.Config, command Command) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprint(r))
		}
	}()
	router := GetRouter(config, command)

	server := &http.Server{Addr: ":" + config.ServerPort, Handler: router}
	go func() {
		log.Println("listening on ", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			debug.PrintStack()
			log.Fatal("FATAL:", err)
		}
	}()
	go func() {
		<-ctx.Done()
		log.Println("api shutdown", server.Shutdown(context.Background()))
	}()
	return
}

func GetRouter(config configuration.Config, command Command) http.Handler {
	router := httprouter.New()
	for _, e := range endpoints {
		log.Println("add endpoint: " + runtime.FuncForPC(reflect.ValueOf(e).Pointer()).Name())
		e(config, router, command)
	}
	handler := util.NewCors(router)
	handler = util.NewLogger(handler)
	return handler
}
