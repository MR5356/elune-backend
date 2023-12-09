package notify

import (
	"github.com/MR5356/elune-backend/pkg/config"
	"github.com/MR5356/elune-backend/pkg/persistence"
	"github.com/MR5356/elune-backend/pkg/persistence/cache"
	"github.com/MR5356/elune-backend/pkg/persistence/database"
	"github.com/MR5356/elune-backend/pkg/utils/fileutil"
	"github.com/MR5356/elune-backend/pkg/utils/structutil"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

type Service struct {
	notifierManager           *NotifierManager
	notifierPluginPersistence *persistence.Persistence[*NotifierPlugin]
	config                    *config.Config
}

func NewService(database *database.Database, cache cache.Cache, cfg *config.Config) *Service {
	return &Service{
		notifierManager:           NewNotifierManager(),
		notifierPluginPersistence: persistence.New(database, cache, &NotifierPlugin{}),
		config:                    cfg,
	}
}

func (s *Service) AddNotifierPlugin(p *NotifierPlugin, filePath string) error {
	if err := s.notifierManager.RegisterPlugin(p.Name, filePath); err != nil {
		return err
	}
	defer os.RemoveAll(filePath)
	p.Version = fileutil.GetFileMd5(filePath)[:6]
	err := os.Rename(filePath, filepath.Join(s.config.Server.RuntimeDirectories[config.PluginDirectoryNotify], p.Name+p.Version+".so"))
	if err != nil {
		return err
	}
	return s.notifierPluginPersistence.Insert(p)
}

func (s *Service) Initialize() error {
	err := s.notifierPluginPersistence.DB.AutoMigrate(&NotifierPlugin{})
	if err != nil {
		return err
	}

	ps, err := s.notifierPluginPersistence.List(&NotifierPlugin{})
	if err != nil {
		return err
	}

	psMap := map[string]*NotifierPlugin{}
	for _, p := range ps {
		if v, ok := psMap[p.Name]; ok {
			if p.CreatedAt.After(v.CreatedAt) {
				psMap[p.Name] = p
			}
		} else {
			psMap[p.Name] = p
		}
	}

	for _, p := range psMap {
		logrus.Infof("load plugin: %+v", p)
		err := s.notifierManager.RegisterPlugin(p.Name, filepath.Join(s.config.Server.RuntimeDirectories[config.PluginDirectoryNotify], p.Name+p.Version+".so"))
		if err != nil {
			logrus.Errorf("plugin register error: %v", err)
			p.Status = err.Error()
		} else {
			p.Status = "success"
		}
		err = s.notifierPluginPersistence.Update(&NotifierPlugin{ID: p.ID}, structutil.Struct2Map(p))
		if err != nil {
			logrus.Errorf("update plugin status error: %v", err)
		}
	}
	return err
}
