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

package eventbatch

import (
	"errors"
	"github.com/SENERGY-Platform/device-command/pkg/auth"
	"github.com/SENERGY-Platform/device-command/pkg/command/dependencies/interfaces"
	"github.com/SENERGY-Platform/external-task-worker/lib/devicerepository/model"
	"log"
	"runtime/debug"
	"sync"
)

type EventBatch struct {
	token           auth.Token
	timescale       interfaces.Timescale
	requests        []interfaces.TimescaleRequest
	responses       []interfaces.TimescaleResponse
	mapping         map[string]map[string][]int //device.id -> service.id -> requests indexes
	mux             sync.Mutex
	wg              sync.WaitGroup
	error           error
	finished        bool
	expectedQueries int64
	queryCount      int64
}

func New(token auth.Token, timescale interfaces.Timescale, expectedQueries int64) *EventBatch {
	result := &EventBatch{
		token:           token,
		timescale:       timescale,
		requests:        []interfaces.TimescaleRequest{},
		mapping:         map[string]map[string][]int{},
		expectedQueries: expectedQueries,
	}
	result.wg.Add(1)
	return result
}

func (this *EventBatch) Query(device model.Device, service model.Service, request []interfaces.TimescaleRequest) (result []interfaces.TimescaleResponse, err error) {
	this.mux.Lock()
	if _, ok := this.mapping[device.Id]; !ok {
		this.mapping[device.Id] = map[string][]int{}
	}
	if _, ok := this.mapping[device.Id][service.Id]; !ok {
		this.mapping[device.Id][service.Id] = []int{}
		for i, _ := range request {
			index := len(this.requests) + i
			this.mapping[device.Id][service.Id] = append(this.mapping[device.Id][service.Id], index)
		}
		this.requests = append(this.requests, request...)
	}
	this.mux.Unlock()
	err = this.Wait()
	if err != nil {
		return result, err
	}
	for _, i := range this.mapping[device.Id][service.Id] {
		if i < len(this.responses) {
			result = append(result, this.responses[i])
		} else {
			err = errors.New("missing value in responses")
			log.Println("ERROR:", err)
			debug.PrintStack()
			return result, err
		}
	}
	return result, err
}

func (this *EventBatch) Wait() error {
	this.mux.Lock()
	this.queryCount++
	if this.queryCount == this.expectedQueries {
		go this.sendRequest()
	}
	this.mux.Unlock()
	this.wg.Wait()
	return this.error
}

func (this *EventBatch) sendRequest() {
	this.mux.Lock()
	defer this.mux.Unlock()
	if this.finished {
		log.Println("WARNING: batch is already finished")
		debug.PrintStack()
		return
	}
	this.finished = true
	this.responses, this.error = this.timescale.Query(this.token, this.requests)
	this.wg.Done()
}
