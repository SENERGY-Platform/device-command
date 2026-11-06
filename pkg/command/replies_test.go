/*
 * Copyright 2026 InfAI (CC SES)
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
	"strconv"
	"sync"
	"testing"

	"github.com/SENERGY-Platform/device-command/pkg/transport"
)

func TestRepliesDeliversToTheWaitingRequest(t *testing.T) {
	replies := NewReplies()
	waiting := replies.expect("a")
	replies.Receive(transport.Reply{TaskId: "b"})
	replies.Receive(transport.Reply{TaskId: "a", Output: map[string]string{"data": "1"}})
	//a second reply for the same task must neither block nor replace the first
	replies.Receive(transport.Reply{TaskId: "a", Output: map[string]string{"data": "2"}})
	reply := <-waiting
	if reply.Output["data"] != "1" {
		t.Error(reply)
	}
	select {
	case reply = <-waiting:
		t.Error("unexpected second reply", reply)
	default:
	}
}

func TestRepliesDropsReplyAfterForget(t *testing.T) {
	replies := NewReplies()
	waiting := replies.expect("a")
	replies.forget("a")
	replies.Receive(transport.Reply{TaskId: "a"})
	select {
	case reply := <-waiting:
		t.Error("unexpected reply", reply)
	default:
	}
	if len(replies.waiting) != 0 {
		t.Error(replies.waiting)
	}
}

func TestRepliesConcurrently(t *testing.T) {
	replies := NewReplies()
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		taskId := strconv.Itoa(i)
		waiting := replies.expect(taskId)
		wg.Add(2)
		go func() {
			defer wg.Done()
			replies.Receive(transport.Reply{TaskId: taskId})
		}()
		go func() {
			defer wg.Done()
			defer replies.forget(taskId)
			if reply := <-waiting; reply.TaskId != taskId {
				t.Error(reply)
			}
		}()
	}
	wg.Wait()
	if len(replies.waiting) != 0 {
		t.Error(replies.waiting)
	}
}
