package main

import (
	dcpsql "github.com/Trendyol/go-dcp-sql"
	_ "github.com/lib/pq"
)

func main() {
	connector, err := dcpsql.NewConnectorBuilder("config.yml").
		Build()
	if err != nil {
		panic(err)
	}

	defer connector.Close()
	connector.Start()
}
