package db

import (
	"database/sql"
	"log"

	"github.com/loghole/db/internal"
)

func Open(driverName, dataSourceName string) (*sql.DB, error) {
	log.Printf("old name: %s", driverName)

	driverName = internal.WrappedDriver(driverName)

	log.Printf("new name: %s", driverName)

	return sql.Open(driverName, dataSourceName)
}
