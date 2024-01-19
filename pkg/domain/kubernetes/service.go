package kubernetes

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/MR5356/elune-backend/pkg/domain/kubernetes/client"
	"github.com/MR5356/elune-backend/pkg/persistence"
	"github.com/MR5356/elune-backend/pkg/persistence/cache"
	"github.com/MR5356/elune-backend/pkg/persistence/database"
	"github.com/MR5356/elune-backend/pkg/utils/structutil"
	"sync"
	"time"
)

const (
	kubernetesReady   = "ready"
	kubernetesError   = "error"
	kubernetesTimeout = "timeout"
)

type Service struct {
	k8sPersistence *persistence.Persistence[*Kubernetes]
	database       *database.Database
	cache          cache.Cache
}

func NewService(database *database.Database, cache cache.Cache) *Service {
	return &Service{
		k8sPersistence: persistence.New(database, cache, &Kubernetes{}),
		database:       database,
		cache:          cache,
	}
}

func (s *Service) AddKubernetes(k *Kubernetes) error {
	if len(k.Title) == 0 {
		return errors.New("集群名称不能为空")
	}
	if len(k.KubeConfig) == 0 {
		return errors.New("集群配置不能为空")
	}
	c, err := client.New(k.KubeConfig)
	if err != nil {
		return errors.New("无效的集群配置")
	}
	ver, err := c.GetVersion()
	if err != nil {
		ver = "unknown"
	}
	k.Version = ver
	k.ID = 0
	k.Md5 = fmt.Sprintf("%x", md5.Sum([]byte(k.KubeConfig)))
	return s.k8sPersistence.Insert(k)
}

func (s *Service) ListKubernetes() ([]*Kubernetes, error) {
	k8ss, err := s.k8sPersistence.List(&Kubernetes{})
	if err != nil {
		return nil, err
	}
	wg := sync.WaitGroup{}
	for _, k8s := range k8ss {
		wg.Add(1)
		k8s := k8s
		go func() {
			defer wg.Done()
			s.setKubernetesInfos(k8s)
			// 屏蔽敏感信息
			k8s.KubeConfig = "******"
		}()
	}
	wg.Wait()
	return k8ss, nil
}

func (s *Service) setKubernetesInfos(k8s *Kubernetes) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	errored := make(chan error)
	finished := make(chan bool)
	go func() {
		c, err := client.New(k8s.KubeConfig)
		if err != nil {
			errored <- err
			return
		}
		nodes, err := c.GetNodes()
		if err != nil {
			k8s.Nodes = make([]*client.NodeInfo, 0)
			errored <- err
			return
		}
		k8s.Nodes = nodes
		ver, err := c.GetVersion()
		if err != nil {
			k8s.Version = "unknown"
			errored <- err
			return
		}
		k8s.Version = ver
		finished <- true
	}()
	select {
	case <-ctx.Done():
		k8s.Status = kubernetesTimeout
	case <-errored:
		k8s.Status = kubernetesError
	case <-finished:
		k8s.Status = kubernetesReady
	}
}

func (s *Service) UpdateKubernetes(k *Kubernetes) error {
	if len(k.Title) == 0 {
		return errors.New("集群名称不能为空")
	}
	if len(k.KubeConfig) == 0 {
		return errors.New("集群配置不能为空")
	}
	c, err := client.New(k.KubeConfig)
	if err != nil {
		return errors.New("无效的集群配置")
	}
	ver, err := c.GetVersion()
	if err != nil {
		ver = "unknown"
	}
	k.Version = ver
	k.Md5 = fmt.Sprintf("%x", md5.Sum([]byte(k.KubeConfig)))
	return s.k8sPersistence.Update(&Kubernetes{ID: k.ID}, structutil.Struct2Map(k))
}

func (s *Service) DeleteKubernetes(id uint) error {
	return s.k8sPersistence.Delete(&Kubernetes{ID: id})
}

func (s *Service) Initialize() error {
	err := s.k8sPersistence.DB.AutoMigrate(&Kubernetes{})
	if err != nil {
		return err
	}

	return nil
}
