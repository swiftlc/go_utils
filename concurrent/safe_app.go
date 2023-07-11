package concurrent

import (
	"context"
	"sync"
	"sync/atomic"

	logger "github.com/sirupsen/logrus"
)

type SafeApp struct {
	*logger.Entry
	Name       string //app name
	ctx        context.Context
	cancelFunc context.CancelFunc
	wg         sync.WaitGroup
	_init      int32
}

func (sa *SafeApp) Init(name string) {
	if atomic.CompareAndSwapInt32(&sa._init, 0, 1) {
		sa.Name = name
		sa.ctx, sa.cancelFunc = context.WithCancel(context.Background())
		sa.Entry = logger.WithField("app_name", name).WithContext(sa.ctx)
		sa.WithContext(sa.ctx).Infof("init")
	}
}

func (sa *SafeApp) SafeRun(taskName string, concurrent int, f func(ctx context.Context)) {
	sa.Infof("run task:%s,concurrent:%d", taskName, concurrent)
	SafeGo(sa.ctx, &sa.wg, concurrent, f)
}

func (sa *SafeApp) SafeHandleGroup(taskName string, concurrent int, input <-chan *EventWithCtx, hanlder func(context.Context, interface{})) {
	sa.Infof("run task:%s,concurrent:%d", taskName, concurrent)
	SafeChanConsume(sa.ctx, &sa.wg, concurrent, input, hanlder)
}

func (sa *SafeApp) Uninit() {
	if atomic.CompareAndSwapInt32(&sa._init, 1, 2) {
		sa.cancelFunc()
		sa.wg.Wait()
		sa.Infof("uninit")
	}
}
