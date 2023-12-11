package notify

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
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

const (
	topicAddNotifierPlugin = "notify.plugin.add"
	topicDelNotifierPlugin = "notify.plugin.del"
)

type Service struct {
	notifierManager           *NotifierManager
	notifierPluginPersistence *persistence.Persistence[*NotifierPlugin]
	config                    *config.Config
	cache                     cache.Cache
}

func NewService(database *database.Database, cache cache.Cache, cfg *config.Config) *Service {
	return &Service{
		notifierManager:           NewNotifierManager(),
		notifierPluginPersistence: persistence.New(database, cache, &NotifierPlugin{}),
		config:                    cfg,
		cache:                     cache,
	}
}

func (s *Service) RemoveNotifierPlugin(id uint) error {
	plugin, err := s.notifierPluginPersistence.Detail(&NotifierPlugin{ID: id})
	if err != nil {
		return errors.New("插件不存在")
	}
	s.cache.Publish(topicDelNotifierPlugin, s.getPluginFilePath(plugin))

	err = s.notifierPluginPersistence.Delete(&NotifierPlugin{ID: id})
	if err != nil {
		return errors.New("删除插件失败")
	}
	return nil
}

func (s *Service) getPluginFilePath(plugin *NotifierPlugin) string {
	return filepath.Join(s.config.Server.RuntimeDirectories[config.PluginDirectoryNotify], plugin.Name+"-"+plugin.Version+".so")
}

func (s *Service) UploadNotifierPlugin(filePath string) error {
	symbol, err := s.notifierManager.GetSymbol(filePath)
	if err != nil {
		return ErrNotifierOpenError
	}
	notifierPlugin, ok := s.notifierManager.Verify(symbol)
	if !ok {
		return ErrNotifierVerifyError
	}
	content, err := fileutil.ReadFromFile(filePath)
	if err != nil {
		return errors.New("读取插件文件失败")
	}

	notifierPlugin.Version = fmt.Sprintf("%x", md5.Sum(content))[:6]
	notifierPlugin.Installed = true

	filename := s.getPluginFilePath(notifierPlugin)
	pi := &PluginInfo{
		Name:     notifierPlugin.Name,
		Filename: filename,
		Content:  content,
	}

	piStr, err := json.Marshal(pi)
	if err != nil {
		return err
	}

	s.cache.Publish(topicAddNotifierPlugin, piStr)

	notifierPlugin.Status = "success"

	return s.notifierPluginPersistence.Insert(notifierPlugin)
}

type PluginInfo struct {
	Name     string `json:"name"`
	Filename string `json:"filename"`
	Content  []byte `json:"content"`
}

func (s *Service) savePluginFile(info string) error {
	pi := new(PluginInfo)
	err := json.Unmarshal([]byte(info), pi)
	if err != nil {
		return err
	}

	logrus.Debugf("save plugin file: %s", pi.Filename)

	err = fileutil.WriteToFile(pi.Filename, pi.Content)
	if err != nil {
		return err
	}
	return s.notifierManager.RegisterPlugin(pi.Name, pi.Filename)
}

func (s *Service) delPluginFile(path string) error {
	logrus.Debugf("del plugin file: %s", path)
	return os.Remove(path)
}

//func (s *Service) AddNotifierPlugin(p *NotifierPlugin, filePath string) error {
//	if err := s.notifierManager.RegisterPlugin(p.Name, filePath); err != nil {
//		return err
//	}
//	defer os.RemoveAll(filePath)
//	p.Version = fileutil.GetFileMd5(filePath)[:6]
//	err := os.Rename(filePath, filepath.Join(s.config.Server.RuntimeDirectories[config.PluginDirectoryNotify], p.Name+p.Version+".so"))
//	if err != nil {
//		return err
//	}
//	return s.notifierPluginPersistence.Insert(p)
//}

func (s *Service) ListNotifierPlugins() ([]*NotifierPlugin, error) {
	return s.notifierPluginPersistence.List(&NotifierPlugin{})
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

	err = s.cache.Subscribe(topicAddNotifierPlugin, s.savePluginFile)
	if err != nil {
		return err
	}

	err = s.cache.Subscribe(topicDelNotifierPlugin, s.delPluginFile)
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
		if p.IsBuiltIn {
			continue
		}
		err := s.notifierManager.RegisterPlugin(p.Name, s.getPluginFilePath(p))
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
