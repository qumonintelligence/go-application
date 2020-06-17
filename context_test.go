package application

import (
	"context"
	"testing"
	"time"
)

func TestContextKey(t *testing.T) {
	key := NewContextKey("hello")
	ctx := context.WithValue(context.Background(), key, "world")

	val := ctx.Value(key).(string)
	if val != "world" {
		t.FailNow()
	}
}

func TestContextKeyToString(t *testing.T) {
	key := NewContextKey("hello")
	ctx := context.WithValue(context.Background(), key, "world")

	ctx2, cancel := context.WithTimeout(ctx, time.Hour)
	defer cancel()
	val := key.ToString(ctx2)
	if val != "world" {
		t.FailNow()
	}
}
