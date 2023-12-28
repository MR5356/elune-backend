package main

import (
	"github.com/MR5356/elune-backend/pkg/config"
	_ "github.com/MR5356/elune-backend/pkg/log"
	"github.com/MR5356/elune-backend/pkg/server"
	"github.com/MR5356/elune-backend/pkg/utils/structutil"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
)

func main() {
	var withs []config.Cfg

	// 是否开启debug
	if os.Getenv(config.EluneEnvDebug) == "true" {
		logrus.SetLevel(logrus.DebugLevel)
		withs = append(withs, config.WithDebug(true))
	}

	// 是否自定义端口
	if len(os.Getenv(config.EluneEnvPort)) > 0 {
		port, err := strconv.Atoi(os.Getenv(config.EluneEnvPort))
		if err == nil {
			withs = append(withs, config.WithPort(port))
		} else {
			logrus.Warnf("invalid port: %s", os.Getenv(config.EluneEnvPort))
		}
	}

	// 是否定义数据库
	if driver := os.Getenv(config.EluneEnvDatabaseDriver); len(driver) > 0 {
		withs = append(withs, config.WithDatabaseDriver(driver))
	}

	if dsn := os.Getenv(config.EluneEnvDatabaseDSN); len(dsn) > 0 {
		withs = append(withs, config.WithDatabaseDsn(dsn))
	}

	// 是否定义Cache
	if driver := os.Getenv(config.EluneEnvCacheDriver); len(driver) > 0 {
		withs = append(withs, config.WithCacheDriver(driver))
	}

	if dsn := os.Getenv(config.EluneEnvCacheDSN); len(dsn) > 0 {
		withs = append(withs, config.WithCacheDsn(dsn))
	}

	cfg := config.New(withs...)
	logrus.Debugf("run with config: %+v", structutil.Struct2String(cfg))

	srv, err := server.New(cfg)
	if err != nil {
		logrus.Fatalf("create server error: %s", err)
	}

	if err := srv.Run(); err != nil {
		logrus.Fatalf("run server error: %s", err)
	}
}
