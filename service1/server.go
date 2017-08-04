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

// Do you think it's secure enough?
const Password = "opensesame"
const ServiceName = "Service1"

type server struct {
	password string
	extract  kithttp.RequestFunc
}

func (s *server) serviceHandler(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	password := values.Get("passwd")

	log.Println("REQUEST CONTEXT IS: ", r.Context())

	ctx := s.extract(
		r.Context(),
		r,
	)
	span := opentracing.SpanFromContext(ctx)

	if span == nil {
		http.Error(w, "Invalid Span data", http.StatusBadRequest)
		return
	}

	defer span.Finish()

	if password != s.password {
		span.SetTag("error", "Invalid Password")
		http.Error(w, "Invalid Password", http.StatusBadRequest)
		return
	}

	span.LogEvent("Password verification successful")

	w.WriteHeader(http.StatusNoContent)
}

func NewServer(password string, tracer opentracing.Tracer) *server {
	return &server{
		password: password,
		extract:  tracekit.HTTPToContext(tracer, ServiceName, kitlog.NewJSONLogger(os.Stdout)),
	}
}

func (s *server) Start() {
	mux := http.NewServeMux()
	mux.Handle("/password", http.HandlerFunc(s.serviceHandler))
	if err := http.ListenAndServe(":4048", mux); err != nil {
		log.Fatalf("Error: %v\n", err)
	}
}
