package postgresSQL

import (
	"database/sql"
	_ "github.com/lib/pq"
	"warehouse-application/pkg/errors"
)

func Connect(dataSourceName string) *sql.DB {
	db, err := sql.Open("postgres", dataSourceName)
	errors.CheckError(err)
	return db
}

func Disconnect(db *sql.DB) {
	err := db.Close()
	errors.CheckWarning(err)
}
