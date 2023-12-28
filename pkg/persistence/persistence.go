package persistence

import (
	"github.com/MR5356/elune-backend/pkg/persistence/cache"
	"github.com/MR5356/elune-backend/pkg/persistence/database"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Persistence[T any] struct {
	DB    *database.Database
	Cache cache.Cache
}

func New[T any](database *database.Database, cache cache.Cache, model T) *Persistence[T] {
	return &Persistence[T]{
		DB:    database,
		Cache: cache,
	}
}

type Pager[T any] struct {
	CurrentPage int64 `json:"currentPage"`
	PageSize    int64 `json:"pageSize"`
	Total       int64 `json:"total"`
	Data        []T   `json:"data"`
}

func (s *Persistence[T]) Insert(entity T) error {
	logrus.Debugf("Insert: %+v", entity)
	return s.DB.Create(entity).Error
}

func (s *Persistence[T]) Delete(entity T) error {
	logrus.Debugf("Delete: %+v", entity)
	return s.DB.Delete(entity).Error
}

func (s *Persistence[T]) Update(entity T, fields map[string]interface{}) error {
	logrus.Debugf("Update: %+v", entity)
	return s.DB.Model(entity).Where(entity).Updates(fields).Error
}

func (s *Persistence[T]) Detail(entity T) (res T, err error) {
	logrus.Debugf("Detail: %+v", entity)
	err = s.DB.First(&res, entity).Error
	return
}

func (s *Persistence[T]) List(entity T) (res []T, err error) {
	logrus.Debugf("List: %+v", entity)
	err = s.DB.Order("updated_at desc").Find(&res, entity).Error
	return
}

func (s *Persistence[T]) Page(entity T, page, size int64) (res *Pager[T], err error) {
	logrus.Debugf("Page: %+v", entity)
	res = new(Pager[T])
	res.CurrentPage = page
	res.PageSize = size
	s.DB.Model(&entity).Where(entity).Count(&res.Total)
	if res.Total == 0 {
		res.Data = make([]T, 0)
	}
	err = s.DB.Model(&entity).Order("updated_at desc").Where(entity).Scopes(Pagination(res)).Find(&res.Data).Error
	return res, err
}

func Pagination[T any](pager *Pager[T]) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		size := pager.PageSize
		page := pager.CurrentPage
		offset := int((page - 1) * size)
		return db.Offset(offset).Limit(int(size))
	}
}
