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

package mgw

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/SENERGY-Platform/device-command/pkg/auth"
	"github.com/SENERGY-Platform/device-command/pkg/command/dependencies/interfaces"
	"github.com/SENERGY-Platform/device-command/pkg/configuration"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Timescale struct {
	TimescaleWrapperUrl string
}

func TimescaleFactory(ctx context.Context, config configuration.Config) (interfaces.Timescale, error) {
	return NewTimescale(config.TimescaleWrapperUrl)
}

func NewTimescale(TimescaleUrl string) (*Timescale, error) {
	if !strings.Contains(TimescaleUrl, "://") {
		TimescaleUrl = "http://" + TimescaleUrl
	}
	parsed, err := url.Parse(TimescaleUrl)
	if err != nil {
		return nil, err
	}
	if parsed.Port() == "" {
		TimescaleUrl = TimescaleUrl + ":8080"
	}
	return &Timescale{TimescaleWrapperUrl: TimescaleUrl}, nil
}

func (this *Timescale) Query(token auth.Token, request []interfaces.TimescaleRequest, timeout time.Duration) (result []interfaces.TimescaleResponse, err error) {
	body := &bytes.Buffer{}
	err = json.NewEncoder(body).Encode(this.castRequest(request))
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", this.TimescaleWrapperUrl+"/last-values", body)
	if err != nil {
		return result, err
	}
	req.Header.Set("Authorization", token.Jwt())
	client := &http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("ERROR: unable to query /last-values", err)
		return result, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		temp, _ := io.ReadAll(resp.Body)
		return result, errors.New(strings.TrimSpace(string(temp)))
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	return
}

func (this *Timescale) castRequest(request []interfaces.TimescaleRequest) (result []TimescaleRequest) {
	for _, r := range request {
		result = append(result, TimescaleRequest{
			DeviceId:   r.Device.LocalId,
			ServiceId:  r.Service.LocalId,
			ColumnName: this.castPath(r.ColumnName),
		})
	}
	return
}

func (this *Timescale) castPath(path string) string {
	if path == "" {
		return path
	}
	parts := strings.Split(path, ".")
	if len(parts) == 0 {
		return path
	}
	return strings.Join(parts[1:], ".")
}

type TimescaleRequest struct {
	DeviceId   string
	ServiceId  string
	ColumnName string
}
