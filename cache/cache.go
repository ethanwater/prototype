package cache

import (
	"context"
	"sync"

	"github.com/ServiceWeaver/weaver"
)

type DataCache interface {
	Put(ctx context.Context, key string, value []string) error
	Get(ctx context.Context, key string) ([]string, error)
}

type Cache struct {
	mu   sync.Mutex
	data map[string][]string
	weaver.Implements[DataCache]
}

func (c *Cache) Put(ctx context.Context, key string, value []string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
	return nil
}

func (c *Cache) Get(ctx context.Context, key string) ([]string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.data[key], nil
}
