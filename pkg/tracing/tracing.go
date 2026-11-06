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

// Package tracing carries OpenTelemetry trace context and baggage through device-command, the way
// the other platform services do it with gin-middleware/otelx: a request continues the trace of its
// caller, and every request to a platform service continues it further.
package tracing

import (
	"context"
	"fmt"
	"net/http"

	"github.com/SENERGY-Platform/gin-middleware/otelx"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// ServiceName is the OpenTelemetry service name of device-command.
const ServiceName = "device-command"

// Off as otel_endpoint disables tracing. A gateway uses it, as it has no collector.
const Off = "off"

// Init sets up the exporter once per process. An empty endpoint selects the collector of the
// platform, the default of otelx; with Off spans are dropped.
func Init(ctx context.Context, endpoint string) error {
	if endpoint == Off {
		return nil
	}
	//otelx initializes only through its middlewares; the returned handler is not used
	_, err := otelx.HTTPOpenTelemetry(ctx, endpoint, ServiceName, http.NotFoundHandler())
	return err
}

// Handler traces incoming requests and continues the trace of the caller; otelx puts user_id and
// username of the token into the baggage. With Off the handler stays unchanged, as otelx would set
// up its exporter on the first call.
func Handler(endpoint string, handler http.Handler) http.Handler {
	if endpoint == Off {
		return handler
	}
	traced, err := otelx.HTTPOpenTelemetry(context.Background(), endpoint, ServiceName, handler)
	if err != nil {
		return handler
	}
	return traced
}

// Inject writes trace context and baggage of ctx into a request to another platform service. A
// request leaving the platform must not carry them, the baggage holds user ids.
func Inject(ctx context.Context, req *http.Request) error {
	return otelx.InjectContextToRequest(ctx, req)
}

// Start begins a span as part of the trace in ctx.
func Start(ctx context.Context, name string, attributes ...attribute.KeyValue) (context.Context, trace.Span) {
	return otel.Tracer(ServiceName).Start(ctx, name, trace.WithAttributes(attributes...))
}

// EndWithStatus ends a span of work that is answered by an http status code. A code from 400 on
// marks the span as failed, with the answer as description.
func EndWithStatus(span trace.Span, code int, answer interface{}) {
	span.SetAttributes(attribute.Int("http.response.status_code", code))
	if code >= http.StatusBadRequest {
		span.SetStatus(codes.Error, fmt.Sprint(answer))
	}
	span.End()
}
