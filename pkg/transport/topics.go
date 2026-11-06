/*
 * Copyright 2026 InfAI (CC SES)
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

package transport

import (
	"net"
	"strconv"
	"strings"

	"github.com/segmentio/kafka-go"
)

// initTopics creates topics with the given partitions, replication factor and the entries of
// configMap that apply to them. It is used instead of the topic creation of service-commons,
// which creates compacted topics without expiry: right for a topic of current states, wrong
// for commands and responses. Taken over from external-task-worker v0.1.1 (lib/com/kafka).
func initTopics(bootstrapUrl string, configMap map[string][]kafka.ConfigEntry, numPartitions int, replicationFactor int, topics ...string) (err error) {
	conn, err := kafka.Dial("tcp", bootstrapUrl)
	if err != nil {
		return err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}
	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return err
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{}
	for _, topic := range topics {
		topicConfigs = append(topicConfigs, kafka.TopicConfig{
			Topic:             topic,
			NumPartitions:     numPartitions,
			ReplicationFactor: replicationFactor,
			ConfigEntries:     getTopicConfig(configMap, topic),
		})
	}
	return controllerConn.CreateTopics(topicConfigs...)
}

// getTopicConfig returns the config entries of topic: those listed under its own name, or else
// those of the first key that is a prefix of it.
func getTopicConfig(configMap map[string][]kafka.ConfigEntry, topic string) []kafka.ConfigEntry {
	if configMap == nil {
		return nil
	}
	result, exists := configMap[topic]
	if exists {
		return result
	}
	for key, conf := range configMap {
		if strings.HasPrefix(topic, key) {
			return conf
		}
	}
	return nil
}
