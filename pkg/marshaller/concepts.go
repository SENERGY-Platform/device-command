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
	"sync"
	"time"

	"github.com/SENERGY-Platform/device-command/pkg/configuration"
	"github.com/SENERGY-Platform/device-command/pkg/devicerepo"
	marshallermodel "github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
	model "github.com/SENERGY-Platform/models/go/models"
)

// concepts holds every concept, characteristic and function-to-concept relation in memory, as the
// in-process marshaller asks for them on every value it converts. It loads on first use and
// reloads every mgw_concept_repo_refresh_interval seconds.
type concepts struct {
	ctx    context.Context
	repo   *devicerepo.Client
	config configuration.Config

	loadMux sync.Mutex
	loaded  bool

	mux               sync.RWMutex
	concepts          map[string]model.Concept
	characteristics   map[string]model.Characteristic //including every sub characteristic
	conceptOfFunction map[string]string
}

func newConcepts(ctx context.Context, repo *devicerepo.Client, config configuration.Config) *concepts {
	result := &concepts{ctx: ctx, repo: repo, config: config}
	if config.MgwConceptRepoRefreshInterval > 0 {
		go result.refresh(ctx, time.Duration(config.MgwConceptRepoRefreshInterval)*time.Second)
	}
	return result
}

func (this *concepts) refresh(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := this.load()
			if err != nil {
				this.config.GetLogger().Warn("unable to update concepts", "error", err)
			}
		}
	}
}

// ensureLoaded loads once; a failed load is repeated on the next use.
func (this *concepts) ensureLoaded() {
	this.loadMux.Lock()
	defer this.loadMux.Unlock()
	if this.loaded {
		return
	}
	err := this.load()
	if err != nil {
		this.config.GetLogger().Error("unable to load concepts", "error", err)
		return
	}
	this.loaded = true
}

func (this *concepts) load() error {
	conceptIds, err := this.repo.ListConceptIds(this.ctx)
	if err != nil {
		return err
	}
	concepts := map[string]model.Concept{}
	characteristics := map[string]model.Characteristic{}
	for _, conceptId := range conceptIds {
		concept, err := this.repo.GetConcept(this.ctx, conceptId)
		if err != nil {
			return err
		}
		concepts[concept.Id] = concept
		for _, characteristicId := range concept.CharacteristicIds {
			characteristic, err := this.repo.GetCharacteristic(this.ctx, characteristicId)
			if err != nil {
				return err
			}
			addWithDescendants(characteristics, characteristic)
		}
	}
	functions, err := this.repo.ListFunctions(this.ctx)
	if err != nil {
		return err
	}
	conceptOfFunction := map[string]string{}
	for _, function := range functions {
		if function.ConceptId == "" {
			continue
		}
		if _, ok := concepts[function.ConceptId]; !ok {
			this.config.GetLogger().Warn("unable to register function with unknown concept", "functionId", function.Id, "conceptId", function.ConceptId)
			continue
		}
		conceptOfFunction[function.Id] = function.ConceptId
	}

	this.mux.Lock()
	defer this.mux.Unlock()
	this.concepts = concepts
	this.characteristics = characteristics
	this.conceptOfFunction = conceptOfFunction
	this.config.GetLogger().Info("concepts loaded", "concepts", len(concepts), "characteristics", len(characteristics), "functions", len(conceptOfFunction))
	return nil
}

func addWithDescendants(index map[string]model.Characteristic, characteristic model.Characteristic) {
	index[characteristic.Id] = characteristic
	for _, sub := range characteristic.SubCharacteristics {
		addWithDescendants(index, sub)
	}
}

func (this *concepts) GetConcept(id string) (model.Concept, error) {
	this.ensureLoaded()
	this.mux.RLock()
	defer this.mux.RUnlock()
	concept, ok := this.concepts[id]
	if !ok {
		return concept, errors.New("no concept found for id " + id)
	}
	return concept, nil
}

func (this *concepts) GetCharacteristic(id string) (model.Characteristic, error) {
	if id == "" {
		return marshallermodel.NullCharacteristic, nil
	}
	this.ensureLoaded()
	this.mux.RLock()
	defer this.mux.RUnlock()
	characteristic, ok := this.characteristics[id]
	if !ok {
		return characteristic, errors.New("no characteristic found for id " + id)
	}
	return characteristic, nil
}

func (this *concepts) GetConceptIdOfFunction(id string) string {
	this.ensureLoaded()
	this.mux.RLock()
	defer this.mux.RUnlock()
	return this.conceptOfFunction[id]
}
