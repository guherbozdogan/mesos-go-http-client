package client

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
	mesos_v1_scheduler "github.com/guherbozdogan/mesos-go-http-client/client/pb/mesos/v1/scheduler"
	//mesos_v1_scheduler "github.com/guherbozdogan/mesos-go-http-client/client/pb/mesos/v1/scheduler"
	"github.com/guherbozdogan/mesos-go-http-client/client/frame"
)

// Endpoints collects all of the endpoints that compose an add service. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
//
//
//
// In a client, it's useful to collect individually constructed endpoints into a
// single type that implements the Service interface. For example, you might
// construct individual endpoints using transport/http.NewClient, combine them
// into an Endpoints, and return it to the caller as a Service.
type Endpoints struct {
	FrameReadFunc       frame.FrameReadFunc
	FrameErrorFunc      frame.ErrorFunc
	SubscribeEndpoint   endpoint.Endpoint
	AcceptEndpoint      endpoint.Endpoint
	DeclineEndpoint     endpoint.Endpoint
	ReviveEndpoint      endpoint.Endpoint
	KillEndpoint        endpoint.Endpoint
	ShutdownEndpoint    endpoint.Endpoint
	AcknowledgeEndpoint endpoint.Endpoint
	ReconcileEndpoint   endpoint.Endpoint
	MessageEndpoint     endpoint.Endpoint
	RequestEndpoint     endpoint.Endpoint
	SuppressEndpoint    endpoint.Endpoint
}
type CallResponse interface {
}

func (e Endpoints) Subscribe(ctx context.Context, FrameworkId *mesos_v1.FrameworkID,
	Mes *mesos_v1_scheduler.Call_Subscribe) (CallResponse, error) {

	t := mesos_v1_scheduler.Call_Type(mesos_v1_scheduler.Call_SUBSCRIBE)
	request := &mesos_v1_scheduler.Call{
		FrameworkId, &t,
		Mes, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil,
		nil, nil}
	response, err := e.SubscribeEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.(CallResponse), err
}

//where is teardown used (from call_type definition in pb file)?

func (e Endpoints) Accept(ctx context.Context, FrameworkId *mesos_v1.FrameworkID,
	Mes *mesos_v1_scheduler.Call_Accept) (CallResponse, error) {

	t := mesos_v1_scheduler.Call_Type(mesos_v1_scheduler.Call_ACCEPT)
	request := &mesos_v1_scheduler.Call{
		FrameworkId, &t,
		nil, Mes, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil,
		nil, nil}
	response, err := e.AcceptEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.(CallResponse), err
}

func (e Endpoints) Decline(ctx context.Context, FrameworkId *mesos_v1.FrameworkID,
	Mes *mesos_v1_scheduler.Call_Decline) (CallResponse, error) {

	t := mesos_v1_scheduler.Call_Type(
		mesos_v1_scheduler.Call_DECLINE)
	request := &mesos_v1_scheduler.Call{
		FrameworkId, &t,
		nil, nil, Mes, nil, nil,
		nil, nil, nil, nil, nil, nil, nil,
		nil, nil}
	response, err := e.DeclineEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.(CallResponse), err
}

func (e Endpoints) Revive(ctx context.Context, FrameworkId *mesos_v1.FrameworkID,
	Mes *mesos_v1_scheduler.Call_Revive) (CallResponse, error) {

	t := mesos_v1_scheduler.Call_Type(
		mesos_v1_scheduler.Call_REVIVE)
	request := &mesos_v1_scheduler.Call{
		FrameworkId, &t,
		nil, nil, nil, nil, nil,
		Mes, nil, nil, nil, nil, nil, nil,
		nil, nil}
	response, err := e.ReviveEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.(CallResponse), err
}

func (e Endpoints) Kill(ctx context.Context, FrameworkId *mesos_v1.FrameworkID,
	Mes *mesos_v1_scheduler.Call_Kill) (CallResponse, error) {

	t := mesos_v1_scheduler.Call_Type(
		mesos_v1_scheduler.Call_KILL)
	request := &mesos_v1_scheduler.Call{
		FrameworkId, &t,
		nil, nil, nil, nil, nil,
		nil, Mes, nil, nil, nil, nil, nil,
		nil, nil}
	response, err := e.KillEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.(CallResponse), err
}

func (e Endpoints) ShutDown(ctx context.Context, FrameworkId *mesos_v1.FrameworkID,
	Mes *mesos_v1_scheduler.Call_Shutdown) (CallResponse, error) {

	t := mesos_v1_scheduler.Call_Type(
		mesos_v1_scheduler.Call_SHUTDOWN)
	request := &mesos_v1_scheduler.Call{
		FrameworkId, &t,
		nil, nil, nil, nil, nil,
		nil, nil, Mes, nil, nil, nil, nil,
		nil, nil}
	response, err := e.ShutdownEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.(CallResponse), err
}

func (e Endpoints) Acknowledge(ctx context.Context, FrameworkId *mesos_v1.FrameworkID,
	Mes *mesos_v1_scheduler.Call_Acknowledge) (CallResponse, error) {

	t := mesos_v1_scheduler.Call_Type(
		mesos_v1_scheduler.Call_ACKNOWLEDGE)
	request := &mesos_v1_scheduler.Call{
		FrameworkId, &t,
		nil, nil, nil, nil, nil,
		nil, nil, nil, Mes, nil, nil, nil,
		nil, nil}
	response, err := e.AcknowledgeEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.(CallResponse), err
}

func (e Endpoints) Reconcile(ctx context.Context, FrameworkId *mesos_v1.FrameworkID,
	Mes *mesos_v1_scheduler.Call_Reconcile) (CallResponse, error) {

	t := mesos_v1_scheduler.Call_Type(
		mesos_v1_scheduler.Call_RECONCILE)
	request := &mesos_v1_scheduler.Call{
		FrameworkId, &t,
		nil, nil, nil, nil, nil,
		nil, nil, nil, nil, Mes, nil, nil,
		nil, nil}
	response, err := e.ReconcileEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.(CallResponse), err
}

func (e Endpoints) Message(ctx context.Context, FrameworkId *mesos_v1.FrameworkID,
	Mes *mesos_v1_scheduler.Call_Message) (CallResponse, error) {

	t := mesos_v1_scheduler.Call_Type(
		mesos_v1_scheduler.Call_MESSAGE)
	request := &mesos_v1_scheduler.Call{
		FrameworkId, &t,
		nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, Mes, nil,
		nil, nil}
	response, err := e.MessageEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.(CallResponse), err
}

func (e Endpoints) Request(ctx context.Context, FrameworkId *mesos_v1.FrameworkID,
	Mes *mesos_v1_scheduler.Call_Request) (CallResponse, error) {

	t := mesos_v1_scheduler.Call_Type(
		mesos_v1_scheduler.Call_REQUEST)
	request := &mesos_v1_scheduler.Call{
		FrameworkId, &t,
		nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, Mes,
		nil, nil}
	response, err := e.RequestEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.(CallResponse), err
}

func (e Endpoints) Suppress(ctx context.Context, FrameworkId *mesos_v1.FrameworkID,
	Mes *mesos_v1_scheduler.Call_Suppress) (CallResponse, error) {

	t := mesos_v1_scheduler.Call_Type(
		mesos_v1_scheduler.Call_SUPPRESS)
	request := &mesos_v1_scheduler.Call{
		FrameworkId, &t,
		nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil,
		Mes, nil}
	response, err := e.SuppressEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.(CallResponse), err
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
