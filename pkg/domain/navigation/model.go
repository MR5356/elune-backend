package navigation

import (
	"time"
)

type Navigation struct {
	ID       uint   `json:"id" gorm:"autoIncrement;primaryKey"`
	Title    string `json:"title" gorm:"not null"`
	Href     string `json:"href"`
	Logo     string `json:"logo"`
	LogoData string `json:"logoData"`
	Desc     string `json:"desc"`
	Parent   uint   `json:"parent" default:"0"` // 0为parent
	Order    int    `json:"order"`
	Unique   string `json:"-" gorm:"unique;not null"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (n *Navigation) TableName() string {
	return "elune_navigation"
}
