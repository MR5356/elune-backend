package cron

import (
	"errors"
	"github.com/sirupsen/logrus"
	"sync"
)

const (
	TaskRunning  = "running"
	TaskFinished = "finished"
	TaskFailed   = "failed"
)

var (
	taskFactory *TaskFactory
	once        sync.Once
)

type Task interface {
	Run()
	SetParams(params string)
	SetCronInfo(cron *Cron)
}

type TaskFactory struct {
	tasks sync.Map
}

func GetTaskFactory() *TaskFactory {
	once.Do(func() {
		taskFactory = &TaskFactory{
			tasks: sync.Map{},
		}
	})
	return taskFactory
}

func (tf *TaskFactory) GetTask(taskName string) (func() Task, error) {
	if f, ok := tf.tasks.Load(taskName); ok {
		return f.(func() Task), nil
	}
	return nil, errors.New("执行器不存在")
}

func (tf *TaskFactory) AddTask(taskName string, f func() Task) error {
	if _, ok := tf.tasks.Load(taskName); ok {
		return errors.New("执行器已经存在")
	}
	tf.tasks.Store(taskName, f)
	logrus.Infof("register task %s", taskName)
	return nil
}
