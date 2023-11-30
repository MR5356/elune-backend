package machine

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/MR5356/elune-backend/pkg/persistence"
	"github.com/MR5356/elune-backend/pkg/persistence/cache"
	"github.com/MR5356/elune-backend/pkg/persistence/database"
	"github.com/MR5356/elune-backend/pkg/utils/structutil"
	"github.com/MR5356/jietan/pkg/executor"
	"github.com/MR5356/jietan/pkg/executor/api"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"time"
)

type Service struct {
	machinePersistence *persistence.Persistence[*Machine]
	groupPersistence   *persistence.Persistence[*Group]
}

func NewService(database *database.Database, cache cache.Cache) *Service {
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
	res := make([]*Group, 0)
	err := s.groupPersistence.DB.Preload("Machines").Find(&res, &Group{}).Error
	for _, item := range res {
		for _, m := range item.Machines {
			m.HostInfo.Password = "******"
		}
	}
	return res, err
}

func (s *Service) DeleteGroup(id uint) error {
	err := s.machinePersistence.DB.Where("group_id = ?", id).Delete(&Machine{GroupId: id}).Error
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
	for _, item := range res {
		item.HostInfo.Password = "******"
	}
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
		Timeout:         time.Second * 20,
	})
	if err != nil {
		return err
	}
	defer sshClient.Close()
	exec := executor.GetExecutor("remote")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	hostInfo := &api.HostInfo{
		Host:   machine.HostInfo.Host,
		Port:   int(machine.HostInfo.Port),
		User:   machine.HostInfo.Username,
		Passwd: machine.HostInfo.Password,
	}
	a := exec.Execute(ctx, &api.ExecuteParams{
		"hosts": []*api.HostInfo{
			hostInfo,
		},
		"script": "#!/bin/bash\n\n# 系统信息\nos=$(cat /etc/os-release | grep ^ID= | cut -d '=' -f2)\nos=${os//\\\"/} # 去掉双引号\nkernel=$(uname -r)\nhostname=$(hostname)\narch=$(uname -m)\n\n# 硬件信息\ncpu_count=$(lscpu | grep \"^CPU:\\|^CPU(s)\" | cut -d ':' -f2 | awk '{$1=$1;print}')\nmem_size=$(free -h | grep Mem | awk '{print $2}')\n\n# 构建JSON\njson=\"{\\\"os\\\": \\\"$os\\\",\n        \\\"kernel\\\": \\\"$kernel\\\",\n        \\\"hostname\\\": \\\"$hostname\\\",\n        \\\"arch\\\": \\\"$arch\\\",\n        \\\"cpu\\\": \\\"$cpu_count\\\",\n        \\\"mem\\\": \\\"$mem_size\\\"}\"\n\n# 输出\necho $json",
		"params": "",
	})
	metaInfo := new(MetaInfo)
	err = json.Unmarshal([]byte(a.Data["log"].(map[string][]string)[hostInfo.String()][0]), metaInfo)
	if err != nil {
		metaInfo = &MetaInfo{
			OS:       "unknown",
			Kernel:   "unknown",
			Hostname: "unknown",
			Arch:     "unknown",
			Cpu:      "unknown",
			Mem:      "unknown",
		}
	}
	// TODO 增加系统信息检测
	machine.MetaInfo = *metaInfo
	return nil
}

func (s *Service) Initialize() error {
	err := s.machinePersistence.DB.AutoMigrate(&Group{}, &Machine{})
	if err != nil {
		return err
	}
	_ = s.groupPersistence.Insert(&Group{
		ID:    1,
		Title: "默认分组",
	})
	return nil
}
