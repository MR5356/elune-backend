package authentication

import (
	"github.com/MR5356/elune-backend/pkg/utils/fileutil"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"os"
)

type RBACService struct {
	enforcer *casbin.Enforcer
	db       *gorm.DB
}

const defaultRBACModel = "[request_definition]\nr = sub, obj, act\n\n[policy_definition]\np = sub, obj, act\n\n[role_definition]\ng = _, _\n\n[matchers]\nm = g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && keyMatch(r.act, p.act)\n\n[policy_effect]\ne = some(where (p.eft == allow))"

func NewRBACService(db *gorm.DB) (*RBACService, error) {
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}

	// 判断模型是否存在，如果不存在则创建一个
	modelPath := "config/rbac_model.conf"
	_, err = os.Stat(modelPath)
	if os.IsNotExist(err) {
		err := os.MkdirAll("config", os.ModePerm)
		if err != nil {
			return nil, err
		}
		err = fileutil.WriteToFile(modelPath, []byte(defaultRBACModel))
		if err != nil {
			return nil, err
		}
	}
	enforcer, err := casbin.NewEnforcer(modelPath, adapter)
	if err != nil {
		return nil, err
	}

	return &RBACService{
		enforcer: enforcer,
		db:       db,
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
	s.db.Exec("DELETE FROM casbin_rule")
	err := s.enforcer.LoadPolicy()
	if err != nil {
		return err
	}
	// 默认角色
	_, _ = s.enforcer.AddRoleForUser("admin", "administrators")
	_, _ = s.enforcer.AddRoleForUser("devops", "devops")
	_, _ = s.enforcer.AddRoleForUser("guest", "users")
	_, _ = s.enforcer.AddRoleForUser("unknown", "guests")

	// 默认权限
	policies := [][]string{
		{
			"administrators", "*", "*",
		},
		{
			"devops", "*", "GET",
		},
		{
			"devops", "/user/*", "*",
		},
		{
			"devops", "/navigation/*", "*",
		},
		{
			"devops", "/script/*", "*",
		},
		{
			"devops", "/machine/*", "*",
		},
		{
			"devops", "/execute/*", "*",
		},
		{
			"devops", "/cron/*", "*",
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
			"guests", "/user/*", "*",
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
