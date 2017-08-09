package scheduler

// This file contains methods to make individual endpoints from services,
// request and response types to serve those endpoints, as well as encoders and
// decoders for those types, for all of our supported transport serialization
// formats. It also includes endpoint middlewares.

import (
	"context"
	"fmt"
	"time"

	"github.com/guherbozdogan/kit/endpoint"
	"github.com/guherbozdogan/kit/log"
	"github.com/guherbozdogan/kit/metrics"
	//proto "github.com/gogo/protobuf/proto"
	//math "math"
	mesos_v1 "github.com/guherbozdogan/mesos-go-http-client/client/pb/mesos/v1"
	//mesos_v1_scheduler "github.com/guherbozdogan/mesos-go-http-client/client/pb/mesos/v1/scheduler"
	
)

// Endpoints collects all of the endpoints that compose an add service. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
//
// In a server, it's useful for functions that need to operate on a per-endpoint
// basis. For example, you might pass an Endpoints to a function that produces
// an http.Handler, with each method (endpoint) wired up to a specific path. (It
// is probably a mistake in design to invoke the Service methods on the
// Endpoints struct in a server.)
//
// In a client, it's useful to collect individually constructed endpoints into a
// single type that implements the Service interface. For example, you might
// construct individual endpoints using transport/http.NewClient, combine them
// into an Endpoints, and return it to the caller as a Service.
type Endpoints struct {
	SubscribeEndpoint   endpoint.Endpoint
	CallEndPoint  endpoint.Endpoint
}

// Sum implements Service. Primarily useful in a client.
func (e Endpoints) Subscribe(ctx context.Context, a, b int) (int, error) {
	request := sumRequest{A: a, B: b}
	response, err := e.SumEndpoint(ctx, request)
	if err != nil {
		return 0, err
	}
	return response.(sumResponse).V, response.(sumResponse).Err
}

// Concat implements Service. Primarily useful in a client.
func (e Endpoints) Callt(ctx context.Context, a, b string) (string, error) {
	request := concatRequest{A: a, B: b}
	response, err := e.ConcatEndpoint(ctx, request)
	if err != nil {
		return "", err
	}
	return response.(concatResponse).V, response.(concatResponse).Err
}


// EndpointInstrumentingMiddleware returns an endpoint middleware that records
// the duration of each invocation to the passed histogram. The middleware adds
// a single field: "success", which is "true" if no error is returned, and
// "false" otherwise.
func EndpointInstrumentingMiddleware(duration metrics.Histogram) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {

			defer func(begin time.Time) {
				duration.With("success", fmt.Sprint(err == nil)).Observe(time.Since(begin).Seconds())
			}(time.Now())
			return next(ctx, request)

		}
	}
}

// EndpointLoggingMiddleware returns an endpoint middleware that logs the
// duration of each invocation, and the resulting error, if any.
func EndpointLoggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {

			defer func(begin time.Time) {
				logger.Log("error", err, "took", time.Since(begin))
			}(time.Now())
			return next(ctx, request)

		}
	}
}

// These types are unexported because they only exist to serve the endpoint
// domain, which is totally encapsulated in this package. They are otherwise
// opaque to all callers.

type subscriberRequest struct{
	
	
}

type subscribeRespone struct {

	
}
type sumRequest struct{ A, B int }

type sumResponse struct {
	V   int
	Err error
}

		
type concatRequest struct{ A, B string }

type concatResponse struct {
	V   string
	Err error
}
