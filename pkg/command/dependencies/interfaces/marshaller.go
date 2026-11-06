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

package interfaces

import (
	"context"

	"github.com/SENERGY-Platform/device-command/pkg/configuration"
	"github.com/SENERGY-Platform/marshaller/lib/api/messages"
	"github.com/SENERGY-Platform/marshaller/lib/configurables"
	marshallermodel "github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
	model "github.com/SENERGY-Platform/models/go/models"
)

// the request types are the ones of the marshaller itself, so that a field added there does
// not have to be repeated here to be usable

type MarshallingV2RequestData = marshallermodel.MarshallingV2RequestData

type UnmarshallingV2Request = messages.UnmarshallingV2Request

type Configurable = configurables.Configurable

type Marshaller interface {
	MarshalFromServiceAndProtocol(characteristicId string, service model.Service, protocol model.Protocol, characteristicData interface{}, configurables []Configurable) (result map[string]string, err error)
	UnmarshalFromServiceAndProtocol(characteristicId string, service model.Service, protocol model.Protocol, message map[string]string, hints []string) (characteristicData interface{}, err error)

	MarshalV2(service model.Service, protocol model.Protocol, data []MarshallingV2RequestData) (result map[string]string, err error)
	UnmarshalV2(request UnmarshallingV2Request) (characteristicData interface{}, err error)
}

type MarshallerFactory func(ctx context.Context, config configuration.Config, iot Iot) (Marshaller, error)
