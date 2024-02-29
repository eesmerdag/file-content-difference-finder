package cache

import (
	"github.com/patrickmn/go-cache"
)

const ReadinessKey = "ready"

type ICache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{})
	Delete(key string)
	Clear()
}

type imCache struct {
	cache *cache.Cache
}

func InitCache() ICache {
	cacheInstance := newCache()
	cacheInstance.Set(ReadinessKey, 0)

	return cacheInstance
}

func newCache() ICache {
	return &imCache{
		cache: cache.New(cache.NoExpiration, cache.NoExpiration),
	}
}

func (c *imCache) Get(key string) (interface{}, bool) {
	return c.cache.Get(key)
}

func (c *imCache) Set(key string, value interface{}) {
	c.cache.Set(key, value, cache.DefaultExpiration)
}

func (c *imCache) Delete(key string) {
	c.cache.Delete(key)
}

func (c *imCache) Clear() {
	c.cache.Flush()
}
