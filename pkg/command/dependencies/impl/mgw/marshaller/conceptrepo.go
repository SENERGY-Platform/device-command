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
	"encoding/json"
	"errors"
	"github.com/SENERGY-Platform/converter/lib/converter/base"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
)

var FunctionToCharacteristics = map[string][]string{}

func SetCharacteristicsOfFunction(functionId string, characteristics []string) {
	if FunctionToCharacteristics == nil {
		FunctionToCharacteristics = map[string][]string{}
	}
	FunctionToCharacteristics[functionId] = characteristics
}

type ConceptRepo struct{}

func (this ConceptRepo) GetCharacteristicsOfFunction(functionId string) (characteristicIds []string, err error) {
	var ok bool
	characteristicIds, ok = FunctionToCharacteristics[functionId]
	if !ok {
		err = errors.New("not found")
	}
	return
}

func (this ConceptRepo) GetConcept(id string) (result model.Concept, err error) {
	temp, err := base.ConceptRepo.GetConcept(id)
	if err != nil {
		return result, err
	}
	b, _ := json.Marshal(temp)
	json.Unmarshal(b, &result)
	return
}

func (this ConceptRepo) GetConceptOfCharacteristic(characteristicId string) (conceptId string, err error) {
	return base.ConceptRepo.GetConceptOfCharacteristic(characteristicId)
}

func (this ConceptRepo) GetCharacteristic(id marshaller.CharacteristicId) (result model.Characteristic, err error) {
	temp, err := base.ConceptRepo.GetCharacteristic(id)
	if err != nil {
		return result, err
	}
	b, _ := json.Marshal(temp)
	json.Unmarshal(b, &result)
	return
}

func (this ConceptRepo) GetRootCharacteristics(ids []marshaller.CharacteristicId) (result []marshaller.CharacteristicId) {
	temp := base.ConceptRepo.GetRootCharacteristics(ids)
	b, _ := json.Marshal(temp)
	json.Unmarshal(b, &result)
	return
}
