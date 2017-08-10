package frametest

import (
	"context"
	"github.com/guherbozdogan/mesos-go-http-client/client/frame"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = Describe("RecordIO", func() {

	var recIO frame.FrameIO
	var recChan chan string
	var errChan chan error

	var frameReadFunc = func(c context.Context, f frame.Frame, i int64) context.Context {
		recChan <- string(f)
		return c
	}
	var errFunc = func(c context.Context, i interface{}) context.Context {

		var err error = i.(error)
		//var s = err.Error()
		//fmt.Print(s)
		errChan <- err
		return c
	}

	BeforeEach(func() {
		recIO = frame.NewRecordIO()
		recChan = make(chan string, 10)
		errChan = make(chan error, 2)

	})
	Describe("Testing Read func", func() {
		Context("with 1 frames (non waited)", func() {

			It("should read ", func() {

				recIO.Read(context.Background(), inputLst["1 frame with no waiting"], frameReadFunc, errFunc)

				gomega.Eventually(recChan, 10).Should(gomega.Receive())

			})
			It("should fail with error if frame data is incomplete ", func() {

				recIO.Read(context.Background(), inputLst["1 frame with 10s waiting"], frameReadFunc, errFunc)

				gomega.Eventually(recChan, 10).Should(gomega.Receive())

			})

		})
		Context("with 2 frames (non waited)", func() {

			It("should read ", func() {
				//				buf.Write(inputLst[0])

				recIO.Read(context.Background(), inputLst["2 frames with no waiting"], frameReadFunc, errFunc)

				gomega.Eventually(recChan, 10).Should(gomega.Receive())

			})

		})
	})
})
