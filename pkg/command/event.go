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

package command

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/SENERGY-Platform/device-command/pkg/auth"
	"github.com/SENERGY-Platform/device-command/pkg/lastvalue"
	"github.com/SENERGY-Platform/device-command/pkg/marshaller"
	marshallermodel "github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/serialization"
	model "github.com/SENERGY-Platform/models/go/models"
)

// StatusMissingLastValue answers a measuring command on a gateway, if the device has not sent since
// mgw-last-value started. A client may then ask the platform instead.
const StatusMissingLastValue = 513

func (this *Command) GetLastEventValue(ctx context.Context, token auth.Token, device model.Device, service model.Service, protocol model.Protocol, characteristicId string, functionId string, aspectNodes []model.AspectNode, timeout time.Duration) (code int, result interface{}, valueTime string) {
	output, valueTime, err, code := this.getLastEventMessage(ctx, token, device, service, protocol, timeout)
	if err != nil {
		return code, "unable to get event value: " + err.Error(), ""
	}
	req := marshaller.UnmarshalRequest{
		Service:          service,
		Protocol:         protocol,
		CharacteristicId: characteristicId,
		Message:          output,
		FunctionId:       functionId,
		AspectNodes:      aspectNodes,
	}
	temp, err := this.marshaller.Unmarshal(ctx, req)
	if err != nil {
		this.config.GetLogger().Debug("unable to unmarshal event value", "error", err, "UnmarshallingV2Request", fmt.Sprintf("%#v", req))
		return http.StatusInternalServerError, "unable to unmarshal event value: " + err.Error(), ""
	}
	return 200, temp, valueTime
}

func (this *Command) getLastEventMessage(ctx context.Context, token auth.Token, device model.Device, service model.Service, protocol model.Protocol, timeout time.Duration) (result map[string]string, valueTime string, err error, code int) {
	response, err := this.lastValues.LastMessage(ctx, token, device, service, protocol, timeout)
	if errors.Is(err, lastvalue.ErrMissing) {
		return result, "", err, StatusMissingLastValue
	}
	if err != nil {
		return result, "", err, http.StatusInternalServerError
	}
	result, err, code = this.useProtocolSerialization(service, protocol, response.Value)
	return result, response.Time, err, code
}

func (this *Command) useProtocolSerialization(service model.Service, protocol model.Protocol, lastMsg map[string]interface{}) (result map[string]string, err error, code int) {
	result = map[string]string{}
	for _, content := range service.Outputs {
		segmentValue := lastMsg[content.ContentVariable.Name]
		segmentName := getSegmentName(protocol, content.ProtocolSegmentId)
		if segmentName != "" && segmentValue != nil {
			result[segmentName], err = marshalSegmentValue(string(content.Serialization), segmentValue, content.ContentVariable.Name)
			if err != nil {
				return result, err, http.StatusInternalServerError
			}
		}
	}
	return result, err, http.StatusInternalServerError
}

func marshalSegmentValue(serializationTo string, value interface{}, rootName string) (string, error) {
	m, ok := serialization.Get(model.Serialization(serializationTo))
	if !ok {
		err := errors.New("unknown serialization")
		slog.Error("unable to get serialization", "error", err, "serialization", serializationTo)
		return "", err
	}
	return m.Marshal(value, marshallermodel.ContentVariable{Name: rootName})
}

func getSegmentName(protocol model.Protocol, id string) string {
	for _, segment := range protocol.ProtocolSegments {
		if segment.Id == id {
			return segment.Name
		}
	}
	return ""
}
