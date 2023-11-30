package cron

import (
	"errors"
	"github.com/sirupsen/logrus"
	"sync"
)

var (
	taskFactory *TaskFactory
	once        sync.Once
)

type Task interface {
	Run()
	SetParams(params string)
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
	return nil, errors.New("taskFunc not found")
}

func (tf *TaskFactory) AddTask(taskName string, f func() Task) error {
	if _, ok := tf.tasks.Load(taskName); ok {
		return errors.New("task already exists")
	}
	tf.tasks.Store(taskName, f)
	logrus.Infof("register task %s", taskName)
	return nil
}
