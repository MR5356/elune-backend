package cron

import (
	"github.com/MR5356/elune-backend/pkg/persistence"
	"github.com/MR5356/elune-backend/pkg/persistence/cache"
	"github.com/MR5356/elune-backend/pkg/persistence/database"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"sync"
)

type Service struct {
	cronPersistence   *persistence.Persistence[*Cron]
	recordPersistence *persistence.Persistence[*Record]
	cron              *cron.Cron
	jobMap            sync.Map
	taskFactory       *TaskFactory
}

func NewService(database *database.Database, cache cache.Cache) *Service {
	c := cron.New(cron.WithSeconds())
	c.Start()
	return &Service{
		cronPersistence:   persistence.New(database, cache, &Cron{}),
		recordPersistence: persistence.New(database, cache, &Record{}),
		cron:              c,
		taskFactory:       GetTaskFactory(),
	}
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

	err = s.taskFactory.AddTask("test", &TestTask{})
	if err != nil {
		return err
	}

	logrus.Debugf("jobs: %+v", jobs)

	for _, job := range jobs {
		logrus.Debugf("add cron job %+v", job)
		task, err := s.taskFactory.GetTask(job.TaskName)
		if err != nil {
			logrus.Errorf("get task error: %v", err)
			continue
		}
		jobId, err := s.cron.AddJob(job.CronString, task)
		if err != nil {
			logrus.Errorf("add cron job error: %v", err)
		} else {
			s.jobMap.Store(job.ID, jobId)
			logrus.Infof("add cron job %s with params: %s", job.Title, job.Params)
		}
	}

	return nil
}
