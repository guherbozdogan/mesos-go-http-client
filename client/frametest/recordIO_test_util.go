package frametest

import (
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
	if cb.err != nil {
		return 0, cb.err
	} else if cb.bytesRead {
		return 0, nil
	} else if cb.bytes != nil {
		copy(p, cb.bytes)
		cb.bytesRead = true
		return len(p), nil
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
			time.Sleep(cb.period * time.Second)
			timeout <- true
		}()

		select {
		case <-timeout:
			return cb.ReadSub(p)
		}
	}
	//	return 0, nil
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
