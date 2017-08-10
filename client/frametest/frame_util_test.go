package frametest

import (
	//. "github.com/guherbozdogan/mesos-go-http-client/client/frame"

	"context"
	"github.com/guherbozdogan/mesos-go-http-client/client/frame"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Frame Util", func() {
	var (
		obiNilRecordIO *frame.RecordIO
		obictRecordIO  = frame.NewFrameIO(frame.CTRecordIO)
		obictOtherIO   = frame.NewFrameIO(frame.CTOtherType)
	)

	Describe("Testing RecordIO constructor", func() {
		Context("with CTRecordIO enum ", func() {
			It("should return non nil", func() {
				Ω(obictRecordIO).NotTo(Equal(nil))

			})
			It("should return RecordIO type", func() {
				Ω(obictRecordIO).Should(BeAssignableToTypeOf(obiNilRecordIO))
			})

		})
		Context("with Other enum ", func() {
			It("should return nil", func() {
				Ω(obictOtherIO).Should(BeNil())
			})

		})
	})
	Describe("DecorateFrameReadFunc", func() {
		var (
			arrEmptyFuncList = []frame.FrameReadFunc{}
			recChan          = make(chan int, 7)
			arrPreFuncList   = []frame.FrameReadFunc{
				func(c context.Context, f frame.Frame, i int64) context.Context {
					recChan <- 0
					return c
				},
				func(c context.Context, f frame.Frame, i int64) context.Context {
					recChan <- 1
					return c
				}}
			arrPostFuncList = []frame.FrameReadFunc{
				func(c context.Context, f frame.Frame, i int64) context.Context {
					recChan <- 3
					return c
				},
				func(c context.Context, f frame.Frame, i int64) context.Context {
					recChan <- 4
					return c
				}}
			arrDecoratedFunc = func(c context.Context, f frame.Frame, i int64) context.Context {
				recChan <- 2
				return c
			}

			fFinalFunc = frame.DecorateFrameReadFunc(arrPreFuncList, arrPostFuncList, arrDecoratedFunc)

			c     = context.Background()
			bytes = []byte{}
			//func()
		)

		Context("with nonEmpty functions ", func() {
			It("should execute functions ordered", func() {
				fFinalFunc(c, bytes, 0)

				var ivalue int
				Ω(recChan).Should(Receive(&ivalue))
				Ω(ivalue).Should(Equal(0))
				Ω(recChan).Should(Receive(&ivalue))
				Ω(ivalue).Should(Equal(1))
				Ω(recChan).Should(Receive(&ivalue))
				Ω(ivalue).Should(Equal(2))
				Ω(recChan).Should(Receive(&ivalue))
				Ω(ivalue).Should(Equal(3))
				Ω(recChan).Should(Receive(&ivalue))
				Ω(ivalue).Should(Equal(4))

			})

		})
		Context("with Empty pre functions ", func() {
			It("should execute correctly", func() {
				recChan = make(chan int, 7)
				fFinalFuncEmptyPre := frame.DecorateFrameReadFunc(arrEmptyFuncList, arrPostFuncList, arrDecoratedFunc)

				c = context.Background()
				bytes = []byte{}

				fFinalFuncEmptyPre(c, bytes, 0)

				var ivalue int

				Ω(recChan).Should(Receive(&ivalue))
				Ω(ivalue).Should(Equal(2))
				Ω(recChan).Should(Receive(&ivalue))
				Ω(ivalue).Should(Equal(3))
				Ω(recChan).Should(Receive(&ivalue))
				Ω(ivalue).Should(Equal(4))

			})

		})
		Context("with Empty post functions ", func() {
			It("should execute correctly", func() {
				recChan = make(chan int, 5)
				fFinalFuncEmptyPre := frame.DecorateFrameReadFunc(arrPreFuncList, arrEmptyFuncList, arrDecoratedFunc)

				c = context.Background()
				bytes = []byte{}

				fFinalFuncEmptyPre(c, bytes, 0)

				var ivalue int

				Ω(recChan).Should(Receive(&ivalue))
				Ω(ivalue).Should(Equal(0))
				Ω(recChan).Should(Receive(&ivalue))
				Ω(ivalue).Should(Equal(1))
				Ω(recChan).Should(Receive(&ivalue))
				Ω(ivalue).Should(Equal(2))

			})

		})
		Context("with nil functions ", func() {
			It("should execute correctly", func() {
				recChan = make(chan int, 5)
				fFinalFuncEmptyPre := frame.DecorateFrameReadFunc(arrPreFuncList, arrEmptyFuncList, nil)

				c = context.Background()
				bytes = []byte{}

				fFinalFuncEmptyPre(c, bytes, 0)

				var ivalue int

				Ω(recChan).Should(Receive(&ivalue))
				Ω(ivalue).Should(Equal(0))
				Ω(recChan).Should(Receive(&ivalue))
				Ω(ivalue).Should(Equal(1))

			})

		})

	})
})
