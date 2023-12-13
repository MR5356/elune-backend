package notify

import (
	"context"
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
	"github.com/MR5356/notify"
	"github.com/MR5356/notify/notifier/lark"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"time"
)

const (
	topicAddNotifierPlugin  = "notify.plugin.add"
	topicDelNotifierPlugin  = "notify.plugin.del"
	topicAddNotifierChannel = "notify.channel.add"
	topicDelNotifierChannel = "notify.channel.del"
)

type Service struct {
	notifierPluginManager      *NotifierPluginManager
	notifierChannelManager     *NotifierChannelManager
	notifierPluginPersistence  *persistence.Persistence[*NotifierPlugin]
	notifierChannelPersistence *persistence.Persistence[*NotifierChannel]
	messageTemplatePersistence *persistence.Persistence[*MessageTemplate]
	config                     *config.Config
	cache                      cache.Cache
}

func NewService(database *database.Database, cache cache.Cache, cfg *config.Config) *Service {
	return &Service{
		notifierPluginManager:      NewNotifierPluginManager(),
		notifierChannelManager:     NewNotifierChannelManager(),
		notifierPluginPersistence:  persistence.New(database, cache, &NotifierPlugin{}),
		notifierChannelPersistence: persistence.New(database, cache, &NotifierChannel{}),
		messageTemplatePersistence: persistence.New(database, cache, &MessageTemplate{}),
		config:                     cfg,
		cache:                      cache,
	}
}

func (s *Service) ListMessageTemplates() ([]*MessageTemplate, error) {
	return s.messageTemplatePersistence.List(&MessageTemplate{})
}

func (s *Service) AddMessageTemplate(message *MessageTemplate) error {
	if len(message.Title) == 0 {
		return errors.New("消息模版名称不能为空")
	}
	return s.messageTemplatePersistence.Insert(message)
}

func (s *Service) RemoveMessageTemplate(id uint) error {
	return s.messageTemplatePersistence.Delete(&MessageTemplate{ID: id})
}

func (s *Service) UpdateMessageTemplate(message *MessageTemplate) error {
	if len(message.Title) == 0 {
		return errors.New("消息模版名称不能为空")
	}
	return s.messageTemplatePersistence.Update(&MessageTemplate{ID: message.ID}, structutil.Struct2Map(message))
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
	symbol, err := s.notifierPluginManager.GetSymbol(filePath)
	if err != nil {
		return ErrNotifierOpenError
	}
	notifierPlugin, ok := s.notifierPluginManager.Verify(symbol)
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
	notifierPlugin.From = "自定义插件"

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
	return s.notifierPluginManager.RegisterPlugin(pi.Name, pi.Filename)
}

func (s *Service) delPluginFile(path string) error {
	logrus.Debugf("del plugin file: %s", path)
	return os.Remove(path)
}

func (s *Service) ListNotifierPlugins() ([]*NotifierPlugin, error) {
	return s.notifierPluginPersistence.List(&NotifierPlugin{})
}

func (s *Service) AddNotifierChannel(notifierChannel *NotifierChannel) error {
	np, err := s.notifierPluginPersistence.Detail(&NotifierPlugin{ID: notifierChannel.NotifierPluginID})
	if err != nil {
		return errors.New("插件不存在")
	}
	if len(notifierChannel.Params) != len(np.Params) {
		return fmt.Errorf("插件参数不匹配，插件需要%d个参数，当前传入%d个参数", len(np.Params), len(notifierChannel.Params))
	}
	err = s.notifierChannelPersistence.Insert(notifierChannel)
	if err != nil {
		return err
	}
	s.cache.Publish(topicAddNotifierChannel, notifierChannel.ID)
	return nil
}

func (s *Service) ListNotifierChannels() ([]*NotifierChannel, error) {
	return s.notifierChannelPersistence.List(&NotifierChannel{})
}

func (s *Service) RemoveNotifierChannel(id uint) error {
	err := s.notifierChannelPersistence.Delete(&NotifierChannel{ID: id})
	if err != nil {
		return err
	}
	s.cache.Publish(topicDelNotifierChannel, id)
	return nil
}

