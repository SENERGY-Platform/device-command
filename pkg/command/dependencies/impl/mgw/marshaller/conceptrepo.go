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
	"github.com/SENERGY-Platform/device-command/pkg/command/dependencies/interfaces"
	"github.com/SENERGY-Platform/device-command/pkg/configuration"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
	"log"
	"runtime/debug"
	"sync"
	"time"
)

type ConceptRepo struct {
	init                               bool
	iot                                interfaces.Iot
	concepts                           map[string]model.Concept
	characteristics                    map[string]model.Characteristic
	conceptByCharacteristic            map[string][]model.Concept
	rootCharacteristicByCharacteristic map[string]model.Characteristic
	characteristicsOfFunction          map[string][]string
	functionToConcept                  map[string]string
	mux                                sync.Mutex
}

type ConceptRepoDefault struct {
	Concept         model.Concept
	Characteristics []model.Characteristic
}

func NewConceptRepo(ctx context.Context, config configuration.Config, iot interfaces.Iot) (result *ConceptRepo, err error) {
	result = &ConceptRepo{
		iot:                                iot,
		concepts:                           map[string]model.Concept{},
		characteristics:                    map[string]model.Characteristic{},
		conceptByCharacteristic:            map[string][]model.Concept{},
		rootCharacteristicByCharacteristic: map[string]model.Characteristic{},
		characteristicsOfFunction:          map[string][]string{},
		functionToConcept:                  map[string]string{},
	}
	ticker := time.NewTicker(time.Duration(config.MgwConceptRepoRefreshInterval) * time.Second)
	go func() {
		defer ticker.Stop()
		<-ctx.Done()
	}()
	go func() {
		for range ticker.C {
			if result.iot.GetLastUsedToken() != "" {
				err = result.Load()
				if err != nil {
					log.Println("WARNING: unable to update concept repository", err)
				}
			}
		}
	}()
	return result, nil
}

func (this *ConceptRepo) ensureInit() {
	if !this.init {
		err := this.Load()
		if err != nil {
			log.Println("ERROR:", err)
			debug.PrintStack()
		}
	}
}

func (this *ConceptRepo) GetCharacteristicsOfFunction(functionId string) (characteristicIds []string, err error) {
	this.ensureInit()
	this.mux.Lock()
	defer this.mux.Unlock()
	var ok bool
	characteristicIds, ok = this.characteristicsOfFunction[functionId]
	if !ok {
		err = errors.New("unknown function-id")
	}
	return
}

func (this *ConceptRepo) GetConcept(id string) (concept model.Concept, err error) {
	this.ensureInit()
	this.mux.Lock()
	defer this.mux.Unlock()
	concept, ok := this.concepts[id]
	if !ok {
		debug.PrintStack()
		return concept, errors.New("no concept found for id " + id)
	}
	return concept, nil
}

func (this *ConceptRepo) GetConceptIdOfFunction(id string) string {
	this.mux.Lock()
	defer this.mux.Unlock()
	return this.functionToConcept[id]
}

func getCharacteristicDescendents(characteristic model.Characteristic) (result []model.Characteristic) {
	result = []model.Characteristic{characteristic}
	for _, child := range characteristic.SubCharacteristics {
		result = append(result, getCharacteristicDescendents(child)...)
	}
	return result
}

func (this *ConceptRepo) GetConceptsOfCharacteristic(characteristicId string) (conceptIds []string, err error) {
	this.ensureInit()
	this.mux.Lock()
	defer this.mux.Unlock()
	concepts, ok := this.conceptByCharacteristic[this.rootCharacteristicByCharacteristic[characteristicId].Id]
	if !ok {
		debug.PrintStack()
		return conceptIds, errors.New("no concept found for characteristic id " + characteristicId)
	}
	for _, concept := range concepts {
		conceptIds = append(conceptIds, concept.Id)
	}
	return conceptIds, nil
}

func (this *ConceptRepo) GetCharacteristic(id string) (characteristic model.Characteristic, err error) {
	this.ensureInit()
	if id == "" {
		return model.NullCharacteristic, nil
	}
	this.mux.Lock()
	defer this.mux.Unlock()
	characteristic, ok := this.characteristics[id]
	if !ok {
		debug.PrintStack()
		return characteristic, errors.New("no characteristic found for id " + id)
	}
	return characteristic, nil
}

