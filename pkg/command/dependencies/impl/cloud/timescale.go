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

package cloud

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/SENERGY-Platform/device-command/pkg/auth"
	"github.com/SENERGY-Platform/device-command/pkg/command/dependencies/interfaces"
	"github.com/SENERGY-Platform/device-command/pkg/configuration"
	"github.com/SENERGY-Platform/external-task-worker/lib/devicerepository/model"
	"github.com/SENERGY-Platform/models/go/models"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Timescale struct {
	TimescaleWrapperUrl string
}

func TimescaleFactory(ctx context.Context, config configuration.Config) (interfaces.Timescale, error) {
	return NewTimescale(config.TimescaleWrapperUrl), nil
}

func NewTimescale(timescaleWrapperUrl string) *Timescale {
	return &Timescale{TimescaleWrapperUrl: timescaleWrapperUrl}
}

type LastMessageResponse struct {
	Time  string                 `json:"time"`
	Value map[string]interface{} `json:"value"`
}

func (this *Timescale) GetLastMessage(token auth.Token, device models.Device, service models.Service, protocol model.Protocol, timeout time.Duration) (result map[string]interface{}, err error) {
	query := url.Values{}
	query.Set("device_id", device.Id)
	query.Set("service_id", service.Id)
	req, err := http.NewRequest("GET", this.TimescaleWrapperUrl+"/last-message?"+query.Encode(), nil)
	if err != nil {
		return result, err
	}
	req.Header.Set("Authorization", token.Jwt())
	client := &http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		temp, _ := io.ReadAll(resp.Body)
		return result, errors.New(strings.TrimSpace(string(temp)))
	}
	wrapper := LastMessageResponse{}
	err = json.NewDecoder(resp.Body).Decode(&wrapper)
	if err != nil {
		return result, err
	}
	return wrapper.Value, nil
}
