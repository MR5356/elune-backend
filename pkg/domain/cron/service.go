package cron

import (
	"encoding/json"
	"errors"
	"github.com/MR5356/elune-backend/pkg/persistence"
	"github.com/MR5356/elune-backend/pkg/persistence/cache"
	"github.com/MR5356/elune-backend/pkg/persistence/database"
	"github.com/MR5356/elune-backend/pkg/utils/structutil"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"strings"
	"sync"
)

const (
	topicRemoveCron = "cron.remove"
	topicAddCron    = "cron.add"
)

type Service struct {
	cronPersistence   *persistence.Persistence[*Cron]
	recordPersistence *persistence.Persistence[*Record]
	cron              *cron.Cron
	jobMap            sync.Map
	taskFactory       *TaskFactory
	cache             cache.Cache
}

func NewService(database *database.Database, cache cache.Cache) *Service {
	c := cron.New(cron.WithSeconds())
	c.Start()
	return &Service{
		cronPersistence:   persistence.New(database, cache, &Cron{}),
		recordPersistence: persistence.New(database, cache, &Record{}),
		cron:              c,
		taskFactory:       GetTaskFactory(),
		cache:             cache,
	}
}

func (s *Service) ListCron() ([]*Cron, error) {
	return s.cronPersistence.List(&Cron{})
}

func (s *Service) SetEnableCron(id uint, enable bool) error {
	c, err := s.cronPersistence.Detail(&Cron{ID: id})
	if err != nil {
		return errors.New("不存在该定时任务")
	}
	c.Enabled = enable
	if enable {
		err = s.addCron(c, c.CronString, c.TaskName, c.Params)
		if err != nil {
			return err
		}
	} else {
		s.removeCron(c.ID)
	}
	return s.cronPersistence.Update(&Cron{ID: id}, structutil.Struct2Map(c))
}

func (s *Service) AddCron(cron *Cron) error {
	if len(cron.Title) == 0 {
		return errors.New("定时任务名称不可为空")
	}
	if len(strings.Split(cron.CronString, " ")) != 6 {
		return errors.New("定时任务cron格式不正确")
	}
	if len(cron.TaskName) == 0 {
		return errors.New("定时任务执行器不可为空")
	}
	if _, err := s.taskFactory.GetTask(cron.TaskName); err != nil {
		return errors.New("定时任务执行器不存在")
	}
	cron.ID = 0
	err := s.cronPersistence.Insert(cron)
	if err != nil {
		return err
	}
	logrus.Infof("添加定时任务:%d %s %s %s", cron.ID, cron.CronString, cron.TaskName, cron.Params)
	err = s.addCron(cron, cron.CronString, cron.TaskName, cron.Params)
	if err != nil {
		_ = s.DeleteCron(cron.ID)
		return errors.New("定时策略表达式不正确")
	}
	return nil
}

func (s *Service) DeleteCron(id uint) error {
	s.removeCron(id)
	return s.cronPersistence.Delete(&Cron{ID: id})
}

func (s *Service) removeCron(cronId uint) {
	s.cache.Publish(topicRemoveCron, cronId)
}

func (s *Service) rmCronSubscriber(cronId uint) {
	logrus.Infof("rm cron task: %d", cronId)
	jobId, ok := s.jobMap.Load(cronId)
	if !ok {
		return
	}
	// 停止定时任务
	s.cron.Remove(jobId.(cron.EntryID))
	// 删除定时任务记录
	s.jobMap.Delete(cronId)
}

type addCronParams struct {
	Cron       *Cron
	CronString string
	TaskName   string
	Params     string
}

func (s *Service) addCronSubscriber(params string) error {
	logrus.Infof("add cron task: %s", params)
	ps := new(addCronParams)
	err := json.Unmarshal([]byte(params), ps)
	if err != nil {
		return err
	}

	taskFunc, err := s.taskFactory.GetTask(ps.TaskName)
	if err != nil {
		return err
	}
	f := taskFunc()
	f.SetParams(ps.Params)
	f.SetCronInfo(ps.Cron)
	jobId, err := s.cron.AddJob(ps.CronString, f)
	if err != nil {
		return err
	}
	s.jobMap.Store(ps.Cron.ID, jobId)
	return nil
}

func (s *Service) addCron(cron *Cron, cronString, taskName, params string) error {
	ps, err := json.Marshal(&addCronParams{
		Cron:       cron,
		CronString: cronString,
		TaskName:   taskName,
		Params:     params,
	})
	if err != nil {
		return err
	}
	s.cache.Publish(topicAddCron, ps)
	return nil
}

func (s *Service) PageCronRecord(pageNum, pageSize int) (*persistence.Pager[*Record], error) {
	return s.recordPersistence.Page(&Record{}, int64(pageNum), int64(pageSize))
}

func (s *Service) Initialize() error {
	err := s.cronPersistence.DB.AutoMigrate(&Cron{}, &Record{})
	if err != nil {
		return err
	}

	jobs, err := s.cronPersistence.List(&Cron{Enabled: true})
	if err != nil {
		return err
	}

	logrus.Debugf("jobs: %+v", jobs)

	for _, job := range jobs {
		logrus.Debugf("add cron job %+v", job)
		err := s.addCron(job, job.CronString, job.TaskName, job.Params)
		if err != nil {
			logrus.Errorf("add cron job error: %v", err)
		}
	}

	err = s.cache.Subscribe(topicRemoveCron, s.rmCronSubscriber)
	if err != nil {
		return err
	}

	err = s.cache.Subscribe(topicAddCron, s.addCronSubscriber)
	if err != nil {
		return err
	}
	return nil
}
