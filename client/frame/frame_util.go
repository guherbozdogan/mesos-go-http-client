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

		for i := len(pre); i >= 0; i++ {
			ctx = pre[i](ctx, fr, n)
		}
		ctx = f(ctx, fr, n)

		for i := len(post); i >= 0; i++ {
			ctx = post[i](ctx, fr, n)
		}
		return ctx
	}

}
