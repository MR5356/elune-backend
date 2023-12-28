package machine

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type Machine struct {
	ID       uint     `json:"id" gorm:"autoIncrement;primaryKey"`
	Title    string   `json:"title" gorm:"not null"`
	Desc     string   `json:"desc"`
	HostInfo HostInfo `json:"hostInfo" gorm:"unique;not null"`
	MetaInfo MetaInfo `json:"metaInfo"`
	Group    Group    `json:"group"`
	GroupId  uint     `json:"groupId"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Group struct {
	ID       uint       `json:"id" gorm:"autoIncrement;primaryKey"`
	Title    string     `json:"title" gorm:"unique;not null"`
	Machines []*Machine `json:"machines" gorm:"foreignkey:GroupId"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (g *Group) TableName() string {
	return "elune_machine_group"
}

type HostInfo struct {
	Host     string `json:"host"`
	Port     uint16 `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type MetaInfo struct {
	OS       string `json:"os"`
	Kernel   string `json:"kernel"`
	Hostname string `json:"hostname"`
	Arch     string `json:"arch"`
	Cpu      string `json:"cpu"`
	Mem      string `json:"mem"`
}

func (m *Machine) TableName() string {
	return "elune_machine"
}

func (m *MetaInfo) Scan(val interface{}) error {
	s := val.(string)
	err := json.Unmarshal([]byte(s), &m)
	return err
}

func (m MetaInfo) Value() (driver.Value, error) {
	s, err := json.Marshal(m)
	return string(s), err
}

func (h *HostInfo) Scan(val interface{}) error {
	s := val.(string)
	err := json.Unmarshal([]byte(s), &h)
	return err
}

func (h HostInfo) Value() (driver.Value, error) {
	s, err := json.Marshal(h)
	return string(s), err
}
