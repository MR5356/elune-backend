package authentication

import (
	"github.com/MR5356/elune-backend/pkg/persistence"
)

type User struct {
	ID       uint   `json:"id" gorm:"autoIncrement;primaryKey"`
	Username string `json:"username" gorm:"not null;unique"`
	Password string `json:"password" gorm:"not null"`
	Email    string `json:"email" gorm:"not null;unique"`

	persistence.BaseModel
}

func (u *User) TableName() string {
	return "elune_users"
}

func (u *User) Desensitization() *User {
	u.Password = "********"
	return u
}
