package echo

import (
	"context"

	"github.com/ServiceWeaver/weaver"
)

type Echoer interface {
	EchoResponse(context.Context, string) error
}

type impl struct {
	weaver.Implements[Echoer]
}

func (e *impl) EchoResponse(ctx context.Context, query string) error {
	logger := e.Logger(ctx)
	logger.Debug("vivian: ECHO", "results", query)

	return nil
}
