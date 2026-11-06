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

package cloud

import (
	"context"

	"github.com/SENERGY-Platform/device-command/pkg/command/dependencies/interfaces"
	"github.com/SENERGY-Platform/device-command/pkg/configuration"
	"github.com/SENERGY-Platform/marshaller/lib/api/messages"
	marshallerclient "github.com/SENERGY-Platform/marshaller/lib/client"
	model "github.com/SENERGY-Platform/models/go/models"
)

func MarshallerFactory(ctx context.Context, config configuration.Config, iot interfaces.Iot) (interfaces.Marshaller, error) {
	//no auth token: the marshaller is called service internally, not routed over the api-gateway
	return &Marshaller{client: marshallerclient.NewClient(config.MarshallerUrl, nil)}, nil
}

// Marshaller calls the marshaller service over http.
type Marshaller struct {
	client marshallerclient.Interface
}

func (this *Marshaller) MarshalV2(service model.Service, protocol model.Protocol, data []interfaces.MarshallingV2RequestData) (result map[string]string, err error) {
	result, err, _ = this.client.MarshalV2(messages.MarshallingV2Request{
		Service:  service,
		Protocol: protocol,
		Data:     data,
	})
	return result, err
}

func (this *Marshaller) UnmarshalV2(request interfaces.UnmarshallingV2Request) (characteristicData interface{}, err error) {
	characteristicData, err, _ = this.client.UnmarshalV2(request)
	return characteristicData, err
}

func (this *Marshaller) MarshalFromServiceAndProtocol(characteristicId string, service model.Service, protocol model.Protocol, characteristicData interface{}, configurables []interfaces.Configurable) (result map[string]string, err error) {
	result, err, _ = this.client.Marshal(messages.MarshallingRequest{
		CharacteristicId: characteristicId,
		Service:          service,
		Protocol:         &protocol,
		Configurables:    configurables,
		Data:             characteristicData,
	})
	return result, err
}

func (this *Marshaller) UnmarshalFromServiceAndProtocol(characteristicId string, service model.Service, protocol model.Protocol, message map[string]string, hints []string) (characteristicData interface{}, err error) {
	characteristicData, err, _ = this.client.Unmarshal(messages.UnmarshallingRequest{
		CharacteristicId:     characteristicId,
		Service:              service,
		Protocol:             &protocol,
		Message:              message,
		ContentVariableHints: hints,
	})
	return characteristicData, err
}
