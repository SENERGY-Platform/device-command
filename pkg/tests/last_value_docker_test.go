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
	"github.com/ory/dockertest/v3"
	"log"
	"net/http"
	"sync"
)

func lastValueEnv(config configuration.Config, ctx context.Context, wg *sync.WaitGroup) (configuration.Config, error) {
	port, _, err := MgwLastValue(ctx, wg, config.MgwMqttBroker)
	if err != nil {
		return config, err
	}
	config.TimescaleWrapperUrl = "http://localhost:" + port
	return config, nil
}

func MgwLastValue(ctx context.Context, wg *sync.WaitGroup, mqttBroker string) (hostPort string, ipAddress string, err error) {
	log.Println("start mgw-last-value")
	pool, err := dockertest.NewPool("")
	if err != nil {
		return "", "", err
	}
	container, err := pool.Run("ghcr.io/senergy-platform/mgw-last-value", "dev", []string{
		"MQTT_BROKER=" + mqttBroker,
		"DEBUG=true",
	})
	if err != nil {
		return "", "", err
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		log.Println("DEBUG: remove container " + container.Container.Name)
		log.Println(container.Close())
	}()
	//go Dockerlog(pool, ctx, container, "MGW-LAST-VALUE")
	hostPort = container.GetPort("8080/tcp")
	err = pool.Retry(func() error {
		_, err := http.Get("http://localhost:" + hostPort)
		return err
	})
	return hostPort, container.Container.NetworkSettings.IPAddress, err
}
