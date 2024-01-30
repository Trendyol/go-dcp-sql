package main

import (
	"fmt"
	sql "github.com/Trendyol/go-dcp-sql"
	"github.com/Trendyol/go-dcp-sql/couchbase"
	_ "github.com/lib/pq"
)

func mapper(event couchbase.Event) []sql.Model {
	var model = sql.SqlModel{
		Query: fmt.Sprintf(
			"INSERT INTO `example-schema`.`example-table` (key, value) VALUES ('%s', '%s')",
			string(event.Key),
			string(event.Value),
		),
	}

	return []sql.Model{&model}
}

func main() {
	connector, err := sql.NewConnectorBuilder("config.yml").
		SetMapper(mapper).
		Build()
	if err != nil {
		panic(err)
	}

	defer connector.Close()
	connector.Start()
}
