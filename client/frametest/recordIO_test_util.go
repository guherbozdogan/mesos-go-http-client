package frametest

import (
	"io"
	"time"
)

//mock readcloser io implementation's interface
type MockReadCloser interface {
	Close() error
	Read(p []byte) (n int, err error)
}

//multiple instances
type MultiMockFrameReaderCloser struct {
	lst   []*MockFrameReaderCloser
	index int
}

//mock implementation
type MockFrameReaderCloser struct {
	isErrorOccuringAtPrefix bool
	bytes                   []byte
	prefixBytes             []byte
	err                     error
	bytesRead               bool
	period                  time.Duration //nano seconds
}

//constructor
func NewMockFrameReaderCloser() *MockFrameReaderCloser {
	return &MockFrameReaderCloser{bytesRead: false, isErrorOccuringAtPrefix: false}
}

//constructor
func NewMultiMockFrameReaderCloser() *MultiMockFrameReaderCloser {
	return &MultiMockFrameReaderCloser{index: 0}
}

//reset to initialize
func (cb *MockFrameReaderCloser) Reset() error {
	//and the error is initialized to no-error
	cb.bytes = nil
	cb.err = nil
	cb.period = time.Duration(0)
	return nil

}

// ReaderCloser Close implementation
func (cb MultiMockFrameReaderCloser) Close() error {
	//and the error is initialized to no-error
	return nil
}

// ReaderCloser Close implementation
func (cb MockFrameReaderCloser) Close() error {
	//and the error is initialized to no-error
	return nil
}

//helper function for mock read implementation
func (cb *MockFrameReaderCloser) ReadSub(p []byte) (n int, err error) {
	if cb.bytes != nil {

		if cb.bytesRead && cb.err != nil {
			return 0, cb.err
		} else if !cb.bytesRead {
			//!guhu risky copy
			copy(p, cb.bytes)

			cb.bytesRead = true
			return len(p), nil
		} else {
			return 0, nil
		}
	} else if cb.err != nil {
		return 0, cb.err
	} else {
		return 0, nil
	}
}

func (cb *MultiMockFrameReaderCloser) Read(p []byte) (n int, err error) {

	if cb.index < len(cb.lst) {
		cb.index++
		return cb.lst[cb.index-1].Read(p)
	} else {
		return 0, nil
	}

}

// ReaderCloser Read implementation
func (cb *MockFrameReaderCloser) Read(p []byte) (n int, err error) {

	if cb.period == 0 {
		return cb.ReadSub(p)
	} else {

		timeout := make(chan bool, 1)
		go func() {
			time.Sleep(cb.period * time.Nanosecond)
			timeout <- true
		}()

		select {
		case <-timeout:
			return cb.ReadSub(p)
		}
	}
	return 0, nil
}
func timeoutHelper(priod time.Duration, f func()) {
	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(priod * time.Nanosecond)
		timeout <- true
	}()

	select {
	case <-timeout:
		f()
	}
}

var byteLst = map[string][]byte{
	"1 frame with no waiting":                  []byte("{\"type\": \"SUBSCRIBED\",\"subscribed\": {\"framework_id\": {\"value\":\"12220-3440-12532-2345\"},\"heartbeat_interval_seconds\":15.0}"),
	"1 frame with 2s waiting":                  []byte("{\"type\": \"SUBSCRIBED\",\"subscribed\": {\"framework_id\": {\"value\":\"12220-3440-12532-2345\"},\"heartbeat_interval_seconds\":15.0}"),
	"1 incomplete frame with no waiting":       []byte("{\"type\": \"SUBSCRIBED\",\"subscribed\": {\"framework_id\": {\"value\":\"12220-3440-12532-2345\"},\"heartbeat_interval_seconds\""),
	"1 frame with no waiting & eof at end":     []byte("{\"type\": \"SUBSCRIBED\",\"subscribed\": {\"framework_id\": {\"value\":\"12220-3440-12532-2345\"},\"heartbeat_interval_seconds\":15.0}"),
	"1 frame with no waiting & prg err at end": []byte("{\"type\": \"SUBSCRIBED\",\"subscribed\": {\"framework_id\": {\"value\":\"12220-3440-12532-2345\"},\"heartbeat_interval_seconds\":15.0}"),
	"2 frames with no waiting":                 []byte("{\"type\": \"SUBSCRIBED\",\"subscribed\": {\"framework_id\": {\"value\":\"12220-3440-12532-2345\"},\"heartbeat_interval_seconds\":15.0}20\n{\"type\":\"HEARTBEAT\"}"),
	"2 frames with 1s waiting":                 []byte("{\"type\": \"SUBSCRIBED\",\"subscribed\": {\"framework_id\": {\"value\":\"12220-3440-12532-2345\"},\"heartbeat_interval_seconds\":15.0}20\n{\"type\":\"HEARTBEAT\"}"),
	"insufficient frame":                       []byte("124\n")}

// test case map
var inputLst = map[string]MockReadCloser{
	"1 frame with no waiting":            &MockFrameReaderCloser{prefixBytes: []byte("121"), bytes: byteLst["1 frame with no waiting"]},
	"1 frame with 2s waiting":            &MockFrameReaderCloser{prefixBytes: []byte("121"), period: 2 * time.Second, bytes: byteLst["1 frame with 2s waiting"]},
	"1 incomplete frame with no waiting": &MockFrameReaderCloser{prefixBytes: []byte("121"), bytes: byteLst["1 incomplete frame with no waiting"]},
	"1 frame with no waiting & eof at end": &MultiMockFrameReaderCloser{index: 0,
		lst: []*MockFrameReaderCloser{
			&MockFrameReaderCloser{prefixBytes: []byte("121"), bytes: byteLst["1 frame with no waiting & eof at end"]},
			&MockFrameReaderCloser{isErrorOccuringAtPrefix: true, prefixBytes: []byte("121"), err: io.EOF, bytes: byteLst["1 frame with no waiting"]}}},
	"1 frame with no waiting & prg err at end": &MultiMockFrameReaderCloser{index: 0,
		lst: []*MockFrameReaderCloser{
			&MockFrameReaderCloser{prefixBytes: []byte("121"), bytes: byteLst["1 frame with no waiting & prg err at end"]},
			&MockFrameReaderCloser{isErrorOccuringAtPrefix: true, prefixBytes: []byte("121"), err: io.ErrNoProgress, bytes: byteLst["1 frame with no waiting & prg err at end"]}}},
	"2 frames with no waiting": &MockFrameReaderCloser{prefixBytes: []byte("121"), bytes: byteLst["2 frames with no waiting"]},
	"2 frames with 1s waiting": &MockFrameReaderCloser{prefixBytes: []byte("121"), bytes: byteLst["2 frames with 1s waiting"]},
	"insufficient frame":       &MockFrameReaderCloser{prefixBytes: []byte("121"), bytes: byteLst["insufficient frame"]}}
