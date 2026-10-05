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
	"time"

	"github.com/SENERGY-Platform/device-command/pkg/auth"
	"github.com/SENERGY-Platform/device-command/pkg/command/dependencies/impl/cloud"
	"github.com/SENERGY-Platform/device-command/pkg/command/dependencies/interfaces"
	"github.com/SENERGY-Platform/device-command/pkg/configuration"
	"github.com/SENERGY-Platform/device-repository/v2/lib/client"
	"github.com/SENERGY-Platform/service-commons/pkg/cache"
)

func IotFactory(ctx context.Context, config configuration.Config) (result interfaces.Iot, err error) {
	c, err := cache.New(cloud.GetCacheConfig())
	if err != nil {
		return result, err
	}
	cacheExpiration := time.Minute
	if config.CacheExpiration != "" && config.CacheExpiration != "-" {
		cacheExpiration, err = time.ParseDuration(config.CacheExpiration)
		if err != nil {
			return nil, err
		}
	}
	a := &auth.OpenidToken{}
	deviceRepoClient := client.NewClient(config.DeviceRepositoryUrl, func() (token string, err error) {
		return a.EnsureAccess(config)
	})
	overwriteAuthTokens := config.RequestUserIdp != "jwt"
	return cloud.NewIotWithDeviceRepoClient(config, c, true, cacheExpiration, deviceRepoClient, overwriteAuthTokens), nil
}
