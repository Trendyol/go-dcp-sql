package config

import (
	"github.com/Trendyol/go-dcp/config"
)

type Sql struct {
	Host           string `yaml:"host"`
	User           string `yaml:"user"`
	Password       string `yaml:"password"`
	DbName         string `yaml:"dbName"`
	Port           uint16 `yaml:"port"`
	SslMode        string `yaml:"sslMode"`
	DriverName     string `yaml:"driverName"`
	BatchSizeLimit int    `yaml:"batchSizeLimit"`
}

type Connector struct {
	Sql Sql        `yaml:"sql" mapstructure:"sql"`
	Dcp config.Dcp `yaml:",inline" mapstructure:",squash"`
}

func (c *Connector) ApplyDefaults() {
	if c.Sql.SslMode == "" {
		c.Sql.SslMode = "disable"
	}
}
