/*
 * Copyright 2026 InfAI (CC SES)
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

package command

import (
	"hash/maphash"
	"reflect"
	"testing"
)

func TestCommandMessageGetAspectIds(t *testing.T) {
	t.Run("reads the deprecated aspect id as a single element list", func(t *testing.T) {
		actual := CommandMessage{AspectId: "aid"}.GetAspectIds()
		if !reflect.DeepEqual(actual, []string{"aid"}) {
			t.Error(actual)
		}
	})

	t.Run("reads the aspect list", func(t *testing.T) {
		actual := CommandMessage{AspectIds: []string{"aid1", "aid2"}}.GetAspectIds()
		if !reflect.DeepEqual(actual, []string{"aid1", "aid2"}) {
			t.Error(actual)
		}
	})

	t.Run("adds the deprecated aspect id to the list", func(t *testing.T) {
		actual := CommandMessage{AspectId: "aid2", AspectIds: []string{"aid1"}}.GetAspectIds()
		if !reflect.DeepEqual(actual, []string{"aid1", "aid2"}) {
			t.Error(actual)
		}
	})

	t.Run("does not repeat an aspect the list already names", func(t *testing.T) {
		actual := CommandMessage{AspectId: "aid", AspectIds: []string{"aid"}}.GetAspectIds()
		if !reflect.DeepEqual(actual, []string{"aid"}) {
			t.Error(actual)
		}
	})

	t.Run("answers nothing without an aspect", func(t *testing.T) {
		actual := CommandMessage{FunctionId: "fid"}.GetAspectIds()
		if len(actual) != 0 {
			t.Error(actual)
		}
	})
}

func TestCommandMessageSetAspectIds(t *testing.T) {
	t.Run("folds the deprecated aspect id into the list and clears it", func(t *testing.T) {
		cmd := CommandMessage{FunctionId: "fid", AspectId: "aid"}
		cmd.SetAspectIds()
		if !reflect.DeepEqual(cmd.AspectIds, []string{"aid"}) {
			t.Error(cmd.AspectIds)
		}
		if cmd.AspectId != "" {
			t.Error(cmd.AspectId)
		}
	})

	t.Run("is idempotent", func(t *testing.T) {
		cmd := CommandMessage{FunctionId: "fid", AspectId: "aid"}
		cmd.SetAspectIds()
		cmd.SetAspectIds()
		if !reflect.DeepEqual(cmd.AspectIds, []string{"aid"}) {
			t.Error(cmd.AspectIds)
		}
	})

	t.Run("lets a deprecated aspect id and a single element list hash alike", func(t *testing.T) {
		seed := maphash.MakeSeed()
		withAspectId := CommandMessage{FunctionId: "fid", DeviceId: "did", ServiceId: "sid", AspectId: "aid"}
		withAspectIds := CommandMessage{FunctionId: "fid", DeviceId: "did", ServiceId: "sid", AspectIds: []string{"aid"}}
		withAspectId.SetAspectIds()
		withAspectIds.SetAspectIds()
		if withAspectId.Hash(seed) != withAspectIds.Hash(seed) {
			t.Error(withAspectId.Hash(seed), withAspectIds.Hash(seed))
		}
	})

	t.Run("keeps commands naming different aspects apart", func(t *testing.T) {
		seed := maphash.MakeSeed()
		one := CommandMessage{FunctionId: "fid", DeviceId: "did", ServiceId: "sid", AspectIds: []string{"aid1"}}
		two := CommandMessage{FunctionId: "fid", DeviceId: "did", ServiceId: "sid", AspectIds: []string{"aid1", "aid2"}}
		one.SetAspectIds()
		two.SetAspectIds()
		if one.Hash(seed) == two.Hash(seed) {
			t.Error(one.Hash(seed))
		}
	})
}

func TestBatchRequestSetAspectIds(t *testing.T) {
	t.Run("folds the deprecated aspect id of every element", func(t *testing.T) {
		batch := BatchRequest{
			{FunctionId: "fid", AspectId: "aid1"},
			{FunctionId: "fid", AspectIds: []string{"aid2"}},
			{FunctionId: "fid"},
		}
		batch.SetAspectIds()
		expected := BatchRequest{
			{FunctionId: "fid", AspectIds: []string{"aid1"}},
			{FunctionId: "fid", AspectIds: []string{"aid2"}},
			{FunctionId: "fid"},
		}
		if !reflect.DeepEqual(batch, expected) {
			t.Error(batch)
		}
	})
}
