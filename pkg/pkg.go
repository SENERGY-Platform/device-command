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

package pkg

import (
	"context"
	"github.com/SENERGY-Platform/device-command/pkg/tracing"

	"github.com/SENERGY-Platform/device-command/pkg/api"
	"github.com/SENERGY-Platform/device-command/pkg/command"
	"github.com/SENERGY-Platform/device-command/pkg/configuration"
	"github.com/SENERGY-Platform/device-command/pkg/devicerepo"
	"github.com/SENERGY-Platform/device-command/pkg/lastvalue"
	"github.com/SENERGY-Platform/device-command/pkg/marshaller"
	"github.com/SENERGY-Platform/device-command/pkg/metrics"
	"github.com/SENERGY-Platform/device-command/pkg/transport"
)

// Start runs the service until ctx is done. Every dependency picks its cloud or gateway variant by
// its own config switch: iot_impl, marshaller_impl, timescale_impl and com_impl.
func Start(ctx context.Context, config configuration.Config) error {
	config = config.WithScalingSuffix()
	err := tracing.Init(ctx, config.OtelEndpoint)
	if err != nil {
		config.GetLogger().Warn("unable to init open-telemetry, continuing without tracing", "error", err)
	}
	repo, err := devicerepo.New(ctx, config)
	if err != nil {
		return err
	}
	m, err := marshaller.New(ctx, config, repo)
	if err != nil {
		return err
	}
	lastValues, err := lastvalue.New(config)
	if err != nil {
		return err
	}
	replies := command.NewReplies()
	sender, err := transport.New(ctx, config, replies.Receive)
	if err != nil {
		return err
	}
	metrics := metrics.New()
	cmd := command.New(config, command.Dependencies{
		DeviceRepo: repo,
		Marshaller: m,
		LastValues: lastValues,
		Transport:  sender,
		Replies:    replies,
		Metrics:    metrics,
	})
	return api.Start(ctx, config, cmd, metrics)
}
