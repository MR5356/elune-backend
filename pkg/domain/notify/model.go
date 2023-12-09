package notify

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/MR5356/elune-backend/pkg/persistence"
)

type NotifierPlugin struct {
	ID      uint   `json:"id" gorm:"autoIncrement;primaryKey"`
	Name    string `json:"name"`
	Desc    string `json:"desc"`
	Version string `json:"-" gorm:"unique;not null"`
	Params  NotifierParams
	Status  string `json:"status"`

	persistence.BaseModel
}

func (p *NotifierPlugin) TableName() string {
	return "elune_notifier_plugins"
}

type NotifierParams []string

func (p *NotifierParams) Scan(val interface{}) error {
	s := val.(string)
	*p = make(NotifierParams, 0)
	return json.Unmarshal([]byte(s), p)
}

func (p NotifierParams) Value() (driver.Value, error) {
	return json.Marshal(p)
}
