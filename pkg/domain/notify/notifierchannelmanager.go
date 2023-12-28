package notify

import (
	"errors"
	"github.com/MR5356/notify"
	"github.com/sirupsen/logrus"
	"sync"
)

var (
	ErrNotifierChannelAlreadyExist = errors.New("通知渠道名称已存在")
	ErrNotifierChannelNotExist     = errors.New("通知渠道不存在")
)

type NotifierChannelManager struct {
	channels map[uint]notify.Notifier
	lock     sync.Mutex
}

func NewNotifierChannelManager() *NotifierChannelManager {
	return &NotifierChannelManager{
		channels: make(map[uint]notify.Notifier),
		lock:     sync.Mutex{},
	}
}

func (m *NotifierChannelManager) RegisterNotifierChannel(id uint, notifier notify.Notifier) {
	logrus.Debugf("register notifier channel: %d", id)
	m.lock.Lock()
	defer m.lock.Unlock()
	if _, ok := m.channels[id]; ok {
		return
	}
	m.channels[id] = notifier
	logrus.Debugf("register notifier channel success: %d", id)
	logrus.Debugf("channels: %v", m.channels)
}

func (m *NotifierChannelManager) RemoveNotifierChannel(id uint) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.channels, id)
}

func (m *NotifierChannelManager) GetNotifierChannel(id uint) (notify.Notifier, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	notifier, ok := m.channels[id]
	if !ok {
		return nil, ErrNotifierChannelNotExist
	}
	return notifier, nil
}
