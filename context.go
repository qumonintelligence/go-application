package application

import (
	"context"

	"github.com/qumonintelligence/go-common/v2/lang"
)

// IContextKey interface
type IContextKey interface {
	ToString(ctx context.Context) string
}

// contextKey of context's value
type contextKey string

// NewContextKey create a new context key
func NewContextKey(name string) IContextKey {
	return contextKey(name)
}

// ToString try to get a context value as String
func (c contextKey) ToString(ctx context.Context) string {
	value := ctx.Value(c)
	if value == nil {
		return ""
	}

	return lang.StringOf(value)
}
