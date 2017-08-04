package service1

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	kitlog "github.com/go-kit/kit/log"
	tracekit "github.com/go-kit/kit/tracing/opentracing"
	kithttp "github.com/go-kit/kit/transport/http"
	opentracing "github.com/opentracing/opentracing-go"
)

const ServiceURL = "http://localhost:4048/password?passwd=%s"

type client struct {
	URL            string
	internalClient *http.Client
	tracer         opentracing.Tracer
	inject         kithttp.RequestFunc
	//	traceMiddleware middleware.RequestFunc
}

func (c *client) PasswordRequest(password string, ctx context.Context) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Password")
	log.Println("*****\n", span)
	defer span.Finish()

	u := fmt.Sprintf(c.URL+"?passwd=%s", "opensesame")
	fmt.Println("URL for request is: ", u)

	req, err := http.NewRequest("POST", u, nil)
	if err != nil {
		return err
	}

	//middleware.FromHTTPRequest(tracer, operationName)

	//req = c.traceMiddleware(req.WithContext(ctx))

	/*	opentracing.GlobalTracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)*/

	ctx = opentracing.ContextWithSpan(ctx, span)

	c.inject(ctx, req)

	//middleware.ToHTTPRequest(tracer)(req.WithContext(ctx))
	fmt.Println(req.Context())
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
		tracer:         tracer,
		internalClient: &http.Client{},
		inject:         tracekit.ContextToHTTP(tracer, kitlog.NewJSONLogger(os.Stdout)),
		//		traceMiddleware: middleware.ToHTTPRequest(tracer),
	}
}
