/*
 * Copyright 2021 InfAI (CC SES)
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
	"fmt"
	"log/slog"

	converterService "github.com/SENERGY-Platform/converter/lib/converter"
	"github.com/SENERGY-Platform/device-command/pkg/configuration"
	"github.com/SENERGY-Platform/device-command/pkg/devicerepo"
	marshallerconfig "github.com/SENERGY-Platform/marshaller/lib/config"
	v2 "github.com/SENERGY-Platform/marshaller/lib/marshaller/v2"
	model "github.com/SENERGY-Platform/models/go/models"
)

// Local runs the marshaller in process, with concepts and characteristics read from the
// device-repository. A gateway uses it to stay independent of the platform.
type Local struct {
	v2 *v2.Marshaller
}

func NewLocal(ctx context.Context, config configuration.Config, repo *devicerepo.Client) (*Local, error) {
	converter, err := converterService.New()
	if err != nil {
		return nil, err
	}
	concepts := newConcepts(ctx, repo, config)
	return &Local{v2: v2.New(marshallerconfig.Config{}, converter, concepts)}, nil
}

func (this *Local) Marshal(service model.Service, protocol model.Protocol, data []MarshalData) (map[string]string, error) {
	return this.v2.Marshal(protocol, service, data)
}

// Unmarshal picks the output path by function and aspects if the request names none, the way the
// marshaller service does.
func (this *Local) Unmarshal(request UnmarshalRequest) (interface{}, error) {
	if request.Path == "" {
		paths := this.v2.GetOutputPaths(request.Service, request.FunctionId, request.GetAspectNodes())
		if len(paths) == 0 {
			return nil, errors.New("no output path found for criteria")
		}
		if len(paths) > 1 {
			slog.Warn("only first path found by FunctionId and AspectNodes is used for Unmarshal", "paths", fmt.Sprintf("%#v", paths))
		}
		request.Path = paths[0]
	}
	return this.v2.Unmarshal(request.Protocol, request.Service, request.CharacteristicId, request.Path, request.Message, nil)
}
