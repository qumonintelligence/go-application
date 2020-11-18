package application_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/qumonintelligence/go-application/v2"
)

func sleepRunnable5s(ctx context.Context) {
	fmt.Println("RUN")
	time.Sleep(5 * time.Second)
}

func TestWait(t *testing.T) {
	app := application.NewApplication(context.Background())

	app.Start(sleepRunnable5s)
	app.Start(sleepRunnable5s)

	app.Wait()

	fmt.Println("WAIT OK")
}
