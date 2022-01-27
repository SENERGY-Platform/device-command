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
	"github.com/SENERGY-Platform/external-task-worker/lib/camunda/interfaces"
	"github.com/SENERGY-Platform/external-task-worker/lib/com"
	"github.com/SENERGY-Platform/external-task-worker/lib/messages"
	"github.com/SENERGY-Platform/external-task-worker/util"
	"net/http"
)

//Get implements github.com/SENERGY-Platform/external-task-worker/lib/camunda/interfaces.FactoryInterface
func (this *Command) Get(_ util.Config, _ com.ProducerInterface) (interfaces.CamundaInterface, error) {
	return this, nil
}

func (this *Command) GetTasks() (tasks []messages.CamundaExternalTask, err error) {
	return []messages.CamundaExternalTask{}, nil
}

func (this *Command) CompleteTask(taskInfo messages.TaskInfo, _ string, output interface{}) (err error) {
	this.register.Complete(taskInfo.TaskId, http.StatusOK, output)
	return nil
}

func (this *Command) SetRetry(taskid string, tenantId string, number int64) {}

func (this *Command) Error(taskId string, processInstanceId string, processDefinitionId string, msg string, tenantId string) {
	this.register.Complete(taskId, http.StatusInternalServerError, msg)
}

func (this *Command) GetWorkerId() string {
	return "device-command"
}

func (this *Command) UnlockTask(taskInfo messages.TaskInfo) (err error) {
	this.register.Complete(taskInfo.TaskId, http.StatusInternalServerError, "")
	return nil
}
