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

package api

import (
	"bytes"
	"encoding/gob"
	"errors"
	"strconv"
)

type CommandMessage struct {
	FunctionId string      `json:"function_id"` //mandatory
	Input      interface{} `json:"input"`

	//device command
	DeviceId  string `json:"device_id,omitempty"`
	ServiceId string `json:"service_id,omitempty"`

	//group command
	GroupId       string `json:"group_id,omitempty"`
	AspectId      string `json:"aspect_id,omitempty"`
	DeviceClassId string `json:"device_class_id,omitempty"`
}

func (this CommandMessage) Validate() error {
	if this.FunctionId == "" {
		return errors.New("expect function_id in body")
	}

	if this.DeviceId != "" && this.ServiceId != "" {
		return nil
	}

	if this.GroupId != "" {
		return nil
	}

	return errors.New("missing device_id, service_id or group_id")
}

func (this CommandMessage) Hash() string {
	var b bytes.Buffer
	gob.NewEncoder(&b).Encode(this)
	return string(b.Bytes())
}

type BatchRequest []CommandMessage

func (this BatchRequest) Validate() error {
	for i, req := range this {
		err := req.Validate()
		if err != nil {
			return errors.New("[" + strconv.Itoa(i) + "]: " + err.Error())
		}
	}
	return nil
}

type BatchResultElement struct {
	StatusCode int         `json:"status_code"`
	Message    interface{} `json:"message"`
}
