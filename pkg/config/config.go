package config

import (
	"github.com/mcuadros/go-defaults"
	"time"
)

const (
	EluneEnvDebug = "ELUNE_DEBUG"
	EluneEnvPort  = "ELUNE_PORT"
)

type Config struct {
	Server      Server      `json:"server" yaml:"server"`
	Persistence Persistence `json:"persistence" yaml:"persistence"`
}

func New(cfgs ...Cfg) *Config {
	config := new(Config)
	defaults.SetDefaults(config)

	for _, cfg := range cfgs {
		cfg(config)
	}
	return config
}

type Server struct {
	Port        int    `json:"port" yaml:"port" default:"5678"`
	Prefix      string `json:"prefix" yaml:"prefix" default:"/api/v1"`
	Debug       bool   `json:"debug" yaml:"debug" default:"false"`
	GracePeriod int    `json:"gracePeriod" yaml:"gracePeriod" default:"30"`
}

type Persistence struct {
	Database Database `json:"database" yaml:"database"`
	Cache    Cache    `json:"cache" yaml:"cache"`
}

type Database struct {
	Driver string `json:"driver" yaml:"driver" default:"sqlite"`
	DSN    string `json:"dsn" yaml:"dsn" default:"db.sqlite"`

	MaxIdleConn int           `json:"maxIdleConn" yaml:"maxIdleConn" default:"10"`
	MaxOpenConn int           `json:"maxOpenConn" yaml:"maxOpenConn" default:"40"`
	ConnMaxLift time.Duration `json:"connMaxLift" yaml:"connMaxLift" default:"0s"`
	ConnMaxIdle time.Duration `json:"connMaxIdle" yaml:"connMaxIdle" default:"0s"`
}

type Cache struct {
	Driver string `json:"driver" yaml:"driver" default:"memory"`
	DSN    string `json:"dsn" yaml:"dsn" default:""`
}

func WithDebug(debug bool) Cfg {
	return func(c *Config) {
		c.Server.Debug = debug
	}
}

func WithPort(port int) Cfg {
	return func(c *Config) {
		c.Server.Port = port
	}
}

type Cfg func(c *Config)
