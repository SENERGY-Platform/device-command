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
	"hash/maphash"
	"net/http"
	"sync"

	"github.com/SENERGY-Platform/device-command/pkg/auth"
)

func (this *Command) Batch(ctx context.Context, token auth.Token, batch BatchRequest, timeout string, preferEventValue bool) []BatchResultElement {
	if len(batch) == 0 {
		return []BatchResultElement{}
	}
	batch.SetAspectIds()
	result := make([]BatchResultElement, len(batch))
	wg := sync.WaitGroup{}
	mux := sync.Mutex{}

	hashSeed := maphash.MakeSeed()
	resultIndexMap := map[uint64][]int{}
	for i, cmd := range batch {
		hash := cmd.Hash(hashSeed)
		resultIndexMap[hash] = append(resultIndexMap[hash], i)
	}
	isAlreadySend := map[uint64]bool{}

	for _, cmd := range batch {
		hash := cmd.Hash(hashSeed)
		if !isAlreadySend[hash] {
			isAlreadySend[hash] = true
			resultIndexesOuter := resultIndexMap[hash]
			wg.Add(1)
			go func(resultIndexes []int, cmd CommandMessage) {
				defer wg.Done()
				var code int
				var temp interface{}
				if cmd.DeviceId != "" && cmd.ServiceId != "" {
					code, temp = this.DeviceCommand(ctx, token, cmd.DeviceId, cmd.ServiceId, cmd.FunctionId, cmd.AspectIds, cmd.Input, timeout, preferEventValue, cmd.CharacteristicId)
				} else if cmd.GroupId != "" {
					code, temp = this.GroupCommand(ctx, token, cmd.GroupId, cmd.FunctionId, cmd.AspectIds, cmd.DeviceClassId, cmd.Input, timeout, preferEventValue, cmd.CharacteristicId)
				}
				if code != http.StatusOK {
					this.config.GetLogger().Warn("error batch response element", "user", token.GetUserId(), "code", code, "response", fmt.Sprintf("%#v", temp))
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
			}(resultIndexesOuter, cmd)
		}
	}
	wg.Wait()
	return result
}
