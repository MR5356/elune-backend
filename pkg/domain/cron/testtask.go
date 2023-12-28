package cron

import "github.com/sirupsen/logrus"

type TestTask struct {
}

func (t *TestTask) Run() {
	logrus.Infof("test task run")
}

func (t *TestTask) SetParams(params string) {

}

func (t *TestTask) SetCronInfo(cron *Cron) {

}

func init() {
	err := GetTaskFactory().AddTask("test", func() Task {
		return &TestTask{}
	})
	if err != nil {
		logrus.Errorf("register test task error: %v", err)
	}
}
