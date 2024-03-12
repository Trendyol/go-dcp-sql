package bulk

import (
	"github.com/Trendyol/go-dcp-sql/sql"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPrepareBulkSQLQueries(t *testing.T) {
	queries := prepareBulkSQLQueries(
		[]sql.Model{
			&sql.Raw{Query: "INSERT INTO `example-schema`.`example-table` (key, value) VALUES (123, 'ABC')"},
			&sql.Raw{Query: "INSERT INTO `example-schema`.`example-table` (key, value) VALUES (456, 'DEF')"},
			&sql.Raw{Query: "INSERT INTO `example-schema`.`example-table2` (key, value) VALUES (456, 'DEF')"},
			&sql.Raw{Query: "UPDATE userinfo SET created = 123 WHERE uid = test"},
		},
	)

	expectedQueries := []string{
		"UPDATE userinfo SET created = 123 WHERE uid = test",
		"INSERT INTO `example-schema`.`example-table2` (key, value)  VALUES  (456, 'DEF')",
		"INSERT INTO `example-schema`.`example-table` (key, value)  VALUES  (123, 'ABC'), (456, 'DEF')",
	}

	assert.ElementsMatch(t, queries, expectedQueries, "batched queries should be merged.")
}
