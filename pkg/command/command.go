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

package command

import (
	"context"
	"github.com/SENERGY-Platform/device-command/pkg/command/dependencies/impl/cloud"
	"github.com/SENERGY-Platform/device-command/pkg/command/dependencies/interfaces"
	"github.com/SENERGY-Platform/device-command/pkg/configuration"
	"github.com/SENERGY-Platform/device-command/pkg/register"
	"github.com/SENERGY-Platform/external-task-worker/lib/com"
	"github.com/SENERGY-Platform/external-task-worker/lib/devicerepository/model"
	"github.com/SENERGY-Platform/external-task-worker/lib/marshaller"
	"strings"
)

type Command struct {
	iot        interfaces.Iot
	timescale  interfaces.Timescale
	register   *register.Register
	config     configuration.Config
	marshaller marshaller.Interface
	producer   com.ProducerInterface
}

func New(ctx context.Context, config configuration.Config) (cmd *Command, err error) {
	return NewWithFactories(ctx, config, cloud.ComFactory, cloud.MarshallerFactory, cloud.IotFactory, cloud.TimescaleFactory)
}

func NewWithFactories(ctx context.Context, config configuration.Config, comFactory interfaces.ComFactory, marshallerFactory interfaces.MarshallerFactory, iotFactory interfaces.IotFactory, timescaleFactory interfaces.TimescaleFactory) (cmd *Command, err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		if err != nil {
			cancel()
		}
	}()
	cmd = &Command{
		config:   config,
		register: register.New(config.DefaultTimeoutDuration, config.Debug),
	}
	cmd.iot, err = iotFactory(ctx, config)
	if err != nil {
		return cmd, err
	}
	cmd.timescale, err = timescaleFactory(ctx, config)
	if err != nil {
		return cmd, err
	}
	cmd.marshaller, err = marshallerFactory(ctx, config)
	if err != nil {
		return cmd, err
	}
	cmd.producer, err = comFactory(ctx, config, cmd.HandleTaskResponse, cmd.ErrorMessageHandler)
	if err != nil {
		return cmd, err
	}
	return cmd, nil
}

func isMeasuringFunctionId(id string) bool {
	if strings.HasPrefix(id, model.MEASURING_FUNCTION_PREFIX) {
		return true
	}
	return false
}
