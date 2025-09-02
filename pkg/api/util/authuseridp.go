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
	"net/http"

	"github.com/SENERGY-Platform/device-command/pkg/auth"
	"github.com/SENERGY-Platform/device-command/pkg/configuration"
)

func NewAuthIdp(config configuration.Config, handler http.Handler) (http.Handler, error) {
	return &AuthIdp{
		handler: handler,
		config:  config,
		auth:    &auth.OpenidToken{},
	}, nil
}

type AuthIdp struct {
	handler http.Handler
	auth    *auth.OpenidToken
	config  configuration.Config
}

func (this *AuthIdp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if this.handler == nil {
		http.Error(w, "Missing Handler in AuthIdp", 500)
	}
	token, err := this.GetToken()
	if err != nil {
		http.Error(w, "unable to use AuthIdp: "+err.Error(), 500)
		return
	}
	r.Header.Set("Authorization", token)
	this.handler.ServeHTTP(w, r)
}

func (this *AuthIdp) GetToken() (token string, err error) {
	return this.auth.EnsureAccess(this.config)
}
