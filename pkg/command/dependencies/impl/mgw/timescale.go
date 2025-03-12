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
	"github.com/SENERGY-Platform/external-task-worker/lib/devicerepository/model"
	"github.com/SENERGY-Platform/models/go/models"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Timescale struct {
	TimescaleWrapperUrl string
	ProtocolSegmentName string
}

func TimescaleFactory(ctx context.Context, config configuration.Config) (interfaces.Timescale, error) {
	return NewTimescale(config.TimescaleWrapperUrl, config.MgwProtocolSegment)
}

func NewTimescale(timescaleUrl string, protocolSegmentName string) (*Timescale, error) {
	if !strings.Contains(timescaleUrl, "://") {
		timescaleUrl = "http://" + timescaleUrl
	}
	parsed, err := url.Parse(timescaleUrl)
	if err != nil {
		return nil, err
	}
	if parsed.Port() == "" {
		timescaleUrl = timescaleUrl + ":8080"
	}
	return &Timescale{TimescaleWrapperUrl: timescaleUrl, ProtocolSegmentName: protocolSegmentName}, nil
}

func (this *Timescale) GetLastMessage(token auth.Token, device models.Device, service models.Service, protocol model.Protocol, timeout time.Duration) (result map[string]interface{}, err error) {
	list, err := this.Query(token, []Request{{
		DeviceId:   device.LocalId,
		ServiceId:  service.LocalId,
		ColumnName: "",
	}}, timeout)
	if err != nil {
		return result, err
	}
	if len(list) != 1 {
		return result, interfaces.ErrMissingLastValue
	}
	element := list[0]
	if element.Time == nil {
		return result, interfaces.ErrMissingLastValue
	}
	protocolSegmentId := ""
	for _, segment := range protocol.ProtocolSegments {
		if segment.Name == this.ProtocolSegmentName {
			protocolSegmentId = segment.Id
			break
		}
	}
	contentVariableName := ""
	for _, output := range service.Outputs {
		if output.ProtocolSegmentId == protocolSegmentId {
			contentVariableName = output.ContentVariable.Name
			break
		}
	}
	result = map[string]interface{}{contentVariableName: element.Value}
	return result, nil
}

func (this *Timescale) Query(token auth.Token, request []Request, timeout time.Duration) (result []Response, err error) {
	body := &bytes.Buffer{}
	err = json.NewEncoder(body).Encode(request)
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

type Request struct {
	DeviceId   string
	ServiceId  string
	ColumnName string
}

type Response struct {
	Time  *string     `json:"time"`
	Value interface{} `json:"value"`
}
