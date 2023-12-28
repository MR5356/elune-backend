package notify

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/MR5356/elune-backend/pkg/persistence"
	"github.com/MR5356/notify"
)

type MessageTemplate struct {
	ID      uint    `json:"id" gorm:"autoIncrement;primaryKey"`
	Title   string  `json:"title" gorm:"length:64;not null"`
	Message Message `json:"message"`

	persistence.BaseModel
}

func (n *MessageTemplate) TableName() string {
	return "elune_notify_message_templates"
}

type Message notify.Message

func (m *Message) Scan(val interface{}) error {
	s := val.([]byte)
	*m = Message{}
	return json.Unmarshal([]byte(s), m)
}

func (m Message) Value() (driver.Value, error) {
	return json.Marshal(m)
}

type NotifierChannel struct {
	ID               uint                  `json:"id" gorm:"autoIncrement;primaryKey"`
	NotifierPluginID uint                  `json:"notifierPluginId"`
	Name             string                `json:"name" gorm:"length:64;not null"`
	Desc             string                `json:"desc"`
	Params           NotifierChannelParams `json:"params"`

	persistence.BaseModel
}

func (c *NotifierChannel) TableName() string {
	return "elune_notifier_channels"
}

type NotifierChannelParams []string

func (p *NotifierChannelParams) Scan(val interface{}) error {
	s := val.([]byte)
	*p = make([]string, 0)
	return json.Unmarshal([]byte(s), p)
}

func (p NotifierChannelParams) Value() (driver.Value, error) {
	return json.Marshal(p)
}

type NotifierPlugin struct {
	ID        uint                                   `json:"id" gorm:"autoIncrement;primaryKey"`
	Name      string                                 `json:"name" gorm:"length:32;not null"`
	Version   string                                 `json:"-" gorm:"unique;not null"`
	Params    NotifierParams                         `json:"params"`
	Status    string                                 `json:"status"`
	URL       string                                 `json:"url"`
	IsBuiltIn bool                                   `json:"isBuiltIn"`
	Installed bool                                   `json:"installed"`
	From      string                                 `json:"from"`
	Symbol    func(params ...string) notify.Notifier `json:"-" gorm:"-"`

	persistence.BaseModel
}

func (p *NotifierPlugin) TableName() string {
	return "elune_notifier_plugins"
}

type NotifierParams []notify.Param

func (p *NotifierParams) Scan(val interface{}) error {
	s := val.([]byte)
	*p = make([]notify.Param, 0)
	return json.Unmarshal([]byte(s), p)
}

func (p NotifierParams) Value() (driver.Value, error) {
	return json.Marshal(p)
}
