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

import "github.com/SENERGY-Platform/external-task-worker/lib/devicerepository"

type CacheImpl struct {
	parent *devicerepository.Cache
}

func (this *CacheImpl) Use(key string, getter func() (interface{}, error), result interface{}) (err error) {
	return this.parent.Use(key, getter, result)
}

func (this *CacheImpl) Set(key string, value []byte) {
	this.parent.Set(key, value)
}

func (this *CacheImpl) Get(key string) (value []byte, err error) {
	item, err := this.parent.Get(key)
	if err != nil {
		return value, err
	}
	return item.Value, nil
}
