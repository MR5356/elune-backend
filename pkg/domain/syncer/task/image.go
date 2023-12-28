package task

import (
	"github.com/MR5356/elune-backend/pkg/domain/cron"
	"github.com/MR5356/elune-backend/pkg/persistence"
	"github.com/MR5356/elune-backend/pkg/persistence/cache"
	"github.com/MR5356/elune-backend/pkg/persistence/database"
	"github.com/sirupsen/logrus"
)

type ImageSyncTask struct {
	cronRecordPersistence *persistence.Persistence[*cron.Record]
	cache                 cache.Cache

	params string
}

func NewImageSyncTask(db *database.Database, cc cache.Cache) *ImageSyncTask {
	return &ImageSyncTask{
		cronRecordPersistence: persistence.New(db, cc, &cron.Record{}),
		cache:                 cc,
	}
}

func (t *ImageSyncTask) Run() {
	err := t.cache.TryLock("image-sync-" + t.params)
	if err != nil {
		return
	}
	defer func(cache cache.Cache, key string) {
		err := cache.Unlock(key)
		if err != nil {
			logrus.Errorf("unlock key %s error: %v", key, err)
		}
	}(t.cache, "image-sync-"+t.params)
	logrus.Infof("image-sync task run")
}

func (t *ImageSyncTask) SetParams(params string) {
	t.params = params
}
