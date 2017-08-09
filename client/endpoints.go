package  client

import (
	"context"
	"fmt"
	"time"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
    "github.com/
 )

//implemented according to http://mesos.apache.org/documentation/latest/scheduler-http-api/
type Endpoints struct {
    //subscribe call handler endpoint
	Subscribe    endpoint.Endpoint
    //other calls handler endpoint
	Call endpoint.Endpoint
}

// Subscribe implements Service. 
func (e Endpoints) Subscribe(ctx context.Context, user string,name string  ) (int, error) {
	request := sumRequest{A: a, B: b}
	response, err := e.SumEndpoint(ctx, request)
	if err != nil {
		return 0, err
	}
	return response.(sumResponse).V, response.(sumResponse).Err
}

// Concat implements Service. Primarily useful in a client.
func (e Endpoints) Concat(ctx context.Context, a, mesosStreamId  string) (string, error) {
	request := concatRequest{A: a, B: b}
	response, err := e.ConcatEndpoint(ctx, request)
	if err != nil {
		return "", err
	}
	return response.(concatResponse).V, response.(concatResponse).Err
}

// MakeSumEndpoint returns an endpoint that invokes Sum on the service.
// Primarily useful in a server.
func MakeSumEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		sumReq := request.(sumRequest)
		v, err := s.Sum(ctx, sumReq.A, sumReq.B)
		if err == ErrIntOverflow {
			return nil, err // special case; see comment on ErrIntOverflow
		}
		return sumResponse{
			V:   v,
			Err: err,
		}, nil
	}
}

// MakeConcatEndpoint returns an endpoint that invokes Concat on the service.
// Primarily useful in a server.
func MakeConcatEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		concatReq := request.(concatRequest)
		v, err := s.Concat(ctx, concatReq.A, concatReq.B)
		return concatResponse{
			V:   v,
			Err: err,
		}, nil
	}
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
