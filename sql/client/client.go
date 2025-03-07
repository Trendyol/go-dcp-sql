package client

import (
	"database/sql"
	"fmt"

	"github.com/Trendyol/go-dcp-sql/config"
)

func NewSQLClient(cfg config.SQL) (*sql.DB, error) {
	dataSourceName := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SslMode,
	)
	c, err := sql.Open(cfg.DriverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return c, nil
}