func (this *ConceptRepo) GetRootCharacteristics(ids []string) (result []string) {
	this.ensureInit()
	this.mux.Lock()
	defer this.mux.Unlock()
	for _, id := range ids {
		root, ok := this.rootCharacteristicByCharacteristic[id]
		if ok {
			result = append(result, root.Id)
		}
	}
	return
}

func (this *ConceptRepo) registerFunction(f FunctionInfo) {
	concept, ok := this.concepts[f.ConceptId]
	if !ok {
		log.Println("WARNING: unable to register function with unknown concept", f)
		return
	}
	this.characteristicsOfFunction[f.Id] = concept.CharacteristicIds
	this.functionToConcept[f.Id] = f.ConceptId
}

func (this *ConceptRepo) Load() error {
	conceptIds, err := this.loadConceptIds()
	if err != nil {
		return err
	}

	type Temp struct {
		Concept         model.Concept
		Characteristics []model.Characteristic
	}
	temp := []Temp{}

	for _, conceptId := range conceptIds {
		concept, err := this.loadConcept(conceptId)
		if err != nil {
			return err
		}
		element := Temp{
			Concept: concept,
		}
		for _, characteristicId := range concept.CharacteristicIds {
			characteristic, err := this.loadCharacteristic(characteristicId)
			if err != nil {
				return err
			}
			element.Characteristics = append(element.Characteristics, characteristic)
		}
		temp = append(temp, element)
	}

	functionInfos, err := this.loadFunctions()
	if err != nil {
		return err
	}

	this.mux.Lock()
	defer this.mux.Unlock()

	this.resetToDefault()

	for _, element := range temp {
		this.register(element.Concept, element.Characteristics)
	}

	for _, f := range functionInfos {
		this.registerFunction(f)
	}
	this.init = true
	return nil
}

func (this *ConceptRepo) resetToDefault() {
	this.concepts = map[string]model.Concept{}
	this.characteristics = map[string]model.Characteristic{}
	this.conceptByCharacteristic = map[string][]model.Concept{}
	this.rootCharacteristicByCharacteristic = map[string]model.Characteristic{}
}

func (this *ConceptRepo) register(concept model.Concept, characteristics []model.Characteristic) {
	log.Println("load concept", concept.Name, concept.Id)
	for _, characteristic := range characteristics {
		log.Println("    load characteristic", characteristic.Name, characteristic.Id)
		concept.CharacteristicIds = append(concept.CharacteristicIds, characteristic.Id)
		this.characteristics[characteristic.Id] = characteristic
		this.conceptByCharacteristic[characteristic.Id] = append(this.conceptByCharacteristic[characteristic.Id], concept)
		this.rootCharacteristicByCharacteristic[characteristic.Id] = characteristic
		for _, descendent := range getCharacteristicDescendents(characteristic) {
			this.rootCharacteristicByCharacteristic[descendent.Id] = characteristic
		}
	}
	this.concepts[concept.Id] = concept
}

func (this *ConceptRepo) loadConceptIds() (ids []string, err error) {
	return this.iot.GetConceptIds(this.iot.GetLastUsedToken())
}

type FunctionInfo struct {
	Id        string `json:"id"`
	ConceptId string `json:"concept_id"`
}

func (this *ConceptRepo) loadFunctions() (functionInfos []FunctionInfo, err error) {
	temp, err := this.iot.ListFunctions(this.iot.GetLastUsedToken())
	for _, element := range temp {
		functionInfos = append(functionInfos, FunctionInfo{
			Id:        element.Id,
			ConceptId: element.ConceptId,
		})
	}
	return functionInfos, err
}

func (this *ConceptRepo) loadConcept(id string) (result model.Concept, err error) {
	concept, err := this.iot.GetConcept(this.iot.GetLastUsedToken(), id)
	if err != nil {
		return result, err
	}
	err = jsonCast(concept, &result)
	return result, err
}

func (this *ConceptRepo) loadCharacteristic(id string) (result model.Characteristic, err error) {
	characteristic, err := this.iot.GetCharacteristic(this.iot.GetLastUsedToken(), id)
	if err != nil {
		return result, err
	}
	err = jsonCast(characteristic, &result)
	return result, err
}
