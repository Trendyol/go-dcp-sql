package client

import (
	"database/sql"
	"fmt"
	"github.com/Trendyol/go-dcp-sql/config"
)

func NewSqlClient(cfg config.Sql) (*sql.DB, error) {
	driverExist(cfg)

	dataSourceName := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DbName, cfg.SslMode,
	)
	c, err := sql.Open(cfg.DriverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func driverExist(cfg config.Sql) {
	var driverExist = false
	for _, driver := range sql.Drivers() {
		if driver == cfg.DriverName {
			driverExist = true
			break
		}
	}

	if !driverExist {
		panic(fmt.Errorf("driver: %s not found", cfg.DriverName))
	}
}
