package application

import (
	"context"
	"sync"
)

type IThread interface {
	Start()
}

type thread struct {
	wg       *sync.WaitGroup
	ctx      context.Context
	runnable Runnable
}

func (t *thread) Start() {
	t.wg.Add(1)
	go func() {
		defer t.wg.Done()
		t.runnable(t.ctx)
	}()
}
