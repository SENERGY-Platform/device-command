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
	"github.com/SENERGY-Platform/marshaller/lib/api/messages"
	marshallerclient "github.com/SENERGY-Platform/marshaller/lib/client"
	model "github.com/SENERGY-Platform/models/go/models"
)

// Remote calls the marshaller service over http.
type Remote struct {
	client marshallerclient.Interface
}

func NewRemote(url string) *Remote {
	//no auth token: the marshaller is called service internally, not routed over the api-gateway
	return &Remote{client: marshallerclient.NewClient(url, nil)}
}

func (this *Remote) Marshal(service model.Service, protocol model.Protocol, data []MarshalData) (result map[string]string, err error) {
	result, err, _ = this.client.MarshalV2(messages.MarshallingV2Request{Service: service, Protocol: protocol, Data: data})
	return result, err
}

func (this *Remote) Unmarshal(request UnmarshalRequest) (result interface{}, err error) {
	result, err, _ = this.client.UnmarshalV2(request)
	return result, err
}
