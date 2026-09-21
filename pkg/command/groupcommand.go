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
	"fmt"
	"net/http"
	"sort"
	"sync"

	"github.com/SENERGY-Platform/device-command/pkg/auth"
	"github.com/SENERGY-Platform/external-task-worker/lib/devicerepository/model"
	marshallermodel "github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
)

func (this *Command) GroupCommand(token auth.Token, groupId string, functionId string, aspectIds []string, deviceClassId string, input interface{}, timeout string, preferEventValue bool, characteristicId string) (code int, resp interface{}) {
	subTasks, err := this.GetSubTasks(token.Jwt(), groupId, functionId, aspectIds, deviceClassId, input)
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	wg := sync.WaitGroup{}
	results := []interface{}{}
	var lastErr interface{}
	var lastErrCode int
	for _, sub := range subTasks {
		wg.Add(1)
		go func(sub SubCommand) {
			defer wg.Done()
			tempCode, temp := this.deviceCommand(token, sub.DeviceId, sub.ServiceId, sub.FunctionId, sub.AspectIds, input, timeout, preferEventValue, characteristicId)
			this.config.GetLogger().Debug("group sub result", "user", token.GetUserId(), "code", tempCode, "result", fmt.Sprintf("%#v", temp))
			if tempCode == http.StatusOK {
				results = append(results, temp)
			} else {
				lastErr = temp
				lastErrCode = tempCode
			}
		}(sub)
	}
	wg.Wait()
	if len(results) == 0 && len(subTasks) > 0 {
		return lastErrCode, lastErr
	}
	return http.StatusOK, results
}

type SubCommand struct {
	FunctionId string `json:"function_id"` //mandatory
	AspectIds  []string
	Input      interface{} `json:"input"`
	DeviceId   string      `json:"device_id,omitempty"`
	ServiceId  string      `json:"service_id,omitempty"`
}

func (this *Command) GetSubTasks(token string, deviceGroupId string, functionId string, aspectIds []string, deviceClassId string, input interface{}) (result []SubCommand, err error) {
	group, err := this.iot.GetDeviceGroup(token, deviceGroupId)
	if err != nil {
		return nil, err
	}
	for _, deviceId := range group.DeviceIds {
		device, err := this.iot.GetDevice(token, deviceId)
		if err != nil {
			return nil, err
		}

		deviceType, err := this.iot.GetDeviceType(token, device.DeviceTypeId)
		if err != nil {
			return nil, err
		}

		aspectNodes := []model.AspectNode{}
		for _, aspectId := range aspectIds {
			aspectNode, err := this.iot.GetAspectNode(aspectId)
			if err != nil {
				this.config.GetLogger().Warn("unable to find aspect node, use aspect node without descendants", "aspect_id", aspectId, "error", err)
				aspectNode = model.AspectNode{Id: aspectId}
			}
			aspectNodes = append(aspectNodes, aspectNode)
		}

		if deviceClassId == "" || deviceClassId == deviceType.DeviceClassId {
			services := this.getFilteredServices(functionId, aspectNodes, deviceType.Services)
			for _, service := range services {
				result = append(result, SubCommand{
					FunctionId: functionId,
					Input:      input,
					DeviceId:   device.Id,
					ServiceId:  service.Id,
					AspectIds:  aspectIds,
				})
			}
		}
	}
	return result, nil
}

func (this *Command) getFilteredServices(functionId string, aspectNodes []model.AspectNode, services []model.Service) (result []model.Service) {
	serviceIndex := map[string]model.Service{}
	for _, service := range services {
		contents := service.Inputs
		if isMeasuringFunctionId(functionId) {
			contents = service.Outputs
		}
		matchesCriteria := anyContentMatchesCriteria(contents, functionId, aspectNodes)
		if matchesCriteria {
			serviceIndex[service.Id] = service
		}
	}
	for _, service := range serviceIndex {
		result = append(result, service)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Id < result[j].Id
	})
	return result
}

func anyContentMatchesCriteria(contents []model.Content, functionId string, aspectNodes []model.AspectNode) bool {
	for _, content := range contents {
		if contentVariableContainsCriteria(content.ContentVariable, functionId, aspectNodes) {
			return true
		}
	}
	return false
}

// contentVariableContainsCriteria reports whether a content variable serves the function and
// the aspects of the command. Every requested aspect has to be matched, by the aspect itself
// or by one of its descendants, the way the device-repository reads a filter-criteria that
// names several aspects; without a requested aspect the function alone decides.
func contentVariableContainsCriteria(variable model.ContentVariable, functionId string, aspectNodes []model.AspectNode) bool {
	if variable.FunctionId == functionId &&
		marshallermodel.AspectMatchLevel(marshallermodel.ContentVariableAspectIds(variable), aspectNodes) >= 0 {
		return true
	}
	for _, sub := range variable.SubContentVariables {
		if contentVariableContainsCriteria(sub, functionId, aspectNodes) {
			return true
		}
	}
	return false
}
