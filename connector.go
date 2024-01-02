package dcpsql

import (
	"github.com/Trendyol/go-dcp"
	"github.com/Trendyol/go-dcp-sql/config"
)

type Connector interface {
	Start()
	Close()
}

type connector struct {
	dcp    dcp.Dcp
	mapper Mapper
	config *config.Connector
}
