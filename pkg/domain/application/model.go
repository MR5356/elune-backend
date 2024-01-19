package application

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/MR5356/elune-backend/pkg/persistence"
)

type Application struct {
	ID    uint   `json:"id" gorm:"autoIncrement;primaryKey"`
	Title string `json:"title" gorm:"not null"`
	Desc  string `json:"desc"`
	Logo  string `json:"logo"`
	Order int    `json:"order"`
	Link  string `json:"link"`
	OS    OS     `json:"os"`
	Note  string `json:"note" gorm:"length:16"`

	persistence.BaseModel
}

func (*Application) TableName() string {
	return "elune_application"
}

type OS []string

func (o *OS) Scan(val interface{}) error {
	s := val.([]byte)
	*o = OS{}
	return json.Unmarshal([]byte(s), o)
}

func (o OS) Value() (driver.Value, error) {
	return json.Marshal(o)
}
