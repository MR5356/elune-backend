package cache

import "github.com/MR5356/elune-backend/pkg/config"

type Cache interface {
	TryLock(key string) error
	Unlock(key string) error
}

func New(cfg *config.Config) (cache Cache, err error) {
	return NewMemoryCache(), nil
}
