/*
 * Copyright 2026 InfAI (CC SES)
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

// Package marshaller translates between characteristic values and protocol messages: remotely
// through the marshaller service in the cloud, or in process on a gateway.
package marshaller

import (
	"context"

	"github.com/SENERGY-Platform/device-command/pkg/configuration"
	"github.com/SENERGY-Platform/device-command/pkg/devicerepo"
	"github.com/SENERGY-Platform/marshaller/lib/api/messages"
	marshallermodel "github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
	model "github.com/SENERGY-Platform/models/go/models"
)

// The request types are the ones of the marshaller itself, so that a field added there is usable
// here without being repeated.
type (
	MarshalData      = marshallermodel.MarshallingV2RequestData
	UnmarshalRequest = messages.UnmarshallingV2Request
)

type Marshaller interface {
	Marshal(ctx context.Context, service model.Service, protocol model.Protocol, data []MarshalData) (map[string]string, error)
	Unmarshal(ctx context.Context, request UnmarshalRequest) (interface{}, error)
}

// New returns the marshaller service client, or the in-process marshaller with
// marshaller_impl "mgw".
func New(ctx context.Context, config configuration.Config, repo *devicerepo.Client) (Marshaller, error) {
	if config.MarshallerImpl == "mgw" {
		return NewLocal(ctx, config, repo)
	}
	return NewRemote(config.MarshallerUrl), nil
}
