package notify

import (
	"errors"
	"github.com/MR5356/notify"
	"github.com/sirupsen/logrus"
	"plugin"
	"sync"
)

var (
	ErrNotifierVerifyError  = errors.New("通知插件校验失败")
	ErrNotifierAlreadyExist = errors.New("通知插件名称已存在")
	ErrNotifierNotExist     = errors.New("通知插件不存在")
)

type NotifierManager struct {
	notifiers map[string]func(params ...string) notify.Notifier
	lock      sync.Mutex
}

func NewNotifierManager() *NotifierManager {
	return &NotifierManager{
		notifiers: make(map[string]func(params ...string) notify.Notifier),
		lock:      sync.Mutex{},
	}
}

func (m *NotifierManager) Verify(symbol plugin.Symbol) bool {
	if _, ok := symbol.(func(params ...string) notify.Notifier); !ok {
		return false
	}
	return true
}

func (m *NotifierManager) register(name string, symbol plugin.Symbol, params ...string) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	if !m.Verify(symbol) {
		return ErrNotifierVerifyError
	}
	if _, ok := m.notifiers[name]; ok {
		return ErrNotifierAlreadyExist
	}
	m.notifiers[name] = symbol.(func(params ...string) notify.Notifier)
	return nil
}

func (m *NotifierManager) Get(name string) (func(params ...string) notify.Notifier, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	notifier, ok := m.notifiers[name]
	if !ok {
		return nil, ErrNotifierNotExist
	}
	return notifier, nil
}

func (m *NotifierManager) RegisterPlugin(name, path string) error {
	logrus.Infof("load plugin: %s", path)
	p, err := plugin.Open(path)
	if err != nil {
		return err
	}
	symbol, err := p.Lookup("New")
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
