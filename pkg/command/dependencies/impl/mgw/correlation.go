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

package mgw

import (
	"errors"
	"github.com/SENERGY-Platform/external-task-worker/lib/messages"
	"sync"
)

type Correlation interface {
	Get(id string) (messages.ProtocolMsg, error)
	Set(id string, element messages.ProtocolMsg) error
}

type CorrelationImpl struct {
	correlations map[string]messages.ProtocolMsg
	mux          sync.Mutex
}

func (this *CorrelationImpl) Get(id string) (messages.ProtocolMsg, error) {
	this.mux.Lock()
	defer this.mux.Unlock()
	if this.correlations == nil {
		this.correlations = map[string]messages.ProtocolMsg{}
	}
	var err error
	result, found := this.correlations[id]
	if !found {
		err = errors.New("correlation id not found")
	}
	return result, err
}

func (this *CorrelationImpl) Set(id string, element messages.ProtocolMsg) error {
	this.mux.Lock()
	defer this.mux.Unlock()
	if this.correlations == nil {
		this.correlations = map[string]messages.ProtocolMsg{}
	}
	this.correlations[id] = element
	return nil
}

var DefaultCorrelation = &CorrelationImpl{}
