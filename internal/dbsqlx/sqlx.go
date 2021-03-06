package dbsqlx

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

func NewSQLx(driverName, dataSourceName string) (*sqlx.DB, error) {
	stdDB, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("can't open db: %w", err)
	}

	if err := stdDB.PingContext(context.TODO()); err != nil {
		return nil, fmt.Errorf("can't ping db: %w", err)
	}

	db := sqlx.NewDb(stdDB, strings.Split(driverName, "-")[0])

	return db, nil
}
