/*
 * Copyright 2024 InfAI (CC SES)
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

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func New() *Metrics {
	reg := prometheus.NewRegistry()

	result := &Metrics{
		httphandler: promhttp.HandlerFor(
			reg,
			promhttp.HandlerOpts{
				Registry: reg,
			},
		),
		commandsSendCountVec: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "device_command_commands_send_count_vec",
			Help: "counter vec for commands send to devices",
		}, []string{"user_id", "device_id", "service_id", "function_id"}),
		lastEventValueRequestCountVec: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "device_command_last_event_value_request_count_vec",
			Help: "counter vec for last-event-value requests",
		}, []string{"user_id", "device_id", "service_id", "function_id"}),
		requestsCountVec: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "device_command_requests_count_vec",
			Help: "counter vec for requests",
		}, []string{"user_id", "endpoint"}),
	}

	reg.MustRegister(
		result.commandsSendCountVec,
		result.lastEventValueRequestCountVec,
		result.requestsCountVec,
	)

	return result
}

type Metrics struct {
	httphandler http.Handler

	commandsSendCountVec          *prometheus.CounterVec
	lastEventValueRequestCountVec *prometheus.CounterVec
	requestsCountVec              *prometheus.CounterVec
}

func (this *Metrics) LogCommandSend(userId string, deviceId string, serviceId string, functionId string) {
	if this == nil {
		return
	}
	this.commandsSendCountVec.WithLabelValues(userId, deviceId, serviceId, functionId).Inc()
}

func (this *Metrics) LogGetLastEventValue(userId string, deviceId string, serviceId string, functionId string) {
	if this == nil {
		return
	}
	this.lastEventValueRequestCountVec.WithLabelValues(userId, deviceId, serviceId, functionId).Inc()
}

func (this *Metrics) LogRequest(userId string, endpoint string) {
	if this == nil {
		return
	}
	this.requestsCountVec.WithLabelValues(userId, endpoint).Inc()
}
