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

	var byteLst = map[string][]byte{
		"1 frame with no waiting":                    []byte("{\"type\": \"SUBSCRIBED\",\"subscribed\": {\"framework_id\": {\"value\":\"12220-3440-12532-2345\"},\"heartbeat_interval_seconds\":15.0}"),
		"1 frame with waiting":                       []byte("{\"type\": \"SUBSCRIBED\",\"subscribed\": {\"framework_id\": {\"value\":\"12220-3440-12532-2345\"},\"heartbeat_interval_seconds\":15.0}"),
		"2 frame with no waiting & eof at prefix":    []byte("{\"type\": \"SUBSCRIBED\",\"subscribed\": {\"framework_id\": {\"value\":\"12220-3440-12532-2345\"},\"heartbeat_interval_seconds\":15.0}"),
		"2 frame with no waiting & eof during frame": []byte("{\"type\": \"SUBSCRIBED\",\"subscribed\": {\"framework_id\": {\"value\":\"12220-3440-12532-2345\"},\"heartbeat_interval_seconds\":15.0}"),

		"2 frame with no waiting & prg err at prefix":    []byte("{\"type\": \"SUBSCRIBED\",\"subscribed\": {\"framework_id\": {\"value\":\"12220-3440-12532-2345\"},\"heartbeat_interval_seconds\":15.0}"),
		"2 frame with no waiting & prg err during frame": []byte("{\"type\": \"SUBSCRIBED\",\"subscribed\": {\"framework_id\": {\"value\":\"12220-3440-12532-2345\"},\"heartbeat_interval_seconds\":15.0}"),
		"2 frame with waiting & eof at prefix":           []byte("{\"type\": \"SUBSCRIBED\",\"subscribed\": {\"framework_id\": {\"value\":\"12220-3440-12532-2345\"},\"heartbeat_interval_seconds\":15.0}"),
		"2 frame with waiting & eof during frame":        []byte("{\"type\": \"SUBSCRIBED\",\"subscribed\": {\"framework_id\": {\"value\":\"12220-3440-12532-2345\"},\"heartbeat_interval_seconds\":15.0}"),

		"2 frame with waiting & prg err at prefix":    []byte("{\"type\": \"SUBSCRIBED\",\"subscribed\": {\"framework_id\": {\"value\":\"12220-3440-12532-2345\"},\"heartbeat_interval_seconds\":15.0}"),
		"2 frame with waiting & prg err during frame": []byte("{\"type\": \"SUBSCRIBED\",\"subscribed\": {\"framework_id\": {\"value\":\"12220-3440-12532-2345\"},\"heartbeat_interval_seconds\":15.0}"),

		"Context cancelled during prefix": []byte(""),
		"Context cancelled during frame":  []byte("trial"),
		"Context cancelled before prefix": []byte("")}

	// test case map
	var inputLst = map[string]MockReadCloser{
		"1 frame with no waiting": &MultiMockFrameReaderCloser{index: 0,
			lst: []*MockFrameReaderCloser{
				&MockFrameReaderCloser{prefixBytes: []byte("121"), bytes: byteLst["1 frame with no waiting"]},
			}},

		"1 frame with waiting": &MultiMockFrameReaderCloser{index: 0,
			lst: []*MockFrameReaderCloser{
				&MockFrameReaderCloser{prefixBytes: []byte("121"), bytes: byteLst["1 frame with waiting"]},
			}},
		//	"1 incomplete frame with no waiting": &MockFrameReaderCloser{prefixBytes: []byte("121"), bytes: byteLst["1 incomplete frame with no waiting"]},
		"2 frame with no waiting & eof at prefix": &MultiMockFrameReaderCloser{index: 0,
			lst: []*MockFrameReaderCloser{
				&MockFrameReaderCloser{prefixBytes: []byte("121"), bytes: byteLst["2 frame with no waiting & eof at prefix"]},
				&MockFrameReaderCloser{isErrorOccuringAtPrefix: true, err: io.EOF}}},

		"2 frame with waiting & eof at prefix": &MultiMockFrameReaderCloser{index: 0,
			lst: []*MockFrameReaderCloser{
				&MockFrameReaderCloser{period: 2, prefixBytes: []byte("121"), bytes: byteLst["2 frame with waiting & eof at prefix"]},
				&MockFrameReaderCloser{period: 2, isErrorOccuringAtPrefix: true, err: io.EOF}}},

		"2 frame with no waiting & eof during frame": &MultiMockFrameReaderCloser{index: 0,
			lst: []*MockFrameReaderCloser{
				&MockFrameReaderCloser{prefixBytes: []byte("121"), bytes: byteLst["2 frame with no waiting & eof during frame"]},
				&MockFrameReaderCloser{isErrorOccuringAtPrefix: false, prefixBytes: []byte("121"), err: io.EOF}}},

		"2 frame with waiting & eof during frame": &MultiMockFrameReaderCloser{index: 0,
			lst: []*MockFrameReaderCloser{
				&MockFrameReaderCloser{period: 2, prefixBytes: []byte("121"), bytes: byteLst["2 frame with waiting & eof during frame"]},
				&MockFrameReaderCloser{period: 2, isErrorOccuringAtPrefix: false, prefixBytes: []byte("121"), err: io.EOF}}},

		"2 frame with no waiting & prg err at prefix": &MultiMockFrameReaderCloser{index: 0,
			lst: []*MockFrameReaderCloser{
				&MockFrameReaderCloser{prefixBytes: []byte("121"), bytes: byteLst["2 frame with no waiting & prg err at prefix"]},
				&MockFrameReaderCloser{isErrorOccuringAtPrefix: true, err: io.ErrNoProgress}}},

		"2 frame with no waiting & prg err during frame": &MultiMockFrameReaderCloser{index: 0,
			lst: []*MockFrameReaderCloser{
				&MockFrameReaderCloser{prefixBytes: []byte("121"), bytes: byteLst["2 frame with no waiting & prg err during frame"]},
				&MockFrameReaderCloser{isErrorOccuringAtPrefix: false, prefixBytes: []byte("121"), err: io.ErrNoProgress}}},

		"2 frame with waiting & prg err at prefix": &MultiMockFrameReaderCloser{index: 0,
			lst: []*MockFrameReaderCloser{
				&MockFrameReaderCloser{period: 1, prefixBytes: []byte("121"), bytes: byteLst["2 frame with waiting & prg err at prefix"]},
				&MockFrameReaderCloser{period: 1, isErrorOccuringAtPrefix: true, err: io.ErrNoProgress}}},
		"2 frame with waiting & prg err during frame": &MultiMockFrameReaderCloser{index: 0,
			lst: []*MockFrameReaderCloser{
				&MockFrameReaderCloser{period: 1, prefixBytes: []byte("121"), bytes: byteLst["2 frame with waiting & prg err during frame"]},
				&MockFrameReaderCloser{period: 1, isErrorOccuringAtPrefix: false, prefixBytes: []byte("121"), err: io.ErrNoProgress}}},

		"Context cancelled during prefix": &MultiMockFrameReaderCloser{index: 0, lst: []*MockFrameReaderCloser{
			&MockFrameReaderCloser{prefixBytes: []byte(""), bytes: byteLst["Context cancelled during prefix"]}}},

		"Context cancelled during frame": &MultiMockFrameReaderCloser{index: 0, lst: []*MockFrameReaderCloser{
			&MockFrameReaderCloser{period: 7, prefixBytes: []byte("121"), bytes: byteLst["Context cancelled during frame"]}}},

		"Context cancelled before prefix": &MultiMockFrameReaderCloser{index: 0, lst: []*MockFrameReaderCloser{
			&MockFrameReaderCloser{prefixBytes: []byte("12"), bytes: byteLst["Context cancelled before prefix"]}}}}

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
	var readLnReplacement = func(r *bufio.Reader, rc io.ReadCloser, icd *uint32) (frame.Bytes, error) {
		tmp, isok := rc.(*MultiMockFrameReaderCloser)
		if isok {
			if tmp.index < len(tmp.lst) {
				if tmp.lst[tmp.index].isErrorOccuringAtPrefix {
					return []byte(""), tmp.lst[tmp.index].err
				} else {
					return tmp.lst[tmp.index].prefixBytes, nil
				}

			} else {
				return []byte(""), nil
			}

		} else {
			tmp2, isok2 := rc.(*MockFrameReaderCloser)
			if isok2 {
				if tmp2.isErrorOccuringAtPrefix {
					return []byte(""), tmp2.err
				} else {
					return tmp2.prefixBytes, nil
				}
			} else {
				return []byte(""), nil
			}
		}

	}
	var earlierReadLine func(r *bufio.Reader, rc io.ReadCloser, icd *uint32) (frame.Bytes, error)
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
		Context("with  frames (non waited)", func() {
			It("should read ", func() {

				contextWithCancel, cancelFunc := context.WithCancel(context.Background())
				go timeoutHelper(2500000, cancelFunc)

				recIO.Read(contextWithCancel, inputLst["1 frame with no waiting"], frameReadFunc, errFunc)

				var s []byte
				gomega.Eventually(recChan, 10).Should(gomega.Receive(&s))
				Ω(s).Should(Equal(byteLst["1 frame with no waiting"]))

			})
			It("should handle EOF gracefully during prefix reading", func() {

				contextWithCancel, cancelFunc := context.WithCancel(context.Background())
				go timeoutHelper(22500000, cancelFunc)

				go recIO.Read(contextWithCancel, inputLst["2 frame with no waiting & eof at prefix"], frameReadFunc, errFunc)

				var s []byte
				var err error
				gomega.Eventually(recChan, 2000).Should(gomega.Receive(&s))
				gomega.Eventually(errChan, 2000).Should(gomega.Receive(&err))

				Ω(s).Should(Equal(byteLst["2 frame with no waiting & eof at prefix"]))
				Ω(err).Should(Equal(io.EOF))

			})
			It("should handle EOF gracefully during frame", func() {

				contextWithCancel, cancelFunc := context.WithCancel(context.Background())
				go timeoutHelper(22500000, cancelFunc)

				go recIO.Read(contextWithCancel, inputLst["2 frame with no waiting & eof during frame"], frameReadFunc, errFunc)

				var s []byte
				var err error
				gomega.Eventually(recChan, 2000).Should(gomega.Receive(&s))
				gomega.Eventually(errChan, 2000).Should(gomega.Receive(&err))

				Ω(s).Should(Equal(byteLst["2 frame with no waiting & eof during frame"]))
				Ω(err).Should(Equal(frame.ErrorChannelClosedBeforeReceivingFrame))

			})
			It("should handle specific io error gracefully at prefix", func() {

				contextWithCancel, cancelFunc := context.WithCancel(context.Background())
				go timeoutHelper(22500000000, cancelFunc)

				recIO.Read(contextWithCancel, inputLst["2 frame with no waiting & prg err at prefix"], frameReadFunc, errFunc)

				var s []byte
				var err error
				gomega.Eventually(recChan, 10).Should(gomega.Receive(&s))
				gomega.Eventually(errChan, 10).Should(gomega.Receive(&err))

				Ω(s).Should(Equal(byteLst["2 frame with no waiting & prg err at prefix"]))
				Ω(err).Should(Equal(io.ErrNoProgress))
			})
			It("should handle specific io error gracefully during frame ", func() {

				contextWithCancel, cancelFunc := context.WithCancel(context.Background())
				go timeoutHelper(22500000000, cancelFunc)

				recIO.Read(contextWithCancel, inputLst["2 frame with no waiting & prg err during frame"], frameReadFunc, errFunc)

				var s []byte
				var err error
				gomega.Eventually(recChan, 10).Should(gomega.Receive(&s))
				gomega.Eventually(errChan, 10).Should(gomega.Receive(&err))

				Ω(s).Should(Equal(byteLst["2 frame with no waiting & prg err during frame"]))
				Ω(err).Should(Equal(io.ErrNoProgress))
			})

		})
		Context("with frames ( waited)", func() {

			It("should read ", func() {

				contextWithCancel, cancelFunc := context.WithCancel(context.Background())
				go timeoutHelper(52500000, cancelFunc)

				recIO.Read(contextWithCancel, inputLst["1 frame with waiting"], frameReadFunc, errFunc)

				var s []byte
				gomega.Eventually(recChan, 30).Should(gomega.Receive(&s))
				Ω(s).Should(Equal(byteLst["1 frame with waiting"]))

			})
			It("should handle EOF gracefully during prefix reading", func() {

				contextWithCancel, cancelFunc := context.WithCancel(context.Background())
				go timeoutHelper(52500000000, cancelFunc)

				go recIO.Read(contextWithCancel, inputLst["2 frame with waiting & eof at prefix"], frameReadFunc, errFunc)

				var s []byte
				var err error
				gomega.Eventually(recChan, 20).Should(gomega.Receive(&s))
				gomega.Eventually(errChan, 20).Should(gomega.Receive(&err))

				Ω(s).Should(Equal(byteLst["2 frame with waiting & eof at prefix"]))
				Ω(err).Should(Equal(io.EOF))

			})
			It("should handle EOF gracefully during frame", func() {

				contextWithCancel, cancelFunc := context.WithCancel(context.Background())
				go timeoutHelper(52500000000, cancelFunc)

				go recIO.Read(contextWithCancel, inputLst["2 frame with waiting & eof during frame"], frameReadFunc, errFunc)

				var s []byte
				var err error
				gomega.Eventually(recChan, 32).Should(gomega.Receive(&s))
				gomega.Eventually(errChan, 32).Should(gomega.Receive(&err))

				Ω(s).Should(Equal(byteLst["2 frame with waiting & eof during frame"]))
				Ω(err).Should(Equal(frame.ErrorChannelClosedBeforeReceivingFrame))

			})
			It("should handle specific io error gracefully at prefix", func() {

				contextWithCancel, cancelFunc := context.WithCancel(context.Background())
				go timeoutHelper(52500000000, cancelFunc)

				recIO.Read(contextWithCancel, inputLst["2 frame with waiting & prg err at prefix"], frameReadFunc, errFunc)

				var s []byte
				var err error
				gomega.Eventually(recChan, 32).Should(gomega.Receive(&s))
				gomega.Eventually(errChan, 32).Should(gomega.Receive(&err))

				Ω(s).Should(Equal(byteLst["2 frame with waiting & prg err at prefix"]))
				Ω(err).Should(Equal(io.ErrNoProgress))
			})
			It("should handle specific io error gracefully during frame ", func() {

				contextWithCancel, cancelFunc := context.WithCancel(context.Background())
				go timeoutHelper(52500000000, cancelFunc)

				recIO.Read(contextWithCancel, inputLst["2 frame with waiting & prg err during frame"], frameReadFunc, errFunc)

				var s []byte
				var err error
				gomega.Eventually(recChan, 32).Should(gomega.Receive(&s))
				gomega.Eventually(errChan, 32).Should(gomega.Receive(&err))

				Ω(s).Should(Equal(byteLst["2 frame with waiting & prg err during frame"]))
				Ω(err).Should(Equal(io.ErrNoProgress))
			})

		})

		Context("with cancel called on prefix", func() {

			It("should gracefully handle ", func() {
				//				buf.Write(inputLst[0])

				contextWithCancel, cancelFunc := context.WithCancel(context.Background())
				go timeoutHelper(1500000000, cancelFunc)

				recIO.Read(contextWithCancel, inputLst["Context cancelled during prefix"], frameReadFunc, errFunc)

				var err error
				gomega.Eventually(errChan, 32).Should(gomega.Receive(&err))

				Ω(err).Should(Equal(frame.ErrorParentContextBeenCancelled))

			})

		})
		Context("with cancel called before prefix", func() {

			It("should gracefully handle ", func() {
				//				buf.Write(inputLst[0])

				contextWithCancel, cancelFunc := context.WithCancel(context.Background())
				cancelFunc()

				recIO.Read(contextWithCancel, inputLst["Context cancelled before prefix"], frameReadFunc, errFunc)

				var err error
				gomega.Eventually(errChan, 32).Should(gomega.Receive(&err))

				Ω(err).Should(Equal(frame.ErrorParentContextBeenCancelled))

			})

		})
		Context("with cancel called during frame", func() {

			It("should gracefully handle ", func() {
				//				buf.Write(inputLst[0])
				contextWithCancel, cancelFunc := context.WithCancel(context.Background())
				go timeoutHelper(2500000000, cancelFunc)

				recIO.Read(contextWithCancel, inputLst["Context cancelled during frame"], frameReadFunc, errFunc)

				var err error
				gomega.Eventually(errChan, 32).Should(gomega.Receive(&err))

				Ω(err).Should(Equal(frame.ErrorParentContextBeenCancelled))

			})

		})

	})
})
