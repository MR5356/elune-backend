package notify

import (
	"errors"
	"github.com/MR5356/notify"
	"github.com/sirupsen/logrus"
	"plugin"
	"sync"
)

var (
	ErrNotifierOpenError    = errors.New("通知插件加载失败，文件错误或者已经添加此插件")
	ErrNotifierVerifyError  = errors.New("通知插件校验失败")
	ErrNotifierAlreadyExist = errors.New("通知插件名称已存在")
	ErrNotifierNotExist     = errors.New("通知插件不存在")
)

type NotifierPluginManager struct {
	notifiers map[string]func(params ...string) notify.Notifier
	lock      sync.Mutex
	symbols   sync.Map
}

func NewNotifierPluginManager() *NotifierPluginManager {
	return &NotifierPluginManager{
		notifiers: make(map[string]func(params ...string) notify.Notifier),
		lock:      sync.Mutex{},
		symbols:   sync.Map{},
	}
}

func (m *NotifierPluginManager) Verify(symbol plugin.Symbol) (*NotifierPlugin, bool) {
	if n, ok := symbol.(func(params ...string) notify.Notifier); ok {
		return &NotifierPlugin{
			Name:   n().Name(),
			Params: n().Params(),
		}, true
	} else {
		return nil, false
	}
}

func (m *NotifierPluginManager) register(name string, symbol plugin.Symbol, params ...string) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	if _, ok := m.Verify(symbol); !ok {
		return ErrNotifierVerifyError
	}
	if _, ok := m.notifiers[name]; ok {
		return ErrNotifierAlreadyExist
	}
	m.notifiers[name] = symbol.(func(params ...string) notify.Notifier)
	return nil
}

func (m *NotifierPluginManager) Get(name string) (func(params ...string) notify.Notifier, error) {
	logrus.Debugf("get notifier plugin: %s", name)
	m.lock.Lock()
	defer m.lock.Unlock()
	notifier, ok := m.notifiers[name]
	if !ok {
		return nil, ErrNotifierNotExist
	}
	return notifier, nil
}

func (m *NotifierPluginManager) GetSymbol(filePath string) (plugin.Symbol, error) {
	if v, ok := m.symbols.Load(filePath); ok {
		return v.(plugin.Symbol), nil
	}
	p, err := plugin.Open(filePath)
	if err != nil {
		logrus.Errorf("plugin open error: %s", err)
		return nil, err
	}
	symbol, err := p.Lookup("New")
	if err != nil {
		logrus.Errorf("plugin lookup error: %s", err)
		return nil, err
	}
	m.symbols.Store(filePath, symbol)
	return symbol, nil
}

func (m *NotifierPluginManager) RegisterBuiltIn(name string, symbol func(params ...string) notify.Notifier) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	if _, ok := m.notifiers[name]; ok {
		return ErrNotifierAlreadyExist
	}
	m.notifiers[name] = symbol
	return nil
}

func (m *NotifierPluginManager) RegisterPlugin(name, path string) error {
	logrus.Infof("load plugin: %s", path)

	symbol, err := m.GetSymbol(path)
	if err != nil {
		return err
	}
	err = m.register(name, symbol)
	if err != nil {
		return err
	}
	logrus.Infof("success load plugin: %s", path)
	return nil
}
