package cache

import (
	"errors"
	"github.com/asaskevich/EventBus"
	"sync"
)

type MemoryCache struct {
	mutexLockMap sync.Map
	evbus        EventBus.Bus
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{mutexLockMap: sync.Map{}}
}

func (c *MemoryCache) TryLock(key string) error {
	if _, locked := c.mutexLockMap.LoadOrStore(key, true); locked {
		return errors.New("already locked")
	}
	return nil
}

func (c *MemoryCache) Unlock(key string) error {
	c.mutexLockMap.Delete(key)
	return nil
}
