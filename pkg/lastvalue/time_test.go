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

package lastvalue

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/SENERGY-Platform/device-command/pkg/auth"
	model "github.com/SENERGY-Platform/models/go/models"
)

func serve(t *testing.T, body string) string {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(body))
	}))
	t.Cleanup(server.Close)
	return server.URL
}

func TestCloudLastMessageTime(t *testing.T) {
	cloud := &Cloud{url: serve(t, `{"time":"2026-11-06T16:47:59.183+01:00","value":{"value":42}}`)}
	msg, err := cloud.LastMessage(context.Background(), auth.Token{}, model.Device{Id: "d"}, model.Service{Id: "s"}, model.Protocol{}, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if msg.Time != "2026-11-06T15:47:59.183Z" {
		t.Error(msg.Time)
	}
	if msg.Value["value"] != float64(42) {
		t.Error(msg.Value)
	}
}

func TestCloudLastMessageWithoutTime(t *testing.T) {
	cloud := &Cloud{url: serve(t, `{"value":{"value":42}}`)}
	msg, err := cloud.LastMessage(context.Background(), auth.Token{}, model.Device{Id: "d"}, model.Service{Id: "s"}, model.Protocol{}, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if msg.Time != "" {
		t.Error(msg.Time)
	}
}

func TestMgwLastMessageTime(t *testing.T) {
	mgw, err := NewMgw(serve(t, `[{"time":"2026-11-06T15:43:29+01:00","value":"{\"v\":1}"}]`), "data")
	if err != nil {
		t.Fatal(err)
	}
	protocol := model.Protocol{ProtocolSegments: []model.ProtocolSegment{{Id: "seg", Name: "data"}}}
	service := model.Service{LocalId: "s", Outputs: []model.Content{{ProtocolSegmentId: "seg", ContentVariable: model.ContentVariable{Name: "root"}}}}
	msg, err := mgw.LastMessage(context.Background(), auth.Token{}, model.Device{LocalId: "d"}, service, protocol, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if msg.Time != "2026-11-06T14:43:29Z" {
		t.Error(msg.Time)
	}
	if msg.Value["root"] != `{"v":1}` {
		t.Error(msg.Value)
	}
}
