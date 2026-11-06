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

// Package transport sends protocol messages to the connectors and reports their replies: over
// kafka in the cloud, over the local mqtt broker on a gateway.
package transport

import (
	"context"

	"github.com/SENERGY-Platform/device-command/pkg/configuration"
	"github.com/SENERGY-Platform/device-command/pkg/messages"
)

// Reply is the answer of a connector to a sent command, identified by the task id of the command.
type Reply struct {
	TaskId string
	Output map[string]string
	// Failed is set if the connector reported an error; Output then carries its message.
	Failed bool
}

type Sender interface {
	Send(ctx context.Context, msg messages.ProtocolMsg) error
}

// New connects to kafka, or to the mqtt broker of a gateway with com_impl "mgw". Every reply is
// passed to receive, which must not block.
func New(ctx context.Context, config configuration.Config, receive func(Reply)) (Sender, error) {
	if config.ComImpl == "mgw" {
		return NewMqtt(ctx, config, receive)
	}
	return NewKafka(ctx, config, receive)
}
