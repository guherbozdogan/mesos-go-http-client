package frame

import  "context"


func NewFrameIOType(t FrameIOType, c context.Context) FrameIO {
    switch t {
    case CTRecordIO:
        return &RecordIO{ctx:c} //NewRecordIO(ctx)
    default:
        return &RecordIO{ctx:c}//NewRecordIO(ctx)
    }
}