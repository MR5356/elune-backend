package executor

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/MR5356/elune-backend/pkg/domain/cron"
	"github.com/MR5356/elune-backend/pkg/domain/machine"
	"github.com/MR5356/elune-backend/pkg/domain/script"
	"github.com/MR5356/elune-backend/pkg/persistence"
	"github.com/MR5356/elune-backend/pkg/persistence/cache"
	"github.com/MR5356/elune-backend/pkg/persistence/database"
	"github.com/MR5356/jietan/pkg/executor"
	"github.com/MR5356/jietan/pkg/executor/api"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"time"
)

type Service struct {
	scriptPersistence       *persistence.Persistence[*script.Script]
	scriptTypePersistence   *persistence.Persistence[*script.Type]
	machinePersistence      *persistence.Persistence[*machine.Machine]
	machineGroupPersistence *persistence.Persistence[*machine.Group]
	recordPersistence       *persistence.Persistence[*Record]

	database *database.Database
	cache    cache.Cache

	jobMap map[string]*JobInfo
}

type JobInfo struct {
	exec      executor.Executor
	ctx       context.Context
	ctxCancel context.CancelFunc
}

func NewService(database *database.Database, cache cache.Cache) *Service {
	return &Service{
		scriptPersistence:       persistence.New(database, cache, &script.Script{}),
		scriptTypePersistence:   persistence.New(database, cache, &script.Type{}),
		machinePersistence:      persistence.New(database, cache, &machine.Machine{}),
		machineGroupPersistence: persistence.New(database, cache, &machine.Group{}),
		recordPersistence:       persistence.New(database, cache, &Record{}),

		database: database,
		cache:    cache,

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
	err := s.recordPersistence.DB.Select("id, script_title, script, host, params, status, message, error, created_at, updated_at").Order("created_at desc").Find(&res, &Record{}).Error
	for _, item := range res {
		item.Result = ""
	}
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
	recordHosts := make([]*api.HostInfo, 0)

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
		recordHosts = append(recordHosts, &api.HostInfo{
			Host:   m.HostInfo.Host,
			Port:   int(m.HostInfo.Port),
			User:   m.HostInfo.Username,
			Passwd: "******",
		})
	}

	hostsStr, _ := json.Marshal(recordHosts)

	record := &Record{
		ID:          id,
		ScriptTitle: scriptInfo.Title,
		Script:      scriptInfo.Content,
		Params:      params,
		Host:        string(hostsStr),
		Status:      "running",
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
		log := res.Data["log"].(map[string][]string)
		logStr, _ := json.Marshal(log)
		record.Result = string(logStr)
		if res.Status == 0 {
			record.Status = "finished"
		} else {
			record.Status = "failed"
		}
		record.Message = res.Message
		record.Error = string(errs)
		err = s.recordPersistence.DB.Updates(record).Error
		if err != nil {
			logrus.Errorf("更新记录入库失败1：%v \n%+v", err, record)
		}
		delete(s.jobMap, id)
	}()

	// 更新日志
	go func() {
		for {
			if job, ok := s.jobMap[id]; ok {
				log := job.exec.GetResult(api.ResultFieldLog, nil).(map[string][]string)
				logStr, _ := json.Marshal(log)
				record.Result = string(logStr)
				err = s.recordPersistence.DB.Updates(record).Error
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

func (s *Service) GetJobLog(id string) (map[string][]string, error) {
	record, err := s.recordPersistence.Detail(&Record{ID: id})
	if err != nil {
		return nil, err
	}
	log := make(map[string][]string)
	err = json.Unmarshal([]byte(record.Result), &log)
	return log, err
}

func (s *Service) Initialize() error {
	err := s.recordPersistence.DB.AutoMigrate(&Record{})
	if err != nil {
		return err
	}

	err = cron.GetTaskFactory().AddTask("script", func() cron.Task {
		return NewScriptTask(s, s.database, s.cache)
	})
	if err != nil {
		return err
	}
	return nil
}
