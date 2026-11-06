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

package transport

import (
	"reflect"
	"testing"
)

// The automation-engine sends commands to the same connectors on a gateway, as
// {"command_id": "automation-engine-<id>", ...}, and every client of the broker gets every
// response and error. Only those with the own prefix may reach the waiting requests.
func TestMqttReadsOnlyOwnReplies(t *testing.T) {
	mqtt := &Mqtt{prefix: "device-command-", segment: "data"}

	reply, ok := mqtt.readResponse([]byte(`{"command_id":"device-command-task-1","data":"42"}`))
	if !ok || !reflect.DeepEqual(reply, Reply{TaskId: "task-1", Output: map[string]string{"data": "42"}}) {
		t.Errorf("%v %#v", ok, reply)
	}
	reply, ok = mqtt.readError("error/command/device-command-task-2", []byte("unreachable"))
	if !ok || !reflect.DeepEqual(reply, Reply{TaskId: "task-2", Output: map[string]string{"data": "unreachable"}, Failed: true}) {
		t.Errorf("%v %#v", ok, reply)
	}

	for _, foreign := range []string{"automation-engine-task-1", "task-1", "", "x-device-command-task-1"} {
		if reply, ok := mqtt.readResponse([]byte(`{"command_id":"` + foreign + `","data":"42"}`)); ok {
			t.Errorf("response %q read as own: %#v", foreign, reply)
		}
		if reply, ok := mqtt.readError("error/command/"+foreign, []byte("unreachable")); ok {
			t.Errorf("error %q read as own: %#v", foreign, reply)
		}
	}
	if _, ok := mqtt.readResponse([]byte(`not json`)); ok {
		t.Error("invalid response read")
	}
}
