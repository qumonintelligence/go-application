package application

import (
	"context"
	"time"

	"go.uber.org/atomic"
)

// ICallable callable
type ICallable func(ctx context.Context, data interface{})

// Task task
type task struct {
	callable ICallable // runner
	ctx      context.Context
	data     interface{} // serialized data
}

type IExecutor interface {
	// Submit a callable to be called later
	Submit(ctx context.Context, callable ICallable, data interface{})

	// Execute a function later using the goroutine pool defined
	ExecuteLater(ctx context.Context, callable ICallable, data interface{}, delay time.Duration)

	// Shutdown the executors
	Shutdown()

	// Count number of task submitted + executeLater
	Count() int64
}

type executor struct {
	ctx               context.Context
	cancel            context.CancelFunc
	count             int
	queue             chan *task
	submitCount       atomic.Int64
	executeLaterCount atomic.Int64
}

func (e *executor) Count() int64 {
	return e.submitCount.Load() + e.executeLaterCount.Load()
}

func (e *executor) Submit(ctx context.Context, callable ICallable, data interface{}) {

	e.submitCount.Inc()

	e.queue <- &task{
		callable: callable,
		ctx:      ctx,
		data:     data,
	}
}

func (e *executor) ExecuteLater(ctx context.Context, callable ICallable, data interface{}, delay time.Duration) {

	e.executeLaterCount.Inc()

	time.AfterFunc(delay, func() {
		e.executeLaterCount.Dec()
		e.Submit(ctx, callable, data)
	})
}

func NewExecutor(ctx context.Context, executorCount int) IExecutor {
	e := &executor{
		count: executorCount,
		queue: make(chan *task, 4096),
	}

	e.ctx, e.cancel = context.WithCancel(ctx)
	e.start()
	return e
}

// Shutdown the executor
func (e *executor) Shutdown() {
	e.cancel()
}

func (e *executor) start() {
	if e.count <= 0 {
		e.count = 1
	}

	for i := 0; i < e.count; i++ {
		go func(i int) {
			for {
				select {
				case t := <-e.queue:
					e.submitCount.Dec()
					t.callable(t.ctx, t.data)

				case <-e.ctx.Done():
					e.cancel()
					return
				}
			}
		}(i)
	}
}
