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

package tracing

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

// record installs a tracer provider that keeps finished spans, and the propagators otelx sets
func record(t *testing.T) *tracetest.SpanRecorder {
	t.Helper()
	recorder := tracetest.NewSpanRecorder()
	previousProvider, previousPropagator := otel.GetTracerProvider(), otel.GetTextMapPropagator()
	otel.SetTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder)))
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	t.Cleanup(func() {
		otel.SetTracerProvider(previousProvider)
		otel.SetTextMapPropagator(previousPropagator)
	})
	return recorder
}

// A request continues the trace of its caller, and a span started while handling it is a child.
func TestHandlerContinuesTheTraceOfTheCaller(t *testing.T) {
	callerTrace := trace.TraceID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	handler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, span := Start(request.Context(), "device command")
		EndWithStatus(span, http.StatusRequestTimeout, "timeout")
	})
	//any endpoint other than Off; otelx sets up its provider here, so the recorder comes after
	traced := Handler("localhost:1", handler)
	recorder := record(t)

	request := httptest.NewRequest(http.MethodPost, "/commands", nil)
	request.Header.Set("traceparent", "00-"+callerTrace.String()+"-0102030405060708-01")
	traced.ServeHTTP(httptest.NewRecorder(), request)

	ended := recorder.Ended()
	if len(ended) != 2 {
		t.Fatalf("spans %+v", ended)
	}
	command, server := ended[0], ended[1]
	if command.Name() != "device command" || server.SpanContext().TraceID() != callerTrace || command.Parent().SpanID() != server.SpanContext().SpanID() {
		t.Fatalf("command %v in %v, server %v", command.Name(), command.Parent(), server.SpanContext())
	}
	if command.Status().Code != codes.Error || command.Status().Description != "timeout" {
		t.Fatalf("status %+v", command.Status())
	}
}

func TestInjectWritesTheTraceContext(t *testing.T) {
	record(t)
	ctx, span := Start(context.Background(), "work")
	defer span.End()
	req, _ := http.NewRequest(http.MethodGet, "http://example.invalid", nil)
	if err := Inject(ctx, req); err != nil {
		t.Fatal(err)
	}
	if req.Header.Get("traceparent") == "" {
		t.Fatalf("headers %v", req.Header)
	}
}

type plain struct{}

func (*plain) ServeHTTP(http.ResponseWriter, *http.Request) {}

func TestOffNeitherWrapsNorInitializes(t *testing.T) {
	handler := &plain{}
	if got := Handler(Off, handler); got != handler {
		t.Fatalf("handler wrapped: %T", got)
	}
	if err := Init(context.Background(), Off); err != nil {
		t.Fatal(err)
	}
}
