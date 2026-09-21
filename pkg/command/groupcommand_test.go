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
	"testing"

	"github.com/SENERGY-Platform/external-task-worker/lib/devicerepository/model"
)

const testMeasuringFunctionId = model.MEASURING_FUNCTION_PREFIX + "temperature"
const testControllingFunctionId = model.CONTROLLING_FUNCTION_PREFIX + "set-temperature"

// testAspectNodes describes air with the children inside_air and outside_air, and the
// unrelated aspect device. The nodes are what the device-repository answers, so a queried
// node knows its subtree.
var testAspectNodes = map[string]model.AspectNode{
	"air": {
		Id:            "air",
		RootId:        "air",
		ChildIds:      []string{"inside_air", "outside_air"},
		DescendentIds: []string{"inside_air", "outside_air"},
	},
	"inside_air":  {Id: "inside_air", RootId: "air", ParentId: "air", AncestorIds: []string{"air"}},
	"outside_air": {Id: "outside_air", RootId: "air", ParentId: "air", AncestorIds: []string{"air"}},
	"device":      {Id: "device", RootId: "device"},
}

func testNodes(ids ...string) (result []model.AspectNode) {
	for _, id := range ids {
		result = append(result, testAspectNodes[id])
	}
	return result
}

// testService is a service with one output content variable that serves the measuring
// function for the given aspects, and one input content variable that serves the
// controlling function for the same aspects.
func testService(id string, aspectIds ...string) model.Service {
	variable := func(name string, functionId string) model.ContentVariable {
		return model.ContentVariable{
			Name:                name,
			SubContentVariables: []model.ContentVariable{{Name: "value", FunctionId: functionId, AspectIds: aspectIds}},
		}
	}
	return model.Service{
		Id:      id,
		Inputs:  []model.Content{{ContentVariable: variable("input", testControllingFunctionId)}},
		Outputs: []model.Content{{ContentVariable: variable("output", testMeasuringFunctionId)}},
	}
}

func serviceIds(services []model.Service) (result []string) {
	for _, service := range services {
		result = append(result, service.Id)
	}
	return result
}

func equalStrings(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestGetFilteredServices(t *testing.T) {
	cmd := &Command{}
	services := []model.Service{
		testService("s_outside_air", "outside_air"),
		testService("s_inside_air", "inside_air"),
		testService("s_device", "device"),
		testService("s_outside_air_and_device", "outside_air", "device"),
	}

	t.Run("selects every service serving the function if no aspect is asked for", func(t *testing.T) {
		actual := serviceIds(cmd.getFilteredServices(testMeasuringFunctionId, nil, services))
		expected := []string{"s_device", "s_inside_air", "s_outside_air", "s_outside_air_and_device"}
		if !equalStrings(actual, expected) {
			t.Error(actual)
		}
	})

	t.Run("selects the services carrying the asked for aspect", func(t *testing.T) {
		actual := serviceIds(cmd.getFilteredServices(testMeasuringFunctionId, testNodes("outside_air"), services))
		expected := []string{"s_outside_air", "s_outside_air_and_device"}
		if !equalStrings(actual, expected) {
			t.Error(actual)
		}
	})

	t.Run("selects the services carrying a descendant of the asked for aspect", func(t *testing.T) {
		actual := serviceIds(cmd.getFilteredServices(testMeasuringFunctionId, testNodes("air"), services))
		expected := []string{"s_inside_air", "s_outside_air", "s_outside_air_and_device"}
		if !equalStrings(actual, expected) {
			t.Error(actual)
		}
	})

	t.Run("requires one content variable to carry every asked for aspect", func(t *testing.T) {
		actual := serviceIds(cmd.getFilteredServices(testMeasuringFunctionId, testNodes("outside_air", "device"), services))
		expected := []string{"s_outside_air_and_device"}
		if !equalStrings(actual, expected) {
			t.Error(actual)
		}
	})

	t.Run("selects nothing if one of the asked for aspects is unmatched", func(t *testing.T) {
		actual := serviceIds(cmd.getFilteredServices(testMeasuringFunctionId, testNodes("inside_air", "device"), services))
		if len(actual) != 0 {
			t.Error(actual)
		}
	})

	t.Run("selects nothing for an aspect no service carries", func(t *testing.T) {
		actual := serviceIds(cmd.getFilteredServices(testMeasuringFunctionId, []model.AspectNode{{Id: "unknown"}}, services))
		if len(actual) != 0 {
			t.Error(actual)
		}
	})

	t.Run("reads the deprecated aspect of a content variable", func(t *testing.T) {
		deprecated := model.Service{
			Id: "s_deprecated",
			Outputs: []model.Content{{ContentVariable: model.ContentVariable{
				Name:                "output",
				SubContentVariables: []model.ContentVariable{{Name: "value", FunctionId: testMeasuringFunctionId, AspectId: "outside_air"}},
			}}},
		}
		actual := serviceIds(cmd.getFilteredServices(testMeasuringFunctionId, testNodes("air"), []model.Service{deprecated}))
		if !equalStrings(actual, []string{"s_deprecated"}) {
			t.Error(actual)
		}
	})

	t.Run("reads the inputs for a controlling function", func(t *testing.T) {
		actual := serviceIds(cmd.getFilteredServices(testControllingFunctionId, testNodes("outside_air"), services))
		expected := []string{"s_outside_air", "s_outside_air_and_device"}
		if !equalStrings(actual, expected) {
			t.Error(actual)
		}
	})
}
