package executor

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/MR5356/elune-backend/pkg/domain/cron"
	"github.com/MR5356/elune-backend/pkg/persistence"
	"github.com/MR5356/elune-backend/pkg/persistence/cache"
	"github.com/MR5356/elune-backend/pkg/persistence/database"
	"github.com/sirupsen/logrus"
	"time"
)

type Task struct {
	cronRecordPersistence *persistence.Persistence[*cron.Record]
	service               *Service

	cache    cache.Cache
	params   string
	cronInfo *cron.Cron
}

func NewScriptTask(service *Service, db *database.Database, cc cache.Cache) *Task {
	return &Task{
		cronRecordPersistence: persistence.New(db, cc, &cron.Record{}),
		cache:                 cc,
		service:               service,
	}
}

type Params struct {
	ScriptID       uint   `json:"scriptId"`
	Params         string `json:"params"`
	MachineGroupId uint   `json:"machineGroupId"`
	MachineIds     []uint `json:"machineIds"`
}

func (t *Task) Run() {
	params := new(Params)
	err := json.Unmarshal([]byte(t.params), params)
	if err != nil {
		logrus.Errorf("unmarshal params %s error: %v", t.params, err)
		return
	}

	uniqueKey := fmt.Sprintf("script-%x", md5.Sum([]byte(t.params)))
	err = t.cache.TryLock(uniqueKey)
	if err != nil {
		return
	}
	defer func(cache cache.Cache, key string) {
		err := cache.Unlock(key)
		if err != nil {
			logrus.Errorf("unlock key %s error: %v", key, err)
		}
	}(t.cache, uniqueKey)
	logrus.Infof("script task run")
	var id string
	if params.MachineGroupId != 0 {
		id, err = t.service.StartNewJobWithMachineGroup(params.ScriptID, params.MachineGroupId, params.Params)
	} else {
		id, err = t.service.StartNewJob(params.ScriptID, params.MachineIds, params.Params)
	}
	if err != nil {
		logrus.Errorf("start new job error: %v", err)
	}
	record := &cron.Record{
		CronId:   t.cronInfo.ID,
		Title:    t.cronInfo.Title,
		Desc:     t.cronInfo.Desc,
		TaskName: t.cronInfo.TaskName,
		Params:   t.params,
		Status:   cron.TaskRunning,
		Log:      "",
	}
	err = t.cronRecordPersistence.Insert(record)
	if err != nil {
		logrus.Errorf("insert record error: %v", err)
		return
	}
	go func() {
		log, err := t.service.GetJobLog(id)
		if err != nil {
			logrus.Errorf("get job log error: %v", err)
		} else {
			logStr, _ := json.Marshal(log)
			record.Log = string(logStr)
			err = t.cronRecordPersistence.DB.Updates(record).Error
			if err != nil {
				logrus.Warnf("update record error: %v", err)
			}
		}
		time.Sleep(time.Second)
	}()
}

func (t *Task) SetParams(params string) {
	t.params = params
}

func (t *Task) SetCronInfo(ci *cron.Cron) {
	t.cronInfo = ci
}
