# Go Dcp SQL

[![Go Reference](https://pkg.go.dev/badge/github.com/Trendyol/go-dcp-sql.svg)](https://pkg.go.dev/github.com/Trendyol/go-dcp-sql) [![Go Report Card](https://goreportcard.com/badge/github.com/Trendyol/go-dcp-sql)](https://goreportcard.com/report/github.com/Trendyol/go-dcp-sql) [![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/Trendyol/go-dcp-sql/badge)](https://scorecard.dev/viewer/?uri=github.com/Trendyol/go-dcp-sql)


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
```go

import (
    _ "github.com/lib/pq" // DON'T FORGET TO ADD THE DRIVER
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
    SetMapper(mapper).Build()
	
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

| Variable                  | Type          | Required | Default | Description                                                                                        |                                                           
|---------------------------|---------------|----------|---------|----------------------------------------------------------------------------------------------------|
| `sql.host`                | string        | yes      |         | SQL connection urls                                                                                |
| `sql.user`                | string        | yes      |         | SQL username                                                                                       |
| `sql.password`            | string        | yes      |         | SQL password                                                                                       |
| `sql.dbName`              | string        | yes      | 1000    | SQL database name                                                                                  |
| `sql.sslMode`             | string        | no       | disable | Enabling SQL SSL mode                                                                              |
| `sql.driverName`          | string        | yes      |         | Driver name                                                                                        |
| `sql.port`                | int           | yes      |         | SQL port                                                                                           |
| `sql.batchSizeLimit`      | int           | no       | 1000    | Maximum message count for batch, if exceed flush will be triggered                                 |
| `sql.batchTickerDuration` | time.Duration | no       | 10s     | Batch is being flushed automatically at specific time intervals for long waiting messages in batch |

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
