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
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/SENERGY-Platform/device-command/pkg/configuration"
	"github.com/golang-jwt/jwt"
)

func NewPresetUserIdp(config configuration.Config, handler http.Handler) (http.Handler, error) {
	userId := strings.TrimPrefix(config.RequestUserIdp, "user:")
	return &PresetUserIdp{
		handler: handler,
		config:  config,
		userId:  userId,
	}, nil
}

type PresetUserIdp struct {
	handler http.Handler
	config  configuration.Config
	userId  string
	mux     sync.Mutex
}

func (this *PresetUserIdp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if this.handler == nil {
		http.Error(w, "Missing Handler in PresetUserIdp", 500)
	}
	if r.Header.Get("Authorization") != "" {
		this.handler.ServeHTTP(w, r)
		return
	}
	userId, err := this.GetUserId()
	if err != nil {
		http.Error(w, "unable to use PresetUserIdp: "+err.Error(), 500)
		return
	}
	token, err := this.GenerateUserTokenById(userId)
	if err != nil {
		http.Error(w, "unable to use PresetUserIdp: "+err.Error(), 500)
		return
	}
	r.Header.Set("Authorization", token)
	this.handler.ServeHTTP(w, r)
}

func (this *PresetUserIdp) GenerateUserTokenById(userid string) (token string, err error) {
	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour).Unix(),
		Issuer:    "mgw-device-command",
		Subject:   userid,
	}
	jwtoken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	unsignedTokenString, err := jwtoken.SigningString()
	if err != nil {
		log.Println("ERROR: PresetUserIdp::SigningString()", err, userid)
		return token, err
	}
	tokenString := strings.Join([]string{unsignedTokenString, ""}, ".")
	return "Bearer " + tokenString, nil
}

func (this *PresetUserIdp) GetUserId() (userId string, err error) {
	return this.userId, nil
}
