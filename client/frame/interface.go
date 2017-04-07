package frame


import "io"
import  "context"


type FrameIOType int
type Frame  [] byte

type FrameReadFunc func(context.Context, Frame, int64 ) 
type FrameReadFuncDecorator  func(FrameReadFunc ) FrameReadFunc 
type ErrorFunc func(context.Context,interface{} )

    
type FrameWritten func (context.Context, int64 ) ()

type FrameIO interface {
   
    Read(context.Context,  io.ReadCloser, FrameReadFunc ,ErrorFunc ) 
    Write(context.Context,  io.WriteCloser, FrameWritten ,ErrorFunc) 
        
 }
