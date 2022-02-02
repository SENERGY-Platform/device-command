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
	"github.com/SENERGY-Platform/external-task-worker/lib/devicerepository/model"
	"github.com/SENERGY-Platform/external-task-worker/lib/messages"
	"net/http"
)

func (this *Command) HandleTaskResponse(message messages.ProtocolMsg) (err error) {
	var output interface{}
	if message.Metadata.OutputCharacteristic != model.NullCharacteristic.Id {
		output, err = this.marshaller.UnmarshalFromServiceAndProtocol(message.Metadata.OutputCharacteristic, message.Metadata.Service, message.Metadata.Protocol, message.Response.Output, message.Metadata.ContentVariableHints)
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
