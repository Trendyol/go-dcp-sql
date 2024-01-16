package main

import (
	sql "github.com/Trendyol/go-dcp-sql"
	"github.com/Trendyol/go-dcp-sql/couchbase"
)

func mapper(event couchbase.Event) []sql.Model {
	var result []sql.Model
	//TODO User should handle events and create related sql.Model model array includes SQL DML queries.
	return result
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
