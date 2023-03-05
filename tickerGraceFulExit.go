package tickerGraceFulExit

import (
	"errors"
	"sync"
	"time"
)

type TickerGraceFulExit struct {
	paused bool
	wg     sync.WaitGroup
	ch     chan int
	rwLock *sync.RWMutex
	opts   options //可选参数变量
}

type options struct {
	MaxTask int
}

func SetMaxTask(a int) Option {
	return func(o *options) {
		o.MaxTask = a
	}
}

type Option func(*options)

func TickerManage(opt ...Option) *TickerGraceFulExit {
	var opts options
	for _, o := range opt {
		o(&opts)
	}

	if opts.MaxTask == 0 {
		opts.MaxTask = 1000
	}

	return &TickerGraceFulExit{
		opts:   opts,
		ch:     make(chan int, opts.MaxTask),
		rwLock: new(sync.RWMutex),
	}
}

func (j *TickerGraceFulExit) IsStop() (stop bool) {
	j.rwLock.RLock()
	defer j.rwLock.RUnlock()

	if j.paused {
		stop = true
	}
	return
}

func (j *TickerGraceFulExit) Add() (err error) {
	j.rwLock.Lock()
	defer j.rwLock.Unlock()

	if j.paused {
		err = errors.New("Work is stopping")
		return
	}
	j.ch <- 1
	return
}

func (j *TickerGraceFulExit) Done() {
	<-j.ch
}

func (j *TickerGraceFulExit) WaitGraceFulExit() {
	j.rwLock.Lock()
	defer j.rwLock.Unlock()

	j.paused = true
	close(j.ch)
	for {
		if len(j.ch) == 0 {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	return
}
