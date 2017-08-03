package service1

import (
	"context"
	"net/http"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go-opentracing/examples/middleware"
)

const ServiceURL = "http://localhost:4048/password?passwd=%s"

type client struct {
	URL             string
	internalClient  *http.Client
	tracer          opentracing.Tracer
	traceMiddleware middleware.RequestFunc
}

func (c *client) PasswordRequest(url string, password string, ctx context.Context) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Password")
	defer span.Finish()

	req, err := http.NewRequest("POST", c.URL, nil)

	if err != nil {
		return err
	}

	req = c.traceMiddleware(req.WithContext(ctx)) //middleware.ToHTTPRequest(tracer)(req.WithContext(ctx))

	resp, err := c.internalClient.Do(req)
	if err != nil {
		span.SetTag("error", err.Error())
		return err
	}

	defer resp.Body.Close()
	return nil
}

func NewClient(url string, tracer opentracing.Tracer) *client {
	return &client{
		URL:             url,
		tracer:          tracer,
		internalClient:  &http.Client{},
		traceMiddleware: middleware.ToHTTPRequest(tracer),
	}
}
