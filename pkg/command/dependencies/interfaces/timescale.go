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
	"errors"
	"github.com/SENERGY-Platform/device-command/pkg/auth"
	"github.com/SENERGY-Platform/device-command/pkg/configuration"
	"github.com/SENERGY-Platform/external-task-worker/lib/devicerepository/model"
	"github.com/SENERGY-Platform/models/go/models"
	"time"
)

type Timescale interface {
	GetLastMessage(token auth.Token, device models.Device, service models.Service, protocol model.Protocol, timeout time.Duration) (result map[string]interface{}, err error)
}
type TimescaleFactory func(ctx context.Context, config configuration.Config) (Timescale, error)

var ErrMissingLastValue = errors.New("missing last value in mgw-last-value")
var ErrMissingLastValueCode = 513 //custom code to signify missing last-value in mgw-last-value
