package authentication

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RBACService struct {
	enforcer *casbin.Enforcer
}

func NewRBACService(db *gorm.DB) (*RBACService, error) {
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}

	enforcer, err := casbin.NewEnforcer("config/rbac_model.conf", adapter)
	if err != nil {
		return nil, err
	}

	return &RBACService{
		enforcer: enforcer,
	}, nil
}

func (s *RBACService) GetRoles() []string {
	return s.enforcer.GetAllRoles()
}

func (s *RBACService) GetObjects() []string {
	return s.enforcer.GetAllObjects()
}

func (s *RBACService) GetActions() []string {
	return s.enforcer.GetAllActions()
}

func (s *RBACService) GetUsersForRole(role string) ([]string, error) {
	return s.enforcer.GetRoleManager().GetUsers(role)
}

func (s *RBACService) GetRolesForUser(user string) ([]string, error) {
	return s.enforcer.GetRoleManager().GetRoles(user)
}

func (s *RBACService) HasRoleForUser(user, obj, role string) (bool, error) {
	logrus.Debugf("HasRoleForUser: %s, %s, %s", user, obj, role)
	return s.enforcer.Enforce(user, obj, role)
}

func (s *RBACService) Initialize() error {
	err := s.enforcer.LoadPolicy()
	if err != nil {
		return err
	}

	// 默认角色
	_, _ = s.enforcer.AddRoleForUser("admin", "administrators")
	_, _ = s.enforcer.AddRoleForUser("guest", "users")
	_, _ = s.enforcer.AddRoleForUser("unknown", "guests")

	// 默认权限
	policies := [][]string{
		{
			"administrators", "*", "*",
		},
		{
			"users", "*", "GET",
		},
		{
			"guests", "*", "GET",
		},
		{
			"users", "/user/*", "*",
		},
		{
			"users", "/navigation/*", "*",
		},
	}

	for _, policy := range policies {
		_, _ = s.enforcer.AddPolicy(policy)
	}

	return nil
}
