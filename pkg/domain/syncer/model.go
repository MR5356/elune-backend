package syncer

import (
	"github.com/MR5356/elune-backend/pkg/persistence"
	"time"
)

type Syncer struct {
	ID     uint   `json:"id" gorm:"autoIncrement;primaryKey"`
	Title  string `json:"title" gorm:"not null"`
	Desc   string `json:"desc"`
	Config string `json:"config" gorm:"not null"`
	Type   Type   `json:"type"`
	TypeId uint   `json:"typeId"`

	persistence.BaseModel
}

func (s *Syncer) TableName() string {
	return "elune_syncer"
}

type Type struct {
	ID      uint      `json:"id" gorm:"autoIncrement;primaryKey"`
	Title   string    `json:"title" gorm:"unique;not null"`
	Syncers []*Syncer `json:"syncers" gorm:"foreignkey:TypeId"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (t *Type) TableName() string {
	return "elune_syncer_type"
}
