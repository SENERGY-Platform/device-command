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
	"encoding/json"
	"github.com/SENERGY-Platform/external-task-worker/lib/devicerepository/model"
	"github.com/SENERGY-Platform/external-task-worker/lib/messages"
	"log"
	"net/http"
	"runtime/debug"
)

func (this *Command) HandleTaskResponse(msg string) (err error) {
	var message messages.ProtocolMsg
	err = json.Unmarshal([]byte(msg), &message)
	if err != nil {
		log.Println("ERROR:", err)
		debug.PrintStack()
		this.register.Complete(message.TaskInfo.TaskId, http.StatusInternalServerError, "unable interpret response message")
		return nil
	}

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

func (this *Command) ErrorMessageHandler(msg string) error {
	var message messages.ProtocolMsg
	err := json.Unmarshal([]byte(msg), &message)
	if err != nil {
		log.Println("ERROR:", err)
		debug.PrintStack()
		return nil
	}
	this.register.Complete(message.TaskInfo.TaskId, http.StatusInternalServerError, msg)
	return nil
}
