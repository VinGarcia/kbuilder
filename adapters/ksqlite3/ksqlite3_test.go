package ksqlite3

import (
	"database/sql"
	"io"
	"testing"

	"github.com/vingarcia/kbuilder"
	"github.com/vingarcia/kbuilder/sqldialect"
)

func TestAdapter(t *testing.T) {
	ksql.RunTestsForAdapter(t, "ksqlite3", sqldialect.Sqlite3Dialect{}, "/tmp/ksql.db", func(t *testing.T) (ksql.DBAdapter, io.Closer) {
		db, err := sql.Open("sqlite3", "/tmp/ksql.db")
		if err != nil {
			t.Fatal(err.Error())
		}
		return SQLAdapter{db}, db
	})
}
