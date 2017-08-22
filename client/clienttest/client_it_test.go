package clienttest

import (
	//. "github.com/guherbozdogan/mesos-go-http-client/client/frame"

	"context"
	"fmt"
	"github.com/guherbozdogan/mesos-go-http-client/client/client"
	"github.com/guherbozdogan/mesos-go-http-client/client/frame"
	"github.com/guherbozdogan/mesos-go-http-client/client/pb/mesos/v1"
	"github.com/guherbozdogan/mesos-go-http-client/client/pb/mesos/v1/scheduler"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = Describe("It testing with running mesos", func() {

	Describe("dummy test", func() {
		Context("dummy test ", func() {
			It("dummy test", func() {
				//Ω(obictRecordIO).NotTo(Equal(nil))

				var recChan chan []byte
				var errChan chan error
				recChan = make(chan []byte, 10)
				errChan = make(chan error, 2)
				clients, err := client.NewHAClient("127.0.0.1/scheduler:5050;127.0.0.2/scheduler:5050",
					nil, nil, nil, nil)

				if err != nil {
					fmt.Println(err.Error())
				} else {
					clients.Clients[0].EndpointsofClient.FrameReadFunc = func(c context.Context, f frame.Frame, i int64) context.Context {
						fmt.Println(f)
						return c
					}

					clients.Clients[0].EndpointsofClient.FrameErrorFunc =
						func(c context.Context, i interface{}) context.Context {

							var err error = i.(error)

							fmt.Println(err.Error())
							return c
						}
					s1 := string("guhu")
					s2 := string("Example1")
					s3 := string("")
					b := bool(false)
					fl := float64(1000000000)
					t1 := mesos_v1.FrameworkInfo_Capability_Type(1)
					clients.Clients[0].EndpointsofClient.Subscribe(context.Background(),
						//&mesos_v1.FrameworkID,
						nil,

						&mesos_v1_scheduler.Call_Subscribe{
							&mesos_v1.FrameworkInfo{
								User: &s1, Name: &s2,
								Id: nil, FailoverTimeout: &fl,
								Checkpoint: &b,
								Role:       &s3,
								Roles:      []string{string("test")}, Capabilities: []*mesos_v1.FrameworkInfo_Capability{
									&mesos_v1.FrameworkInfo_Capability{
										&t1,
										nil}}},
							nil, nil})

					var s []byte
					var err error
					gomega.Eventually(recChan).Should(gomega.Receive(&s))
					gomega.Eventually(errChan).Should(gomega.Receive(&err))

				}
			})

		})
	})
})
