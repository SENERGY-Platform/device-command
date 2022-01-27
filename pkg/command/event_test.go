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

package command

import (
	"github.com/SENERGY-Platform/external-task-worker/lib/devicerepository/model"
	"reflect"
	"strings"
	"testing"
)

func Test_createTimescaleRequest(t *testing.T) {
	actualResult := createTimescaleRequest(model.Device{
		Id:      "deviceid",
		LocalId: "localdeviceid",
	}, model.Service{
		Id:      "serviceid",
		LocalId: "localserviceid",
		Outputs: []model.Content{
			{Id: "contentid_meta", ContentVariable: model.ContentVariable{}, ProtocolSegmentId: "meta_segment_id"},
			{Id: "contentid_data", ProtocolSegmentId: "data_segment_id", ContentVariable: model.ContentVariable{
				Name: "root",
				SubContentVariables: []model.ContentVariable{
					{Name: "foo"},
					{Name: "bar"},
					{Name: "batz", SubContentVariables: []model.ContentVariable{
						{Name: "blub"},
					}},
				},
			}},
		},
	})
	expectedResult := []TimescaleRequest{
		{
			DeviceId:   "deviceid",
			ServiceId:  "serviceid",
			ColumnName: "root.foo",
		},
		{
			DeviceId:   "deviceid",
			ServiceId:  "serviceid",
			ColumnName: "root.bar",
		},
		{
			DeviceId:   "deviceid",
			ServiceId:  "serviceid",
			ColumnName: "root.batz.blub",
		},
	}
	if !reflect.DeepEqual(actualResult, expectedResult) {
		t.Error(actualResult)
	}
}

func Test_setPath(t *testing.T) {
	paths := map[string]interface{}{
		"foo":              42,
		"root.bar":         "batz",
		"root.blub.bla":    13,
		"root.blub.blabla": "leaf",
	}
	result := map[string]interface{}{}
	for path, value := range paths {
		result = setPath(result, strings.Split(path, "."), value)
	}

	var expected interface{} = map[string]interface{}{
		"foo": 42,
		"root": map[string]interface{}{
			"bar": "batz",
			"blub": map[string]interface{}{
				"bla":    13,
				"blabla": "leaf",
			},
		},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Error(result)
	}
}

func Test_getTimescaleValue(t *testing.T) {

	result := getTimescaleValue([]TimescaleRequest{
		{
			DeviceId:   "d1",
			ServiceId:  "s1",
			ColumnName: "foo",
		},
		{
			DeviceId:   "d1",
			ServiceId:  "s1",
			ColumnName: "root.bar",
		},
		{
			DeviceId:   "d1",
			ServiceId:  "s1",
			ColumnName: "root.blub.bla",
		},
		{
			DeviceId:   "d1",
			ServiceId:  "s1",
			ColumnName: "root.blub.blabla",
		},
	}, []TimescaleResponse{
		{
			Value: 42,
		},
		{
			Value: "batz",
		},
		{
			Value: 13,
		},
		{
			Value: "leaf",
		},
	})

	var expected interface{} = map[string]interface{}{
		"foo": 42,
		"root": map[string]interface{}{
			"bar": "batz",
			"blub": map[string]interface{}{
				"bla":    13,
				"blabla": "leaf",
			},
		},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Error(result)
	}
}

func Test_createEventValueFromTimescaleValues(t *testing.T) {
	result, err := createEventValueFromTimescaleValues(model.Service{
		Outputs: []model.Content{
			{
				ContentVariable: model.ContentVariable{
					Name: "root",
					SubContentVariables: []model.ContentVariable{
						{
							Name: "bar",
							Type: model.String,
						},
						{
							Name: "blub",
							SubContentVariables: []model.ContentVariable{
								{Name: "bla"},
								{Name: "blabla"},
							},
						},
					},
				},
				Serialization:     "json",
				ProtocolSegmentId: "bodyId",
			},
		},
	}, model.Protocol{
		ProtocolSegments: []model.ProtocolSegment{
			{
				Id:   "metaid",
				Name: "meta",
			},
			{
				Id:   "bodyId",
				Name: "body",
			},
		},
	}, []TimescaleRequest{
		{
			DeviceId:   "d1",
			ServiceId:  "s1",
			ColumnName: "foo",
		},
		{
			DeviceId:   "d1",
			ServiceId:  "s1",
			ColumnName: "root.bar",
		},
		{
			DeviceId:   "d1",
			ServiceId:  "s1",
			ColumnName: "root.blub.bla",
		},
		{
			DeviceId:   "d1",
			ServiceId:  "s1",
			ColumnName: "root.blub.blabla",
		},
	}, []TimescaleResponse{
		{
			Value: 42,
		},
		{
			Value: "batz",
		},
		{
			Value: 13,
		},
		{
			Value: "leaf",
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
	expected := map[string]string{
		"body": `{"bar":"batz","blub":{"bla":13,"blabla":"leaf"}}`,
	}
	if !reflect.DeepEqual(result, expected) {
		t.Error(result)
		return
	}
}
