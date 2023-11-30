package cache

import "github.com/MR5356/elune-backend/pkg/config"

type Cache interface {
	TryLock(key string) error
	Unlock(key string) error
	Subscribe(topic string, fn interface{}) error
	Publish(topic string, data interface{})
}

func New(cfg *config.Config) (cache Cache, err error) {
	return NewMemoryCache(), nil
}
