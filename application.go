package application

import (
	"context"
	"os"
	"os/signal"
	"runtime"
	"sync"
)

// waitGroupContextKey contextKey for wait group
const waitGroupContextKey = contextKey("APP:WAIT_GROUP")

// IApplication interface
type IApplication interface {

	// Background run the application in background, and wait for OS signal to terminate
	Background()

	// Run a runnable in a goroutine
	Start(runnable Runnable)

	// Stop the application
	Stop()

	// Run and wait all runnable stopped
	Wait()
}

// WaitGroupFromContext get a sync.WaitGroup from context, if available
func WaitGroupFromContext(ctx context.Context) *sync.WaitGroup {
	value := ctx.Value(waitGroupContextKey)
	if value == nil {
		return nil
	}

	switch wg := value.(type) {
	case *sync.WaitGroup:
		return wg

	default:
		return nil
	}
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

	app.ctx, app.cancel = context.WithCancel(
		context.WithValue(ctx, waitGroupContextKey, app.wg),
	)
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

func (a *application) Wait() {
	defer runtime.GC()

	done := make(chan bool, 1)

	go func() {
		a.wg.Wait()
		done <- true
	}()

	for {
		select {
		case <-done:
			break

		case <-a.ctx.Done():
			break
		}
	}
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