func (s *Service) SendTestMessage(notifierChannelId uint) error {
	notifier, err := s.notifierChannelManager.GetNotifierChannel(notifierChannelId)
	if err != nil {
		return err
	}

	err = notifier.Send(
		context.Background(),
		notify.NewMessage("通道测试消息").
			CardBuilder().
			WithLevel(notify.MessageLevelInfo).
			AddText(notify.MessageLayoutBisected, true, "**时间**\n"+time.Now().Format("2006-01-02 15:04:05"), "**事件**\n测试消息通道").
			AddText(notify.MessageLayoutDefault, true, "如果你看到这条消息，说明通道配置正确").
			WithNote("来自Elune", true).
			Build(),
	)
	return err
}

func (s *Service) addNotifierChannel(id uint) error {
	notifierChannel, err := s.notifierChannelPersistence.Detail(&NotifierChannel{ID: id})
	if err != nil {
		return err
	}
	np, err := s.notifierPluginPersistence.Detail(&NotifierPlugin{ID: notifierChannel.NotifierPluginID})
	if err != nil {
		return err
	}

	for k, v := range s.notifierPluginManager.notifiers {
		logrus.Debugf("plugins k%+v : %+v", k, v)
	}
	p, err := s.notifierPluginManager.Get(np.Name)
	if err != nil {
		return err
	}

	c := p(notifierChannel.Params...)
	s.notifierChannelManager.RegisterNotifierChannel(notifierChannel.ID, c)
	return nil
}

func (s *Service) delNotifierChannel(id uint) error {
	s.notifierChannelManager.RemoveNotifierChannel(id)
	return nil
}

func (s *Service) newLarkNotifierPlugin(params ...string) notify.Notifier {
	return lark.NewWebhookBot(params[0])
}

func (s *Service) Initialize() error {
	err := s.notifierPluginPersistence.DB.AutoMigrate(&NotifierPlugin{}, &NotifierChannel{}, &MessageTemplate{})
	if err != nil {
		return err
	}

	defaultNotifyMessageTemplate := []*MessageTemplate{
		{
			ID:    1,
			Title: "通道测试消息",
			Message: Message(*notify.NewMessage("通道测试消息").
				CardBuilder().
				WithLevel(notify.MessageLevelInfo).
				AddText(notify.MessageLayoutBisected, true, "**时间**\n"+time.Now().Format("2006-01-02 15:04:05"), "**事件**\n测试消息通道").
				AddText(notify.MessageLayoutDefault, true, "如果你看到这条消息，说明通道配置正确").
				WithNote("来自Elune", true).
				Build()),
		},
	}

	for _, message := range defaultNotifyMessageTemplate {
		_ = s.messageTemplatePersistence.Insert(message)
	}

	builtinPlugins := []*NotifierPlugin{
		{
			ID:      1,
			Name:    "larkNotifier",
			Version: "9ae861",
			Params: []notify.Param{
				{
					Name: "url",
					Type: "string",
					Desc: "飞书通知机器人 webhook 地址",
				},
			},
			Symbol: s.newLarkNotifierPlugin,
		},
	}

	for _, p := range builtinPlugins {
		p.From = "内置插件"
		p.IsBuiltIn = true
		p.Installed = true
		p.Status = "success"
		_ = s.notifierPluginPersistence.Insert(p)
		err := s.notifierPluginManager.RegisterBuiltIn(p.Name, p.Symbol)
		if err != nil {
			return err
		}
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

	err = s.cache.Subscribe(topicAddNotifierChannel, s.addNotifierChannel)
	if err != nil {
		return err
	}

	err = s.cache.Subscribe(topicDelNotifierChannel, s.delNotifierChannel)
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
		err := s.notifierPluginManager.RegisterPlugin(p.Name, s.getPluginFilePath(p))
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

	notifierChannels, err := s.notifierChannelPersistence.List(&NotifierChannel{})
	if err != nil {
		return err
	}
	for _, notifierChannel := range notifierChannels {
		err := s.addNotifierChannel(notifierChannel.ID)
		if err != nil {
			return err
		}
	}

	return err
}
