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

// Package devicerepo reads devices, services and the semantic model from the device-repository.
package devicerepo

import (
	"context"
	"errors"
	"time"

	"github.com/SENERGY-Platform/device-command/pkg/auth"
	"github.com/SENERGY-Platform/device-command/pkg/configuration"
	"github.com/SENERGY-Platform/device-repository/v3/lib/client"
	repomodel "github.com/SENERGY-Platform/device-repository/v3/lib/model"
	model "github.com/SENERGY-Platform/models/go/models"
)

// Client reads through the cache of the device-repository client, which keeps a response per
// request and per user, and drops it when the repository announces a change on kafka.
type Client struct {
	repo   client.Interface
	cached client.Interface
	// devices answers device reads: uncached in the cloud, where each read is a fresh access
	// check of the requesting user, cached on a gateway, which must not depend on the platform
	devices client.Interface
	// serviceToken is set if every read uses the credentials of the service instead of the token
	// of the requesting user: on a gateway that identifies its user without a jwt
	serviceToken bool
}

func New(ctx context.Context, config configuration.Config) (*Client, error) {
	ttl := time.Minute
	if config.CacheExpiration != "" && config.CacheExpiration != "-" {
		var err error
		ttl, err = time.ParseDuration(config.CacheExpiration)
		if err != nil {
			return nil, err
		}
	}
	options := client.Options{Timeout: 10 * time.Second}
	gateway := config.IotImpl == "mgw"
	if gateway {
		openid := &auth.OpenidToken{}
		//sent with every request that carries no token of its own
		options.OptionalAuthTokenForApiGatewayRequest = func() (string, error) { return openid.EnsureAccess(config) }
	}
	repo := client.NewClient(config.DeviceRepositoryUrl, options)
	result := &Client{
		repo:         repo,
		cached:       client.WithCache(repo, ttl),
		serviceToken: gateway && config.RequestUserIdp != "jwt",
	}
	result.devices = result.repo
	if gateway {
		result.devices = result.cached
	}
	unsubscribe := client.InvalidateOnSignals(repo, nil)
	go func() {
		<-ctx.Done()
		unsubscribe()
	}()
	err := startCacheInvalidator(ctx, config)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// token is the token a read is sent with; an empty one is replaced by the token of the service
func (this *Client) token(userToken string) string {
	if this.serviceToken {
		return ""
	}
	return userToken
}

func (this *Client) GetDevice(ctx context.Context, token string, id string) (result model.Device, err error) {
	result, err, _ = this.devices.ReadDevice(ctx, id, this.token(token), model.UnsetPermissionFlag)
	return result, err
}

func (this *Client) GetDeviceType(ctx context.Context, token string, id string) (result model.DeviceType, err error) {
	result, err, _ = this.cached.ReadDeviceType(ctx, id, this.token(token))
	return result, err
}

// GetService finds the service in the device type of the device, so that a service of another
// device type is refused.
func (this *Client) GetService(ctx context.Context, token string, device model.Device, id string) (model.Service, error) {
	deviceType, err := this.GetDeviceType(ctx, token, device.DeviceTypeId)
	if err != nil {
		return model.Service{}, err
	}
	for _, service := range deviceType.Services {
		if service.Id == id {
			return service, nil
		}
	}
	return model.Service{}, errors.New("service not found")
}

func (this *Client) GetProtocol(ctx context.Context, token string, id string) (result model.Protocol, err error) {
	result, err, _ = this.cached.ReadProtocol(ctx, id, this.token(token))
	return result, err
}

func (this *Client) GetDeviceGroup(ctx context.Context, token string, id string) (result model.DeviceGroup, err error) {
	result, err, _ = this.cached.ReadDeviceGroup(ctx, id, this.token(token), false)
	return result, err
}

func (this *Client) GetFunction(ctx context.Context, id string) (result model.Function, err error) {
	result, err, _ = this.cached.GetFunction(ctx, id)
	return result, err
}

func (this *Client) GetConcept(ctx context.Context, id string) (result model.Concept, err error) {
	result, err, _ = this.cached.GetConceptWithoutCharacteristics(ctx, id)
	return result, err
}

func (this *Client) GetCharacteristic(ctx context.Context, id string) (result model.Characteristic, err error) {
	result, err, _ = this.cached.GetCharacteristic(ctx, id)
	return result, err
}

func (this *Client) GetAspectNode(ctx context.Context, id string) (result model.AspectNode, err error) {
	result, err, _ = this.cached.GetAspectNode(ctx, id)
	return result, err
}

const listLimit = 1000

func (this *Client) ListConceptIds(ctx context.Context) (ids []string, err error) {
	for offset := int64(0); ; offset += listLimit {
		concepts, _, err, _ := this.cached.ListConcepts(ctx, repomodel.ConceptListOptions{Limit: listLimit, Offset: offset, SortBy: "name.asc"})
		if err != nil {
			return ids, err
		}
		for _, concept := range concepts {
			ids = append(ids, concept.Id)
		}
		if len(concepts) < listLimit {
			return ids, nil
		}
	}
}

func (this *Client) ListFunctions(ctx context.Context) (result []model.Function, err error) {
	for offset := int64(0); ; offset += listLimit {
		functions, _, err, _ := this.cached.ListFunctions(ctx, repomodel.FunctionListOptions{Limit: listLimit, Offset: offset, SortBy: "name.asc"})
		if err != nil {
			return result, err
		}
		result = append(result, functions...)
		if len(functions) < listLimit {
			return result, nil
		}
	}
}
