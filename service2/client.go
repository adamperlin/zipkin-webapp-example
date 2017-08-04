package service2

import (
	"context"
	"fmt"
	"net/http"
	"os"

	kitlog "github.com/go-kit/kit/log"
	tracekit "github.com/go-kit/kit/tracing/opentracing"
	kithttp "github.com/go-kit/kit/transport/http"
	opentracing "github.com/opentracing/opentracing-go"
)

type client struct {
	URL            string
	internalClient *http.Client
	tracer         opentracing.Tracer
	inject         kithttp.RequestFunc
}

func (c *client) DatabaseRequest(ctx context.Context, dbname string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Database")
	defer span.Finish()

	u := fmt.Sprintf(c.URL+"?table=%s", dbname)
	fmt.Println("URL for request is: ", u)

	req, err := http.NewRequest("POST", u, nil)
	if err != nil {
		return err
	}

	ctx = opentracing.ContextWithSpan(ctx, span)
	ctx = c.inject(ctx, req)

	resp, err := c.internalClient.Do(req.WithContext(ctx))
	if err != nil {
		span.SetTag("error", err.Error())
		return err
	}

	defer resp.Body.Close()
	return nil
}

func NewClient(url string, tracer opentracing.Tracer) *client {
	return &client{
		URL:            url,
		internalClient: &http.Client{},
		inject:         tracekit.ContextToHTTP(tracer, kitlog.NewJSONLogger(os.Stdout)),
	}
}
