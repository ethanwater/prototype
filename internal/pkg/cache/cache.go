package cache

import (
	"context"
	"sync"

	"github.com/ServiceWeaver/weaver"
)

type Cache interface {
	Get(context.Context, string) (string, error)
	Put(context.Context, string, string) error
}

type impl struct {
	weaver.Implements[Cache]
	weaver.WithRouter[route]
	mu   sync.Mutex
	data map[string]string
}

func (c *impl) Init(context.Context) error {
	c.data = make(map[string]string)
	return nil
}

type route struct{}

func (route) Get(ctx context.Context, key string) string        { return key }
func (route) Put(ctx context.Context, key, value string) string { return key }

func (c *impl) Get(ctx context.Context, key string) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Logger(ctx).Debug("vivian: [ok] cache tunneled")

	return c.data[key], nil
}

func (c *impl) Put(ctx context.Context, key, value string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
	c.Logger(ctx).Debug("vivian: [ok] cache successful")
	return nil
}

func (c *impl) Clear(ctx context.Context) error {
	c.data = make(map[string]string)
	c.Logger(ctx).Debug("vivian: [ok] cache cleared")
	return nil
}
