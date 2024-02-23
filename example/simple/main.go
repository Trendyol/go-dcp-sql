package main

import (
	"fmt"
	"github.com/Trendyol/go-dcp-sql"
	"github.com/Trendyol/go-dcp-sql/couchbase"
	"github.com/Trendyol/go-dcp-sql/sql"
	_ "github.com/lib/pq"
)

func mapper(event couchbase.Event) []sql.Model {
	var raw = sql.Raw{
		Query: fmt.Sprintf(
			"INSERT INTO `example-schema`.`example-table` (key, value) VALUES ('%s', '%s')",
			string(event.Key),
			string(event.Value),
		),
	}

	return []sql.Model{&raw}
}

func main() {
	connector, err := dcpsql.NewConnectorBuilder("config.yml").
		SetMapper(mapper).
		Build()
	if err != nil {
		panic(err)
	}

	defer connector.Close()
	connector.Start()
}
