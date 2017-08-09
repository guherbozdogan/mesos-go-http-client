package frame

import "io"
import  "context"


type FrameIOType int
type Frame  [] byte

type FrameReadFunc func(context.Context, Frame, int64 )  context.Context


    
    

                            
type ErrorFunc func(context.Context,interface{} )  context.Context

    
type FrameWritten func (context.Context, int64 ) context.Context

type FrameIO interface {
   
    Read(context.Context,  io.ReadCloser, FrameReadFunc ,ErrorFunc ) 
    Write(context.Context,  io.WriteCloser, FrameWritten ,ErrorFunc) 
        
 }
