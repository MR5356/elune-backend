package navigation

import "github.com/MR5356/elune-backend/pkg/persistence"

type Navigation struct {
	ID     uint   `json:"id" gorm:"autoIncrement;primaryKey"`
	Title  string `json:"title" gorm:"not null"`
	Href   string `json:"href"`
	Logo   string `json:"logo"`
	Desc   string `json:"desc"`
	Parent uint   `json:"parent" default:"0"` // -1为parent
	Order  int    `json:"order"`
	Unique string `json:"-" gorm:"unique;not null"`

	persistence.BaseModel
}

func (n *Navigation) TableName() string {
	return "elune_navigation"
}
