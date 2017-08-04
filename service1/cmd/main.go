package main

import (
	"fmt"
	"os"

	"github.com/go_zipkin_example/service1"
	opentracing "github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
)

const (
	zipkinEndpoint = "http://localhost:9411/api/v1/spans"
	Debug          = false
	serviceName    = "service1"
)

func main() {

	// Create the zipkin http collector. This is what sends data
	// to the zipkin server
	collector, err := zipkin.NewHTTPCollector(zipkinEndpoint)
	if err != nil {
		fmt.Printf("unable to create Zipkin HTTP collector: %+v\n", err)
		os.Exit(-1)
	}

	// The recorder is an abstraction layer over the collector
	recorder := zipkin.NewRecorder(collector, Debug, zipkinEndpoint, "service1")

	// The zipkin tracer implements the opentracing.Tracer interface, so that
	// it can be used with our server
	tracer, err := zipkin.NewTracer(recorder)

	if err != nil {
		fmt.Printf("unable to create Zipkin tracer: %+v\n", err)
		os.Exit(-1)
	}

	// Set our global tracer to the tracer implementation we're using
	// Not necessary in our case, but good practice.
	// This allows any method in the program to call opentracting.GlobalTracer(),
	// and have access to our tracer. Good for expansion.
	opentracing.InitGlobalTracer(tracer)
	// Create our server with a generic password and our opentracing.Tracer
	// implementation
	server := service1.NewServer("opensesame", tracer)
	server.Start()

	// Make sure to close the zipkin collector so that it flushes out all
	// buffered spans
	collector.Close()
}
