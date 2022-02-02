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

package cloud

import (
	"context"
	"errors"
	"fmt"
	"github.com/SENERGY-Platform/device-command/pkg/command/dependencies/interfaces"
	"github.com/SENERGY-Platform/device-command/pkg/configuration"
	"github.com/SENERGY-Platform/external-task-worker/lib/com/kafka"
	"github.com/SENERGY-Platform/external-task-worker/lib/devicerepository"
	"github.com/SENERGY-Platform/external-task-worker/util"
	"log"
)

func ComFactory(ctx context.Context, config configuration.Config, responseListener func(msg string) error, errorListener func(msg string) error) (producer interfaces.Producer, err error) {
	comFactory := kafka.Factory
	libConfig := createLibConfig(config)
	if config.ResponseWorkerCount > 1 {
		err = comFactory.NewConsumer(ctx, libConfig, getQueuedResponseHandler(ctx, config.ResponseWorkerCount, config.ResponseWorkerCount, responseListener), errorListener)
	} else {
		err = comFactory.NewConsumer(ctx, libConfig, responseListener, errorListener)
	}
	if err != nil {
		return producer, err
	}
	producer, err = comFactory.NewProducer(ctx, libConfig)
	if err != nil {
		return producer, err
	}
	return producer, nil
}

func getQueuedResponseHandler(ctx context.Context, workerCount int64, queueSize int64, respHandler func(msg string) error) func(msg string) (err error) {
	queue := make(chan string, queueSize)
	for i := int64(0); i < workerCount; i++ {
		go func() {
			for msg := range queue {
				err := respHandler(msg)
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
		MetadataResponseTo:              config.MetadataResponseTo,
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
