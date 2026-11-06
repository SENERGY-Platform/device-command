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

package lastvalue

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/SENERGY-Platform/device-command/pkg/auth"
	"github.com/SENERGY-Platform/device-command/pkg/tracing"
	model "github.com/SENERGY-Platform/models/go/models"
)

// Cloud reads from the timescale-wrapper of the platform.
type Cloud struct {
	url string
}

func (this *Cloud) LastMessage(ctx context.Context, token auth.Token, device model.Device, service model.Service, protocol model.Protocol, timeout time.Duration) (map[string]interface{}, error) {
	query := url.Values{}
	query.Set("device_id", device.Id)
	query.Set("service_id", service.Id)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, this.url+"/last-message?"+query.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", token.Jwt())
	err = tracing.Inject(ctx, req)
	if err != nil {
		return nil, err
	}
	resp, err := (&http.Client{Timeout: timeout}).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		temp, _ := io.ReadAll(resp.Body)
		return nil, errors.New(strings.TrimSpace(string(temp)))
	}
	result := struct {
		Value map[string]interface{} `json:"value"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	return result.Value, err
}
