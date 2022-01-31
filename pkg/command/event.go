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
	"errors"
	"github.com/SENERGY-Platform/device-command/pkg/auth"
	"github.com/SENERGY-Platform/device-command/pkg/command/dependencies/interfaces"
	"github.com/SENERGY-Platform/external-task-worker/lib/devicerepository/model"
	marshallermodel "github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/serialization"
	"log"
	"net/http"
	"sort"
	"strings"
)

func (this *Command) GetLastEventValue(token auth.Token, device model.Device, service model.Service, protocol model.Protocol, characteristicId string) (code int, result interface{}) {
	output, err := this.getLastEventMessage(token, device, service, protocol)
	if err != nil {
		return http.StatusInternalServerError, "unable to get event value: " + err.Error()
	}
	temp, err := this.marshaller.UnmarshalFromServiceAndProtocol(characteristicId, service, protocol, output, nil)
	if err != nil {
		return http.StatusInternalServerError, "unable to unmarshal event value: " + err.Error()
	}
	return 200, temp
}

func (this *Command) getLastEventMessage(token auth.Token, device model.Device, service model.Service, protocol model.Protocol) (result map[string]string, err error) {
	request := createTimescaleRequest(device, service)
	response, err := this.timescale.Query(token, request)
	if err != nil {
		return result, err
	}
	if len(request) != len(response) {
		return result, errors.New("timescale response has less elements then the request")
	}
	result, err = createEventValueFromTimescaleValues(service, protocol, request, response)
	return result, err
}

func createEventValueFromTimescaleValues(service model.Service, protocol model.Protocol, request []interfaces.TimescaleRequest, response []interfaces.TimescaleResponse) (result map[string]string, err error) {
	result = map[string]string{}
	timescaleValue := getTimescaleValue(request, response)
	for _, content := range service.Outputs {
		segmentValue := timescaleValue[content.ContentVariable.Name]
		segmentName := getSegmentName(protocol, content.ProtocolSegmentId)
		if segmentName != "" && segmentValue != nil {
			result[segmentName], err = marshalSegmentValue(content.Serialization, segmentValue, content.ContentVariable.Name)
			if err != nil {
				return result, err
			}
		}
	}
	return result, err
}

func marshalSegmentValue(serializationTo string, value interface{}, rootName string) (string, error) {
	m, ok := serialization.Get(serializationTo)
	if !ok {
		err := errors.New("unknown serialization")
		log.Println("ERROR: unknown serialization", serializationTo)
		return "", err
	}
	return m.Marshal(value, marshallermodel.ContentVariable{Name: rootName})
}

func getTimescaleValue(timescaleRequests []interfaces.TimescaleRequest, timescaleResponses []interfaces.TimescaleResponse) (result map[string]interface{}) {
	pathToValue := map[string]interface{}{}
	paths := []string{}
	for i, request := range timescaleRequests {
		pathToValue[request.ColumnName] = timescaleResponses[i].Value
		paths = append(paths, request.ColumnName)
	}

	sort.Strings(paths)

	result = map[string]interface{}{}
	for _, path := range paths {
		result = setPath(result, strings.Split(path, "."), pathToValue[path])
	}

	return result
}

func setPath(orig map[string]interface{}, path []string, value interface{}) map[string]interface{} {
	if len(path) == 0 {
		return orig
	}
	first := path[0]
	if len(path) == 1 {
		orig[first] = value
		return orig
	}
	rest := path[1:]
	sub, ok := orig[first]
	if !ok {
		orig[first] = setPath(map[string]interface{}{}, rest, value)
	} else {
		subMap, okCast := sub.(map[string]interface{})
		if !okCast {
			log.Println("ERROR: setPath() expect map in path", orig, sub, strings.Join(path, "."))
			return orig
		}
		orig[first] = setPath(subMap, rest, value)
	}
	return orig
}

func getSegmentName(protocol model.Protocol, id string) string {
	for _, segment := range protocol.ProtocolSegments {
		if segment.Id == id {
			return segment.Name
		}
	}
	return ""
}

func createTimescaleRequest(device model.Device, service model.Service) (result []interfaces.TimescaleRequest) {
	paths := []string{}
	for _, content := range service.Outputs {
		paths = append(paths, getContentPaths([]string{}, content.ContentVariable)...)
	}
	for _, path := range paths {
		result = append(result, interfaces.TimescaleRequest{
			DeviceId:   device.Id,
			ServiceId:  service.Id,
			ColumnName: path,
		})
	}
	return result
}

func getContentPaths(current []string, variable model.ContentVariable) (result []string) {
	//skip empty content
	if variable.Name == "" {
		return result
	}
	//skip list
	if variable.Type == model.List {
		return result
	}
	if len(variable.SubContentVariables) == 0 {
		return []string{strings.Join(append(current, variable.Name), ".")}
	}
	for _, sub := range variable.SubContentVariables {
		result = append(result, getContentPaths(append(current, variable.Name), sub)...)
	}
	return result
}
