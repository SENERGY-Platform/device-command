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

package command

import (
	"encoding/json"
	"errors"
	"github.com/SENERGY-Platform/device-command/pkg/configuration"
	"github.com/SENERGY-Platform/external-task-worker/lib/devicerepository"
	"github.com/SENERGY-Platform/external-task-worker/lib/devicerepository/model"
	"github.com/SENERGY-Platform/external-task-worker/util"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Iot struct {
	cache  *devicerepository.Cache
	config configuration.Config
}

//Get implements devicerepository.FactoryInterface for use in external-task-worker
func (this *Iot) Get(_ util.Config) devicerepository.RepoInterface {
	return this
}

func NewIot(config configuration.Config) *Iot {
	return &Iot{config: config, cache: devicerepository.NewCache()}
}

func (this *Iot) GetFunction(token string, id string) (result model.Function, err error) {
	err = this.cache.Use("function."+id, func() (interface{}, error) {
		return this.getFunction(token, id)
	}, &result)
	return
}

func (this *Iot) getFunction(token string, id string) (result model.Function, err error) {
	req, err := http.NewRequest("GET", this.config.DeviceManagerUrl+"/functions/"+url.PathEscape(id), nil)
	if err != nil {
		return result, err
	}
	req.Header.Set("Authorization", token)
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		temp, _ := io.ReadAll(resp.Body)
		return result, errors.New(strings.TrimSpace(string(temp)))
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	return
}

func (this *Iot) GetConcept(token string, id string) (result model.Concept, err error) {
	err = this.cache.Use("concept."+id, func() (interface{}, error) {
		return this.getConcept(token, id)
	}, &result)
	return
}

func (this *Iot) getConcept(token string, id string) (result model.Concept, err error) {
	req, err := http.NewRequest("GET", this.config.DeviceManagerUrl+"/concepts/"+url.PathEscape(id), nil)
	if err != nil {
		return result, err
	}
	req.Header.Set("Authorization", token)
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		temp, _ := io.ReadAll(resp.Body)
		return result, errors.New(string(temp))
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	return
}

func (this *Iot) GetToken(user string) (devicerepository.Impersonate, error) {
	item, err := this.cache.Get("user_token." + user)
	return devicerepository.Impersonate(item.Value), err
}

//StoreToken stores a user token which allows the external-task-worker to use GetToken()
func (this *Iot) StoreToken(user string, token string) {
	this.cache.SetWithExpiration("user_token."+user, []byte(token), int(this.config.TimeoutDuration.Seconds()))
}

func (this *Iot) GetDevice(token devicerepository.Impersonate, id string) (result model.Device, err error) {
	err = this.cache.Use("device."+id, func() (interface{}, error) {
		return this.getDevice(token, id)
	}, &result)
	return
}

func (this *Iot) getDevice(token devicerepository.Impersonate, id string) (result model.Device, err error) {
	err = token.GetJSON(this.config.DeviceRepositoryUrl+"/devices/"+url.QueryEscape(id), &result)
	return
}

func (this *Iot) GetProtocol(token devicerepository.Impersonate, id string) (result model.Protocol, err error) {
	err = this.cache.Use("protocol."+id, func() (interface{}, error) {
		return this.getProtocol(token, id)
	}, &result)
	return
}

func (this *Iot) getProtocol(token devicerepository.Impersonate, id string) (result model.Protocol, err error) {
	err = token.GetJSON(this.config.DeviceRepositoryUrl+"/protocols/"+url.QueryEscape(id), &result)
	return
}

func (this *Iot) GetService(token devicerepository.Impersonate, device model.Device, id string) (result model.Service, err error) {
	result, err = this.getServiceFromCache(id)
	if err != nil {
		dt, err := this.GetDeviceType(token, device.DeviceTypeId)
		if err != nil {
			log.Println("ERROR: unable to load device-type", device.DeviceTypeId, token)
			return result, err
		}
		for _, service := range dt.Services {
			if service.Id == id {
				this.saveServiceToCache(service)
				return service, nil
			}
		}
		log.Println("ERROR: unable to find service in device-type", device.DeviceTypeId, id)
		return result, errors.New("service not found")
	}
	return
}

func (this *Iot) getServiceFromCache(id string) (service model.Service, err error) {
	item, err := this.cache.Get("service." + id)
	if err != nil {
		return service, err
	}
	err = json.Unmarshal(item.Value, &service)
	return
}

func (this *Iot) saveServiceToCache(service model.Service) {
	buffer, _ := json.Marshal(service)
	this.cache.Set("service."+service.Id, buffer)
}

func (this *Iot) GetDeviceType(token devicerepository.Impersonate, id string) (result model.DeviceType, err error) {
	err = this.cache.Use("deviceType."+id, func() (interface{}, error) {
		return this.getDeviceType(token, id)
	}, &result)
	return
}

func (this *Iot) getDeviceType(token devicerepository.Impersonate, id string) (result model.DeviceType, err error) {
	err = token.GetJSON(this.config.DeviceRepositoryUrl+"/device-types/"+url.QueryEscape(id), &result)
	return
}

func (this *Iot) GetDeviceGroup(token devicerepository.Impersonate, id string) (result model.DeviceGroup, err error) {
	err = this.cache.Use("deviceGroup."+id, func() (interface{}, error) {
		return this.getDeviceGroup(token, id)
	}, &result)
	return
}

func (this *Iot) getDeviceGroup(token devicerepository.Impersonate, id string) (result model.DeviceGroup, err error) {
	err = token.GetJSON(this.config.DeviceRepositoryUrl+"/device-groups/"+url.QueryEscape(id), &result)
	return
}
