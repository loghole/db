package db

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"

	"github.com/loghole/db/internal"
	"github.com/loghole/db/wrapper"
)

func Open(tracer opentracing.Tracer, driverName, dataSourceName string) (db *sql.DB, err error) {
	driverName, err = internal.WrappedDriver(wrapper.NewWrapper(tracer, dataSourceName), driverName)
	if err != nil {
		return nil, err
	}

	return sql.Open(driverName, dataSourceName)
}

func OpenSQLx(tracer opentracing.Tracer, driverName, dataSourceName string) (db *sqlx.DB, err error) {
	driverName, err = internal.WrappedDriver(wrapper.NewWrapper(tracer, dataSourceName), driverName)
	if err != nil {
		return nil, err
	}

	return sqlx.Open(driverName, dataSourceName)
}
