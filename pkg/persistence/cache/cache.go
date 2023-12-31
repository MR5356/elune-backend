package cache

import "github.com/MR5356/elune-backend/pkg/config"

type Cache interface {
	TryLock(key string) error
	Unlock(key string) error
	Subscribe(topic string, fn interface{}) error
	Publish(topic string, data interface{})
}

func New(cfg *config.Config) (cache Cache, err error) {
	if cfg.Persistence.Cache.Driver == "redis" {
		return NewRedisCache(cfg.Persistence.Cache.DSN)
	} else {
		return NewMemoryCache(), nil
	}
}
