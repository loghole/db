package internal

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
)

const wrappedAlias = "%s-wrapped"

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
	BeforeQuery(ctx context.Context, action string) context.Context
	AfterQuery(ctx context.Context, err error)
}

type Driver struct {
	driver.Driver
	Wrapper
}

func WrappedDriver(wrapper Wrapper, driverName string) (newName string, err error) {
	db, err := sql.Open(driverName, "")
	if err != nil {
		return "", err
	}

	newName = fmt.Sprintf(wrappedAlias, driverName)

	sql.Register(newName, &Driver{Driver: db.Driver(), Wrapper: wrapper})

	return newName, nil
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
