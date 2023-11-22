package machine

import (
	"errors"
	"fmt"
	"github.com/MR5356/elune-backend/pkg/persistence"
	"github.com/MR5356/elune-backend/pkg/persistence/cache"
	"github.com/MR5356/elune-backend/pkg/persistence/database"
	"github.com/MR5356/elune-backend/pkg/utils/structutil"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"time"
)

type Service struct {
	machinePersistence *persistence.Persistence[*Machine]
	groupPersistence   *persistence.Persistence[*Group]
}

func NewService(database *database.Database, cache *cache.Cache) *Service {
	return &Service{
		machinePersistence: persistence.New(database, cache, &Machine{}),
		groupPersistence:   persistence.New(database, cache, &Group{}),
	}
}

func (s *Service) AddGroup(group *Group) error {
	group.ID = 0
	if len(group.Title) == 0 {
		return errors.New("组名称不可为空")
	}
	return s.groupPersistence.Insert(group)
}

func (s *Service) ListGroup() ([]*Group, error) {
	return s.groupPersistence.List(&Group{})
}

func (s *Service) DeleteGroup(id uint) error {
	err := s.machinePersistence.Delete(&Machine{GroupId: id})
	if err != nil {
		return err
	}
	return s.groupPersistence.Delete(&Group{ID: id})
}

func (s *Service) UpdateGroup(group *Group) error {
	if len(group.Title) == 0 {
		return errors.New("组名称不可为空")
	}
	return s.groupPersistence.Update(&Group{ID: group.ID}, structutil.Struct2Map(&Group{
		ID:    group.ID,
		Title: group.Title,
	}))
}

func (s *Service) AddMachine(machine *Machine) error {
	logrus.Debugf("machine: %+v", machine)
	machine.ID = 0
	if len(machine.Title) == 0 {
		return errors.New("机器名称不可为空")
	}
	if err := s.checkMachine(machine); err != nil {
		return fmt.Errorf("机器连接失败: %s", err)
	}
	return s.machinePersistence.Insert(machine)
}

func (s *Service) ListMachine() ([]*Machine, error) {
	res := make([]*Machine, 0)
	err := s.machinePersistence.DB.Joins("Group").Find(&res).Error
	return res, err
}

func (s *Service) DeleteMachine(id uint) error {
	return s.machinePersistence.Delete(&Machine{ID: id})
}

func (s *Service) UpdateMachine(machine *Machine) error {
	if len(machine.Title) == 0 {
		return errors.New("机器名称不可为空")
	}
	if err := s.checkMachine(machine); err != nil {
		return fmt.Errorf("机器连接失败: %s", err)
	}
	return s.machinePersistence.Update(&Machine{ID: machine.ID}, structutil.Struct2Map(&Machine{
		ID:       machine.ID,
		Title:    machine.Title,
		Desc:     machine.Desc,
		HostInfo: machine.HostInfo,
		MetaInfo: machine.MetaInfo,
		GroupId:  machine.GroupId,
	}))
}

func (s *Service) checkMachine(machine *Machine) error {
	sshClient, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", machine.HostInfo.Host, machine.HostInfo.Port), &ssh.ClientConfig{
		User:            machine.HostInfo.Username,
		Auth:            []ssh.AuthMethod{ssh.Password(machine.HostInfo.Password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * 10,
	})
	if err != nil {
		return err
	}
	defer sshClient.Close()
	// TODO 增加系统信息检测
	machine.MetaInfo = MetaInfo{
		OS:       "centos",
		Kernel:   "3.15.2",
		Hostname: "hw",
		Arch:     "amd64",
		Cpu:      "2",
		Mem:      "2048",
	}
	return err
}

func (s *Service) Initialize() error {
	err := s.machinePersistence.DB.AutoMigrate(&Group{}, &Machine{})
	if err != nil {
		return err
	}
	_ = s.groupPersistence.Insert(&Group{
		ID:    1,
		Title: "未分组",
	})
	return nil
}
