package persistence

import (
	"time"
)

type BaseModel struct {
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	//DeleteAt  gorm.DeletedAt `json:"-"`
}
