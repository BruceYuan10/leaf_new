package timer

import (
	"github.com/BruceYuan10/leaf_new/conf"
	"github.com/BruceYuan10/leaf_new/log"
	"runtime"
	"sync"
	"time"
)

// one dispatcher per goroutine (goroutine not safe)
type Dispatcher struct {
	ChanTimer chan *Timer
	disMutex  *sync.Mutex
	closeFlag bool
}

func NewDispatcher(l int) *Dispatcher {
	disp := new(Dispatcher)
	disp.ChanTimer = make(chan *Timer, l)
	disp.disMutex = new(sync.Mutex)
	return disp
}

// Timer
type Timer struct {
	t  *time.Timer
	cb func()
}

func (t *Timer) Stop() {
	t.t.Stop()
	t.cb = nil
}

func (t *Timer) Cb() {
	defer func() {
		t.cb = nil
		if r := recover(); r != nil {
			if conf.LenStackBuf > 0 {
				buf := make([]byte, conf.LenStackBuf)
				l := runtime.Stack(buf, false)
				log.Error("%v: %s", r, buf[:l])
			} else {
				log.Error("%v", r)
			}
		}
	}()

	if t.cb != nil {
		t.cb()
	}
}

func (disp *Dispatcher) AfterFunc(d time.Duration, cb func()) *Timer {
	t := new(Timer)
	t.cb = cb
	t.t = time.AfterFunc(d, func() {
		disp.ChanTimer <- t
	})
	return t
}

// Cron
type Cron struct {
	t *Timer
}

func (c *Cron) Stop() {
	if c.t != nil {
		c.t.Stop()
	}
}

func (disp *Dispatcher) CronFunc(cronExpr *CronExpr, _cb func()) *Cron {
	c := new(Cron)

	now := time.Now()
	nextTime := cronExpr.Next(now)
	if nextTime.IsZero() {
		return c
	}

	// callback
	var cb func()
	cb = func() {
		defer _cb()

		now := time.Now()
		nextTime := cronExpr.Next(now)
		if nextTime.IsZero() {
			return
		}
		c.t = disp.AfterFunc(nextTime.Sub(now), cb)
	}

	c.t = disp.AfterFunc(nextTime.Sub(now), cb)
	return c
}

func (disp *Dispatcher) OnClose() {
	disp.disMutex.Lock()
	defer disp.disMutex.Unlock()

	if !disp.closeFlag {
		close(disp.ChanTimer)
		disp.closeFlag = true
	}
}
