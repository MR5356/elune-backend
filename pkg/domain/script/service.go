package script

import (
	"errors"
	"github.com/MR5356/elune-backend/pkg/persistence"
	"github.com/MR5356/elune-backend/pkg/persistence/cache"
	"github.com/MR5356/elune-backend/pkg/persistence/database"
	"github.com/MR5356/elune-backend/pkg/utils/structutil"
)

type Service struct {
	scriptPersistence     *persistence.Persistence[*Script]
	scriptTypePersistence *persistence.Persistence[*Type]
}

func NewService(database *database.Database, cache cache.Cache) *Service {
	return &Service{
		scriptPersistence:     persistence.New(database, cache, &Script{}),
		scriptTypePersistence: persistence.New(database, cache, &Type{}),
	}
}

func (s *Service) AddScript(script *Script) error {
	script.ID = 0
	if len(script.Title) == 0 {
		return errors.New("脚本名称不可为空")
	}
	if len(script.Content) == 0 {
		return errors.New("脚本内容不可为空")
	}
	script.TypeId = 1
	return s.scriptPersistence.Insert(script)
}

func (s *Service) ListScript() ([]*Script, error) {
	res := make([]*Script, 0)
	err := s.scriptPersistence.DB.Order("elune_script.created_at").Joins("Type").Find(&res).Error
	return res, err
}

func (s *Service) DeleteScript(id uint) error {
	return s.scriptPersistence.Delete(&Script{ID: id})
}

func (s *Service) UpdateScript(script *Script) error {
	if len(script.Title) == 0 {
		return errors.New("脚本名称不可为空")
	}
	if len(script.Content) == 0 {
		return errors.New("脚本内容不可为空")
	}
	st := &Script{
		ID:      script.ID,
		Title:   script.Title,
		Desc:    script.Desc,
		Content: script.Content,
		TypeId:  1,
	}
	return s.scriptPersistence.Update(&Script{ID: script.ID}, structutil.Struct2Map(st))
}

func (s *Service) Initialize() error {
	err := s.scriptPersistence.DB.AutoMigrate(&Type{}, &Script{})
	if err != nil {
		return err
	}

	_ = s.scriptTypePersistence.Insert(&Type{
		ID:    1,
		Title: "shell",
	})
	return nil
}
