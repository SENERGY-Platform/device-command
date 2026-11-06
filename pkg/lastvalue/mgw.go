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
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/SENERGY-Platform/device-command/pkg/auth"
	"github.com/SENERGY-Platform/device-command/pkg/tracing"
	model "github.com/SENERGY-Platform/models/go/models"
)

// Mgw reads from mgw-last-value on a gateway. It knows devices and services by their local ids
// and stores the raw value of one protocol segment.
type Mgw struct {
	url             string
	protocolSegment string
}

// NewMgw accepts the url of mgw-last-value without scheme and port; they default to http and 8080.
func NewMgw(lastValueUrl string, protocolSegment string) (*Mgw, error) {
	if !strings.Contains(lastValueUrl, "://") {
		lastValueUrl = "http://" + lastValueUrl
	}
	parsed, err := url.Parse(lastValueUrl)
	if err != nil {
		return nil, err
	}
	if parsed.Port() == "" {
		lastValueUrl = lastValueUrl + ":8080"
	}
	return &Mgw{url: lastValueUrl, protocolSegment: protocolSegment}, nil
}

type mgwRequest struct {
	DeviceId   string
	ServiceId  string
	ColumnName string
}

type mgwResponse struct {
	Time  *string     `json:"time"`
	Value interface{} `json:"value"`
}

// LastMessage returns the value under the name of the output content variable that is sent in
// the configured protocol segment.
func (this *Mgw) LastMessage(ctx context.Context, token auth.Token, device model.Device, service model.Service, protocol model.Protocol, timeout time.Duration) (map[string]interface{}, error) {
	list, err := this.query(ctx, token, []mgwRequest{{DeviceId: device.LocalId, ServiceId: service.LocalId}}, timeout)
	if err != nil {
		return nil, err
	}
	if len(list) != 1 || list[0].Time == nil {
		return nil, ErrMissing
	}
	segmentId := ""
	for _, segment := range protocol.ProtocolSegments {
		if segment.Name == this.protocolSegment {
			segmentId = segment.Id
			break
		}
	}
	variableName := ""
	for _, output := range service.Outputs {
		if output.ProtocolSegmentId == segmentId {
			variableName = output.ContentVariable.Name
			break
		}
	}
	return map[string]interface{}{variableName: list[0].Value}, nil
}

func (this *Mgw) query(ctx context.Context, token auth.Token, request []mgwRequest, timeout time.Duration) (result []mgwResponse, err error) {
	body := &bytes.Buffer{}
	err = json.NewEncoder(body).Encode(request)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, this.url+"/last-values", body)
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
		slog.Error("unable to query /last-values", "error", err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		temp, _ := io.ReadAll(resp.Body)
		return nil, errors.New(strings.TrimSpace(string(temp)))
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	return result, err
}
