package frame

import (
	"bufio"
	//	"bytes"
	"context"
	//binary "encoding/binary"
	"errors"
	"io"
	"sync"
	"sync/atomic"

	//	"bytes"
	//	"fmt"
	"strconv"
	//zap logger to be added
	//byte allocation buffers to be limited
)

//const BufferLen = 2048 //to make
var (
	ErrorChannelClosedBeforeReceivingFrame = errors.New("Channel closed before receiving frame")
	ErrorFormatErr                         = errors.New("Format Error")

	ErrorParentContextBeenCancelled = errors.New("Parent context been cancelled")
	ErrorGeneralIOErr               = errors.New("General IO Error")
	ErrorEOF                        = errors.New("EOF Received")
	ErrorInternal                   = errors.New("Module Internal Error")
)

type RecordIO struct {
	ReadLine func(r *bufio.Reader, rc io.ReadCloser) (Bytes, error)
}

func NewRecordIO() *RecordIO {
	return &RecordIO{ReadLine: Readln}
}

type Bytes []byte

//helper functionto read a line with bufio
func Readln(r *bufio.Reader, rc io.ReadCloser) (Bytes, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		if line != nil {
			ln = append(ln, line...)
		} else {
			return []byte(""), ErrorInternal
		}

	}
	return ln, err
}

//helper function to check whether a flag is set or not in atomic sense
func isAtomicValSet(v *uint32) bool {
	opsFinal := atomic.LoadUint32(v)
	if opsFinal > 0 {
		return true
	} else {
		return false
	}
}

//interface function for Read
func (rc *RecordIO) Read(ctx context.Context, reader io.ReadCloser, f FrameReadFunc, erf ErrorFunc) {

	//flag to be set when context is done
	var icd uint32 = 0

	errc := make(chan error) //channel for sending error out

	//initialize the bufio
	r := bufio.NewReader(reader)
	//r := reader

	//to be able to wait for the graceful exit of below function
	var wg sync.WaitGroup
	wg.Add(1)
	//note: check how channels receive arrays, if they receive as copies, is it possible to send a channel an address of bytes from a local variable? if so, use the pointer version instead of copy by value

	// go function
	//go func (r bufio.Reader,  rfc chan bytes[], erc chan error, ops * uint32,)   {
	go func() {
		//graceful exit
		defer wg.Done()

		//check if context is already cancelled/done, if so, return
		fcc := func() bool {
			opsFinal := atomic.LoadUint32(&icd)
			if opsFinal > 0 {
				return true
			} else {
				return false
			}
		}

		for !fcc() {
			b, err := rc.ReadLine(r, reader)

			//check if context is already cancelled/done, if so, return
			if fcc() {
				return
			}
			if err != nil {
				errc <- err
				//log error !
				return
			}

			//b = append(b, byte(' '))
			if b == nil || len(b) == 0 {
				continue
			}
			s1 := string(b[:])

			//buf := bytes.NewBuffer(b) // b is []byte
			si, errCn := strconv.Atoi(s1)
			//add err handling here!

			if errCn != nil {
				errc <- ErrorFormatErr
				//log error !
				return
			}

			s := int64(si)

			//add some check constraints here to the read buffer size!!!!!!!!!!!  might be illegitimate!

			//check for panics of memory alloc later

			var trb []byte
			trb = make([]byte, s, s)

			trbr := trb[:]
			l := int64(0) //length of read bytes
			for !fcc() {
				ni, err := reader.Read(trbr)
				n := int64(ni)
				l += n
				//check if context is already cancelled/done, if so, return
				if fcc() {
					return
				}

				//read entire data
				// add check if received n is smaller and or n is larger but eof returned
				if err == io.EOF {
					if l < s {
						errc <- ErrorChannelClosedBeforeReceivingFrame
						//log error!
						return
					} else {
						trbTmp := make([]byte, s, s)
						copy(trbTmp, trb)
						f(ctx, Frame(trbTmp), s)

						errc <- ErrorEOF
						return
					}

				}
				if l == s { //if all entire frame read
					trbTmp := make([]byte, s, s)
					copy(trbTmp, trb)
					f(ctx, Frame(trbTmp), s)

					if err != nil {
						errc <- err
						return

					}
					break

				} else if err != nil {
					errc <- err
					return

				} else { //if entire frame is not read yet, continue reading in next loop

					trbr = trb[l:]
				}
			}
		}
	}()

	done := ctx.Done()

	select {
	case <-done:

		atomic.AddUint32(&icd, 1)
		//wait graceful close
		wg.Wait()
		erf(ctx, ctx.Err())
		close(errc)
		return

	case e := <-errc:

		wg.Wait()
		erf(ctx, e)
		close(errc)

		return
	}
}

func (*RecordIO) Write(c context.Context, writer io.WriteCloser, f FrameWritten, erf ErrorFunc) {

	//not need to implement right now:)
	//to be added later
	return

}
