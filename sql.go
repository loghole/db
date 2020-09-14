package db

import (
	"database/sql"
	"log"

	"github.com/opentracing/opentracing-go"

	"github.com/loghole/db/internal"
)

func Open(driverName, dataSourceName string) (*sql.DB, error) {
	log.Printf("old name: %s", driverName)

	driverName, err = internal.WrappedDriver(driverName, )

	log.Printf("new name: %s", driverName)

	return sql.Open(driverName, dataSourceName)
}
