package application_test

import (
	"context"
	"testing"
	"time"

	"github.com/qumonintelligence/go-application/v2"
)

func sleep5s(ctx context.Context, data interface{}) {
	time.Sleep(time.Second)
}

func TestExecutor10(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	ex := application.NewExecutor(ctx, 3)

	for i := 0; i < 15; i++ {
		ex.Submit(ctx, sleep5s, nil)
	}

	time.Sleep(15 * time.Second)
	if ex.Count() > 0 {
		t.Fail()
	}
}
