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
	"log"
	"sync"

	"github.com/SENERGY-Platform/device-command/pkg/configuration"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func lastValueEnv(config configuration.Config, ctx context.Context, wg *sync.WaitGroup) (configuration.Config, error) {
	port, err := MgwLastValue(ctx, wg, config.MgwMqttBroker)
	if err != nil {
		return config, err
	}
	config.TimescaleWrapperUrl = "http://localhost:" + port
	return config, nil
}

func MgwLastValue(ctx context.Context, wg *sync.WaitGroup, mqttBroker string) (hostPort string, err error) {
	log.Println("start mgw-last-value")
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "ghcr.io/senergy-platform/mgw-last-value:dev",
			ExposedPorts: []string{"8080/tcp"},
			Env: map[string]string{
				"MQTT_BROKER": mqttBroker,
				"DEBUG":       "true",
			},
			WaitingFor: wait.ForHTTP("/").WithPort("8080/tcp").WithStatusCodeMatcher(func(int) bool { return true }),
		},
		Started: true,
	})
	if err != nil {
		return "", err
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		log.Println("DEBUG: remove container mgw-last-value", container.Terminate(context.Background()))
	}()
	port, err := container.MappedPort(ctx, "8080/tcp")
	if err != nil {
		return "", err
	}
	return port.Port(), nil
}
