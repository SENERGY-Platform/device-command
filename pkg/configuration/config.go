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
	DeviceRepositoryUrl string `json:"device_repository_url"`

	CacheExpiration                 string   `json:"cache_expiration"`
	CacheInvalidationAllKafkaTopics []string `json:"cache_invalidation_all_kafka_topics"`
	DeviceKafkaTopic                string   `json:"device_kafka_topic"`
	DeviceGroupKafkaTopic           string   `json:"device_group_kafka_topic"`

	TimescaleWrapperUrl string `json:"timescale_wrapper_url"`
	TimescaleImpl       string `json:"timescale_impl"` //"mgw" || "cloud" defaults to "cloud"

	KafkaUrl               string        `json:"kafka_url"`
	DefaultTimeout         string        `json:"default_timeout"`
	DefaultTimeoutDuration time.Duration `json:"-"`

	KafkaConsumerGroup string `json:"kafka_consumer_group"`
	ResponseTopic      string `json:"response_topic"`
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

	TopicSuffixForScaling string `json:"topic_suffix_for_scaling"` //only for kafka & cloud

	KafkaTopicConfigs map[string][]kafka.ConfigEntry `json:"kafka_topic_configs"`

	MgwCorrelationIdPrefix string `json:"mgw_correlation_id_prefix"`
	MgwProtocolSegment     string `json:"mgw_protocol_segment"`
	MgwMqttBroker          string `json:"mgw_mqtt_broker"`
	MgwMqttClientId        string `json:"mgw_mqtt_client_id"`
	MgwMqttUser            string `json:"mgw_mqtt_user" config:"secret"`
	MgwMqttPw              string `json:"mgw_mqtt_pw" config:"secret"`
	ComImpl                string `json:"com_impl"`        //"mgw" || "cloud" defaults to "cloud"
	MarshallerImpl         string `json:"marshaller_impl"` //"mgw" || "cloud" defaults to "cloud"
	UseIotFallback         bool   `json:"use_iot_fallback"`
	IotFallbackFile        string `json:"iot_fallback_file"`

	MgwConceptRepoRefreshInterval int64 `json:"mgw_concept_repo_refresh_interval"` //in seconds

	OverwriteAuthToken       bool    `json:"overwrite_auth_token"`
	AuthExpirationTimeBuffer float64 `json:"auth_expiration_time_buffer"`
	AuthEndpoint             string  `json:"auth_endpoint"`
	AuthClientId             string  `json:"auth_client_id" config:"secret"`
	AuthUserName             string  `json:"auth_user_name" config:"secret"`
	AuthPassword             string  `json:"auth_password" config:"secret"`
	AuthFallbackToken        string  `json:"auth_fallback_token"`

	InitTopics bool `json:"init_topics"`
}

// loads config from json in location and used environment variables (e.g ZookeeperUrl --> ZOOKEEPER_URL)
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
		fieldConfig := configType.Field(index).Tag.Get("config")
		envName := fieldNameToEnvName(fieldName)
		envValue := os.Getenv(envName)
		if envValue != "" {
			loggedEnvValue := envValue
			if strings.Contains(fieldConfig, "secret") {
				loggedEnvValue = "***"
			}
			fmt.Println("use environment variable: ", envName, " = ", loggedEnvValue)
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
