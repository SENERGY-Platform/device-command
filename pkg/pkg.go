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
	"github.com/SENERGY-Platform/device-command/pkg/api"
	"github.com/SENERGY-Platform/device-command/pkg/command"
	"github.com/SENERGY-Platform/device-command/pkg/configuration"
)

func Start(ctx context.Context, config configuration.Config) error {
	cmd, err := command.New(ctx, config)
	if err != nil {
		return err
	}
	return api.Start(ctx, config, cmd)
}
