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

package cache

import (
	"encoding/json"
	"errors"
	"github.com/SENERGY-Platform/device-command/pkg/command/dependencies/impl/mgw/fallback"
	"github.com/coocood/freecache"
	"log"
	"runtime/debug"
)

var L1Expiration = 60         // 60sec
var L1Size = 40 * 1024 * 1024 //40MB
var Debug = false

type Cache struct {
	l1       *freecache.Cache
	fallback fallback.Fallback
}

type Item struct {
	Key   string
	Value []byte
}

var ErrNotFound = errors.New("key not found in cache")

func NewCacheWithFallback(fallback fallback.Fallback) *Cache {
	return &Cache{l1: freecache.NewCache(L1Size), fallback: fallback}
}

func (this *Cache) Get(key string) (value []byte, err error) {
	value, err = this.l1.Get([]byte(key))
	if err != nil && err != freecache.ErrNotFound {
		log.Println("ERROR: in Cache::l1.Get()", err)
	}
	return
}

func (this *Cache) Set(key string, value []byte) {
	err := this.l1.Set([]byte(key), value, L1Expiration)
	if err != nil {
		log.Println("ERROR: in Cache::l1.Set()", err)
	}
	return
}

func (this *Cache) Use(key string, getter func() (interface{}, error), result interface{}) (err error) {
	value, err := this.Get(key)
	if err == nil {
		err = json.Unmarshal(value, result)
		return
	}
	temp, err := getter()
	if err != nil {
		temp, err = this.fallback.Get(key)
		if err != nil {
			return err
		}
	} else {
		err = this.fallback.Set(key, temp)
		if err != nil {
			log.Println("WARNING: unable to store value in fallback storage", err)
			err = nil
		}
	}
	value, err = json.Marshal(temp)
	if err != nil {
		return err
	}
	this.Set(key, value)
	err = json.Unmarshal(value, &result)
	if err != nil {
		debug.PrintStack()
	}
	return err
}
