package kubernetes

import (
	"github.com/MR5356/elune-backend/pkg/kubernetes/client"
	"github.com/sirupsen/logrus"
)

type Service struct {
	client *client.Client
}

func NewService(config string) *Service {
	c, err := client.New(config)
	if err != nil {
		logrus.Errorf("create kubernetes client error: %+v", err)
		return nil
	}
	return &Service{
		client: c,
	}
}

func (s *Service) GetNodes() ([]client.NodeInfo, error) {
	return s.client.GetNodes()
}

func (s *Service) Initialize() error {
	return nil
}
