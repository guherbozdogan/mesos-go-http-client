package frame


import "io"

type FrameIOType int
type Frame * [] byte

type FrameRead func(Frame, int64 ) ()
type FrameWritten func (int64 ) ()

type FrameIO interface {
   
    Read(io.ReadCloser, FrameRead )  (interface{}, error) 
    Write(io.WriteCloser, FrameWritten )  (interface{}, error) 
        
 }
const (
    CTRecordIO FrameIOType = 1 << iota
  )