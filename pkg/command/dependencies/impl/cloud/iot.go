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
	"encoding/json"
	"errors"
	"github.com/SENERGY-Platform/device-command/pkg/command/dependencies/interfaces"
	"github.com/SENERGY-Platform/device-command/pkg/configuration"
	"github.com/SENERGY-Platform/external-task-worker/lib/devicerepository"
	"github.com/SENERGY-Platform/external-task-worker/lib/devicerepository/model"
	"io"
	"log"
	"net/http"
	"net/url"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

type Iot struct {
	cache         Cache
	config        configuration.Config
	cacheDevices  bool //no caching to ensure access check in repository
	lastUsedToken string
}

type Cache interface {
	Use(key string, getter func() (interface{}, error), result interface{}) (err error)
	Set(key string, value []byte)
	Get(key string) (value []byte, err error)
}

func IotFactory(ctx context.Context, config configuration.Config) (interfaces.Iot, error) {
	devicerepository.L1Size = config.DeviceRepoCacheSizeInMb * 1024 * 1024
	return NewIot(config, &CacheImpl{parent: devicerepository.NewCache()}, false), nil
}

func NewIot(config configuration.Config, cache Cache, cacheDevices bool) *Iot {
	return &Iot{config: config, cache: cache, cacheDevices: cacheDevices}
}

func (this *Iot) GetFunction(token string, id string) (result model.Function, err error) {
	err = this.cache.Use("function."+id, func() (interface{}, error) {
		return this.getFunction(token, id)
	}, &result)
	return
}

func (this *Iot) getFunction(token string, id string) (result model.Function, err error) {
	err = this.GetJson(token, this.config.DeviceManagerUrl+"/functions/"+url.PathEscape(id), &result)
	return
}

func (this *Iot) GetConcept(token string, id string) (result model.Concept, err error) {
	err = this.cache.Use("concept."+id, func() (interface{}, error) {
		return this.getConcept(token, id)
	}, &result)
	return
}

func (this *Iot) getConcept(token string, id string) (result model.Concept, err error) {
	err = this.GetJson(token, this.config.DeviceManagerUrl+"/concepts/"+url.PathEscape(id), &result)
	return
}

func (this *Iot) GetDevice(token string, id string) (result model.Device, err error) {
	if this.cacheDevices {
		err = this.cache.Use("device."+id, func() (interface{}, error) {
			return this.getDevice(token, id)
		}, &result)
		return result, err
	}
	return this.getDevice(token, id)
}

func (this *Iot) getDevice(token string, id string) (result model.Device, err error) {
	err = this.GetJson(token, this.config.DeviceRepositoryUrl+"/devices/"+url.QueryEscape(id), &result)
	return
}

func (this *Iot) GetProtocol(token string, id string) (result model.Protocol, err error) {
	err = this.cache.Use("protocol."+id, func() (interface{}, error) {
		return this.getProtocol(token, id)
	}, &result)
	return
}

func (this *Iot) getProtocol(token string, id string) (result model.Protocol, err error) {
	err = this.GetJson(token, this.config.DeviceRepositoryUrl+"/protocols/"+url.QueryEscape(id), &result)
	return
}

func (this *Iot) GetService(token string, device model.Device, id string) (result model.Service, err error) {
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
	value, err := this.cache.Get("service." + id)
	if err != nil {
		return service, err
	}
	err = json.Unmarshal(value, &service)
	return
}

func (this *Iot) saveServiceToCache(service model.Service) {
	buffer, _ := json.Marshal(service)
	this.cache.Set("service."+service.Id, buffer)
}

func (this *Iot) GetDeviceType(token string, id string) (result model.DeviceType, err error) {
	err = this.cache.Use("deviceType."+id, func() (interface{}, error) {
		return this.getDeviceType(token, id)
	}, &result)
	return
}

func (this *Iot) getDeviceType(token string, id string) (result model.DeviceType, err error) {
	err = this.GetJson(token, this.config.DeviceRepositoryUrl+"/device-types/"+url.QueryEscape(id), &result)
	return
}

func (this *Iot) GetDeviceGroup(token string, id string) (result model.DeviceGroup, err error) {
	err = this.cache.Use("deviceGroup."+id, func() (interface{}, error) {
		return this.getDeviceGroup(token, id)
	}, &result)
	return
}

func (this *Iot) getDeviceGroup(token string, id string) (result model.DeviceGroup, err error) {
	err = this.GetJson(token, this.config.DeviceRepositoryUrl+"/device-groups/"+url.QueryEscape(id), &result)
	return
}

func (this *Iot) GetJson(token string, endpoint string, result interface{}) (err error) {
	this.lastUsedToken = token
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", token)
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		temp, _ := io.ReadAll(resp.Body)
		return errors.New(strings.TrimSpace(string(temp)))
	}
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		log.Println("ERROR:", err.Error())
		debug.PrintStack()
	}
	return
}

func (this *Iot) GetAspectNode(token string, id string) (result model.AspectNode, err error) {
	err = this.cache.Use("aspect-nodes."+id, func() (interface{}, error) {
		return this.getAspectNode(token, id)
	}, &result)
	return
}

func (this *Iot) getAspectNode(token string, id string) (result model.AspectNode, err error) {
	err = this.GetJson(token, this.config.DeviceRepositoryUrl+"/aspect-nodes/"+url.QueryEscape(id), &result)
	return
}

func (this *Iot) GetLastUsedToken() string {
	return this.lastUsedToken
}

type IdWrapper struct {
	Id string `json:"id"`
}

func (this *Iot) GetConceptIds(token string) (ids []string, err error) {
	limit := 100
	offset := 0
	temp := []IdWrapper{}
	for len(temp) == limit || offset == 0 {
		temp = []IdWrapper{}
		err = this.GetJson(token, this.config.PermissionsUrl+"/v3/resources/concepts?limit="+strconv.Itoa(limit)+"&offset="+strconv.Itoa(offset)+"&sort=name.asc&rights=r", &temp)
		if err != nil {
			return ids, err
		}
		for _, wrapper := range temp {
			ids = append(ids, wrapper.Id)
		}
		offset = offset + limit
	}
	return ids, err
}

func (this *Iot) ListFunctions(token string) (functionInfos []model.Function, err error) {
	limit := 100
	offset := 0
	temp := []model.Function{}
	for len(temp) == limit || offset == 0 {
		temp = []model.Function{}
		endpoint := this.config.PermissionsUrl + "/v3/resources/functions?limit=" + strconv.Itoa(limit) + "&offset=" + strconv.Itoa(offset) + "&sort=name.asc&rights=r"
		err = this.GetJson(token, endpoint, &temp)
		if err != nil {
			return functionInfos, err
		}
		functionInfos = append(functionInfos, temp...)
		offset = offset + limit
	}
	return functionInfos, err
}

func (this *Iot) GetCharacteristic(token string, id string) (result model.Characteristic, err error) {
	err = this.cache.Use("characteristics."+id, func() (interface{}, error) {
		return this.getCharacteristic(token, id)
	}, &result)
	return
}

func (this *Iot) getCharacteristic(token string, id string) (result model.Characteristic, err error) {
	err = this.GetJson(token, this.config.DeviceManagerUrl+"/characteristics/"+url.PathEscape(id), &result)
	return
}
