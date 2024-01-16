package dcpsql

import "github.com/Trendyol/go-dcp-sql/couchbase"

type Mapper func(event couchbase.Event) []Model

type Model interface {
	ConvertSql() string
}

//type UpsertModel struct {
//	Table string
//	Data  map[string]any
//}
//
//func (u *UpsertModel) ConvertSql() string {
//	a := "INSERT INTO"
//	for k, v := range u.Data {
//
//	}
//
//	// TODO
//	return "INSERT INTO"
//}
//
//type DeleteModel struct {
//	Table string
//	Data  map[string]any
//}
//
//func (u *DeleteModel) ConvertSql() string {
//	// TODO
//	return "DELETE ..."
//}

type SqlModel struct {
	Query string
}

func (u *SqlModel) ConvertSql() string {
	return u.Query
}

func DefaultMapper(event couchbase.Event) []Model {
	//TODO
	return nil
}
