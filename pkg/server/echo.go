package server

import (
	"context"

	"github.com/ServiceWeaver/weaver"
)

type EchoInterface interface {
	Echo(context.Context, string) error
}

type echo struct {
	weaver.Implements[EchoInterface]
}

func (e *echo) Echo(ctx context.Context, query string) error {
	logger := e.Logger(ctx)
	logger.Debug("echo", "results", query)

	return nil
}
