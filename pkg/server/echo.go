package server

import (
	"context"

	"github.com/ServiceWeaver/weaver"
)

type Echo interface {
	EchoResponse(context.Context, string) error
}

type echo struct {
	weaver.Implements[Echo]
}

func (e *echo) EchoResponse(ctx context.Context, query string) error {
	logger := e.Logger(ctx)
	logger.Debug("vivian: ECHO", "results", query)

	return nil
}
