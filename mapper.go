package dcpsql

import (
	"fmt"
	"github.com/Trendyol/go-dcp-sql/config"
	"github.com/Trendyol/go-dcp-sql/couchbase"
	"github.com/Trendyol/go-dcp-sql/sql"
)

type Mapper func(event couchbase.Event) []sql.Model

var collectionTableMappings *[]config.CollectionTableMapping
var mappingCache = make(map[string]config.CollectionTableMapping)

func SetCollectionTableMappings(mappings *[]config.CollectionTableMapping) {
	collectionTableMappings = mappings
	mappingCache = make(map[string]config.CollectionTableMapping)
}

func DefaultMapper(event couchbase.Event) []sql.Model {
	if event.IsMutated {
		mapping := findCollectionTableMapping(event.CollectionName)
		query := buildUpsertQuery(mapping, event)

		return []sql.Model{&query}
	}

	return nil
}

func findCollectionTableMapping(collectionName string) config.CollectionTableMapping {
	if mapping, exists := mappingCache[collectionName]; exists {
		return mapping
	}

	for _, mapping := range *collectionTableMappings {
		if mapping.Collection == collectionName {
			mappingCache[collectionName] = mapping
			return mapping
		}
	}

	panic(fmt.Sprintf("no mapping found for collection: %s", collectionName))
}

func buildUpsertQuery(mapping config.CollectionTableMapping, event couchbase.Event) sql.Raw {
	var query sql.Raw

	audit := mapping.Audit
	if audit.Enabled && len(audit.CreatedAtColumnName) > 0 && len(audit.UpdatedAtColumnName) > 0 {
		query = sql.Raw{
			Query: fmt.Sprintf(
				"INSERT INTO %s (%s, %s, %s, %s) VALUES($1, $2, NOW(), NOW()) "+
					"ON CONFLICT (%s) DO UPDATE SET "+
					"%s = $2, %s = NOW()",
				mapping.TableName,
				mapping.KeyColumnName,
				mapping.ValueColumnName,
				audit.CreatedAtColumnName,
				audit.UpdatedAtColumnName,
				mapping.KeyColumnName,
				mapping.ValueColumnName,
				audit.UpdatedAtColumnName,
			),
			Args: []interface{}{
				string(event.Key),
				string(event.Value),
			},
		}
	} else {
		query = sql.Raw{
			Query: fmt.Sprintf(
				"INSERT INTO %s (%s, %s) VALUES($1, $2) "+
					"ON CONFLICT (%s) DO UPDATE SET %s = $2",
				mapping.TableName,
				mapping.KeyColumnName,
				mapping.ValueColumnName,
				mapping.KeyColumnName,
				mapping.ValueColumnName,
			),
			Args: []interface{}{
				string(event.Key),
				string(event.Value),
			},
		}
	}

	return query
}
