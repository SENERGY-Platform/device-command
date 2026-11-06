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
	"fmt"
	"github.com/SENERGY-Platform/device-command/pkg/tracing"
	"go.opentelemetry.io/otel/attribute"
	"net/http"
	"sort"
	"sync"

	"github.com/SENERGY-Platform/device-command/pkg/auth"
	marshallermodel "github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
	model "github.com/SENERGY-Platform/models/go/models"
)

func (this *Command) GroupCommand(ctx context.Context, token auth.Token, groupId string, functionId string, aspectIds []string, deviceClassId string, input interface{}, timeout string, preferEventValue bool, characteristicId string) (code int, resp interface{}) {
	ctx, span := tracing.Start(ctx, "group command",
		attribute.String("device_group_id", groupId),
		attribute.String("function_id", functionId),
		attribute.StringSlice("aspect_ids", aspectIds))
	defer func() { tracing.EndWithStatus(span, code, resp) }()
	subTasks, err := this.GetSubTasks(ctx, token.Jwt(), groupId, functionId, aspectIds, deviceClassId, input)
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	span.SetAttributes(attribute.Int("device_commands", len(subTasks)))
	wg := sync.WaitGroup{}
	mux := sync.Mutex{}
	results := []interface{}{}
	var lastErr interface{}
	var lastErrCode int
	for _, sub := range subTasks {
		wg.Add(1)
		go func(sub SubCommand) {
			defer wg.Done()
			tempCode, temp, _ := this.deviceCommand(ctx, token, sub.DeviceId, sub.ServiceId, sub.FunctionId, sub.AspectIds, input, timeout, preferEventValue, characteristicId)
			this.config.GetLogger().Debug("group sub result", "user", token.GetUserId(), "code", tempCode, "result", fmt.Sprintf("%#v", temp))
			mux.Lock()
			defer mux.Unlock()
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

func (this *Command) GetSubTasks(ctx context.Context, token string, deviceGroupId string, functionId string, aspectIds []string, deviceClassId string, input interface{}) (result []SubCommand, err error) {
	group, err := this.repo.GetDeviceGroup(ctx, token, deviceGroupId)
	if err != nil {
		return nil, err
	}
	for _, deviceId := range group.DeviceIds {
		device, err := this.repo.GetDevice(ctx, token, deviceId)
		if err != nil {
			return nil, err
		}

		deviceType, err := this.repo.GetDeviceType(ctx, token, device.DeviceTypeId)
		if err != nil {
			return nil, err
		}

		aspectNodes := []model.AspectNode{}
		for _, aspectId := range aspectIds {
			aspectNode, err := this.repo.GetAspectNode(ctx, aspectId)
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
