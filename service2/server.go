package service2

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	kitlog "github.com/go-kit/kit/log"
	tracekit "github.com/go-kit/kit/tracing/opentracing"
	kithttp "github.com/go-kit/kit/transport/http" /* opentracing*/
	"github.com/opentracing/opentracing-go"
)

const (
	ServiceName = "Service 2"
)

type server struct {
	db      mockDB
	extract kithttp.RequestFunc
}

type mockDB struct{}

func (m mockDB) Query(ctx context.Context, query string) {

	time.Sleep(20 * time.Millisecond)
	span := opentracing.SpanFromContext(ctx)
	defer span.Finish()

	span.LogEvent(fmt.Sprintf("Executing query %s", query))
}

func (s *server) serviceHandler(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	database := values.Get("table")

	/*wireCtx, err := opentracing.GlobalTracer().Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(r.Header),
	)*/

	// extract our new context using some middleware
	ctx := s.extract(
		r.Context(),
		r,
	)

	s.db.Query(ctx, fmt.Sprintf("SELECT * FROM %s WHERE answer=42", database))
	w.WriteHeader(http.StatusNoContent)
}

func (s *server) Start() {
	mux := http.NewServeMux()
	mux.Handle("/database", http.HandlerFunc(s.serviceHandler))
	if err := http.ListenAndServe(":6000", mux); err != nil {
		log.Fatalf("Error: %v\n", err)
	}
}

func NewServer(tracer opentracing.Tracer) *server {
	return &server{
		db:      mockDB{},
		extract: tracekit.HTTPToContext(tracer, ServiceName, kitlog.NewJSONLogger(os.Stdout)),
	}
}
