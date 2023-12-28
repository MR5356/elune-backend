package executor

import "github.com/MR5356/elune-backend/pkg/persistence"

type Record struct {
	ID          string `json:"id" gorm:"primaryKey"`
	ScriptTitle string `json:"scriptTitle"`
	Script      string `json:"script"`
	Host        string `json:"host"`
	Params      string `json:"params"`
	Result      string `json:"result"`
	Status      string `json:"status"`
	Message     string `json:"message"`
	Error       string `json:"error"`

	persistence.BaseModel
}

func (r *Record) TableName() string {
	return "elune_executor_records"
}
