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
	"sync"
	"time"

	"github.com/SENERGY-Platform/device-command/pkg/configuration"
	"github.com/SENERGY-Platform/service-commons/pkg/kafka"
	"github.com/SENERGY-Platform/service-commons/pkg/testing/docker"
)

func kafkaEnv(initialConfig configuration.Config, ctx context.Context, wg *sync.WaitGroup) (config configuration.Config, err error) {
	config = initialConfig
	config.InitTopics = true
	//config.KafkaConsumerGroup = uuid.NewString()

	//kafka
	kafkaUrl, err := docker.Kafka(ctx, wg)
	if err != nil {
		return config, err
	}

	config.KafkaUrl = kafkaUrl

	time.Sleep(5 * time.Second)

	return config, nil
}

// testProducer stands in for a protocol connector that answers on the topic the command names
type testProducer struct {
	ctx      context.Context
	kafkaUrl string
	mux      sync.Mutex
	writers  map[string]*kafka.Producer
}

func newTestProducer(ctx context.Context, kafkaUrl string) *testProducer {
	return &testProducer{ctx: ctx, kafkaUrl: kafkaUrl, writers: map[string]*kafka.Producer{}}
}

func (this *testProducer) ProduceWithKey(topic string, key string, message string) error {
	this.mux.Lock()
	writer, ok := this.writers[topic]
	if !ok {
		var err error
		writer, err = kafka.NewProducer(this.ctx, kafka.Config{KafkaUrl: this.kafkaUrl}, topic)
		if err != nil {
			this.mux.Unlock()
			return err
		}
		this.writers[topic] = writer
	}
	this.mux.Unlock()
	return writer.Produce(key, []byte(message))
}

// testConsumer reads the commands a protocol connector would get. The topic is created by the
// consumer, before device-command produces to it.
func testConsumer(ctx context.Context, kafkaUrl string, groupId string, topic string, listener func(msg []byte) error, onError func(err error)) error {
	return kafka.NewConsumer(ctx, kafka.Config{
		KafkaUrl:      kafkaUrl,
		ConsumerGroup: groupId,
		StartOffset:   kafka.FirstOffset,
		InitTopic:     true,
		OnError:       onError,
	}, topic, listener)
}
