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
	"github.com/SENERGY-Platform/external-task-worker/lib/devicerepository/model"
	"github.com/SENERGY-Platform/external-task-worker/lib/marshaller"
	"github.com/SENERGY-Platform/external-task-worker/lib/messages"
	"github.com/SENERGY-Platform/external-task-worker/util"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (this *Command) DeviceCommand(token auth.Token, deviceId string, serviceId string, functionId string, aspectId string, input interface{}, timeout string) (code int, resp interface{}) {
	code, resp = this.deviceCommand(token, deviceId, serviceId, functionId, aspectId, input, timeout)
	if code == http.StatusOK {
		resp = []interface{}{resp}
	}
	return code, resp
}

func (this *Command) deviceCommand(token auth.Token, deviceId string, serviceId string, functionId string, aspectId string, input interface{}, timeout string) (code int, resp interface{}) {
	timeoutDuration := this.config.DefaultTimeoutDuration
	var err error
	if timeout != "" {
		timeoutDuration, err = time.ParseDuration(timeout)
		if err != nil {
			timeoutDuration = this.config.DefaultTimeoutDuration
		}
	}

	device, err := this.iot.GetDevice(token.Jwt(), deviceId)
	if err != nil {
		return http.StatusInternalServerError, "unable to load device: " + err.Error()
	}
	service, err := this.iot.GetService(token.Jwt(), device, serviceId)
	if err != nil {
		return http.StatusInternalServerError, "unable to load service: " + err.Error()
	}

	function, err := this.iot.GetFunction(token.Jwt(), functionId)
	if err != nil {
		return http.StatusInternalServerError, "unable to load function: " + err.Error()
	}

	characteristicId := ""
	if function.ConceptId != "" {
		concept, err := this.iot.GetConcept(token.Jwt(), function.ConceptId)
		if err != nil {
			return http.StatusInternalServerError, "unable to load concept: " + err.Error()
		}
		characteristicId = concept.BaseCharacteristicId
	}

	protocol, err := this.iot.GetProtocol(token.Jwt(), service.ProtocolId)
	if err != nil {
		return http.StatusInternalServerError, "unable to load protocol: " + err.Error()
	}
	if service.Interaction == model.EVENT && isMeasuringFunctionId(functionId) {
		return this.GetLastEventValue(token, device, service, protocol, characteristicId, functionId)
	}

	var aspectNode *model.AspectNode
	if aspectId != "" {
		temp, err := this.iot.GetAspectNode(token.Jwt(), aspectId)
		if err != nil {
			return http.StatusInternalServerError, "unable to load aspect node: " + err.Error()
		}
		aspectNode = &temp
	}

	var inputCharacteristicId string
	var outputCharacteristicId string

	var inputFunctionId string
	var outputFunctionId string

	var inputAspectNode *model.AspectNode
	var outputAspectNode *model.AspectNode

	data := []marshaller.MarshallingV2RequestData{}

	if isControllingFunction(function) {
		inputCharacteristicId = characteristicId
		inputFunctionId = functionId
		if input != nil {
			data = []marshaller.MarshallingV2RequestData{
				{
					Value:            input,
					CharacteristicId: inputCharacteristicId,
					FunctionId:       inputFunctionId,
					AspectNode:       inputAspectNode,
				},
			}
		}
		inputAspectNode = aspectNode
	} else {
		outputCharacteristicId = characteristicId
		outputFunctionId = functionId
		outputAspectNode = aspectNode
	}

	marshalledInput, err := this.marshaller.MarshalV2(service, protocol, data)
	if err != nil {
		return http.StatusInternalServerError, "unable to marshal input: " + err.Error()
	}

	taskId := uuid.New().String()
	this.register.Register(taskId)

	protocolMessage := messages.ProtocolMsg{
		TaskInfo: messages.TaskInfo{
			TaskId:   taskId,
			Time:     strconv.FormatInt(util.TimeNow().Unix(), 10),
			TenantId: token.GetUserId(),
		},
		Request: messages.ProtocolRequest{
			Input: marshalledInput,
		},
		Metadata: messages.Metadata{
			Version:              3,
			Device:               device,
			Service:              service,
			Protocol:             protocol,
			OutputFunctionId:     outputFunctionId,
			OutputAspectNode:     outputAspectNode,
			InputCharacteristic:  inputCharacteristicId,
			OutputCharacteristic: outputCharacteristicId,
			ResponseTo:           this.config.MetadataResponseTo,
			ErrorTo:              this.config.MetadataErrorTo,
		},
		Trace: []messages.Trace{},
	}

	err = this.producer.SendCommand(protocolMessage)
	if err != nil {
		log.Println("ERROR:", err)
		this.register.Complete(taskId, http.StatusInternalServerError, "unable to produce message")
	}
	return this.register.WaitWithTimeout(taskId, timeoutDuration)
}

func isControllingFunction(function model.Function) bool {
	if function.RdfType == model.SES_ONTOLOGY_CONTROLLING_FUNCTION {
		return true
	}
	if strings.HasPrefix(function.Id, "urn:infai:ses:controlling-function:") {
		return true
	}
	return false
}
