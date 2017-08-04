package main

import (
	"fmt"
	"os"

	"github.com/go_zipkin_example/service2"
	opentracing "github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
)

const (
	zipkinEndpoint = "http://localhost:9411/api/v1/spans"
	Debug          = false
	serviceName    = "service2"
)

func main() {

	// create collector.
	collector, err := zipkin.NewHTTPCollector(zipkinEndpoint)
	if err != nil {
		fmt.Printf("unable to create Zipkin HTTP collector: %+v\n", err)
		os.Exit(-1)
	}

	// create recorder.
	recorder := zipkin.NewRecorder(collector, Debug, zipkinEndpoint, "service2")

	// create tracer.
	tracer, err := zipkin.NewTracer(
		recorder,
	)

	if err != nil {
		fmt.Printf("unable to create Zipkin tracer: %+v\n", err)
		os.Exit(-1)
	}

	// explicitly set our tracer to be the default tracer.
	opentracing.InitGlobalTracer(tracer)
	server := service2.NewServer(tracer)
	server.Start()

	collector.Close()

}
