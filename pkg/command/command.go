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

// Package command answers a command for a device or a device group: it marshals the input, sends
// the protocol message, waits for the reply of the connector and unmarshals its output. A
// measuring function of an event service is answered from the last value instead.
package command

import (
	"context"
	"net/http"

	"github.com/SENERGY-Platform/device-command/pkg/auth"
	"github.com/SENERGY-Platform/device-command/pkg/configuration"
	"github.com/SENERGY-Platform/device-command/pkg/devicerepo"
	"github.com/SENERGY-Platform/device-command/pkg/lastvalue"
	"github.com/SENERGY-Platform/device-command/pkg/marshaller"
	"github.com/SENERGY-Platform/device-command/pkg/metrics"
	"github.com/SENERGY-Platform/device-command/pkg/transport"
)

type Dependencies struct {
	DeviceRepo *devicerepo.Client
	Marshaller marshaller.Marshaller
	LastValues lastvalue.Reader
	Transport  transport.Sender
	// Replies receives what Transport reports back; see NewReplies.
	Replies *Replies
	Metrics *metrics.Metrics
}

type Command struct {
	config     configuration.Config
	repo       *devicerepo.Client
	marshaller marshaller.Marshaller
	lastValues lastvalue.Reader
	transport  transport.Sender
	replies    *Replies
	metrics    *metrics.Metrics
}

func New(config configuration.Config, deps Dependencies) *Command {
	return &Command{
		config:     config,
		repo:       deps.DeviceRepo,
		marshaller: deps.Marshaller,
		lastValues: deps.LastValues,
		transport:  deps.Transport,
		replies:    deps.Replies,
		metrics:    deps.Metrics,
	}
}

func (this *Command) Command(ctx context.Context, token auth.Token, cmd CommandMessage, timeout string, preferEventValue bool) (code int, resp interface{}) {
	cmd.SetAspectIds()
	if cmd.DeviceId != "" && cmd.ServiceId != "" {
		return this.DeviceCommand(ctx, token, cmd.DeviceId, cmd.ServiceId, cmd.FunctionId, cmd.AspectIds, cmd.Input, timeout, preferEventValue, cmd.CharacteristicId)
	}
	if cmd.GroupId != "" {
		return this.GroupCommand(ctx, token, cmd.GroupId, cmd.FunctionId, cmd.AspectIds, cmd.DeviceClassId, cmd.Input, timeout, preferEventValue, cmd.CharacteristicId)
	}
	return http.StatusBadRequest, "missing device_id, service_id or group_id"
}
