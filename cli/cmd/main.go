package main

import (
	"context"
	"log"

	"github.com/go_zipkin_example/service1"
	opentracing "github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
)

const (
	ZipkinEndpoint = "http://localhost:9411/api/v1/spans"
	Debug          = false
	HostName       = "0.0.0.0:0"
	ServiceName    = "cli"
)

func main() {
	collector, err := zipkin.NewHTTPCollector(ZipkinEndpoint)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	recorder := zipkin.NewRecorder(collector, Debug, HostName, ServiceName)

	tracer, err := zipkin.NewTracer(recorder)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	// Initialize global tracer with zipkin implementation so we
	// stick to using opentracing from now on.
	opentracing.InitGlobalTracer(tracer)

	// Our root span, the beginning of a trace
	span := opentracing.StartSpan("Main")

	// Inject span into a context.Context so that it can be sent
	ctx := opentracing.ContextWithSpan(context.Background(), span)

	server := service1.NewServer("opensesame")
	server.Start()

}
