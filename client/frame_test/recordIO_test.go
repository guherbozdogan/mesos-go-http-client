package frame_test

import (
	//. "github.com/guherbozdogan/mesos-go-http-client/client/frame"

	. "github.com/onsi/ginkgo"
	//. "github.com/onsi/gomega"
	"context"
	//	"fmt"
	"bytes"
	"github.com/guherbozdogan/mesos-go-http-client/client/frame"
	"github.com/onsi/gomega"
	//"os"
	"io"
)

type ClosingBuffer struct {
	*bytes.Buffer
}

func (cb *ClosingBuffer) Close() error {
	//and the error is initialized to no-error
	return nil
}

var _ = Describe("RecordIO", func() {

	//rp, rw := io.Pipe()
	var inputLst = [][]byte{
		[]byte("120\n{\"type\": \"SUBSCRIBED\",\"subscribed\": {\"framework_id\": {\"value\":\"12220-3440-12532-2345\"},\"heartbeat_interval_seconds\":15.0}"),
		[]byte(`124\n`)}
	Describe("Testing Read func", func() {
		Context("with correct param ", func() {

			recChan := make(chan string, 10)
			var frameReadFunc = func(c context.Context, f frame.Frame, i int64) context.Context {
				recChan <- string(f)
				return c
			}
			var errFunc = func(c context.Context, i interface{}) context.Context {

				return c
			}

			var recIO = frame.NewRecordIO()
			var tmpRC = &ClosingBuffer{bytes.NewBuffer(inputLst[0])}
			var rc io.ReadCloser
			rc = tmpRC
			recIO.Read(context.Background(), rc, frameReadFunc, errFunc)

			It("should read ", func() {
				//				buf.Write(inputLst[0])

				gomega.Eventually(recChan, 10).Should(gomega.Receive())

			})
		})
	})
})
