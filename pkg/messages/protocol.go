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

// Package messages holds the protocol message this service sends to a protocol connector and
// reads back from it. It is the contract with every connector of the platform, and this
// service is not its only writer: the automation-engine sends the same message for device
// tasks of a process. The types were taken over from external-task-worker v0.1.1
// (lib/messages), which was the shared definition until the worker was archived. A field
// changes its meaning for every connector at once, so the json form stays as it was there;
// protocol_test.go pins it.
package messages

import (
	"slices"
	"strings"

	marshallermodel "github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
	"github.com/SENERGY-Platform/models/go/models"
)

type TaskInfo struct {
	WorkerId            string `json:"worker_id"`
	TaskId              string `json:"task_id"`
	ProcessInstanceId   string `json:"process_instance_id"`
	ProcessDefinitionId string `json:"process_definition_id"`
	CompletionStrategy  string `json:"completion_strategy"`
	Time                string `json:"time"`
	TenantId            string `json:"tenant_id"`
	BusinessKey         string `json:"business_key"`
}

type ProtocolRequest struct {
	Input map[string]string `json:"input"`
}

type ProtocolResponse struct {
	Output map[string]string `json:"output"`
}

type Metadata struct {
	Version              int64               `json:"version,omitempty"`
	Device               models.Device       `json:"device"`
	Service              models.Service      `json:"service"`
	Protocol             models.Protocol     `json:"protocol"`
	OutputPath           string              `json:"output_path,omitempty"`         //only for version >= 3
	OutputFunctionId     string              `json:"output_function_id,omitempty"`  //only for version >= 3 if no OutputPath is known
	OutputAspectNode     *models.AspectNode  `json:"output_aspect_node,omitempty"`  //deprecated: please use OutputAspectNodes
	OutputAspectNodes    []models.AspectNode `json:"output_aspect_nodes,omitempty"` //only for version >= 3 if no OutputPath is known
	InputCharacteristic  string              `json:"input_characteristic,omitempty"`
	OutputCharacteristic string              `json:"output_characteristic,omitempty"`
	ContentVariableHints []string            `json:"content_variable_hints,omitempty"` //only for version < 3
	ResponseTo           string              `json:"response_to"`
	ErrorTo              string              `json:"error_to,omitempty"`
}

type ProtocolMsg struct {
	Request  ProtocolRequest  `json:"request"`
	Response ProtocolResponse `json:"response"`
	TaskInfo TaskInfo         `json:"task_info"`
	Metadata Metadata         `json:"metadata"`
	Trace    []Trace          `json:"trace,omitempty"`
}

type Trace struct {
	TimeUnit  string `json:"time_unit"`
	Timestamp int64  `json:"timestamp"`
	Location  string `json:"location"`
}

// GetOutputAspectNodes returns the aspects of the response the metadata describes. A
// connector that does not know the list yet echoes only the deprecated single node.
func (this Metadata) GetOutputAspectNodes() []models.AspectNode {
	if this.OutputAspectNode == nil {
		return this.OutputAspectNodes
	}
	return marshallermodel.AspectNodesAlias(*this.OutputAspectNode, this.OutputAspectNodes)
}

// SetOutputAspectNodes sets the aspects of the expected response and keeps the deprecated
// single node in sync, for a reader of the message that does not know the list yet.
func (this *Metadata) SetOutputAspectNodes(aspectNodes []models.AspectNode) {
	this.OutputAspectNodes = aspectNodes
	this.OutputAspectNode = deprecatedAspectNode(aspectNodes)
}

// deprecatedAspectNode picks the node a written message reports in its deprecated single
// aspect field. It is the node with the alphabetically first id, the way the
// device-repository fills the deprecated field of a path option.
func deprecatedAspectNode(aspectNodes []models.AspectNode) *models.AspectNode {
	if len(aspectNodes) == 0 {
		return nil
	}
	first := slices.MinFunc(aspectNodes, func(a models.AspectNode, b models.AspectNode) int {
		return strings.Compare(a.Id, b.Id)
	})
	return &first
}
