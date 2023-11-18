package echo

import (
	"context"

	"github.com/ServiceWeaver/weaver"
)

type T interface {
	EchoResponse(context.Context, string) error
}

type impl struct {
	weaver.Implements[T]
}

func (e *impl) EchoResponse(ctx context.Context, query string) error {
	logger := e.Logger(ctx)
	logger.Debug("vivian:", "echo", query)

	return nil
}
