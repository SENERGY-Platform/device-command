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
	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func mqttEnv(config configuration.Config, ctx context.Context, wg *sync.WaitGroup) (configuration.Config, error) {
	_, ip, err := Mqtt(ctx, wg)
	if err != nil {
		return config, err
	}
	config.MgwMqttBroker = "tcp://" + ip + ":1883"
	return config, nil
}

func Mqtt(ctx context.Context, wg *sync.WaitGroup) (hostPort string, ipAddress string, err error) {
	log.Println("start mqtt broker")
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:           "eclipse-mosquitto:1.6.12",
			ExposedPorts:    []string{"1883/tcp"},
			WaitingFor:      wait.ForListeningPort("1883/tcp"),
			AlwaysPullImage: true,
		},
		Started: true,
	})
	if err != nil {
		return "", "", err
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			log.Println("DEBUG: remove container", container.Terminate(context.Background()))
		}()
		<-ctx.Done()
		/*
			reader, err := container.Logs(context.Background())
			if err != nil {
				log.Println("ERROR: unable to get container log")
				return
			}
			buf := new(strings.Builder)
			io.Copy(buf, reader)
			fmt.Println("MOSQUITTO LOGS: ------------------------------------------")
			fmt.Println(buf.String())
			fmt.Println("\n---------------------------------------------------------------")
		//*/
	}()

	ipAddress, err = container.ContainerIP(ctx)
	if err != nil {
		return "", "", err
	}
	temp, err := container.MappedPort(ctx, "1883/tcp")
	if err != nil {
		return "", "", err
	}
	hostPort = temp.Port()

	return hostPort, ipAddress, err
}

// testMqttClient stands in for a connector on a gateway
type testMqttClient struct {
	client paho.Client
}

func newTestMqttClient(ctx context.Context, broker string, clientId string) (*testMqttClient, error) {
	client := paho.NewClient(paho.NewClientOptions().AddBroker(broker).SetClientID(clientId).SetCleanSession(true).SetOrderMatters(false))
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	go func() {
		<-ctx.Done()
		client.Disconnect(0)
	}()
	return &testMqttClient{client: client}, nil
}

func (this *testMqttClient) Subscribe(topic string, qos byte, handler func(topic string, payload []byte)) error {
	token := this.client.Subscribe(topic, qos, func(_ paho.Client, message paho.Message) {
		handler(message.Topic(), message.Payload())
	})
	token.Wait()
	return token.Error()
}

func (this *testMqttClient) Publish(topic string, qos byte, retained bool, payload []byte) error {
	token := this.client.Publish(topic, qos, retained, payload)
	token.Wait()
	return token.Error()
}
