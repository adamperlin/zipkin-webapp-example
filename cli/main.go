package main

import (
	"context"
	"log"

	"github.com/go_zipkin_example/service1"
	"github.com/go_zipkin_example/service2"
	opentracing "github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
)

const (
	ZipkinEndpoint   = "http://localhost:9411/api/v1/spans"
	Debug            = false
	Service1Endpoint = "http://localhost:4048/password"
	Service2Endpoint = "http://localhost:6000/database"
	ServiceName      = "cli"
)

func main() {
	collector, err := zipkin.NewHTTPCollector(ZipkinEndpoint)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	recorder := zipkin.NewRecorder(collector, Debug, ZipkinEndpoint, ServiceName)

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
	log.Println(ctx)

	//Initialize our client with a tracer for middleware to use
	client := service1.NewClient(Service1Endpoint, tracer)
	// Add an annotation to the span...
	span.LogEvent("Call Service 1")

	// Use the password "opensesame" for the request
	err = client.PasswordRequest("opensesame", ctx)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	client2 := service2.NewClient(Service2Endpoint, tracer)
	// Annotate the span so we can trace behavior
	span.LogEvent("Call Service 2")

	// Make a request in the table "lifetheuniverseandeverything"
	err = client2.DatabaseRequest(ctx, "lifetheuniverseandeverything")
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	// End the span
	span.Finish()
	// Collector must be closed to flush out all spans
	collector.Close()
}
