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
	"github.com/SENERGY-Platform/device-command/pkg/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/SENERGY-Platform/device-command/pkg/auth"
	"github.com/SENERGY-Platform/device-command/pkg/marshaller"
	"github.com/SENERGY-Platform/device-command/pkg/messages"
	marshallermodel "github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
	model "github.com/SENERGY-Platform/models/go/models"
	"github.com/google/uuid"
)

func (this *Command) DeviceCommand(ctx context.Context, token auth.Token, deviceId string, serviceId string, functionId string, aspectIds []string, input interface{}, timeout string, preferEventValue bool, characteristicId string) (code int, resp interface{}) {
	code, resp, _ = this.deviceCommandWithTime(ctx, token, deviceId, serviceId, functionId, aspectIds, input, timeout, preferEventValue, characteristicId)
	return code, resp
}

// deviceCommandWithTime also returns when an answer read from the last values was stored (RFC3339),
// empty for an answer of the device itself.
func (this *Command) deviceCommandWithTime(ctx context.Context, token auth.Token, deviceId string, serviceId string, functionId string, aspectIds []string, input interface{}, timeout string, preferEventValue bool, characteristicId string) (code int, resp interface{}, valueTime string) {
	code, resp, valueTime = this.deviceCommand(ctx, token, deviceId, serviceId, functionId, aspectIds, input, timeout, preferEventValue, characteristicId)
	if code == http.StatusOK {
		resp = []interface{}{resp}
	}
	return code, resp, valueTime
}

func (this *Command) deviceCommand(ctx context.Context, token auth.Token, deviceId string, serviceId string, functionId string, aspectIds []string, input interface{}, timeout string, preferEventValue bool, characteristicId string) (code int, resp interface{}, valueTime string) {
	ctx, span := tracing.Start(ctx, "device command",
		attribute.String("device_id", deviceId),
		attribute.String("service_id", serviceId),
		attribute.String("function_id", functionId),
		attribute.StringSlice("aspect_ids", aspectIds))
	defer func() { tracing.EndWithStatus(span, code, resp) }()

	timeoutDuration := this.config.DefaultTimeoutDuration
	if timeout != "" {
		parsed, err := time.ParseDuration(timeout)
		if err == nil {
			timeoutDuration = parsed
		}
	}

	device, err := this.repo.GetDevice(ctx, token.Jwt(), deviceId)
	if err != nil {
		return http.StatusInternalServerError, "unable to load device: " + err.Error(), ""
	}
	service, err := this.repo.GetService(ctx, token.Jwt(), device, serviceId)
	if err != nil {
		return http.StatusInternalServerError, "unable to load service: " + err.Error(), ""
	}
	function, err := this.repo.GetFunction(ctx, functionId)
	if err != nil {
		return http.StatusInternalServerError, "unable to load function: " + err.Error(), ""
	}
	if characteristicId == "" && function.ConceptId != "" {
		concept, err := this.repo.GetConcept(ctx, function.ConceptId)
		if err != nil {
			return http.StatusInternalServerError, "unable to load concept: " + err.Error(), ""
		}
		characteristicId = concept.BaseCharacteristicId
	}
	protocol, err := this.repo.GetProtocol(ctx, token.Jwt(), service.ProtocolId)
	if err != nil {
		return http.StatusInternalServerError, "unable to load protocol: " + err.Error(), ""
	}
	aspectNodes := []model.AspectNode{}
	for _, aspectId := range aspectIds {
		aspectNode, err := this.repo.GetAspectNode(ctx, aspectId)
		if err != nil {
			return http.StatusInternalServerError, "unable to load aspect node: " + err.Error(), ""
		}
		aspectNodes = append(aspectNodes, aspectNode)
	}

	if isMeasuringFunctionId(functionId) && (service.Interaction == model.EVENT || (preferEventValue && service.Interaction == model.EVENT_AND_REQUEST)) {
		this.metrics.LogGetLastEventValue(token.GetUserId(), device.Id, service.Id, functionId)
		span.SetAttributes(attribute.String("answered_from", "last value"))
		return this.GetLastEventValue(ctx, token, device, service, protocol, characteristicId, functionId, aspectNodes, timeoutDuration)
	}

	var inputCharacteristicId, outputCharacteristicId, outputFunctionId string
	var outputAspectNodes []model.AspectNode
	data := []marshaller.MarshalData{}
	if isControllingFunction(function) {
		inputCharacteristicId = characteristicId
		if input != nil {
			data = append(data, marshaller.MarshalData{
				Value:            input,
				CharacteristicId: inputCharacteristicId,
				FunctionId:       functionId,
				AspectNodes:      aspectNodes,
			})
		}
	} else {
		outputCharacteristicId = characteristicId
		outputFunctionId = functionId
		outputAspectNodes = aspectNodes
	}

	marshalledInput, err := this.marshaller.Marshal(ctx, service, protocol, data)
	if err != nil {
		return http.StatusInternalServerError, "unable to marshal input: " + err.Error(), ""
	}

	taskId := uuid.New().String()
	span.SetAttributes(attribute.String("answered_from", "connector"), attribute.String("task_id", taskId))
	msg := messages.ProtocolMsg{
		TaskInfo: messages.TaskInfo{
			TaskId:   taskId,
			Time:     strconv.FormatInt(time.Now().Unix(), 10),
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
			InputCharacteristic:  inputCharacteristicId,
			OutputCharacteristic: outputCharacteristicId,
			ResponseTo:           this.config.MetadataResponseTo,
			ErrorTo:              this.config.MetadataErrorTo,
		},
		Trace: []messages.Trace{},
	}
	//keeps the deprecated single node filled, for a protocol handler that does not know the list yet
	msg.Metadata.SetOutputAspectNodes(outputAspectNodes)

	this.metrics.LogCommandSend(token.GetUserId(), device.Id, service.Id, functionId)

	replies := this.replies.expect(taskId)
	defer this.replies.forget(taskId)
	err = this.transport.Send(ctx, msg)
	if err != nil {
		this.config.GetLogger().Error("unable to produce command message", "error", err)
		return http.StatusInternalServerError, "unable to produce message", ""
	}
	span.AddEvent("command sent")

	select {
	case reply := <-replies:
		span.AddEvent("reply received", trace.WithAttributes(attribute.Bool("failed", reply.Failed)))
		if reply.Failed {
			return http.StatusInternalServerError, reply.Output, ""
		}
		if outputCharacteristicId == "" || outputCharacteristicId == marshallermodel.NullCharacteristic.Id {
			return http.StatusOK, nil, ""
		}
		output, err := this.marshaller.Unmarshal(ctx, marshaller.UnmarshalRequest{
			Service:          service,
			Protocol:         protocol,
			CharacteristicId: outputCharacteristicId,
			Message:          reply.Output,
			FunctionId:       outputFunctionId,
			AspectNodes:      outputAspectNodes,
		})
		if err != nil {
			return http.StatusInternalServerError, err.Error(), ""
		}
		return http.StatusOK, output, ""
	case <-time.After(timeoutDuration):
		return http.StatusRequestTimeout, "timeout", ""
	}
}

func isControllingFunction(function model.Function) bool {
	return function.RdfType == model.SES_ONTOLOGY_CONTROLLING_FUNCTION ||
		strings.HasPrefix(function.Id, model.CONTROLLING_FUNCTION_PREFIX)
}

func isMeasuringFunctionId(id string) bool {
	return strings.HasPrefix(id, model.MEASURING_FUNCTION_PREFIX)
}
