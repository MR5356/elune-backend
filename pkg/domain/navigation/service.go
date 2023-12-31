package navigation

import (
	"errors"
	"fmt"
	"github.com/MR5356/elune-backend/pkg/persistence"
	"github.com/MR5356/elune-backend/pkg/persistence/cache"
	"github.com/MR5356/elune-backend/pkg/persistence/database"
	"github.com/MR5356/elune-backend/pkg/utils/imgutil"
	"github.com/MR5356/elune-backend/pkg/utils/structutil"
	"github.com/sirupsen/logrus"
)

type Service struct {
	persistence *persistence.Persistence[*Navigation]
}

func NewService(database *database.Database, cache cache.Cache) *Service {
	return &Service{
		persistence.New(database, cache, &Navigation{}),
	}
}

func (s *Service) ListNavigation() ([]*Navigation, error) {
	res, err := s.persistence.List(&Navigation{})
	if err != nil {
		return nil, err
	}
	for _, item := range res {
		if len(item.Logo) == 0 || len(item.LogoData) > 0 {
			continue
		}
		logrus.Debugf("imgutil.ImgLinkToBase64: %s", item.Logo)
		data, err := imgutil.ImgLinkToBase64(item.Logo)
		if err == nil {
			item.LogoData = data
			s.UpdateNavigation(item)
		} else {
			item.LogoData = item.Logo
			s.UpdateNavigation(item)
			logrus.Errorf("imgutil.ImgLinkToBase64 err: %v", err)
		}
	}
	return res, nil
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

	if len(navigation.Logo) > 0 {
		data, err := imgutil.ImgLinkToBase64(navigation.Logo)
		if err == nil {
			navigation.LogoData = data
		} else {
			navigation.LogoData = navigation.Logo
			logrus.Errorf("imgutil.ImgLinkToBase64 err: %v", err)
		}
	}

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

	if len(navigation.Logo) > 0 {
		data, err := imgutil.ImgLinkToBase64(navigation.Logo)
		if err == nil {
			navigation.LogoData = data
		} else {
			navigation.LogoData = navigation.Logo
			logrus.Errorf("imgutil.ImgLinkToBase64 err: %v", err)
		}
	}

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
