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
	"github.com/SENERGY-Platform/device-command/pkg/command/dependencies/interfaces"
	"github.com/SENERGY-Platform/device-command/pkg/messages"
	marshallermodel "github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
	"net/http"
)

func (this *Command) HandleTaskResponse(message messages.ProtocolMsg) (err error) {
	var output interface{}
	if message.Metadata.OutputCharacteristic != marshallermodel.NullCharacteristic.Id && message.Metadata.OutputCharacteristic != "" {
		output, err = this.marshaller.UnmarshalV2(interfaces.UnmarshallingV2Request{
			Service:          message.Metadata.Service,
			Protocol:         message.Metadata.Protocol,
			CharacteristicId: message.Metadata.OutputCharacteristic,
			Message:          message.Response.Output,
			FunctionId:       message.Metadata.OutputFunctionId,
			AspectNodes:      message.Metadata.GetOutputAspectNodes(),
		})
		if err != nil {
			this.register.Complete(message.TaskInfo.TaskId, http.StatusInternalServerError, err.Error())
			return nil
		}
	}
	this.register.Complete(message.TaskInfo.TaskId, http.StatusOK, output)
	return
}

func (this *Command) ErrorMessageHandler(message messages.ProtocolMsg) error {
	this.register.Complete(message.TaskInfo.TaskId, http.StatusInternalServerError, message.Response.Output)
	return nil
}
