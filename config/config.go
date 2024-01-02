package config

import (
	"github.com/Trendyol/go-dcp/config"
)

type DB struct {
	//TODO string or []string?
	Host       string `yaml:"host"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	DriverName string `yaml:"driverName"`
	Port       uint16 `yaml:"port"`
}

type Connector struct {
	DB  DB         `yaml:"db" mapstructure:"db"`
	Dcp config.Dcp `yaml:",inline" mapstructure:",squash"`
}

func (c *Connector) ApplyDefaults() {
}
