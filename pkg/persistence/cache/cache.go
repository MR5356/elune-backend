package cache

import "github.com/MR5356/elune-backend/pkg/config"

type Cache interface {
}

func New(cfg *config.Config) (cache *Cache, err error) {
	return
}
