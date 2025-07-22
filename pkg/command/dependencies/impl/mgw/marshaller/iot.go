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

package marshaller

import (
	"context"
	"errors"
	"github.com/SENERGY-Platform/device-command/pkg/auth"
	"github.com/SENERGY-Platform/device-command/pkg/command/dependencies/interfaces"
	"github.com/SENERGY-Platform/device-command/pkg/configuration"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
	"net/http"
)

func NewMarshallerIot(ctx context.Context, conf configuration.Config, a *auth.OpenidToken, iot interfaces.Iot) (result *MarshallerIot, err error) {
	return &MarshallerIot{iot: iot, auth: a, config: conf}, nil
}

type MarshallerIot struct {
	iot    interfaces.Iot
	auth   *auth.OpenidToken
	config configuration.Config
}

func (this *MarshallerIot) GetAspectNode(id string) (result model.AspectNode, err error) {
	token, err := this.auth.EnsureAccess(this.config)
	if err != nil {
		return result, err
	}
	temp, err := this.iot.GetAspectNode(token, id)
	if err != nil {
		return result, err
	}
	err = jsonCast(temp, &result)
	return result, err
}

// this method should only be needed for old marshal/unmarshal requests
func (this MarshallerIot) GetDeviceType(id string) (result model.DeviceType, err error, code int) {
	return result, errors.New("not implemented"), http.StatusInternalServerError
}
