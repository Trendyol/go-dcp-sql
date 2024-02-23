package dcpsql

import (
	"github.com/Trendyol/go-dcp-sql/couchbase"
	"github.com/Trendyol/go-dcp-sql/sql"
)

type Mapper func(event couchbase.Event) []sql.Model

func DefaultMapper(event couchbase.Event) []sql.Model {
	return nil
}
