package executor

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/MR5356/elune-backend/pkg/domain/machine"
	"github.com/MR5356/elune-backend/pkg/domain/script"
	"github.com/MR5356/elune-backend/pkg/persistence"
	"github.com/MR5356/elune-backend/pkg/persistence/cache"
	"github.com/MR5356/elune-backend/pkg/persistence/database"
	"github.com/MR5356/jietan/pkg/executor"
	"github.com/MR5356/jietan/pkg/executor/api"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

type Service struct {
	scriptPersistence       *persistence.Persistence[*script.Script]
	scriptTypePersistence   *persistence.Persistence[*script.Type]
	machinePersistence      *persistence.Persistence[*machine.Machine]
	machineGroupPersistence *persistence.Persistence[*machine.Group]
	recordPersistence       *persistence.Persistence[*Record]

	jobMap map[string]*JobInfo
}

type JobInfo struct {
	exec      executor.Executor
	ctx       context.Context
	ctxCancel context.CancelFunc
}

func NewService(database *database.Database, cache *cache.Cache) *Service {
	return &Service{
		scriptPersistence:       persistence.New(database, cache, &script.Script{}),
		scriptTypePersistence:   persistence.New(database, cache, &script.Type{}),
		machinePersistence:      persistence.New(database, cache, &machine.Machine{}),
		machineGroupPersistence: persistence.New(database, cache, &machine.Group{}),
		recordPersistence:       persistence.New(database, cache, &Record{}),

		jobMap: make(map[string]*JobInfo),
	}
}

func (s *Service) StartNewJobWithMachineGroup(scriptId, machineGroupId uint, params string) (id string, err error) {
	list, err := s.machinePersistence.List(&machine.Machine{GroupId: machineGroupId})
	if err != nil {
		return "", err
	}

	ms := make([]uint, 0)
	for _, item := range list {
		ms = append(ms, item.ID)
	}

	return s.StartNewJob(scriptId, ms, params)
}

func (s *Service) ListJob() ([]*Record, error) {
	res := make([]*Record, 0)
	err := s.recordPersistence.DB.Order("created_at desc").Find(&res, &Record{}).Error
	return res, err
}

func (s *Service) StartNewJob(scriptId uint, machineId []uint, params string) (id string, err error) {
	id = uuid.NewString()
	ctx, cancel := context.WithCancel(context.Background())
	exec := executor.GetExecutor("remote")
	s.jobMap[id] = &JobInfo{
		exec:      exec,
		ctx:       ctx,
		ctxCancel: cancel,
	}

	scriptInfo, err := s.scriptPersistence.Detail(&script.Script{ID: scriptId})
	if err != nil {
		return "", err
	}

	hosts := make([]*api.HostInfo, 0)
	logrus.Debugf("machineId: %+v", machineId)
	if len(machineId) == 0 {
		return "", errors.New("至少包含一台主机")
	}
	for _, item := range machineId {
		m, err := s.machinePersistence.Detail(&machine.Machine{ID: item})
		if err != nil {
			return "", errors.New("机器不存在")
		}
		hosts = append(hosts, &api.HostInfo{
			Host:   m.HostInfo.Host,
			Port:   int(m.HostInfo.Port),
			User:   m.HostInfo.Username,
			Passwd: m.HostInfo.Password,
		})
	}

	hostsStr, _ := json.Marshal(hosts)

	record := &Record{
		ID:          id,
		ScriptTitle: scriptInfo.Title,
		Script:      scriptInfo.Content,
		Params:      params,
		Host:        string(hostsStr),
		Status:      "RUNNING",
	}

	logrus.Debugf("record: %+v", record)
	// 记录入库
	err = s.recordPersistence.Insert(record)

	if err != nil {
		logrus.Errorf("记录入库失败：%v \n%+v", err, record)
	}

	// 主进程
	go func() {
		res := exec.Execute(ctx, &api.ExecuteParams{
			"hosts":  hosts,
			"script": scriptInfo.Content,
			"params": params,
		})
		errs, _ := json.Marshal(res.Data["error"])
		record.Result = strings.Join(res.Data["log"].([]string), "\n")
		record.Status = strconv.Itoa(int(res.Status))
		record.Message = res.Message
		record.Error = string(errs)
		err := s.recordPersistence.DB.Updates(record)
		//err := s.recordPersistence.Update(&Record{ID: record.ID}, structutil.Struct2Map(&Record{
		//	Result:  strings.Join(res.Data["log"].([]string), "\n"),
		//	Status:  strconv.Itoa(int(res.Status)),
		//	Message: res.Message,
		//	Error:   string(errs),
		//}))
		if err != nil {
			logrus.Errorf("更新记录入库失败：%v \n%+v", err, record)
		}
		delete(s.jobMap, id)
	}()

	// 更新日志
	go func() {
		for {
			if job, ok := s.jobMap[id]; ok {
				record.Result = strings.Join(job.exec.GetResult(api.ResultFieldLog, nil).([]string), "\n")
				err := s.recordPersistence.DB.Updates(record)
				//err := s.recordPersistence.Update(&Record{ID: record.ID}, structutil.Struct2Map(&Record{
				//	Result: strings.Join(job.exec.GetResult(api.ResultFieldLog, nil).([]string), "\n"),
				//}))
				if err != nil {
					logrus.Errorf("更新记录入库失败：%v \n%+v", err, record)
				}
			}
			time.Sleep(time.Second)
		}
	}()
	return id, nil
}

func (s *Service) StopJob(id string) error {
	jobInfo, ok := s.jobMap[id]
	if !ok {
		return nil
	}
	jobInfo.ctxCancel()
	delete(s.jobMap, id)
	return nil
}

func (s *Service) GetJobLog(id string) ([]string, error) {
	record, err := s.recordPersistence.Detail(&Record{ID: id})
	if err != nil {
		return nil, err
	}
	return strings.Split(record.Result, "\n"), nil
}

func (s *Service) Initialize() error {
	err := s.recordPersistence.DB.AutoMigrate(&Record{})
	if err != nil {
		return err
	}
	return nil
}
