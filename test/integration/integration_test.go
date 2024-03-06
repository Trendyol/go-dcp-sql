package integration

import (
	"context"
	"encoding/json"
	"fmt"
	dcpsql "github.com/Trendyol/go-dcp-sql"
	"github.com/Trendyol/go-dcp-sql/config"
	"github.com/Trendyol/go-dcp-sql/couchbase"
	"github.com/Trendyol/go-dcp-sql/sql"
	"github.com/Trendyol/go-dcp-sql/sql/client"
	_ "github.com/lib/pq"
	"sync"
	"testing"
	"time"
)

type AirlineEvent struct {
	name string
}

func Mapper(event couchbase.Event) []sql.Model {

	var airlineEvent AirlineEvent

	err := json.Unmarshal(event.Value, &airlineEvent)
	if err != nil {
		panic(err)
	}
	var raw = sql.Raw{
		Query: fmt.Sprintf(
			"INSERT INTO example_table (id, name) VALUES ('%s', '%s')",
			string(event.Key),
			airlineEvent.name,
		),
	}
	return []sql.Model{&raw}
}

func TestSql(t *testing.T) {
	connector, err := dcpsql.NewConnectorBuilder("config.yml").SetMapper(Mapper).Build()
	if err != nil {
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		connector.Start()
	}()

	time.Sleep(1 * time.Second)

	go func() {
		sql, err := client.NewSQLClient(config.SQL{
			Host:       "localhost",
			User:       "user",
			Password:   "password",
			DBName:     "example",
			Port:       5432,
			DriverName: "postgres",
			SslMode:    "disable",
		})
		if err != nil {
			t.Fatalf("could not open connection to sql %s", err)
		}

		ctx, _ := context.WithTimeout(context.Background(), 3*time.Minute)

	CountCheckLoop:
		for {
			select {
			case <-ctx.Done():
				t.Fatalf("deadline exceed")
			default:
				var count int
				err := sql.QueryRow("SELECT COUNT(*) FROM example_table").Scan(&count)
				if err != nil {
					t.Fatalf("sql query error %s", err)
				}
				if count == 31591 {
					connector.Close()
					goto CountCheckLoop
				}
				time.Sleep(2 * time.Second)
			}
		}

	}()

	wg.Wait()
}
