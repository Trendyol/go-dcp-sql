package config

import (
	"time"

	"github.com/Trendyol/go-dcp/config"
)

type SQL struct {
	Host                string        `yaml:"host"`
	User                string        `yaml:"user"`
	Password            string        `yaml:"password"`
	DBName              string        `yaml:"dbName"`
	Port                uint16        `yaml:"port"`
	SslMode             string        `yaml:"sslMode"`
	DriverName          string        `yaml:"driverName"`
	BatchSizeLimit      int           `yaml:"batchSizeLimit"`
	BatchTickerDuration time.Duration `yaml:"batchTickerDuration"`
}

type Connector struct {
	SQL SQL        `yaml:"sql" mapstructure:"sql"`
	Dcp config.Dcp `yaml:",inline" mapstructure:",squash"`
}

func (c *Connector) ApplyDefaults() {
	if c.SQL.SslMode == "" {
		c.SQL.SslMode = "disable"
	}
}
