package blog

import (
	"github.com/MR5356/elune-backend/pkg/persistence"
	"github.com/MR5356/elune-backend/pkg/persistence/cache"
	"github.com/MR5356/elune-backend/pkg/persistence/database"
)

type Service struct {
	persistence *persistence.Persistence[*Blog]
}

func NewService(database *database.Database, cache *cache.Cache) *Service {
	return &Service{
		persistence.New(database, cache, &Blog{}),
	}
}

func (s *Service) Initialize() error {
	return s.persistence.DB.AutoMigrate(&Blog{})
}
