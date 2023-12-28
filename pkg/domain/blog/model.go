package blog

import (
	"database/sql/driver"
	"github.com/MR5356/elune-backend/pkg/persistence"
	"strings"
)

type Blog struct {
	ID       uint     `json:"id" gorm:"autoIncrement;primaryKey"`
	Title    string   `json:"title" gorm:"not null"`
	Desc     string   `json:"desc"`
	Author   string   `json:"author"`
	Content  string   `json:"content"`
	Likes    int      `json:"likes"`
	Reads    int      `json:"reads"`
	Category Category `json:"categories"`

	persistence.BaseModel
}

func (c *Blog) TableName() string {
	return "elune_blog"
}

type Category []string

func (c *Category) Scan(val interface{}) error {
	s := val.(string)
	ss := strings.Split(s, ",,,")
	*c = ss
	return nil
}

func (c Category) Value() (driver.Value, error) {
	return strings.Join(c, ",,,"), nil
}
