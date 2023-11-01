package site

import "github.com/MR5356/elune-backend/pkg/persistence"

type SiteConfig struct {
	ID    uint   `json:"id" gorm:"autoIncrement;primaryKey"`
	Key   string `json:"key" gorm:"unique;not null"`
	Value string `json:"value"`

	persistence.BaseModel
}

func (c *SiteConfig) TableName() string {
	return "site_config"
}
