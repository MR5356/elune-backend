package notify

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/MR5356/elune-backend/pkg/persistence"
	"github.com/MR5356/notify"
)

type NotifierPlugin struct {
	ID      uint   `json:"id" gorm:"autoIncrement;primaryKey"`
	Name    string `json:"name" gorm:"length:32;unique;not null"`
	Version string `json:"-" gorm:"unique;not null"`
	Params  NotifierParams
	Status  string `json:"status"`

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
