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
	"errors"
	"fmt"
	"github.com/SENERGY-Platform/device-command/pkg/command/iot"
	"github.com/SENERGY-Platform/device-command/pkg/configuration"
	"github.com/SENERGY-Platform/device-command/pkg/register"
	"github.com/SENERGY-Platform/external-task-worker/lib/com"
	"github.com/SENERGY-Platform/external-task-worker/lib/com/comswitch"
	"github.com/SENERGY-Platform/external-task-worker/lib/devicerepository"
	"github.com/SENERGY-Platform/external-task-worker/lib/devicerepository/model"
	"github.com/SENERGY-Platform/external-task-worker/lib/marshaller"
	"github.com/SENERGY-Platform/external-task-worker/util"
	"log"
	"strings"
)

type Command struct {
	iot        Iot
	register   *register.Register
	config     configuration.Config
	marshaller marshaller.Interface
	producer   com.ProducerInterface
}

type Iot interface {
	GetFunction(token string, id string) (result model.Function, err error)
	GetConcept(token string, id string) (result model.Concept, err error)
	GetDevice(token string, id string) (result model.Device, err error)
	GetProtocol(token string, id string) (result model.Protocol, err error)
	GetService(token string, device model.Device, id string) (result model.Service, err error)
	GetDeviceType(token string, id string) (result model.DeviceType, err error)
	GetDeviceGroup(token string, id string) (result model.DeviceGroup, err error)
}

type IotFactory func(config configuration.Config) (Iot, error)

func DefaultIotFactory(config configuration.Config) (Iot, error) {
	return iot.NewIot(config), nil
}

func New(ctx context.Context, config configuration.Config) (cmd *Command, err error) {
	return NewWithFactories(ctx, config, comswitch.Factory, marshaller.Factory, DefaultIotFactory)
}

func NewWithFactories(ctx context.Context, config configuration.Config, comFactory com.FactoryInterface, marshallerFactory marshaller.FactoryInterface, iotFactory IotFactory) (cmd *Command, err error) {
	cmd = &Command{
		config:     config,
		register:   register.New(config.DefaultTimeoutDuration, config.Debug),
		marshaller: marshallerFactory.New(config.MarshallerUrl),
	}
	cmd.iot, err = iotFactory(config)
	if err != nil {
		return cmd, err
	}
	libConfig := createLibConfig(config)
	if config.ResponseWorkerCount > 1 {
		err = comFactory.NewConsumer(ctx, libConfig, cmd.GetQueuedResponseHandler(ctx, config.ResponseWorkerCount, config.ResponseWorkerCount), cmd.ErrorMessageHandler)
	} else {
		err = comFactory.NewConsumer(ctx, libConfig, cmd.HandleTaskResponse, cmd.ErrorMessageHandler)
	}
	if err != nil {
		return cmd, err
	}
	cmd.producer, err = comFactory.NewProducer(ctx, libConfig)
	if err != nil {
		return cmd, err
	}

	return cmd, nil
}

func (this *Command) GetQueuedResponseHandler(ctx context.Context, workerCount int64, queueSize int64) func(msg string) (err error) {
	queue := make(chan string, queueSize)
	for i := int64(0); i < workerCount; i++ {
		go func() {
			for msg := range queue {
				err := this.HandleTaskResponse(msg)
				if err != nil {
					log.Println("ERROR: ", err)
				}
			}
		}()
	}
	go func() {
		<-ctx.Done()
		close(queue)
	}()
	return func(msg string) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = errors.New(fmt.Sprint(r))
			}
		}()
		queue <- msg
		return err
	}
}

func createLibConfig(config configuration.Config) util.Config {
	devicerepository.L1Size = config.DeviceRepoCacheSizeInMb * 1024 * 1024
	return util.Config{
		Debug:                           config.Debug,
		DeviceRepoUrl:                   config.DeviceRepositoryUrl,
		CompletionStrategy:              util.PESSIMISTIC,
		OptimisticTaskCompletionTimeout: 100,
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
		ResponseWorkerCount:             config.ResponseWorkerCount,
		MetadataErrorTo:                 config.MetadataErrorTo,
		ErrorTopic:                      config.ErrorTopic,
		KafkaTopicConfigs:               config.KafkaTopicConfigs,
	}
}

func isMeasuringFunctionId(id string) bool {
	if strings.HasPrefix(id, model.MEASURING_FUNCTION_PREFIX) {
		return true
	}
	return false
}
