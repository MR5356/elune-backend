package script

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/MR5356/elune-backend/pkg/persistence"
)

type Script struct {
	ID      uint   `json:"id" gorm:"autoIncrement;primaryKey"`
	Title   string `json:"title" gorm:"not null"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
	Params  Params `json:"params"`
	Type    Type   `json:"type"`
	TypeId  uint   `json:"typeId"`

	persistence.BaseModel
}

func (s *Script) TableName() string {
	return "elune_script"
}

type Type struct {
	ID      uint      `json:"id" gorm:"autoIncrement;primaryKey"`
	Title   string    `json:"title" gorm:"not null"`
	Scripts []*Script `json:"scripts" gorm:"foreignkey:TypeId"`

	persistence.BaseModel
}

func (t *Type) TableName() string {
	return "elune_script_type"
}

type Params map[string]any

func (p *Params) Scan(val interface{}) error {
	s := val.(string)
	err := json.Unmarshal([]byte(s), &p)
	return err
}

func (p Params) Value() (driver.Value, error) {
	s, err := json.Marshal(p)
	return string(s), err
}
