package syncer

import (
	"github.com/MR5356/elune-backend/pkg/persistence"
	"github.com/MR5356/elune-backend/pkg/persistence/cache"
	"github.com/MR5356/elune-backend/pkg/persistence/database"
)

type Service struct {
	syncerPersistence *persistence.Persistence[*Syncer]
	typePersistence   *persistence.Persistence[*Type]

	database *database.Database
	cache    cache.Cache
}

func NewService(database *database.Database, cache cache.Cache) *Service {
	return &Service{
		syncerPersistence: persistence.New(database, cache, &Syncer{}),
		typePersistence:   persistence.New(database, cache, &Type{}),
	}
}

func (s *Service) Initialize() error {
	err := s.syncerPersistence.DB.AutoMigrate(&Type{}, &Syncer{})
	if err != nil {
		return err
	}

	_ = s.typePersistence.Insert(&Type{
		ID:    1,
		Title: "image",
	})
	_ = s.typePersistence.Insert(&Type{
		ID:    2,
		Title: "git",
	})

	//err = cron.GetTaskFactory().AddTask("image-sync", func() cron.Task {
	//	return task.NewImageSyncTask(s.database, s.cache)
	//})
	//if err != nil {
	//	return err
	//}
	//
	//err = cron.GetTaskFactory().AddTask("git-sync", func() cron.Task {
	//	return task.NewGitSyncTask()
	//})
	//if err != nil {
	//	return err
	//}

	return nil
}
