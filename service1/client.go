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
	// now we have some middleware, but for span injection instead of span
	// extraction
	inject kithttp.RequestFunc
}

// We'll make our request function a method of the client
func (c *client) PasswordRequest(password string, ctx context.Context) error {
	// We start the span back up from the context passed to the function:
	span, ctx := opentracing.StartSpanFromContext(ctx, "Password")
	log.Println("*****\n", span)
	defer span.Finish()

	// Build our custom url for the password request
	u := fmt.Sprintf(c.URL+"?passwd=%s", "opensesame")
	fmt.Println("URL for request is: ", u)

	req, err := http.NewRequest("POST", u, nil)
	if err != nil {
		return err
	}

	// we inject the spans found in the context into our http request:
	c.inject(ctx, req)

	fmt.Println(req.Context())
	// And we make the request, passing in our context as well.
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
		// Once again, we use our middleware to return a kithttp.RequestFunc
		// which injects tracing data found in a context.Context into an
		// http request's headers.
		inject: tracekit.ContextToHTTP(tracer, kitlog.NewJSONLogger(os.Stdout)),
	}
}
