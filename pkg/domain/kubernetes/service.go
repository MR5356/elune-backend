package kubernetes

import (
	client2 "github.com/MR5356/elune-backend/pkg/domain/kubernetes/client"
	"github.com/sirupsen/logrus"
)

type Service struct {
	client *client2.Client
}

func NewService(config string) *Service {
	c, err := client2.New(config)
	if err != nil {
		logrus.Errorf("create kubernetes client error: %+v", err)
		return nil
	}
	return &Service{
		client: c,
	}
}

func (s *Service) GetNodes() ([]client2.NodeInfo, error) {
	return s.client.GetNodes()
}

func (s *Service) Initialize() error {
	return nil
}
