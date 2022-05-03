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
	"sync"
)

func (this *Command) Batch(token auth.Token, batch BatchRequest, timeout string, preferEventValue bool) []BatchResultElement {
	if len(batch) == 0 {
		return []BatchResultElement{}
	}

	result := make([]BatchResultElement, len(batch))
	wg := sync.WaitGroup{}
	mux := sync.Mutex{}

	resultIndexMap := map[string][]int{}
	for i, cmd := range batch {
		hash := cmd.Hash()
		resultIndexMap[hash] = append(resultIndexMap[hash], i)
	}
	isAlreadySend := map[string]bool{}

	ebatch, err := this.GetEventBatch(token, batch, preferEventValue)
	if err != nil {
		log.Println("WARNING: unable to create event batch --> run without batching", err)
		ebatch = nil
	}

	for _, cmd := range batch {
		hash := cmd.Hash()
		if !isAlreadySend[hash] {
			isAlreadySend[hash] = true
			wg.Add(1)
			go func(resultIndexes []int, cmd CommandMessage) {
				defer wg.Done()
				var code int
				var temp interface{}
				if cmd.DeviceId != "" && cmd.ServiceId != "" {
					code, temp = this.DeviceCommand(token, cmd.DeviceId, cmd.ServiceId, cmd.FunctionId, cmd.AspectId, cmd.Input, timeout, preferEventValue, ebatch)
				} else if cmd.GroupId != "" {
					code, temp = this.GroupCommand(token, cmd.GroupId, cmd.FunctionId, cmd.AspectId, cmd.DeviceClassId, cmd.Input, timeout, preferEventValue, ebatch)
				}
				mux.Lock()
				defer mux.Unlock()
				for _, index := range resultIndexes {
					result[index] = BatchResultElement{
						StatusCode: code,
						Message:    temp,
					}
				}
				return
			}(resultIndexMap[hash], cmd)
		}
	}
	wg.Wait()
	return result
}

func (this *Command) GetEventBatch(token auth.Token, tasks BatchRequest, preferEventValue bool) (batch *eventbatch.EventBatch, err error) {
	count, err := this.expectedEventRequests(token, tasks, preferEventValue)
	if err != nil {
		return nil, err
	}
	return eventbatch.New(token, this.timescale, count), nil
}

func (this *Command) expectedEventRequests(token auth.Token, batch []CommandMessage, preferEventValue bool) (count int64, err error) {
	isAlreadySend := map[string]bool{}
	for _, cmd := range batch {
		hash := cmd.Hash()
		if !isAlreadySend[hash] {
			isAlreadySend[hash] = true
			if cmd.DeviceId != "" && cmd.ServiceId != "" {
				device, err := this.iot.GetDevice(token.Jwt(), cmd.DeviceId)
				if err != nil {
					return count, err
				}
				service, err := this.iot.GetService(token.Jwt(), device, cmd.ServiceId)
				if err != nil {
					return count, err
				}
				var aspectError error
				if cmd.AspectId != "" {
					_, aspectError = this.iot.GetAspectNode(token.Jwt(), cmd.AspectId)
				}
				if aspectError == nil && isMeasuringFunctionId(cmd.FunctionId) && (service.Interaction == model.EVENT || (preferEventValue && service.Interaction == model.EVENT_AND_REQUEST)) {
					count = count + 1
				}
			} else if cmd.GroupId != "" {
				subTasks, err := this.GetSubTasks(token.Jwt(), cmd.GroupId, cmd.FunctionId, cmd.AspectId, cmd.DeviceClassId, cmd.Input)
				if err != nil {
					return count, err
				}
				for _, sub := range subTasks {
					device, err := this.iot.GetDevice(token.Jwt(), sub.DeviceId)
					if err != nil {
						return count, err
					}
					service, err := this.iot.GetService(token.Jwt(), device, sub.ServiceId)
					if err != nil {
						return count, err
					}
					if isMeasuringFunctionId(sub.FunctionId) && (service.Interaction == model.EVENT || (preferEventValue && service.Interaction == model.EVENT_AND_REQUEST)) {
						count = count + 1
					}
				}
			}
		}
	}
	return count, nil
}
