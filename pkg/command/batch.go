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

	ebatch, err := this.GetEventBatch(token)
	if err != nil {
		log.Println("WARNING: unable to create event batch --> run without batching", err)
		ebatch = nil
	}

	ebatch.CountCommands(len(batch))
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
				} else {
					ebatch.CountWait() // cancel invalid count
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
		} else {
			ebatch.CountWait() // cancel count of reused command
		}
	}
	wg.Wait()
	return result
}

func (this *Command) GetEventBatch(token auth.Token) (batch *eventbatch.EventBatch, err error) {
	return eventbatch.New(token, this.timescale), nil
}
