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

package register

import (
	"log"
	"net/http"
	"sync"
	"time"
)

func New(defaultTimeout time.Duration, debug bool) *Register {
	return &Register{register: map[string]*State{}, defaultTimeout: defaultTimeout, debug: debug}
}

type Register struct {
	debug          bool
	register       map[string]*State
	mux            sync.Mutex
	defaultTimeout time.Duration
}

type State struct {
	wg   *sync.WaitGroup
	code int
	resp interface{}
}

func (this *Register) Register(id string) {
	this.mux.Lock()
	defer this.mux.Unlock()
	wg := &sync.WaitGroup{}
	wg.Add(1)
	this.register[id] = &State{
		wg:   wg,
		code: 0,
		resp: nil,
	}
}

func (this *Register) Complete(id string, code int, value interface{}) {
	if this.debug {
		log.Println("complete", id, code, value)
	}
	this.mux.Lock()
	defer this.mux.Unlock()
	state, ok := this.register[id]
	if !ok {
		return
	}
	state.resp, state.code = value, code
	state.wg.Done()
}

func (this *Register) Wait(id string) (int, interface{}) {
	return this.WaitWithTimeout(id, this.defaultTimeout)
}

func (this *Register) WaitWithTimeout(id string, timeout time.Duration) (int, interface{}) {
	this.mux.Lock()
	state, ok := this.register[id]
	this.mux.Unlock()
	if !ok {
		return http.StatusInternalServerError, "unregistered correlation id"
	}
	defer delete(this.register, id)

	t := time.AfterFunc(timeout, func() {
		this.Complete(id, http.StatusRequestTimeout, "timeout")
	})
	defer func() {
		t.Stop()
	}()

	state.wg.Wait()
	return state.code, state.resp
}
