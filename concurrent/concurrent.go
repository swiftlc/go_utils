package concurrent

import (
	"context"
	"sync"

	"go.opencensus.io/trace"
)

func SafeGo(ctx context.Context, w *sync.WaitGroup, concurrent int,
	f func(ctx context.Context)) {
	if concurrent < 1 {
		concurrent = 1
	}
	if w != nil {
		w.Add(concurrent)
	}
	for i := 0; i < concurrent; i++ {
		go func(index int) {
			defer func() {
				if e := recover(); e != nil {
				}
				if w != nil {
					w.Done()
				}
			}()
			f(ctx)
		}(i)
	}
}

func CopyCtx(ctxFrom context.Context, ctxTo context.Context) context.Context {
	//trace
	ctxTo = trace.NewContext(ctxTo, trace.FromContext(ctxFrom))

	return ctxTo
}

func SafeChanConsume(ctx context.Context, w *sync.WaitGroup, concurrent int, ch <-chan *EventWithCtx, f func(context.Context, interface{})) {
	SafeGo(ctx, w, concurrent, func(ctx context.Context) {
		for v := range ch {
			newCtx := CopyCtx(v.Ctx, ctx)
			f(newCtx, v.Event)
		}
	})
}
