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

func (tf *TaskFactory) GetTask(taskName string) (Task, error) {
	if task, ok := tf.tasks.Load(taskName); ok {
		return task.(Task), nil
	}
	return nil, errors.New("task not found")
}

func (tf *TaskFactory) AddTask(taskName string, task Task) error {
	if _, ok := tf.tasks.Load(taskName); ok {
		return errors.New("task already exists")
	}
	tf.tasks.Store(taskName, task)
	logrus.Infof("register task %s", taskName)
	return nil
}

type TestTask struct {
}

func (t *TestTask) Run() {
	logrus.Infof("test task run")
}

func (t *TestTask) SetParams(params string) {

}
