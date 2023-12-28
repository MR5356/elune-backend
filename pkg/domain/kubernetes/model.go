package kubernetes

import (
	"github.com/MR5356/elune-backend/pkg/domain/kubernetes/client"
	"time"
)

type Kubernetes struct {
	ID         uint               `json:"id" gorm:"autoIncrement;primaryKey"`
	Title      string             `json:"title" gorm:"not null"`
	Desc       string             `json:"desc"`
	KubeConfig string             `json:"kubeConfig" gorm:"not null"`
	Md5        string             `json:"md5" gorm:"unique;not null"`
	Version    string             `json:"version"`
	Status     string             `json:"status" gorm:"-"`
	Nodes      []*client.NodeInfo `json:"nodes" gorm:"-"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (k *Kubernetes) TableName() string {
	return "elune_kubernetes"
}
