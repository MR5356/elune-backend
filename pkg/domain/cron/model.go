package cron

import (
	"github.com/MR5356/elune-backend/pkg/persistence"
	"time"
)

type Cron struct {
	ID         uint      `json:"id" gorm:"autoIncrement;primaryKey"`
	Title      string    `json:"title" gorm:"not null"`
	Desc       string    `json:"desc"`
	CronString string    `json:"cronString"`
	TaskName   string    `json:"taskName"`
	NextTime   time.Time `json:"nextTime" gorm:"-"`
	Params     string    `json:"params"`
	Enabled    bool      `json:"enabled"`

	persistence.BaseModel
}

func (c *Cron) TableName() string {
	return "elune_cron"
}

type Record struct {
	ID       uint   `json:"id" gorm:"autoIncrement;primaryKey"`
	CronId   uint   `json:"cronId"`
	Title    string `json:"title"`
	Desc     string `json:"desc"`
	TaskName string `json:"taskName"`
	Params   string `json:"params"`
	Log      string `json:"log"`
	Status   string `json:"status"`

	persistence.BaseModel
}

func (r *Record) TableName() string {
	return "elune_cron_records"
}
