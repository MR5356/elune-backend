package authentication

import (
	"github.com/MR5356/elune-backend/pkg/persistence"
	"github.com/MR5356/elune-backend/pkg/persistence/cache"
	"github.com/MR5356/elune-backend/pkg/persistence/database"
)

type Service struct {
	userPersistence         *persistence.Persistence[*User]
	groupPersistence        *persistence.Persistence[*Group]
	rolePersistence         *persistence.Persistence[*Role]
	roleRelationPersistence *persistence.Persistence[*RoleRelation]
	userGroupPersistence    *persistence.Persistence[*UserGroup]
}

func NewService(database *database.Database, cache cache.Cache) *Service {
	return &Service{
		userPersistence:         persistence.New(database, cache, &User{}),
		groupPersistence:        persistence.New(database, cache, &Group{}),
		rolePersistence:         persistence.New(database, cache, &Role{}),
		roleRelationPersistence: persistence.New(database, cache, &RoleRelation{}),
		userGroupPersistence:    persistence.New(database, cache, &UserGroup{}),
	}
}

func (s *Service) Initialize() error {
	err := s.userPersistence.DB.AutoMigrate(&User{}, &Group{}, &UserGroup{}, &Role{}, &RoleRelation{})
	if err != nil {
		return err
	}

	defaultUsers := []*User{
		{
			ID:       1,
			Username: "admin",
			Nickname: "系统管理员",
			Password: "admin",
			Email:    "admin@example.com",
		},
		{
			ID:       2,
			Username: "devops",
			Nickname: "运维工程师",
			Password: "devops",
			Email:    "devops@example.com",
		},
	}

	defaultGroups := []*Group{
		{
			ID:    1,
			Title: "administrators",
			Desc:  "系统管理员组",
		},
		{
			ID:    2,
			Title: "devops",
			Desc:  "运维工程师组",
		},
		{
			ID:    3,
			Title: "users",
			Desc:  "普通用户",
		},
		{
			ID:    4,
			Title: "guests",
			Desc:  "游客",
		},
	}

	defaultUserGroups := []*UserGroup{
		{
			ID:      1,
			UserID:  1,
			GroupID: 1,
		},
		{
			ID:      2,
			UserID:  2,
			GroupID: 2,
		},
	}

	for _, user := range defaultUsers {
		_ = s.userPersistence.Insert(user)
	}

	for _, group := range defaultGroups {
		_ = s.groupPersistence.Insert(group)
	}

	for _, ug := range defaultUserGroups {
		_ = s.userGroupPersistence.Insert(ug)
	}
	return nil
}
