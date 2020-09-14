package db

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"

	"github.com/loghole/db/internal"
	"github.com/loghole/db/wrapper"
)

func Open(tracer opentracing.Tracer, driverName, dataSourceName string) (*sql.DB, error) {
	newName, err := internal.WrappedDriver(wrapper.NewWrapper(tracer, dataSourceName), driverName)
	if err != nil {
		return nil, err
	}

	return sql.Open(newName, dataSourceName)
}

func OpenSQLx(tracer opentracing.Tracer, driverName, dataSourceName string) (*sqlx.DB, error) {
	db, err := Open(tracer, driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	return sqlx.NewDb(db, driverName), nil
}
