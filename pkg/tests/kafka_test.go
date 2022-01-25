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

package tests

import (
	"context"
	"github.com/SENERGY-Platform/device-command/pkg/configuration"
	"github.com/SENERGY-Platform/external-task-worker/lib/test/docker"
	"github.com/ory/dockertest/v3"
	"sync"
)

func kafkaEnv(initialConfig configuration.Config, ctx context.Context, wg *sync.WaitGroup) (config configuration.Config, err error) {
	config = initialConfig

	pool, err := dockertest.NewPool("")
	if err != nil {
		return config, err
	}

	closeZk, _, zkIp, err := docker.Zookeeper(pool)
	if err != nil {
		return config, err
	}
	wg.Add(1)
	go func() {
		<-ctx.Done()
		closeZk()
		wg.Done()
	}()

	zookeeperUrl := zkIp + ":2181"

	//kafka
	kafkaUrl, closeKafka, err := docker.Kafka(pool, zookeeperUrl)
	if err != nil {
		return config, err
	}
	wg.Add(1)
	go func() {
		<-ctx.Done()
		closeKafka()
		wg.Done()
	}()

	config.KafkaUrl = kafkaUrl

	return config, nil
}
