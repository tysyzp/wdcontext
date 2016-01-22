//
// User: zhangpeng
// Date: 01/22/2016
// Time: 13:11
//

package wdcontext

import (
	"golang.org/x/net/context"
	"time"
)

type key int

var contextKey key = 0

type WDContext struct {
	context.Context
	feed  chan struct{}
	timer *time.Timer
}

func WithWatchDogTimeout(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc, chan struct{}) {
	cancelCtx, cancelFun := context.WithCancel(parent)
	c := &WDContext{}
	cancelCtx = context.WithValue(cancelCtx, contextKey, c)
	c.Context = cancelCtx
	c.feed = make(chan struct{})

	cancel := func() {
		c.timer.Stop()
		cancelFun()
	}
	c.timer = time.AfterFunc(timeout, func() {
		select {
		case <-c.feed:
		default:
			cancel()
		}
	})

	return cancelCtx, cancel, c.feed
}
