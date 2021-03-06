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

package mgw

import (
	"context"
	"github.com/SENERGY-Platform/device-command/pkg/command/dependencies/impl/cloud"
	"github.com/SENERGY-Platform/device-command/pkg/command/dependencies/impl/mgw/cache"
	"github.com/SENERGY-Platform/device-command/pkg/command/dependencies/impl/mgw/fallback"
	"github.com/SENERGY-Platform/device-command/pkg/command/dependencies/interfaces"
	"github.com/SENERGY-Platform/device-command/pkg/configuration"
)

func IotFactory(ctx context.Context, config configuration.Config) (result interfaces.Iot, err error) {
	fallback, err := fallback.NewFallback(config.IotFallbackFile)
	if err != nil {
		return result, err
	}
	return cloud.NewIot(config, cache.NewCacheWithFallback(fallback), true), nil
}
