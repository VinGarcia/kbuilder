package ksqlserver

import (
	"context"
	"database/sql"

	"github.com/vingarcia/kbuilder"
	"github.com/vingarcia/kbuilder/sqldialect"

	// This is imported here so the user don't
	// have to worry about it when he uses it.
	_ "github.com/denisenkom/go-mssqldb"
)

// NewFromSQLDB builds a ksql.DB from a *sql.DB instance
func NewFromSQLDB(db *sql.DB) (ksql.DB, error) {
	return ksql.NewWithAdapter(NewSQLAdapter(db), sqldialect.SqlserverDialect{})
}

// New instantiates a new KSQL client using the "sqlserver" driver
func New(
	_ context.Context,
	connectionString string,
	config ksql.Config,
) (ksql.DB, error) {
	config.SetDefaultValues()

	db, err := sql.Open("sqlserver", connectionString)
	if err != nil {
		return ksql.DB{}, err
	}
	if err = db.Ping(); err != nil {
		return ksql.DB{}, err
	}

	db.SetMaxOpenConns(config.MaxOpenConns)

	return ksql.NewWithAdapter(NewSQLAdapter(db), sqldialect.SqlserverDialect{})
}
