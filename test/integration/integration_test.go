package integration

import (
	"context"
	dcpsql "github.com/Trendyol/go-dcp-sql"
	"github.com/Trendyol/go-dcp-sql/config"
	"github.com/Trendyol/go-dcp-sql/sql/client"
	"github.com/Trendyol/go-dcp/logger"
	_ "github.com/lib/pq"
	"sync"
	"testing"
	"time"
)

type AirlineEvent struct {
	name string
}

func TestDefaultMapper(t *testing.T) {
	t.Run("TestDefaultMapperInsert", testDefaultMapperInsert)
	t.Run("TestDefaultMapperDelete", testDefaultMapperDelete)
}

func testDefaultMapperInsert(t *testing.T) {
	time.Sleep(time.Second * 30)

	connector, err := dcpsql.NewConnectorBuilder("config.yml").Build()
	if err != nil {
		t.Fatal(err)
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
			CollectionTableMapping: []config.CollectionTableMapping{
				{
					Collection:      "_default",
					TableName:       "public.example_table",
					KeyColumnName:   "id",
					ValueColumnName: "name",
					Audit: struct {
						Enabled             bool   `yaml:"enabled"`
						CreatedAtColumnName string `yaml:"createdAtColumnName"`
						UpdatedAtColumnName string `yaml:"updatedAtColumnName"`
					}{},
				},
			},
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
					logger.Log.Info("done")
					connector.Close()
					goto CountCheckLoop
				}
				time.Sleep(2 * time.Second)
			}
		}

	}()

	wg.Wait()
}

func testDefaultMapperDelete(t *testing.T) {
	time.Sleep(time.Second * 30)

	connector, err := dcpsql.NewConnectorBuilder("config.yml").Build()
	if err != nil {
		t.Fatal(err)
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
			CollectionTableMapping: []config.CollectionTableMapping{
				{
					Collection:      "_default",
					TableName:       "public.example_table",
					KeyColumnName:   "id",
					ValueColumnName: "name",
					Audit: struct {
						Enabled             bool   `yaml:"enabled"`
						CreatedAtColumnName string `yaml:"createdAtColumnName"`
						UpdatedAtColumnName string `yaml:"updatedAtColumnName"`
					}{},
				},
			},
		})
		if err != nil {
			t.Fatalf("could not open connection to sql %s", err)
		}

		ctx, _ := context.WithTimeout(context.Background(), 3*time.Minute)

	DeleteCheckLoop:
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
				if count == 0 {
					logger.Log.Info("all records deleted successfully")
					connector.Close()
					goto DeleteCheckLoop
				}
				time.Sleep(2 * time.Second)
			}
		}
	}()

	wg.Wait()
}
