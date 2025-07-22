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
	"context"
	"fmt"
	"github.com/SENERGY-Platform/device-command/pkg/configuration"
	"github.com/SENERGY-Platform/mgw-cloud-proxy/cert-manager/lib/client"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

func NewMgwRequestUserIdp(config configuration.Config, handler http.Handler) (http.Handler, error) {
	mgwUserIdUrl := strings.TrimPrefix(config.RequestUserIdp, "mgw:")
	_, err := url.Parse(mgwUserIdUrl)
	if err != nil {
		return nil, fmt.Errorf("invalid RequestUserIdp, expect 'mgw:<url>', err=%w", err)
	}
	return &MgwRequestUserIdpMiddleWare{
		handler: handler,
		config:  config,
		client:  client.New(http.DefaultClient, mgwUserIdUrl),
	}, nil
}

type MgwRequestUserIdpMiddleWare struct {
	handler               http.Handler
	config                configuration.Config
	client                *client.Client
	lastUserId            string
	lastUserIdRequestTime time.Time
	mux                   sync.Mutex
}

func (this *MgwRequestUserIdpMiddleWare) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if this.handler == nil {
		http.Error(w, "Missing Handler in MgwRequestUserIdpMiddleWare", 500)
	}
	userId, err := this.GetUserId()
	if err != nil {
		http.Error(w, "unable to use MgwRequestUserIdpMiddleWare: "+err.Error(), 500)
		return
	}
	token, err := this.GenerateUserTokenById(userId)
	if err != nil {
		http.Error(w, "unable to use MgwRequestUserIdpMiddleWare: "+err.Error(), 500)
		return
	}
	r.Header.Set("Authorization", token)
	this.handler.ServeHTTP(w, r)
}

func (this *MgwRequestUserIdpMiddleWare) GenerateUserTokenById(userid string) (token string, err error) {
	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour).Unix(),
		Issuer:    "mgw-device-command",
		Subject:   userid,
	}
	jwtoken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	unsignedTokenString, err := jwtoken.SigningString()
	if err != nil {
		log.Println("ERROR: GenerateUserTokenById::SigningString()", err, userid)
		return token, err
	}
	tokenString := strings.Join([]string{unsignedTokenString, ""}, ".")
	return "Bearer " + tokenString, nil
}

func (this *MgwRequestUserIdpMiddleWare) GetUserId() (userId string, err error) {
	this.mux.Lock()
	defer this.mux.Unlock()
	if this.lastUserId != "" && time.Since(this.lastUserIdRequestTime) < 10*time.Minute {
		return this.lastUserId, nil
	}
	network, err := this.client.NetworkInfo(context.Background(), "")
	if err != nil {
		return "", err
	}
	this.lastUserId = network.UserID
	this.lastUserIdRequestTime = time.Now()
	return this.lastUserId, nil
}
