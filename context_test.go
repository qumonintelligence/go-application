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

type simple struct {
	val int
}

func TestContextWithPointer(t *testing.T) {
	key := NewContextKey("hello")
	ptr := &simple{val: 1}

	ctx := context.WithValue(context.Background(), key, ptr)
	ptr.val = 2

	val := ctx.Value(key)
	pval := val.(*simple)
	if pval.val != 2 {
		t.FailNow()
	}
	pval.val = 10

	if ptr.val != 10 {
		t.FailNow()
	}

}
