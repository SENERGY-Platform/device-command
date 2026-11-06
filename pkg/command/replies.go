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
	"sync"

	"github.com/SENERGY-Platform/device-command/pkg/transport"
)

// Replies hands the reply of a connector to the request that waits for it. It is created before
// the transport, which delivers every reply to Receive, and is then passed to New.
type Replies struct {
	mux     sync.Mutex
	waiting map[string]chan transport.Reply
}

func NewReplies() *Replies {
	return &Replies{waiting: map[string]chan transport.Reply{}}
}

// Receive passes the reply on to the waiting request. A reply that nobody waits for any more,
// because the request timed out or a reply already arrived, is dropped.
func (this *Replies) Receive(reply transport.Reply) {
	this.mux.Lock()
	defer this.mux.Unlock()
	waiting, ok := this.waiting[reply.TaskId]
	if !ok {
		return
	}
	delete(this.waiting, reply.TaskId)
	waiting <- reply
}

// expect has to be called before the command is sent, so that an early reply is not lost.
func (this *Replies) expect(taskId string) <-chan transport.Reply {
	this.mux.Lock()
	defer this.mux.Unlock()
	result := make(chan transport.Reply, 1)
	this.waiting[taskId] = result
	return result
}

func (this *Replies) forget(taskId string) {
	this.mux.Lock()
	defer this.mux.Unlock()
	delete(this.waiting, taskId)
}
