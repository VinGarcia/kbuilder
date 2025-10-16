package kpostgres

import (
	"database/sql"

	"github.com/vingarcia/kbuilder"
	"github.com/vingarcia/kbuilder/sqldialect"
)

// NewFromSQLDB builds a ksql.DB from a *sql.DB instance
func NewFromSQLDB(db *sql.DB) (ksql.DB, error) {
	return ksql.NewWithAdapter(NewSQLAdapter(db), sqldialect.PostgresDialect{})
}
