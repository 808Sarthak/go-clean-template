package cache

import (
	"context"
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
)

type Cache struct {
	cache *cache.Cache
}

func New(redisClient *redis.Client) *Cache {
	return &Cache{
		cache: cache.New(&cache.Options{
			Redis: redisClient,
		}),
	}
}

func (c *Cache) Get(ctx context.Context, key string, value interface{}) error {
	return c.cache.Get(ctx, key, value)
}

func (c *Cache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return c.cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: value,
		TTL:   ttl,
	})
}

func (c *Cache) Delete(ctx context.Context, key string) error {
	return c.cache.Delete(ctx, key)
}
