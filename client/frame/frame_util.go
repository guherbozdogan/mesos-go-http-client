package frame

import "context"

const (
	CTRecordIO  FrameIOType = 1 << iota
	CTOtherType FrameIOType = 1 << iota
)

func NewFrameIO(t FrameIOType) FrameIO {
	switch t {
	case CTRecordIO:
		return &RecordIO{} //NewRecordIO(ctx)
	default:
		return nil //NewRecordIO(ctx)
	}
}

func DecorateFrameReadFunc(pre []FrameReadFunc, post []FrameReadFunc, f FrameReadFunc) FrameReadFunc {
	return func(ctx context.Context, fr Frame, n int64) context.Context {

		for i := 0; i < len(pre); i++ {
			if pre[i] != nil {
				ctx = pre[i](ctx, fr, n)
			}
		}
		if f != nil {
			ctx = f(ctx, fr, n)
		}
		for i := 0; i < len(post); i++ {
			if post[i] != nil {
				ctx = post[i](ctx, fr, n)
			}
		}
		return ctx
	}

}
