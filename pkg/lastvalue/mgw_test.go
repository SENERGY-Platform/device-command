/*
 * Copyright 2025 InfAI (CC SES)
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

package lastvalue

import (
	"testing"
)

func TestNewMgw(t *testing.T) {
	obj, err := NewMgw("foobar", "data")
	if err != nil {
		t.Error(err)
		return
	}
	if obj.url != "http://foobar:8080" {
		t.Error(obj.url)
		return
	}

	obj, err = NewMgw("https://foobar:80", "data")
	if err != nil {
		t.Error(err)
		return
	}
	if obj.url != "https://foobar:80" {
		t.Error(obj.url)
		return
	}
}
