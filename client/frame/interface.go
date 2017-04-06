package frame


import "io"
import  "context"

type FrameIOType int
type Frame * [] byte

type FrameRead func(Frame, int64 ) ()
type FrameWritten func (int64 ) ()

type FrameIO interface {
   
    Read(context.Context,  io.ReadCloser, FrameRead )  (interface{}, error) 
    Write(context.Context,  io.WriteCloser, FrameWritten )  (interface{}, error) 
        
 }
const (
    CTRecordIO FrameIOType = 1 << iota
  )