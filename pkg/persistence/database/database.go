package database

import (
	"encoding/json"
	"errors"
	"github.com/MR5356/elune-backend/pkg/config"
	"github.com/glebarez/sqlite"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	*gorm.DB
}

var (
	DBDriverNotSupport = errors.New("database driver not support")
)

func New(cfg *config.Config) (database *Database, err error) {
	var driver gorm.Dialector
	logrus.Debugf("database driver: %s", cfg.Persistence.Database.Driver)
	switch cfg.Persistence.Database.Driver {
	case "sqlite":
		driver = sqlite.Open(cfg.Persistence.Database.DSN)
	case "mysql":
		driver = mysql.Open(cfg.Persistence.Database.DSN)
	case "postgres":
		driver = postgres.Open(cfg.Persistence.Database.DSN)
	default:
		return nil, DBDriverNotSupport
	}

	var dbLogLevel = logger.Error
	if cfg.Server.Debug {
		dbLogLevel = logger.Info
	}
	logrus.Debugf("database log level: %+v", dbLogLevel)

	client, err := gorm.Open(driver, &gorm.Config{
		Logger: logger.Default.LogMode(dbLogLevel),
	})
	if err != nil {
		return nil, err
	}

	db, err := client.DB()
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(cfg.Persistence.Database.MaxIdleConn)
	db.SetMaxOpenConns(cfg.Persistence.Database.MaxOpenConn)
	db.SetConnMaxLifetime(cfg.Persistence.Database.ConnMaxLift)
	db.SetConnMaxIdleTime(cfg.Persistence.Database.ConnMaxIdle)

	dbStat, _ := json.Marshal(db.Stats())
	logrus.Debugf("database stats: %s", dbStat)
	return &Database{client}, nil
}
