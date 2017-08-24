package client

import (
	"net/url"
	"strings"
	"time"

	jujuratelimit "github.com/juju/ratelimit"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/sony/gobreaker"

	"github.com/guherbozdogan/kit/circuitbreaker"
	. "github.com/guherbozdogan/kit/endpoint"

	"github.com/guherbozdogan/kit/log"
	"github.com/guherbozdogan/kit/ratelimit"
	//	"github.com/guherbozdogan/kit/tracing/opentracing"
	httptransport "github.com/guherbozdogan/kit/transport/http"
	"github.com/guherbozdogan/mesos-go-http-client/client/frame"
)

//yeppp after break/rest, check:
//-- this should create a client  for endpoints struct
//---the client currently would provide grpc data structs for passing data and also json
//-- the client has frame reading capa for which response events are also returned to client
//--- investigate case where two active redundant schedulers exist
//--- frame work id, executor id, all relation to be read
//--- add backoff strategy for network partition or failures
//---add stream id header as with a middleware
//---add opentrace support
//---add direct zookeeper support
//
//
//--then scheduler implementation with yaml files

type Client struct {
	URL               *url.URL
	EndpointsofClient *Endpoints
}
type Clients struct {
	Clients     []*Client
	masterIndex int32
}

func (c *Clients) initLeaderInfo() error {
	//function to set masterinfo
	return nil
}

//non thread safe
func (c *Clients) getClientofLeader() (lc *Client, err error) {
	if c.masterIndex == -1 {
		err = c.initLeaderInfo()
		if err != nil {
			return nil, err
		} else {
			return c.Clients[c.masterIndex], nil
		}

	} else {
		return c.Clients[c.masterIndex], nil
	}
}

// New returns an AddService backed by an HTTP server living at the remote
// instance. We expect instance to come from a service discovery system, so
// likely of the form "host:port;host:port".

func NewHAClient(instances string, tracer stdopentracing.Tracer, logger log.Logger, f frame.FrameReadFunc, e frame.ErrorFunc) (*Clients, error) {

	if tracer == nil {

	}
	var hostList []string
	hostList = strings.Split(instances, ";")

	var clientList *Clients = &Clients{masterIndex: -1}
	clientList.Clients = make([]*Client, len(hostList))
	for i, item := range hostList {

		if !strings.HasPrefix(item, "http") {
			item = string("http://") + item
		}

		u, err := url.Parse(item)
		if err != nil {
			return nil, err
		} else {
			client := &Client{u, &Endpoints{FrameReadFunc: f, FrameErrorFunc: e}}
			clientList.Clients[i] = client
			//clientList.Clients = append(clientList.Clients, client)
		}
	}

	// We construct a single ratelimiter middleware, to limit the total outgoing
	// QPS from this client to all methods on the remote instance. We also
	// construct per-endpoint circuitbreaker middlewares to demonstrate how
	// that's done, although they could easily be combined into a single breaker
	// for the entire remote instance, too.

	i := 0
	for i < len(clientList.Clients) {

		var c *Client = clientList.Clients[i]
		u := c.URL

		limiter := ratelimit.NewTokenBucketLimiter(jujuratelimit.NewBucketWithRate(100, 100))

		var subscribeEndPoint Endpoint = nil

		subscribeEndPoint = httptransport.NewBSClient(
			"POST",
			copyURL(u, "/api/v1/scheduler"),
			EncodeHTTPSubscribeRequestToJSON,
			DecodeHTTPSubscribeResponse,
			c.EndpointsofClient.FrameReadFunc,
			frame.CTRecordIO,
			c.EndpointsofClient.FrameErrorFunc,
			//httptransport.ClientBefore(opentracing.ToHTTPRequest(tracer, logger)),
		).Endpoint()
		//	subscribeEndPoint = opentracing.TraceClient(tracer, "Subscribe")(subscribeEndPoint)
		subscribeEndPoint = limiter(subscribeEndPoint)
		subscribeEndPoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "Subscribe",
			Timeout: 30 * time.Second,
		}))(subscribeEndPoint)

		//var t *Endpoints = c.EndpointsofClient

		c.EndpointsofClient.SubscribeEndpoint = subscribeEndPoint
		//t.SubscribeEndpoint = subscribeEndPoint

		i++
	}
	return clientList, nil
}

func copyURL(base *url.URL, path string) *url.URL {
	next := *base
	next.Path = path
	return &next
}
