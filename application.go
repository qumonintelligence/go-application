package application

import (
	"context"
	"os"
	"os/signal"
	"runtime"
	"sync"
)

// IApplication interface
type IApplication interface {

	// Background run the application in background, and wait for OS signal to terminate
	Background()

	// Run a runnable in a goroutine
	Start(runnable Runnable)

	// Stop the application
	Stop()
}

type application struct {
	wg     *sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

// NewApplication create a new application that can run in background
func NewApplication(ctx context.Context) IApplication {

	app := &application{
		wg: &sync.WaitGroup{},
	}

	if ctx == nil {
		ctx = context.Background()
	}

	app.ctx, app.cancel = context.WithCancel(ctx)
	return app
}

func (a *application) Start(runnable Runnable) {
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		runnable(a.ctx)
	}()
}

func (a *application) Stop() {
	a.cancel()
}

func (a *application) Background() {
	// force GC so that everything is cleared
	defer runtime.GC()

	a.wg.Add(1)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, os.Kill)

	go func() {
		defer a.wg.Done()
		select {
		case <-signalCh:
			a.cancel()
			return

		case <-a.ctx.Done():
			return
		}
	}()

	a.wg.Wait()
}
