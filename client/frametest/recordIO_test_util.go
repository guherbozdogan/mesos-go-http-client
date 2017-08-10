package frametest

import (
	"io"
	"time"
)

type MockFrameReaderCloser struct {
	bytes  []byte
	err    error
	period time.Duration //nano seconds
}

func NewMockFrameReaderCloser() *MockFrameReaderCloser {
	return &MockFrameReaderCloser{}
}

func (cb *MockFrameReaderCloser) Reset() error {
	//and the error is initialized to no-error
	cb.bytes = nil
	cb.err = nil
	cb.period = time.Duration(0)
	return nil
}

func (cb *MockFrameReaderCloser) Close() error {
	//and the error is initialized to no-error
	return nil
}
func (cb *MockFrameReaderCloser) ReadSub(p []byte) (n int, err error) {
	if cb.err != nil {
		return 0, err
	} else if cb.bytes != nil {

		//!guhu risky copy
		copy(p, cb.bytes)
		return len(p), nil
	} else {
		return 0, nil
	}

}

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

var inputLst = map[string]*MockFrameReaderCloser{
	"1 frame with no waiting":                  &MockFrameReaderCloser{bytes: []byte("121\n{\"type\": \"SUBSCRIBED\",\"subscribed\": {\"framework_id\": {\"value\":\"12220-3440-12532-2345\"},\"heartbeat_interval_seconds\":15.0}")},
	"1 frame with 10s waiting":                 &MockFrameReaderCloser{period: 2 * time.Second, bytes: []byte("121\n{\"type\": \"SUBSCRIBED\",\"subscribed\": {\"framework_id\": {\"value\":\"12220-3440-12532-2345\"},\"heartbeat_interval_seconds\":15.0}")},
	"1 incomplete frame with no waiting":       &MockFrameReaderCloser{bytes: []byte("121\n{\"type\": \"SUBSCRIBED\",\"subscribed\": {\"framework_id\": {\"value\":\"12220-3440-12532-2345\"},\"heartbeat_interval_seconds\"")},
	"1 frame with no waiting & eof at end":     &MockFrameReaderCloser{err: io.EOF, bytes: []byte("121\n{\"type\": \"SUBSCRIBED\",\"subscribed\": {\"framework_id\": {\"value\":\"12220-3440-12532-2345\"},\"heartbeat_interval_seconds\":15.0}")},
	"1 frame with no waiting & prg err at end": &MockFrameReaderCloser{err: io.ErrNoProgress, bytes: []byte("121\n{\"type\": \"SUBSCRIBED\",\"subscribed\": {\"framework_id\": {\"value\":\"12220-3440-12532-2345\"},\"heartbeat_interval_seconds\":15.0}")},
	"2 frames with no waiting":                 &MockFrameReaderCloser{bytes: []byte("121\n{\"type\": \"SUBSCRIBED\",\"subscribed\": {\"framework_id\": {\"value\":\"12220-3440-12532-2345\"},\"heartbeat_interval_seconds\":15.0}20\n{\"type\":\"HEARTBEAT\"}")},
	"2 frames with 1s waiting":                 &MockFrameReaderCloser{bytes: []byte("121\n{\"type\": \"SUBSCRIBED\",\"subscribed\": {\"framework_id\": {\"value\":\"12220-3440-12532-2345\"},\"heartbeat_interval_seconds\":15.0}20\n{\"type\":\"HEARTBEAT\"}")},
	"insufficient frame":                       &MockFrameReaderCloser{bytes: []byte("124\n")}}
