package frame

type FrameIOType int
type Frame bytes []

type FrameRead func(bytes[], int )
type FrameWritten func (int )

type FrameIO interface {
   
    Read(io.ReadCloser, FrameRead )  (interface{}, error) 
    Write(io.WriteCloser, FrameWitten )  (interface{}, error) 
        
 }

const (
    RecordIO FrameIOType = 1 << iota
    
)

func NewFrameIOType(t FrameIOType) FrameIO {
    switch t {
    case RecordIO:
        return newRecordIO()
    default:
        return newRecordIO()
    }
}