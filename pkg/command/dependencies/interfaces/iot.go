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
	"github.com/SENERGY-Platform/external-task-worker/lib/devicerepository/model"
)

type Iot interface {
	ListFunctions(token string) (functionInfos []model.Function, err error)
	GetFunction(token string, id string) (result model.Function, err error)
	GetConcept(token string, id string) (result model.Concept, err error)
	GetCharacteristic(token string, id string) (result model.Characteristic, err error)
	GetDevice(token string, id string) (result model.Device, err error)
	GetProtocol(token string, id string) (result model.Protocol, err error)
	GetService(token string, device model.Device, id string) (result model.Service, err error)
	GetDeviceType(token string, id string) (result model.DeviceType, err error)
	GetDeviceGroup(token string, id string) (result model.DeviceGroup, err error)
	GetAspectNode(token string, id string) (model.AspectNode, error)
	GetConceptIds(token string) ([]string, error)
}

type IotFactory func(ctx context.Context, config configuration.Config) (Iot, error)
