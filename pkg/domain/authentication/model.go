package authentication

import (
	"github.com/MR5356/elune-backend/pkg/persistence"
)

type User struct {
	ID       uint   `json:"id" gorm:"autoIncrement;primaryKey"`
	Nickname string `json:"nickname" gorm:"not null;default:''"`
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

type Group struct {
	ID    uint   `json:"id" gorm:"autoIncrement;primaryKey"`
	Title string `json:"title" gorm:"not null"`
	Desc  string `json:"desc"`

	persistence.BaseModel
}

func (g *Group) TableName() string {
	return "elune_user_groups"
}

type UserGroup struct {
	ID      uint `json:"id" gorm:"autoIncrement;primaryKey"`
	UserID  uint `json:"userId" gorm:"not null"`
	GroupID uint `json:"groupId" gorm:"not null"`

	persistence.BaseModel
}

func (ug *UserGroup) TableName() string {
	return "elune_user_group_relations"
}

type Role struct {
	ID   uint   `json:"id" gorm:"autoIncrement;primaryKey"`
	Name string `json:"name" gorm:"not null;unique"`
	Desc string `json:"desc"`

	persistence.BaseModel
}

func (r *Role) TableName() string {
	return "elune_roles"
}

type RoleRelation struct {
	ID     uint   `json:"id" gorm:"autoIncrement;primaryKey"`
	RoleID uint   `json:"roleId" gorm:"not null"`
	UID    uint   `json:"uid" gorm:"not null"`
	Type   string `json:"type"` // group/user
}

func (rr *RoleRelation) TableName() string {
	return "elune_role_relations"
}
