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
	"github.com/SENERGY-Platform/device-command/pkg/auth"
	"github.com/SENERGY-Platform/device-command/pkg/configuration"
	"github.com/SENERGY-Platform/device-command/pkg/register"
	"github.com/SENERGY-Platform/external-task-worker/lib"
	"github.com/SENERGY-Platform/external-task-worker/lib/com"
	"github.com/SENERGY-Platform/external-task-worker/lib/com/comswitch"
	"github.com/SENERGY-Platform/external-task-worker/lib/devicegroups"
	"github.com/SENERGY-Platform/external-task-worker/lib/devicerepository"
	"github.com/SENERGY-Platform/external-task-worker/lib/devicerepository/model"
	"github.com/SENERGY-Platform/external-task-worker/lib/marshaller"
	"github.com/SENERGY-Platform/external-task-worker/lib/messages"
	"github.com/SENERGY-Platform/external-task-worker/util"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

type Command struct {
	taskWorker *lib.CmdWorker
	iot        *Iot
	register   *register.Register
	config     configuration.Config
	marshaller marshaller.Interface
}

func New(ctx context.Context, config configuration.Config) *Command {
	return NewWithFactories(ctx, config, comswitch.Factory, marshaller.Factory)
}

func NewWithFactories(ctx context.Context, config configuration.Config, comFactory com.FactoryInterface, marshallerFactory marshaller.FactoryInterface) *Command {
	cmd := &Command{
		config:     config,
		iot:        NewIot(config),
		register:   register.New(config.TimeoutDuration, config.Debug),
		marshaller: marshallerFactory.New(config.MarshallerUrl),
	}
	cmd.taskWorker = lib.New(ctx, workerConfig(config), comFactory, cmd.iot, cmd, marshallerFactory)
	return cmd
}

func workerConfig(config configuration.Config) util.Config {
	devicegroups.LocalDbSize = config.PartialResultStoreSizeInMb * 1024 * 1024
	return util.Config{
		Debug:                           config.Debug,
		DeviceRepoUrl:                   config.DeviceRepositoryUrl,
		CompletionStrategy:              util.PESSIMISTIC,
		OptimisticTaskCompletionTimeout: 100,
		CamundaFetchLockDuration:        config.TimeoutDuration.Milliseconds(),
		KafkaUrl:                        config.KafkaUrl,
		KafkaConsumerGroup:              config.KafkaConsumerGroup,
		ResponseTopic:                   config.ResponseTopic,
		PermissionsUrl:                  config.PermissionsUrl,
		MarshallerUrl:                   config.MarshallerUrl,
		GroupScheduler:                  config.GroupScheduler,
		HttpCommandConsumerPort:         config.HttpCommandConsumerPort,
		HttpCommandConsumerSync:         config.HttpCommandConsumerSync,
		MetadataResponseTo:              config.MetadataResponseTo,
		DisableKafkaConsumer:            config.DisableKafkaConsumer,
		DisableHttpConsumer:             config.DisableHttpConsumer,
		AsyncFlushFrequency:             config.AsyncFlushFrequency,
		AsyncCompression:                config.AsyncCompression,
		SyncCompression:                 config.SyncCompression,
		Sync:                            config.Sync,
		SyncIdempotent:                  config.SyncIdempotent,
		PartitionNum:                    config.PartitionNum,
		ReplicationFactor:               config.ReplicationFactor,
		AsyncFlushMessages:              config.AsyncFlushMessages,
		KafkaConsumerMaxWait:            config.KafkaConsumerMaxWait,
		KafkaConsumerMinBytes:           config.KafkaConsumerMinBytes,
		KafkaConsumerMaxBytes:           config.KafkaConsumerMaxBytes,
		SubResultExpirationInSeconds:    config.SubResultExpirationInSeconds,
		SubResultDatabaseUrls:           config.SubResultDatabaseUrls,
		MemcachedTimeout:                config.MemcachedTimeout,
		MemcachedMaxIdleConns:           config.MemcachedMaxIdleConns,
		ResponseWorkerCount:             config.ResponseWorkerCount,
		MetadataErrorTo:                 config.MetadataErrorTo,
		ErrorTopic:                      config.ErrorTopic,
		KafkaTopicConfigs:               config.KafkaTopicConfigs,
	}
}

func (this *Command) DeviceCommand(token auth.Token, deviceId string, serviceId string, functionId string, input interface{}) (code int, resp interface{}) {
	this.iot.StoreToken(token.GetUserId(), token.Jwt())

	iotToken := devicerepository.Impersonate(token.Jwt())
	device, err := this.iot.GetDevice(iotToken, deviceId)
	if err != nil {
		return http.StatusInternalServerError, "unable to load device: " + err.Error()
	}
	service, err := this.iot.GetService(iotToken, device, serviceId)
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

	var protocol *model.Protocol
	if service.Interaction == model.EVENT && isMeasuringFunctionId(functionId) {
		temp, err := this.iot.GetProtocol(iotToken, service.ProtocolId)
		if err != nil {
			return http.StatusInternalServerError, "unable to load protocol: " + err.Error()
		}
		protocol = &temp
		return this.GetLastEventValue(token, device, service, temp, characteristicId)
	}

	taskId := uuid.New().String()
	this.register.Register(taskId)

	this.taskWorker.ExecuteCommand(messages.Command{
		Version:          2,
		Function:         function,
		CharacteristicId: characteristicId,
		Device:           &device,
		Service:          &service,
		Input:            input,
		Protocol:         protocol,
		Retries:          0,
	}, messages.CamundaExternalTask{
		Id:        taskId,
		Variables: nil,
		Retries:   0,
		TenantId:  token.GetUserId(),
	}, util.CALLER_CAMUNDA_LOOP)

	return this.register.Wait(taskId)
}

func (this *Command) GroupCommand(token auth.Token, groupId string, functionId string, input interface{}) (code int, resp interface{}) {
	this.iot.StoreToken(token.GetUserId(), token.Jwt())
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

	taskId := uuid.New().String()
	this.register.Register(taskId)

	this.taskWorker.ExecuteCommand(messages.Command{
		Version:          2,
		Function:         function,
		CharacteristicId: characteristicId,
		DeviceGroupId:    groupId,
		Input:            input,
		Retries:          0,
	}, messages.CamundaExternalTask{
		Id:        taskId,
		Variables: nil,
		Retries:   0,
		TenantId:  token.GetUserId(),
	}, util.CALLER_CAMUNDA_LOOP)

	return this.register.Wait(taskId)
}

func isMeasuringFunctionId(id string) bool {
	if strings.HasPrefix(id, model.MEASURING_FUNCTION_PREFIX) {
		return true
	}
	return false
}
