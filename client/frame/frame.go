package frame

import "context"
const (
    CTRecordIO FrameIOType = 1 << iota
  )


func NewFrameIO(t FrameIOType) FrameIO {
    switch t {
    case CTRecordIO:
        return &RecordIO{} //NewRecordIO(ctx)
    default:
        return &RecordIO{}//NewRecordIO(ctx)
    }
}

func DecorateFrameReadFunc(b [] FrameReadFunc,e []  FrameReadFunc, f FrameReadFunc) FrameReadFunc  {
    return func(ctx context.Context, fr Frame,n int64 )  context.Context {
    
        for i:= len(b); i>=0; i++ {
                ctx=b[i](ctx,fr, n );    
        }
        ctx=f(ctx,fr, n );
        
        for i:= len(e); i>=0; i++ {
                ctx=e[i](ctx,fr, n );    
        }
        return ctx;
    }
    
}