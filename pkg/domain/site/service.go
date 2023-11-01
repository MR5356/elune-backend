package site

import (
	"errors"
	"github.com/MR5356/elune-backend/pkg/persistence"
	"github.com/MR5356/elune-backend/pkg/persistence/cache"
	"github.com/MR5356/elune-backend/pkg/persistence/database"
	"github.com/MR5356/elune-backend/pkg/utils/structutil"
)

type Service struct {
	persistence *persistence.Persistence[*SiteConfig]
}

func NewService(database *database.Database, cache *cache.Cache) *Service {
	return &Service{
		persistence.New(database, cache, &SiteConfig{}),
	}
}

func (s *Service) GetKey(key string) (string, error) {
	entity, err := s.persistence.Detail(&SiteConfig{Key: key})
	if err != nil {
		return "", err
	} else {
		return entity.Value, nil
	}
}

func (s *Service) SetKey(key string, value string) error {
	if len(key) == 0 {
		return errors.New("key cannot be empty")
	}
	entity, err := s.persistence.Detail(&SiteConfig{Key: key})
	if err != nil {
		err := s.persistence.Insert(&SiteConfig{Key: key, Value: value})
		return err
	} else {
		err := s.persistence.Update(entity, structutil.Struct2Map(&SiteConfig{Key: key, Value: value}))
		return err
	}
}

func (s *Service) Initialize() error {
	err := s.persistence.DB.AutoMigrate(&SiteConfig{})
	if err != nil {
		return err
	}

	defaultSiteConfigs := []*SiteConfig{
		{Key: "title", Value: "Elune"},
		{Key: "description", Value: "Website of Elune"},
		{Key: "logo", Value: "/logo.svg"},
		{Key: "favicon", Value: "/favicon.ico"},
		{Key: "copyright", Value: "© 2022 Elune"},
		{Key: "beian", Value: "冀公网安备 13112202000250号"},
		{Key: "beianMiit", Value: "冀ICP备20003324号-3"},
	}

	for _, siteConfig := range defaultSiteConfigs {
		_ = s.persistence.Insert(siteConfig)
	}
	return nil
}
