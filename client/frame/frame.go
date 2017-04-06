package frame




func NewFrameIOType(t FrameIOType) FrameIO {
    switch t {
    case CTRecordIO:
        return &RecordIO{} //NewRecordIO(ctx)
    default:
        return &RecordIO{}//NewRecordIO(ctx)
    }
}