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
	"errors"
	"github.com/SENERGY-Platform/device-command/pkg/command/dependencies/impl/mgw/marshaller"
	"github.com/SENERGY-Platform/device-command/pkg/command/dependencies/interfaces"
	"github.com/SENERGY-Platform/device-command/pkg/configuration"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
	"net/http"
)

func MarshallerFactory(ctx context.Context, config configuration.Config) (interfaces.Marshaller, error) {
	return marshaller.NewMarshaller(config, MarshallerIot{}), nil
}

type MarshallerIot struct{}

//this method should only be needed for a path options request, which is not relevant for this project
func (this MarshallerIot) GetDeviceType(id string) (result model.DeviceType, err error, code int) {
	return result, errors.New("not implemented"), http.StatusInternalServerError
}