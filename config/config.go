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
	SslMode             string        `yaml:"sslMode"`
	DriverName          string        `yaml:"driverName"`
	BatchTickerDuration time.Duration `yaml:"batchTickerDuration"`
	Port                uint16        `yaml:"port"`
}

type CollectionTableMapping struct {
	Collection      string `yaml:"collection"`
	TableName       string `yaml:"tableName"`
	KeyColumnName   string `yaml:"keyColumnName"`
	ValueColumnName string `yaml:"valueColumnName"`
	Audit           struct {
		Enabled             bool   `yaml:"enabled"`
		CreatedAtColumnName string `yaml:"createdAtColumnName"`
		UpdatedAtColumnName string `yaml:"updatedAtColumnName"`
	} `yaml:"audit,omitempty"`
}

type Connector struct {
	SQL                    SQL                      `yaml:"sql" mapstructure:"sql"`
	CollectionTableMapping []CollectionTableMapping `yaml:"collectionTableMapping,omitempty"`
	Dcp                    config.Dcp               `yaml:",inline" mapstructure:",squash"`
}

func (c *Connector) ApplyDefaults() {
	if c.SQL.SslMode == "" {
		c.SQL.SslMode = "disable"
	}

	if c.SQL.BatchTickerDuration == 0 {
		c.SQL.BatchTickerDuration = 10 * time.Second
	}
}
