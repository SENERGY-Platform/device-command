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

// Package lastvalue reads the last message a device sent for a service. It answers commands for
// measuring functions of event services, which cannot be asked.
package lastvalue

import (
	"context"
	"errors"
	"time"

	"github.com/SENERGY-Platform/device-command/pkg/auth"
	"github.com/SENERGY-Platform/device-command/pkg/configuration"
	model "github.com/SENERGY-Platform/models/go/models"
)

// ErrMissing is returned by the gateway variant if the device has not sent since mgw-last-value
// started. A client may then ask the platform instead.
var ErrMissing = errors.New("missing last value in mgw-last-value")

// Message is the last message of a service. Value maps the content variable name to its value;
// Time is when the message was stored, RFC3339 in UTC, empty if the source does not say.
type Message struct {
	Value map[string]interface{}
	Time  string
}

// Reader returns the last message of a service.
type Reader interface {
	LastMessage(ctx context.Context, token auth.Token, device model.Device, service model.Service, protocol model.Protocol, timeout time.Duration) (Message, error)
}

// normalizeTime returns t as RFC3339 in UTC, or t unchanged if it is no RFC3339 time.
func normalizeTime(t string) string {
	parsed, err := time.Parse(time.RFC3339Nano, t)
	if err != nil {
		return t
	}
	return parsed.UTC().Format(time.RFC3339Nano)
}

// New returns the timescale-wrapper of the platform, or mgw-last-value with timescale_impl "mgw".
func New(config configuration.Config) (Reader, error) {
	if config.TimescaleImpl == "mgw" {
		return NewMgw(config.TimescaleWrapperUrl, config.MgwProtocolSegment)
	}
	return &Cloud{url: config.TimescaleWrapperUrl}, nil
}
