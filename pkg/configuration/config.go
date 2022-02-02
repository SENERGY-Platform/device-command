/*
 * Copyright 2021 InfAI (CC SES)
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

package configuration

import (
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	ServerPort          string `json:"server_port"`
	Debug               bool   `json:"debug"`
	ResponseWorkerCount int64  `json:"response_worker_count"`
	MarshallerUrl       string `json:"marshaller_url"`
	DeviceManagerUrl    string `json:"device_manager_url"`
	DeviceRepositoryUrl string `json:"device_repository_url"`
	TimescaleWrapperUrl string `json:"timescale_wrapper_url"`

	KafkaUrl               string        `json:"kafka_url"`
	DefaultTimeout         string        `json:"default_timeout"`
	DefaultTimeoutDuration time.Duration `json:"-"`

	KafkaConsumerGroup string `json:"kafka_consumer_group"`
	ResponseTopic      string `json:"response_topic"`
	PermissionsUrl     string `json:"permissions_url"`
	GroupScheduler     string `json:"group_scheduler"`

	MetadataResponseTo string `json:"metadata_response_to"`

	AsyncFlushFrequency string `json:"async_flush_frequency"`
	AsyncCompression    string `json:"async_compression"`
	SyncCompression     string `json:"sync_compression"`
	Sync                bool   `json:"sync"`
	SyncIdempotent      bool   `json:"sync_idempotent"`
	PartitionNum        int64  `json:"partition_num"`
	ReplicationFactor   int64  `json:"replication_factor"`
	AsyncFlushMessages  int64  `json:"async_flush_messages"`

	KafkaConsumerMaxWait  string `json:"kafka_consumer_max_wait"`
	KafkaConsumerMinBytes int64  `json:"kafka_consumer_min_bytes"`
	KafkaConsumerMaxBytes int64  `json:"kafka_consumer_max_bytes"`

	MetadataErrorTo string `json:"metadata_error_to"`
	ErrorTopic      string `json:"error_topic"`

	TopicSuffixForScaling string `json:"topic_suffix_for_scaling"`

	DeviceRepoCacheSizeInMb int                            `json:"device_repo_cache_size_in_mb"`
	KafkaTopicConfigs       map[string][]kafka.ConfigEntry `json:"kafka_topic_configs"`
}

//loads config from json in location and used environment variables (e.g ZookeeperUrl --> ZOOKEEPER_URL)
func Load(location string) (config Config, err error) {
	file, error := os.Open(location)
	if error != nil {
		log.Println("error on config load: ", error)
		return config, error
	}
	decoder := json.NewDecoder(file)
	error = decoder.Decode(&config)
	if error != nil {
		log.Println("invalid config json: ", error)
		return config, error
	}
	handleEnvironmentVars(&config)
	config.DefaultTimeoutDuration, err = time.ParseDuration(config.DefaultTimeout)
	return config, err
}

var camel = regexp.MustCompile("(^[^A-Z]*|[A-Z]*)([A-Z][^A-Z]+|$)")

func fieldNameToEnvName(s string) string {
	var a []string
	for _, sub := range camel.FindAllStringSubmatch(s, -1) {
		if sub[1] != "" {
			a = append(a, sub[1])
		}
		if sub[2] != "" {
			a = append(a, sub[2])
		}
	}
	return strings.ToUpper(strings.Join(a, "_"))
}

// preparations for docker
func handleEnvironmentVars(config *Config) {
	configValue := reflect.Indirect(reflect.ValueOf(config))
	configType := configValue.Type()
	for index := 0; index < configType.NumField(); index++ {
		fieldName := configType.Field(index).Name
		envName := fieldNameToEnvName(fieldName)
		envValue := os.Getenv(envName)
		if envValue != "" {
			fmt.Println("use environment variable: ", envName, " = ", envValue)
			if configValue.FieldByName(fieldName).Kind() == reflect.Int64 {
				i, _ := strconv.ParseInt(envValue, 10, 64)
				configValue.FieldByName(fieldName).SetInt(i)
			}
			if configValue.FieldByName(fieldName).Kind() == reflect.String {
				configValue.FieldByName(fieldName).SetString(envValue)
			}
			if configValue.FieldByName(fieldName).Kind() == reflect.Bool {
				b, _ := strconv.ParseBool(envValue)
				configValue.FieldByName(fieldName).SetBool(b)
			}
			if configValue.FieldByName(fieldName).Kind() == reflect.Float64 {
				f, _ := strconv.ParseFloat(envValue, 64)
				configValue.FieldByName(fieldName).SetFloat(f)
			}
			if configValue.FieldByName(fieldName).Kind() == reflect.Slice {
				val := []string{}
				for _, element := range strings.Split(envValue, ",") {
					val = append(val, strings.TrimSpace(element))
				}
				configValue.FieldByName(fieldName).Set(reflect.ValueOf(val))
			}
			if configValue.FieldByName(fieldName).Kind() == reflect.Map {
				value := map[string]string{}
				for _, element := range strings.Split(envValue, ",") {
					keyVal := strings.Split(element, ":")
					key := strings.TrimSpace(keyVal[0])
					val := strings.TrimSpace(keyVal[1])
					value[key] = val
				}
				configValue.FieldByName(fieldName).Set(reflect.ValueOf(value))
			}
		}
	}
}
