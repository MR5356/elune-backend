package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/MR5356/elune-backend/pkg/config"
	"github.com/MR5356/elune-backend/pkg/controller"
	"github.com/MR5356/elune-backend/pkg/domain/authentication"
	"github.com/MR5356/elune-backend/pkg/domain/blog"
	"github.com/MR5356/elune-backend/pkg/domain/navigation"
	"github.com/MR5356/elune-backend/pkg/domain/site"
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

func New(config *config.Config) (server *Server, err error) {
	if config.Server.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(
		middleware.Record(),
		middleware.CORS(),
		gzip.Gzip(gzip.DefaultCompression),
		middleware.Recovery(),
	)

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
	rbacService, err := authentication.NewRBACService(db.DB)
	if err != nil {
		return nil, err
	}
	// jwt
	jwtService := authentication.NewJWTService(config.Server.Secret, config.Server.Issuer, config.Server.Expire)

	siteService := site.NewService(db, cc)
	navigationService := navigation.NewService(db, cc)
	userService := authentication.NewService(db, cc)
	//
	//selfKubeconfig, _ := siteService.GetKey("kubeconfig")
	//kubernetesService := kubernetes.NewService(selfKubeconfig)

	services := []service.Service{
		siteService,
		navigationService,
		//kubernetesService,
		rbacService,
		jwtService,
		userService,
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
		//kubernetes.NewController(kubernetesService),
		authentication.NewController(rbacService, jwtService, userService),
		blog.NewController(),
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
