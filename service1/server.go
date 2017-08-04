package service1

import (
	"log"
	"net/http"
	"os"

	opentracing "github.com/opentracing/opentracing-go"

	kitlog "github.com/go-kit/kit/log"
	tracekit "github.com/go-kit/kit/tracing/opentracing"
	kithttp "github.com/go-kit/kit/transport/http"
)

const ServiceName = "Service1"

type server struct {
	password string
	// some middleware for extracting trace data from http headers
	extract kithttp.RequestFunc
}

func (s *server) serviceHandler(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	password := values.Get("passwd")

	// our extract middleware returns a context.Context with tracing data injected
	// into it. Opentracing uses context.Contexts heavily for intermediate trace
	// storage
	ctx := s.extract(
		r.Context(),
		r,
	)

	// Now, we extract the actual span from the context.Context...
	span := opentracing.SpanFromContext(ctx)

	// In case the request was malformed, we check for error
	if span == nil {
		http.Error(w, "Invalid Span data", http.StatusBadRequest)
		return
	}
	// span.Finish must be called to send off the spans data
	defer span.Finish()

	// Perform our password check
	if password != s.password {
		// Annotate the span if there was an error
		span.SetTag("error", "Invalid Password")
		http.Error(w, "Invalid Password", http.StatusBadRequest)
		return
	}

	// Add an annotation to signify there were no problems
	span.LogEvent("Password verification successful")

	w.WriteHeader(http.StatusNoContent)
}

// Our constructor must now take an opentracing.Tracer, in order to decorate our middleware
func NewServer(password string, tracer opentracing.Tracer) *server {
	return &server{
		password: password,
		// Use go kit middleware for span extraction.Returns a RequestFunc which decorates a
		// context.Context with tracing data from an http request
		extract: tracekit.HTTPToContext(tracer, ServiceName, kitlog.NewJSONLogger(os.Stdout)),
	}
}

// No changes here
func (s *server) Start() {
	mux := http.NewServeMux()
	mux.Handle("/password", http.HandlerFunc(s.serviceHandler))
	if err := http.ListenAndServe(":4048", mux); err != nil {
		log.Fatalf("Error: %v\n", err)
	}
}
