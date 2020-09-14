package internal

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"database/sql/driver"
	"encoding/hex"
)

type driverConnBeginTx interface {
	driver.Conn
	driver.ConnBeginTx
}

type driverConnNamedValue interface {
	driverConnBeginTx
	driver.NamedValueChecker
}

type driverConnQueryExec interface {
	driverConnBeginTx
	driver.QueryerContext
	driver.ExecerContext
}

type driverConnQueryExecAndNamedValue interface {
	driverConnQueryExec
	driver.NamedValueChecker
}

type Wrapper interface {
	BeforeQuery(ctx context.Context) context.Context
	AfterQuery(ctx context.Context, err error)
}

type Driver struct {
	driver.Driver
}

func WrappedDriver(driverName string) string {
	db, err := sql.Open(driverName, "")
	if err != nil {
		panic(err)
	}

	hash := sha256.New()
	hash.Write([]byte(driverName))

	newName := hex.EncodeToString(hash.Sum(nil))

	sql.Register(newName, &Driver{Driver: db.Driver()})

	return newName
}

func (d *Driver) Open(name string) (driver.Conn, error) {
	conn, err := d.Driver.Open(name)
	if err != nil {
		return conn, err
	}

	switch conn := conn.(type) {
	case driverConnQueryExecAndNamedValue:
		return NewConnQueryExecAndNamedValue(conn), nil
	case driverConnQueryExec:
		return NewConnQueryExec(conn), nil
	case driverConnNamedValue:
		return NewConnNamedValue(conn), nil
	case driverConnBeginTx:
		return NewConnBeginTx(conn), nil
	}

	return conn, nil
}
