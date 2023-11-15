package utils

import (
	"context"
	"time"

	"github.com/ServiceWeaver/weaver"
)

type T interface {
	Time(context.Context) ([]byte, error)
	LoggerSocket(context.Context) error
}

type impl struct {
	weaver.Implements[T]
}

func (i *impl) LoggerSocket(ctx context.Context) error {
	i.Logger(ctx).Debug("vivian: socket: connected")
	return nil
}

func (i *impl) Time(ctx context.Context) ([]byte, error) {
	return []byte(time.Now().Format(time.RFC3339)), nil
}
