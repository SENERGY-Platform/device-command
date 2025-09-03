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

package tests

import (
	"context"
	_ "embed"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"time"

	"github.com/SENERGY-Platform/device-command/pkg/configuration"
	"github.com/SENERGY-Platform/device-repository/lib/api"
	"github.com/SENERGY-Platform/device-repository/lib/client"
	repoconf "github.com/SENERGY-Platform/device-repository/lib/configuration"
	"github.com/SENERGY-Platform/device-repository/lib/database"
	"github.com/SENERGY-Platform/models/go/models"
	"github.com/SENERGY-Platform/service-commons/pkg/jwt"
)

func iotEnvSetExport[T any, F func(ctx context.Context, e T, sf func(T) error) error](ctx context.Context, key string, value interface{}, prefix string, setter F, permCb func(T) error) error {
	if strings.HasPrefix(key, prefix) {
		e := new(T)
		temp, err := json.Marshal(value)
		if err != nil {
			log.Println("ERROR: unable to marshal", prefix, err)
			return err
		}
		err = json.Unmarshal(temp, e)
		if err != nil {
			log.Println("ERROR: unable to unmarshal", prefix, err)
			return err
		}
		err = setter(ctx, *e, permCb)
		if err != nil {
			log.Println("ERROR: unable to set", prefix, err)
			return err
		}
	}
	return nil
}

func CreateAspectNodes(db database.Database, aspect models.Aspect, rootId string, parentId string, ancestors []string) (descendents []string, err error) {
	descendents = []string{}
	children := []string{}
	for _, sub := range aspect.SubAspects {
		children = append(children, sub.Id)
		temp, err := CreateAspectNodes(db, sub, rootId, aspect.Id, append(ancestors, aspect.Id))
		if err != nil {
			return descendents, err
		}
		descendents = append(descendents, temp...)
	}
	err = db.SetAspectNode(context.Background(), models.AspectNode{
		Id:            aspect.Id,
		Name:          aspect.Name,
		RootId:        rootId,
		ParentId:      parentId,
		ChildIds:      children,
		AncestorIds:   ancestors,
		DescendentIds: descendents,
	})
	return append(descendents, aspect.Id), err
}

func NilCallback[T any](T) error {
	return nil
}

func iotEnv(initialConfig configuration.Config, ctx context.Context, wg *sync.WaitGroup, export []byte) (config configuration.Config, c client.Interface, db database.Database, err error) {
	config = initialConfig

	config, err = authMock(config, ctx, wg)
	if err != nil {
		return config, c, db, err
	}

	c, db, err = client.NewTestClient()
	if err != nil {
		return config, c, db, err
	}
	time.Sleep(time.Second)

	exportStruct := map[string]interface{}{}
	err = json.Unmarshal(export, &exportStruct)
	if err != nil {
		return config, c, db, err
	}

	for k, v := range exportStruct {
		err = iotEnvSetExport(ctx, k, v, "/characteristics/", db.SetCharacteristic, NilCallback)
		if err != nil {
			return config, c, db, err
		}
		err = iotEnvSetExport(ctx, k, v, "/concepts/", db.SetConcept, NilCallback)
		if err != nil {
			return config, c, db, err
		}
		err = iotEnvSetExport(ctx, k, v, "/aspects/", func(ctx context.Context, aspect models.Aspect, sf func(a2 models.Aspect) error) error {
			err := db.SetAspect(ctx, aspect, sf)
			if err != nil {
				return err
			}
			_, err = CreateAspectNodes(db, aspect, aspect.Id, "", []string{})
			return err
		}, NilCallback)
		if err != nil {
			return config, c, db, err
		}
		err = iotEnvSetExport(ctx, k, v, "/functions/", db.SetFunction, NilCallback)
		if err != nil {
			return config, c, db, err
		}
		err = iotEnvSetExport(ctx, k, v, "/device-classes/", db.SetDeviceClass, NilCallback)
		if err != nil {
			return config, c, db, err
		}
		err = iotEnvSetExport(ctx, k, v, "/protocols/", db.SetProtocol, NilCallback)
		if err != nil {
			return config, c, db, err
		}
		err = iotEnvSetExport(ctx, k, v, "/device-types/", db.SetDeviceType, NilCallback)
		if err != nil {
			return config, c, db, err
		}

	}

	router := api.GetRouter(repoconf.Config{}, c)
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if config.ComImpl == "mgw" {
			token, err := jwt.Parse(request.Header.Get("Authorization"))
			if err != nil {
				log.Println("ERROR: unable to parse token", err)
				writer.WriteHeader(http.StatusUnauthorized)
				return
			}
			if token.GetUserId() == "mgw-fallback-token" {
				log.Println("ERROR: used fallback-token in", request.URL.String())
				writer.WriteHeader(http.StatusForbidden)
				return
			}
		}
		router.ServeHTTP(writer, request)
	}))
	wg.Add(1)
	go func() {
		<-ctx.Done()
		server.Close()
		wg.Done()
	}()
	config.DeviceRepositoryUrl = server.URL
	return config, c, db, err
}

//go:embed test_export_1.json
var export1 []byte

//go:embed test_export_2.json
var export2 []byte
