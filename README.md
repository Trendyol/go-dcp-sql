# Go Dcp SQL [![Go Reference](https://pkg.go.dev/badge/github.com/Trendyol/go-dcp-sql.svg)](https://pkg.go.dev/github.com/Trendyol/go-dcp-sql) [![Go Report Card](https://goreportcard.com/badge/github.com/Trendyol/go-dcp-sql)](https://goreportcard.com/report/github.com/Trendyol/go-dcp-sql) [![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/Trendyol/go-dcp-sql/badge)](https://scorecard.dev/viewer/?uri=github.com/Trendyol/go-dcp-sql)

**Go Dcp SQL** streams documents from Couchbase Database Change Protocol (DCP) and writes to
SQL tables in near real-time.

## Features

* Custom SQL queries **per** DCP event.
* **Update multiple documents** for a DCP event(see [Example](#example)).
* Handling different DCP events such as **expiration, deletion and mutation**(see [Example](#example)).
* **Managing batch configurations** such as maximum batch size, batch ticker durations.
* **Scale up and down** by custom membership algorithms(Couchbase, KubernetesHa, Kubernetes StatefulSet or
  Static, see [examples](https://github.com/Trendyol/go-dcp#examples)).
* **Easily manageable configurations**.

## Example

**Note:** If you prefer to use the default mapper by entering the configuration instead of creating a custom mapper, please refer to [this](#collection-table-mapping-configuration) topic.
Otherwise, you can refer to the example provided below:

```go
package main

import (
  _ "github.com/lib/pq" // DON'T FORGET TO ADD THE DRIVER
)

func mapper(event couchbase.Event) []sql.Model {
  var raw = sql.Raw{
    Query: fmt.Sprintf(
      "INSERT INTO `example-schema`.`example-table` (key, value) VALUES ($1, $2);",
    ),
    Args: []interface{}{
      string(event.Key),
      string(event.Value),
    },
  }

  return []sql.Model{&raw}
}

func main() {
  connector, err := dcpsql.NewConnectorBuilder("config.yml").
    SetMapper(mapper). // NOT NEEDED IF YOU'RE USING DEFAULT MAPPER. JUST CALL Build() FUNCTION
    Build()
  if err != nil {
    panic(err)
  }

  defer connector.Close()
  connector.Start()
}
```

## Configuration

### Dcp Configuration

Check out on [go-dcp](https://github.com/Trendyol/go-dcp#configuration)

### SQL Specific Configuration

| Variable                                                | Type                     | Required | Default | Description                                                                                        |                                                           
|---------------------------------------------------------|--------------------------|----------|---------|----------------------------------------------------------------------------------------------------|
| `sql.host`                                              | string                   | yes      |         | SQL connection urls                                                                                |
| `sql.user`                                              | string                   | yes      |         | SQL username                                                                                       |
| `sql.password`                                          | string                   | yes      |         | SQL password                                                                                       |
| `sql.dbName`                                            | string                   | yes      |         | SQL database name                                                                                  |
| `sql.sslMode`                                           | string                   | no       | disable | Enabling SQL SSL mode                                                                              |
| `sql.driverName`                                        | string                   | yes      |         | Driver name                                                                                        |
| `sql.port`                                              | int                      | yes      |         | SQL port                                                                                           |
| `sql.batchTickerDuration`                               | time.Duration            | no       | 10s     | Batch is being flushed automatically at specific time intervals for long waiting messages in batch |
| `sql.collectionTableMapping`                            | []CollectionTableMapping | no       | 10s     | Will be used for default mapper. Please read the next topic.                                       |

### Collection Table Mapping Configuration

Collection table mapping configuration is optional. This configuration should only be provided if you are using the default mapper. If you are implementing your own custom mapper function, this configuration is not needed.

| Variable                                                 | Type    | Required | Default | Description                                                                  |                                                           
|----------------------------------------------------------|---------|----------|---------|------------------------------------------------------------------------------|
| `sql.collectionTableMapping[].collection`                | string  | yes      |         | Couchbase collection name                                                    |
| `sql.collectionTableMapping[].tableName`                 | string  | yes      |         | Target SQL table name                                                        |
| `sql.collectionTableMapping[].keyColumnName`             | string  | yes      |         | Column name for document key in SQL table                                    |
| `sql.collectionTableMapping[].valueColumnName`           | string  | yes      |         | Column name for document value in SQL table                                  |
| `sql.collectionTableMapping[].audit.enabled`             | bool    | no       |         | Enable audit columns for tracking document changes                           |
| `sql.collectionTableMapping[].audit.createdAtColumnName` | string  | no       |         | Column name for tracking document creation time                              |
| `sql.collectionTableMapping[].audit.updatedAtColumnName` | string  | no       |         | Column name for tracking document update time                                |

## Exposed metrics

| Metric Name                                   | Description                   | Labels | Value Type |
|-----------------------------------------------|-------------------------------|--------|------------|
| sql_connector_latency_ms                      | Time to adding to the batch.  | N/A    | Gauge      |
| sql_connector_bulk_request_process_latency_ms | Time to process bulk request. | N/A    | Gauge      |

You can also use all DCP-related metrics explained [here](https://github.com/Trendyol/go-dcp#exposed-metrics).
All DCP-related metrics are automatically injected. It means you don't need to do anything. 

## Contributing

Go Dcp SQL is always open for direct contributions. For more information please check
our [Contribution Guideline document](./CONTRIBUTING.md).

## License

Released under the [MIT License](LICENSE).
