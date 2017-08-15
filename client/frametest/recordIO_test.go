package frametest

import (
	"context"
	"github.com/guherbozdogan/mesos-go-http-client/client/frame"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	. "github.com/onsi/gomega"
	//	"io"
	"bufio"
	"io"
)

var _ = Describe("RecordIO", func() {

	var recIO *frame.RecordIO
	var recChan chan []byte
	var errChan chan error

	var frameReadFunc = func(c context.Context, f frame.Frame, i int64) context.Context {
		recChan <- f
		return c
	}
	var errFunc = func(c context.Context, i interface{}) context.Context {

		var err error = i.(error)
		//var s = err.Error()
		//fmt.Print(s)
		errChan <- err
		return c
	}

	//mock read line replacement
	var readLnReplacement = func(r *bufio.Reader, rc io.ReadCloser) (frame.Bytes, error) {
		tmp, isok := rc.(*MultiMockFrameReaderCloser)
		if isok {
			if tmp.index < len(tmp.lst) {
				return tmp.lst[tmp.index].prefixBytes, nil
			} else {
				if tmp.lst[tmp.index-1].err != nil {
					if tmp.lst[tmp.index-1].isErrorOccuringAtPrefix {
						return []byte(""), tmp.lst[tmp.index-1].err
					} else {

						tmp.index--
						return tmp.lst[tmp.index].prefixBytes, nil

					}

				} else {
					return []byte(""), nil
				}

			}

		} else {
			tmp2, isok2 := rc.(*MockFrameReaderCloser)
			if isok2 {
				if !tmp2.bytesRead {
					return tmp2.prefixBytes, nil
				} else if tmp2.err != nil {
					if tmp2.isErrorOccuringAtPrefix {
						return []byte(""), nil
					} else {
						return tmp2.prefixBytes, nil
					}

				} else {
					return []byte(""), nil
				}
			} else {
				return []byte(""), nil
			}
		}
	}
	var earlierReadLine func(r *bufio.Reader, rc io.ReadCloser) (frame.Bytes, error)
	BeforeEach(func() {
		//frame.Readln =
		recIO = frame.NewRecordIO()
		earlierReadLine = recIO.ReadLine

		recIO.ReadLine = readLnReplacement
		recChan = make(chan []byte, 10)
		errChan = make(chan error, 2)

	})
	AfterEach(func() {
		recIO.ReadLine = earlierReadLine

	})
	Describe("Testing Read func", func() {
		Context("with 1 frames (non waited)", func() {

			It("should read ", func() {

				contextWithCancel, cancelFunc := context.WithCancel(context.Background())
				go timeoutHelper(2500000, cancelFunc)

				recIO.Read(contextWithCancel, inputLst["1 frame with no waiting"], frameReadFunc, errFunc)

				var s []byte
				gomega.Eventually(recChan, 10).Should(gomega.Receive(&s))
				Ω(s).Should(Equal(byteLst["1 frame with no waiting"]))

			})
			It("should handle EOF gracefully ", func() {

				contextWithCancel, cancelFunc := context.WithCancel(context.Background())
				go timeoutHelper(22500000, cancelFunc)

				go recIO.Read(contextWithCancel, inputLst["1 frame with no waiting & eof at end"], frameReadFunc, errFunc)

				var s []byte
				var err error
				gomega.Eventually(recChan, 2000).Should(gomega.Receive(&s))
				gomega.Eventually(errChan, 2000).Should(gomega.Receive(&err))

				Ω(s).Should(Equal(byteLst["1 frame with no waiting & eof at end"]))
				Ω(err).Should(Equal(io.EOF))

			})
			It("should handle specific io error gracefully ", func() {

				contextWithCancel, cancelFunc := context.WithCancel(context.Background())
				go timeoutHelper(22500000, cancelFunc)

				recIO.Read(contextWithCancel, inputLst["1 frame with no waiting & prg err at end"], frameReadFunc, errFunc)

				var s []byte
				var err error
				gomega.Eventually(recChan, 10).Should(gomega.Receive(&s))
				gomega.Eventually(errChan, 10).Should(gomega.Receive(&err))

				Ω(s).Should(Equal(byteLst["1 frame with no waiting & prg err at end"]))
				Ω(err).Should(Equal(io.ErrNoProgress))
			})

		}) /*
			Context("with 2 frames (non waited)", func() {

				It("should read ", func() {
					//				buf.Write(inputLst[0])

					recIO.Read(context.Background(), inputLst["2 frames with no waiting"], frameReadFunc, errFunc)

					gomega.Eventually(recChan, 10).Should(gomega.Receive())

				})

			})
		*/

	})
})
