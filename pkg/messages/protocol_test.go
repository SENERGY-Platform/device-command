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

package messages

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/SENERGY-Platform/models/go/models"
)

// The expected json was produced by external-task-worker v0.1.1 (lib/messages) for the same
// message, before its types were taken over here. Only the parts this package defines are
// written out; device, service, protocol and aspect nodes are models/go and are inserted as
// models/go marshals them, so that a models update does not break this test.
func TestProtocolMsgJsonIsUnchanged(t *testing.T) {
	device := models.Device{Id: "device-1", LocalId: "local-1", DeviceTypeId: "dt-1"}
	service := models.Service{Id: "service-1", LocalId: "s-local", ProtocolId: "protocol-1", Interaction: models.REQUEST}
	protocol := models.Protocol{Id: "protocol-1", Handler: "connector"}
	nodeA := models.AspectNode{Id: "a", Name: "A"}
	nodeB := models.AspectNode{Id: "b", Name: "B"}

	msg := ProtocolMsg{
		Request:  ProtocolRequest{Input: map[string]string{"data": `{"level":1}`}},
		Response: ProtocolResponse{Output: map[string]string{"data": `{"level":2}`}},
		TaskInfo: TaskInfo{TaskId: "task-1", Time: "1700000000", TenantId: "user-1"},
		Metadata: Metadata{
			Version:              3,
			Device:               device,
			Service:              service,
			Protocol:             protocol,
			OutputFunctionId:     "f-1",
			InputCharacteristic:  "c-in",
			OutputCharacteristic: "c-out",
			ResponseTo:           "device-command-response",
			ErrorTo:              "device-command-errors",
		},
		Trace: []Trace{{TimeUnit: "unix_nano", Timestamp: 42, Location: "connector"}},
	}
	msg.Metadata.SetOutputAspectNodes([]models.AspectNode{nodeB, nodeA})

	expected := fmt.Sprintf(`{"request":{"input":{"data":"{\"level\":1}"}},"response":{"output":{"data":"{\"level\":2}"}},"task_info":{"worker_id":"","task_id":"task-1","process_instance_id":"","process_definition_id":"","completion_strategy":"","time":"1700000000","tenant_id":"user-1","business_key":""},"metadata":{"version":3,"device":%s,"service":%s,"protocol":%s,"output_function_id":"f-1","output_aspect_node":%s,"output_aspect_nodes":[%s,%s],"input_characteristic":"c-in","output_characteristic":"c-out","response_to":"device-command-response","error_to":"device-command-errors"},"trace":[{"time_unit":"unix_nano","timestamp":42,"location":"connector"}]}`,
		mustJson(t, device), mustJson(t, service), mustJson(t, protocol), mustJson(t, nodeA), mustJson(t, nodeB), mustJson(t, nodeA))

	if actual := mustJson(t, msg); actual != expected {
		t.Errorf("\nexpected %s\nactual   %s", expected, actual)
	}
}

func TestEmptyProtocolMsgJsonIsUnchanged(t *testing.T) {
	expected := fmt.Sprintf(`{"request":{"input":null},"response":{"output":null},"task_info":{"worker_id":"","task_id":"","process_instance_id":"","process_definition_id":"","completion_strategy":"","time":"","tenant_id":"","business_key":""},"metadata":{"device":%s,"service":%s,"protocol":%s,"response_to":""}}`,
		mustJson(t, models.Device{}), mustJson(t, models.Service{}), mustJson(t, models.Protocol{}))
	if actual := mustJson(t, ProtocolMsg{}); actual != expected {
		t.Errorf("\nexpected %s\nactual   %s", expected, actual)
	}
}

// A connector that does not know the list answers with the deprecated single node only.
func TestOutputAspectNodesReadsDeprecatedNode(t *testing.T) {
	msg := ProtocolMsg{}
	err := json.Unmarshal([]byte(`{"metadata":{"output_aspect_node":{"id":"a"}}}`), &msg)
	if err != nil {
		t.Fatal(err)
	}
	actual := msg.Metadata.GetOutputAspectNodes()
	if !reflect.DeepEqual(actual, []models.AspectNode{{Id: "a"}}) {
		t.Errorf("%#v", actual)
	}
}

func TestSetOutputAspectNodesWithoutNodes(t *testing.T) {
	metadata := Metadata{OutputAspectNode: &models.AspectNode{Id: "old"}}
	metadata.SetOutputAspectNodes(nil)
	if metadata.OutputAspectNode != nil || metadata.OutputAspectNodes != nil {
		t.Errorf("%#v", metadata)
	}
}

func mustJson(t *testing.T, value any) string {
	t.Helper()
	result, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return string(result)
}
