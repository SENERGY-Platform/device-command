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

package util

import (
	"github.com/SENERGY-Platform/device-command/pkg/auth"
	"github.com/SENERGY-Platform/device-command/pkg/configuration"
	"net/http"
)

func NewTokenOverwrite(config configuration.Config, handler http.Handler) http.Handler {
	return &TokenOverwriteMiddleWare{
		handler: handler,
		config:  config,
		auth:    &auth.OpenidToken{},
	}
}

type TokenOverwriteMiddleWare struct {
	handler http.Handler
	config  configuration.Config
	auth    *auth.OpenidToken
}

func (this *TokenOverwriteMiddleWare) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if this.handler != nil {
		token, err := this.auth.EnsureAccess(this.config)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		r.Header.Set("Authorization", token)
		this.handler.ServeHTTP(w, r)
	} else {
		http.Error(w, "Forbidden", 403)
	}
}
