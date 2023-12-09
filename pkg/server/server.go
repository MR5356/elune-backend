package server

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"github.com/MR5356/elune-backend/pkg/config"
	"github.com/MR5356/elune-backend/pkg/controller"
	"github.com/MR5356/elune-backend/pkg/domain/authentication"
	"github.com/MR5356/elune-backend/pkg/domain/blog"
	"github.com/MR5356/elune-backend/pkg/domain/cron"
	"github.com/MR5356/elune-backend/pkg/domain/executor"
	"github.com/MR5356/elune-backend/pkg/domain/kubernetes"
	"github.com/MR5356/elune-backend/pkg/domain/machine"
	"github.com/MR5356/elune-backend/pkg/domain/navigation"
	"github.com/MR5356/elune-backend/pkg/domain/notify"
	"github.com/MR5356/elune-backend/pkg/domain/script"
	"github.com/MR5356/elune-backend/pkg/domain/site"
	"github.com/MR5356/elune-backend/pkg/domain/syncer"
	"github.com/MR5356/elune-backend/pkg/middleware"
	"github.com/MR5356/elune-backend/pkg/persistence/cache"
	"github.com/MR5356/elune-backend/pkg/persistence/database"
	"github.com/MR5356/elune-backend/pkg/response"
	"github.com/MR5356/elune-backend/pkg/service"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	engine *gin.Engine
	config *config.Config
}

//go:embed static
var fs embed.FS

func New(config *config.Config) (server *Server, err error) {
	if config.Server.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// 准备目录
	for name, path := range config.Server.RuntimeDirectories {
		logrus.Infof("准备 %s 目录: %s", name, path)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			err = os.MkdirAll(path, os.ModePerm)
			if err != nil {
				return nil, err
			}
		}
	}

	// jwt
	jwtService := authentication.NewJWTService(config.Server.Secret, config.Server.Issuer, config.Server.Expire)

	engine := gin.New()
	engine.MaxMultipartMemory = 8 << 20
	engine.Use(
		middleware.Record(jwtService),
		middleware.CORS(),
		gzip.Gzip(gzip.DefaultCompression),
		middleware.Recovery(),
	)

	// 前端代理接口
	engine.Use(middleware.Static("/", middleware.NewStaticFileSystem(fs, "static")))

	// 后端接口
	api := engine.Group(config.Server.Prefix)

	api.GET("/health", func(c *gin.Context) {
		response.Success(c, nil)
	})

	engine.NoRoute(func(c *gin.Context) {
		response.New(c, http.StatusNotFound, response.CodeNotFound, response.MsgNotFound, nil)
	})

	// 数据库驱动
	db, err := database.New(config)
	if err != nil {
		return nil, err
	}
	logrus.Debugf("database: %+v", db)

	// 缓存驱动
	cc, err := cache.New(config)
	if err != nil {
		return nil, err
	}
	logrus.Debugf("cache: %+v", cc)

	// rbac
	userService := authentication.NewService(db, cc)
	rbacService, err := authentication.NewRBACService(db.DB, userService)
	if err != nil {
		return nil, err
	}

	siteService := site.NewService(db, cc)
	navigationService := navigation.NewService(db, cc)
	scriptService := script.NewService(db, cc)
	machineService := machine.NewService(db, cc)
	execService := executor.NewService(db, cc)
	syncerService := syncer.NewService(db, cc)
	cronService := cron.NewService(db, cc)
	kubernetesService := kubernetes.NewService(db, cc)
	notifyService := notify.NewService(db, cc, config)

	services := []service.Service{
		siteService,
		navigationService,
		userService,
		rbacService,
		jwtService,
		scriptService,
		machineService,
		execService,
		syncerService,
		cronService,
		kubernetesService,
		notifyService,
	}
	for _, srv := range services {
		err := srv.Initialize()
		if err != nil {
			return nil, err
		}
	}

	api.Use(middleware.Authentication(config, rbacService, jwtService))

	controllers := []controller.Controller{
		site.NewController(siteService),
		navigation.NewController(navigationService),
		authentication.NewController(rbacService, jwtService, userService, config),
		blog.NewController(),
		script.NewController(scriptService),
		machine.NewController(machineService),
		executor.NewController(execService),
		syncer.NewController(syncerService),
		cron.NewController(cronService),
		kubernetes.NewController(kubernetesService),
		notify.NewController(notifyService),
	}
	for _, ctrl := range controllers {
		ctrl.RegisterRoute(api)
	}

	server = &Server{
		engine: engine,
		config: config,
	}

	return
}

func (s *Server) Run() error {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.Server.Port),
		Handler: s.engine,
	}

	go func() {
		logrus.Infof("start server on port %d", s.config.Server.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Fatalf("start server error: %s", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.config.Server.GracePeriod)*time.Second)
	defer cancel()

	ch := <-sig
	logrus.Infof("receive signal: %s", ch)
	return server.Shutdown(ctx)
}
