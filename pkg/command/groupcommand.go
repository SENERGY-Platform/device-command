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
	"github.com/SENERGY-Platform/device-command/pkg/auth"
	"github.com/SENERGY-Platform/device-command/pkg/command/eventbatch"
	"github.com/SENERGY-Platform/external-task-worker/lib/devicerepository/model"
	"log"
	"net/http"
	"sort"
	"sync"
)

func (this *Command) GroupCommand(token auth.Token, groupId string, functionId string, aspectId string, deviceClassId string, input interface{}, timeout string, preferEventValue bool, batch *eventbatch.EventBatch) (code int, resp interface{}) {
	subTasks, err := this.GetSubTasks(token.Jwt(), groupId, functionId, aspectId, deviceClassId, input)
	if err != nil {
		batch.CountWait() //error -> cancel count of group
		return http.StatusInternalServerError, err.Error()
	}
	wg := sync.WaitGroup{}
	results := []interface{}{}
	var lastErr interface{}
	var lastErrCode int
	batch.CountCommands(len(subTasks)) //add count for device commands
	batch.CountWait()                  //cancel count of group command
	for _, sub := range subTasks {
		wg.Add(1)
		go func(sub SubCommand) {
			defer wg.Done()
			tempCode, temp := this.deviceCommand(token, sub.DeviceId, sub.ServiceId, sub.FunctionId, sub.AspectId, input, timeout, preferEventValue, batch)
			if this.config.Debug {
				log.Println("DEBUG: group sub result:", tempCode, temp)
			}
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
	AspectId   string
	Input      interface{} `json:"input"`
	DeviceId   string      `json:"device_id,omitempty"`
	ServiceId  string      `json:"service_id,omitempty"`
}

func (this *Command) GetSubTasks(token string, deviceGroupId string, functionId string, aspectId string, deviceClassId string, input interface{}) (result []SubCommand, err error) {
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

		aspect := model.AspectNode{}
		if aspectId != "" {
			aspect, err = this.iot.GetAspectNode(token, aspectId)
			if err != nil {
				log.Println("WARNING: unable to find aspect node, use aspect node without descendants", err)
				aspect.Id = aspectId
				err = nil
			}
		}

		if deviceClassId == "" || deviceClassId == deviceType.DeviceClassId {
			services := this.getFilteredServices(functionId, aspect, deviceType.Services)
			for _, service := range services {
				result = append(result, SubCommand{
					FunctionId: functionId,
					Input:      input,
					DeviceId:   device.Id,
					ServiceId:  service.Id,
					AspectId:   aspectId,
				})
			}
		}
	}
	return result, nil
}

func (this *Command) getFilteredServices(functionId string, aspect model.AspectNode, services []model.Service) (result []model.Service) {
	serviceIndex := map[string]model.Service{}
	for _, service := range services {
		contents := service.Inputs
		if isMeasuringFunctionId(functionId) {
			contents = service.Outputs
		}
		matchesCriteria := anyContentMatchesCriteria(contents, model.DeviceGroupFilterCriteria{FunctionId: functionId, AspectId: aspect.Id}, aspect)
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

func anyContentMatchesCriteria(contents []model.Content, criteria model.DeviceGroupFilterCriteria, aspectNode model.AspectNode) bool {
	for _, content := range contents {
		if contentVariableContainsCriteria(content.ContentVariable, criteria, aspectNode) {
			return true
		}
	}
	return false
}

func contentVariableContainsCriteria(variable model.ContentVariable, criteria model.DeviceGroupFilterCriteria, aspectNode model.AspectNode) bool {
	if variable.FunctionId == criteria.FunctionId &&
		(criteria.AspectId == "" ||
			variable.AspectId == criteria.AspectId ||
			listContains(aspectNode.DescendentIds, variable.AspectId)) {
		return true
	}
	for _, sub := range variable.SubContentVariables {
		if contentVariableContainsCriteria(sub, criteria, aspectNode) {
			return true
		}
	}
	return false
}

func listContains(list []string, search string) bool {
	for _, element := range list {
		if element == search {
			return true
		}
	}
	return false
}
