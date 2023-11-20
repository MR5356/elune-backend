package navigation

import (
	"errors"
	"fmt"
	"github.com/MR5356/elune-backend/pkg/persistence"
	"github.com/MR5356/elune-backend/pkg/persistence/cache"
	"github.com/MR5356/elune-backend/pkg/persistence/database"
	"github.com/MR5356/elune-backend/pkg/utils/structutil"
	"github.com/sirupsen/logrus"
)

type Service struct {
	persistence *persistence.Persistence[*Navigation]
}

func NewService(database *database.Database, cache *cache.Cache) *Service {
	return &Service{
		persistence.New(database, cache, &Navigation{}),
	}
}

func (s *Service) ListNavigation() ([]*Navigation, error) {
	return s.persistence.List(&Navigation{})
}

func (s *Service) AddNavigation(navigation *Navigation) error {
	navigation.ID = 0
	if len(navigation.Title) == 0 {
		return errors.New("菜单名称不可为空")
	}

	// 检查是否已存在同名的菜单
	//res, err := s.persistence.Detail(&Navigation{Title: navigation.Title, Parent: navigation.Parent})
	//if err == nil && res.Parent == navigation.Parent {
	//	return errors.New("菜单名称重复")
	//}

	navigation.Unique = fmt.Sprintf("%s-%d", navigation.Title, navigation.Parent)
	return s.persistence.Insert(navigation)
}

func (s *Service) UpdateNavigation(navigation *Navigation) error {
	if len(navigation.Title) == 0 {
		return errors.New("navigation title cannot be empty")
	}

	//res, err := s.persistence.Detail(&Navigation{Title: navigation.Title, Parent: navigation.Parent})
	//if err == nil && res.Parent == navigation.Parent && (res.Title != navigation.Title || res.Parent != navigation.Parent) {
	//	return errors.New("菜单名称重复")
	//}

	navigation.Unique = fmt.Sprintf("%s-%d", navigation.Title, navigation.Parent)
	return s.persistence.Update(&Navigation{ID: navigation.ID}, structutil.Struct2Map(navigation))
}

func (s *Service) DeleteNavigation(id uint) error {
	return s.persistence.Delete(&Navigation{ID: id})
}

func (s *Service) Initialize() error {
	err := s.persistence.DB.AutoMigrate(&Navigation{})
	if err != nil {
		return err
	}

	_ = s.persistence.Insert(&Navigation{
		ID:     1,
		Title:  "默认分类",
		Parent: 0,
		Unique: "默认分类-0",
	})

	err = s.persistence.DB.Migrator().DropColumn(&Navigation{}, "delete_at")
	if err != nil {
		logrus.Errorf("drop column err: %+v", err)
	}

	return nil
}
