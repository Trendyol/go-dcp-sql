package dcpsql

import (
	"github.com/Trendyol/go-dcp-sql/couchbase"
	"github.com/Trendyol/go-dcp-sql/model"
)

type Mapper func(event couchbase.Event) []model.DcpSqlItem

func DefaultMapper(event couchbase.Event) []model.DcpSqlItem {
	return nil
}
