package authentication

import (
	"github.com/MR5356/elune-backend/pkg/persistence"
	"github.com/MR5356/elune-backend/pkg/persistence/cache"
	"github.com/MR5356/elune-backend/pkg/persistence/database"
)

type Service struct {
	persistence *persistence.Persistence[*User]
}

func NewService(database *database.Database, cache cache.Cache) *Service {
	return &Service{
		persistence.New(database, cache, &User{}),
	}
}

func (s *Service) Initialize() error {
	err := s.persistence.DB.AutoMigrate(&User{})
	if err != nil {
		return err
	}
	_ = s.persistence.Insert(&User{
		ID:       1,
		Username: "admin",
		Password: "admin",
		Email:    "admin@example.com",
	})
	_ = s.persistence.Insert(&User{
		ID:       2,
		Username: "guest",
		Password: "guest",
		Email:    "guest@example.com",
	})
	_ = s.persistence.Insert(&User{
		ID:       3,
		Username: "devops",
		Password: "devops",
		Email:    "devops@example.com",
	})
	return nil
}
