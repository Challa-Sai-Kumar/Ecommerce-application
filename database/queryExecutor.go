package database

import "database/sql"

type QueryExecutor interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
}
