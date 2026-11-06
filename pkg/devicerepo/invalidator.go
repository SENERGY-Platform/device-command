/*
 * Copyright 2023 InfAI (CC SES)
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

package devicerepo

import (
	"context"
	"time"

	"github.com/SENERGY-Platform/device-command/pkg/configuration"
	"github.com/SENERGY-Platform/service-commons/pkg/cache/invalidator"
	"github.com/SENERGY-Platform/service-commons/pkg/kafka"
)

// startCacheInvalidator drops cached entries when the repository announces a change on kafka.
// Without kafka, as on a gateway, entries live until they expire.
func startCacheInvalidator(ctx context.Context, config configuration.Config) error {
	if config.KafkaUrl == "" || config.KafkaUrl == "-" {
		return nil
	}
	kafkaConf := kafka.Config{
		KafkaUrl:               config.KafkaUrl,
		StartOffset:            kafka.LastOffset,
		Debug:                  config.Debug,
		PartitionWatchInterval: time.Minute,
		OnError: func(err error) {
			config.GetLogger().Error("kafka consumer error for cache invalidator", "error", err)
		},
		InitTopic: config.InitTopics,
	}
	if len(config.CacheInvalidationAllKafkaTopics) > 0 {
		err := invalidator.StartCacheInvalidatorAll(ctx, kafkaConf, config.CacheInvalidationAllKafkaTopics, nil)
		if err != nil {
			return err
		}
	}
	return invalidator.StartKnownCacheInvalidators(ctx, kafkaConf, invalidator.KnownTopics{
		DeviceTopic:      config.DeviceKafkaTopic,
		DeviceGroupTopic: config.DeviceGroupKafkaTopic,
	}, nil)
}
