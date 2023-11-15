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
	"github.com/SENERGY-Platform/external-task-worker/lib/devicerepository/model"
	"github.com/SENERGY-Platform/service-commons/pkg/cache"
	"github.com/SENERGY-Platform/service-commons/pkg/signal"
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
	cache           *cache.Cache
	config          configuration.Config
	cacheDevices    bool //no caching to ensure access check in repository
	lastUsedToken   string
	cacheExpiration time.Duration
}

func GetCacheConfig() cache.Config {
	return cache.Config{
		CacheInvalidationSignalHooks: map[cache.Signal]cache.ToKey{
			signal.Known.CacheInvalidationAll: nil,
			signal.Known.FunctionCacheInvalidation: func(signalValue string) (cacheKey string) {
				return "function." + signalValue
			},
			signal.Known.ConceptCacheInvalidation: func(signalValue string) (cacheKey string) {
				return "concept." + signalValue
			},
			signal.Known.DeviceCacheInvalidation: func(signalValue string) (cacheKey string) {
				return "device." + signalValue
			},
			signal.Known.ProtocolInvalidation: func(signalValue string) (cacheKey string) {
				return "protocol." + signalValue
			},
			signal.Known.DeviceTypeCacheInvalidation: func(signalValue string) (cacheKey string) {
				return "device-type." + signalValue
			},
			signal.Known.DeviceGroupInvalidation: func(signalValue string) (cacheKey string) {
				return "device-group." + signalValue
			},
			signal.Known.CharacteristicCacheInvalidation: func(signalValue string) (cacheKey string) {
				return "characteristics." + signalValue
			},
		},
	}
}

func IotFactory(ctx context.Context, config configuration.Config) (interfaces.Iot, error) {
	c, err := cache.New(GetCacheConfig())
	if err != nil {
		return nil, err
	}
	cacheExpiration := time.Minute
	if config.CacheExpiration != "" && config.CacheExpiration != "-" {
		cacheExpiration, err = time.ParseDuration(config.CacheExpiration)
		if err != nil {
			return nil, err
		}
	}
	return NewIot(config, c, false, cacheExpiration), nil
}

func NewIot(config configuration.Config, cache *cache.Cache, cacheDevices bool, cacheExpiration time.Duration) *Iot {
	return &Iot{config: config, cache: cache, cacheDevices: cacheDevices, cacheExpiration: cacheExpiration}
}

func (this *Iot) GetFunction(token string, id string) (result model.Function, err error) {
	return cache.Use(this.cache, "function."+id, func() (model.Function, error) {
		return this.getFunction(token, id)
	}, this.cacheExpiration)
}

func (this *Iot) getFunction(token string, id string) (result model.Function, err error) {
	err = this.GetJson(token, this.config.DeviceManagerUrl+"/functions/"+url.PathEscape(id), &result)
	return
}

func (this *Iot) GetConcept(token string, id string) (result model.Concept, err error) {
	return cache.Use(this.cache, "concept."+id, func() (model.Concept, error) {
		return this.getConcept(token, id)
	}, this.cacheExpiration)
}

func (this *Iot) getConcept(token string, id string) (result model.Concept, err error) {
	err = this.GetJson(token, this.config.DeviceManagerUrl+"/concepts/"+url.PathEscape(id), &result)
	return
}

func (this *Iot) GetDevice(token string, id string) (result model.Device, err error) {
	if this.cacheDevices {
		return cache.Use(this.cache, "device."+id, func() (model.Device, error) {
			return this.getDevice(token, id)
		}, this.cacheExpiration)
	}
	return this.getDevice(token, id)
}

func (this *Iot) getDevice(token string, id string) (result model.Device, err error) {
	err = this.GetJson(token, this.config.DeviceRepositoryUrl+"/devices/"+url.QueryEscape(id), &result)
	return
}

func (this *Iot) GetProtocol(token string, id string) (result model.Protocol, err error) {
	return cache.Use(this.cache, "protocol."+id, func() (model.Protocol, error) {
		return this.getProtocol(token, id)
	}, this.cacheExpiration)
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
	item, err := this.cache.Get("service." + id)
	if err != nil {
		return service, err
	}
	var ok bool
	service, ok = item.(model.Service)
	if !ok {
		err = errors.New("unable to interpret cache value as model.Service")
	}
	return service, err
}

func (this *Iot) saveServiceToCache(service model.Service) {
	_ = this.cache.Set("service."+service.Id, service, this.cacheExpiration)
}

func (this *Iot) GetDeviceType(token string, id string) (result model.DeviceType, err error) {
	return cache.Use(this.cache, "device-type."+id, func() (model.DeviceType, error) {
		return this.getDeviceType(token, id)
	}, this.cacheExpiration)
}

func (this *Iot) getDeviceType(token string, id string) (result model.DeviceType, err error) {
	err = this.GetJson(token, this.config.DeviceRepositoryUrl+"/device-types/"+url.QueryEscape(id), &result)
	return
}

func (this *Iot) GetDeviceGroup(token string, id string) (result model.DeviceGroup, err error) {
	return cache.Use(this.cache, "device-group."+id, func() (model.DeviceGroup, error) {
		return this.getDeviceGroup(token, id)
	}, this.cacheExpiration)
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
	return cache.Use(this.cache, "aspect-nodes."+id, func() (model.AspectNode, error) {
		return this.getAspectNode(token, id)
	}, this.cacheExpiration)
}

func (this *Iot) getAspectNode(token string, id string) (result model.AspectNode, err error) {
	err = this.GetJson(token, this.config.DeviceRepositoryUrl+"/aspect-nodes/"+url.QueryEscape(id), &result)
	return
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
	return cache.Use(this.cache, "characteristics."+id, func() (model.Characteristic, error) {
		return this.getCharacteristic(token, id)
	}, this.cacheExpiration)
}

func (this *Iot) getCharacteristic(token string, id string) (result model.Characteristic, err error) {
	err = this.GetJson(token, this.config.DeviceManagerUrl+"/characteristics/"+url.PathEscape(id), &result)
	return
}
