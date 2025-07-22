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
	"github.com/SENERGY-Platform/device-repository/lib/client"
	"github.com/SENERGY-Platform/external-task-worker/lib/devicerepository/model"
	"github.com/SENERGY-Platform/service-commons/pkg/cache"
	"github.com/SENERGY-Platform/service-commons/pkg/signal"
	"io"
	"log"
	"net/http"
	"net/url"
	"runtime/debug"
	"strings"
	"time"
)

type Iot struct {
	cache           *cache.Cache
	config          configuration.Config
	cacheDevices    bool //no caching to ensure access check in repository
	lastUsedToken   string
	cacheExpiration time.Duration
	client          client.Interface
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
	return NewIotWithDeviceRepoClient(config, cache, cacheDevices, cacheExpiration, client.NewClient(config.DeviceRepositoryUrl, nil))
}

func NewIotWithDeviceRepoClient(config configuration.Config, cache *cache.Cache, cacheDevices bool, cacheExpiration time.Duration, client client.Interface) *Iot {
	return &Iot{config: config, cache: cache, cacheDevices: cacheDevices, cacheExpiration: cacheExpiration, client: client}
}

func (this *Iot) GetFunction(token string, id string) (result model.Function, err error) {
	return cache.Use(this.cache, "function."+id, func() (model.Function, error) {
		return this.getFunction(token, id)
	}, func(function model.Function) error {
		if function.Id == "" {
			return errors.New("invalid function loaded from cache")
		}
		return nil
	}, this.cacheExpiration)
}

func (this *Iot) getFunction(token string, id string) (result model.Function, err error) {
	result, err, _ = this.client.GetFunction(id)
	return
}

func (this *Iot) GetConcept(token string, id string) (result model.Concept, err error) {
	return cache.Use(this.cache, "concept."+id, func() (model.Concept, error) {
		return this.getConcept(token, id)
	}, func(concept model.Concept) error {
		if concept.Id == "" {
			return errors.New("invalid concept loaded from cache")
		}
		return nil
	}, this.cacheExpiration)
}

func (this *Iot) getConcept(token string, id string) (result model.Concept, err error) {
	result, err, _ = this.client.GetConceptWithoutCharacteristics(id)
	return
}

func (this *Iot) GetDevice(token string, id string) (result model.Device, err error) {
	if this.cacheDevices {
		return cache.Use(this.cache, "device."+id, func() (model.Device, error) {
			return this.getDevice(token, id)
		}, func(device model.Device) error {
			if device.Id == "" {
				return errors.New("invalid device loaded from cache")
			}
			return nil
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
	}, func(protocol model.Protocol) error {
		if protocol.Id == "" {
			return errors.New("invalid protocol loaded from cache")
		}
		return nil
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
			log.Println("ERROR: unable to load device-type", device.DeviceTypeId)
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
	service, err = cache.Get[model.Service](this.cache, "service."+id, func(service model.Service) error {
		if service.Id == "" {
			return errors.New("invalid service loaded from cache")
		}
		return nil
	})
	return service, err
}

func (this *Iot) saveServiceToCache(service model.Service) {
	_ = this.cache.Set("service."+service.Id, service, this.cacheExpiration)
}

func (this *Iot) GetDeviceType(token string, id string) (result model.DeviceType, err error) {
	return cache.Use(this.cache, "device-type."+id, func() (model.DeviceType, error) {
		return this.getDeviceType(token, id)
	}, func(deviceType model.DeviceType) error {
		if deviceType.Id == "" {
			return errors.New("invalid device-type loaded from cache")
		}
		return nil
	}, this.cacheExpiration)
}

func (this *Iot) getDeviceType(token string, id string) (result model.DeviceType, err error) {
	err = this.GetJson(token, this.config.DeviceRepositoryUrl+"/device-types/"+url.QueryEscape(id), &result)
	return
}

func (this *Iot) GetDeviceGroup(token string, id string) (result model.DeviceGroup, err error) {
	return cache.Use(this.cache, "device-group."+id, func() (model.DeviceGroup, error) {
		return this.getDeviceGroup(token, id)
	}, func(group model.DeviceGroup) error {
		if group.Id == "" {
			return errors.New("invalid device-group loaded from cache")
		}
		return nil
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
	if this.config.AuthEnabled() {
		req.Header.Set("Authorization", token)
	}
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := c.Do(req)
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
	}, func(node model.AspectNode) error {
		if node.Id == "" {
			return errors.New("invalid aspect-node loaded from cache")
		}
		return nil
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
	limit := 1000
	offset := 0
	temp := []model.Concept{}
	for len(temp) == limit || offset == 0 {
		temp, _, err, _ = this.client.ListConcepts(client.ConceptListOptions{
			Limit:  int64(limit),
			Offset: int64(offset),
			SortBy: "name.asc",
		})
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
	limit := 1000
	offset := 0
	temp := []model.Function{}
	for len(temp) == limit || offset == 0 {
		temp, _, err, _ = this.client.ListFunctions(client.FunctionListOptions{
			Limit:  int64(limit),
			Offset: int64(offset),
			SortBy: "name.asc",
		})
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
	}, func(characteristic model.Characteristic) error {
		if characteristic.Id == "" {
			return errors.New("invalid characteristic loaded from cache")
		}
		return nil
	}, this.cacheExpiration)
}

func (this *Iot) getCharacteristic(token string, id string) (result model.Characteristic, err error) {
	result, err, _ = this.client.GetCharacteristic(id)
	return
}
