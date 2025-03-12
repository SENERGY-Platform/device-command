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
	"encoding/json"
	"errors"
	"github.com/SENERGY-Platform/device-command/pkg/auth"
	"github.com/SENERGY-Platform/device-command/pkg/command/dependencies/interfaces"
	"github.com/SENERGY-Platform/external-task-worker/lib/devicerepository/model"
	"github.com/SENERGY-Platform/external-task-worker/lib/marshaller"
	marshallermodel "github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/serialization"
	"github.com/SENERGY-Platform/models/go/models"
	"log"
	"net/http"
	"time"
)

func (this *Command) GetLastEventValue(token auth.Token, device model.Device, service model.Service, protocol model.Protocol, characteristicId string, functionId string, aspect model.AspectNode, timeout time.Duration) (code int, result interface{}) {
	output, err, code := this.getLastEventMessage(token, device, service, protocol, timeout)
	if err != nil {
		return code, "unable to get event value: " + err.Error()
	}
	temp, err := this.marshaller.UnmarshalV2(marshaller.UnmarshallingV2Request{
		Service:          service,
		Protocol:         protocol,
		CharacteristicId: characteristicId,
		Message:          output,
		FunctionId:       functionId,
		AspectNode:       aspect,
		AspectNodeId:     aspect.Id,
	})
	if err != nil {
		if this.config.Debug {
			log.Println("ERROR:", err)
			marshalRequestStr, _ := json.Marshal(marshaller.UnmarshallingV2Request{
				Service:          service,
				Protocol:         protocol,
				CharacteristicId: characteristicId,
				Message:          output,
				FunctionId:       functionId,
				AspectNode:       aspect,
				AspectNodeId:     aspect.Id,
			})
			log.Println("ERROR: unmarshal request", string(marshalRequestStr))
		}
		return http.StatusInternalServerError, "unable to unmarshal event value: " + err.Error()
	}
	return 200, temp
}

func (this *Command) getLastEventMessage(token auth.Token, device model.Device, service model.Service, protocol model.Protocol, timeout time.Duration) (result map[string]string, err error, code int) {
	response, err := this.timescale.GetLastMessage(token, device, service, protocol, timeout)
	if errors.Is(err, interfaces.ErrMissingLastValue) {
		return result, err, interfaces.ErrMissingLastValueCode
	}
	if err != nil {
		return result, err, http.StatusInternalServerError
	}
	result, err, code = this.useProtocolSerialization(service, protocol, response)
	return result, err, code
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
	m, ok := serialization.Get(models.Serialization(serializationTo))
	if !ok {
		err := errors.New("unknown serialization")
		log.Println("ERROR: unknown serialization", serializationTo)
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
