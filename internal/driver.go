package internal

import (
	"database/sql"
	"database/sql/driver"
	"fmt"

	"github.com/loghole/db/wrapper"
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

type Driver struct {
	driver.Driver
	wrapper.Wrapper
}

func WrappedDriver(wrapper wrapper.Wrapper, driverName string) (newName string, err error) {
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
		return NewConnQueryExecAndNamedValue(conn, d.Wrapper), nil
	case driverConnQueryExec:
		return NewConnQueryExec(conn, d.Wrapper), nil
	case driverConnNamedValue:
		return NewConnNamedValue(conn, d.Wrapper), nil
	case driverConnBeginTx:
		return NewConnBeginTx(conn, d.Wrapper), nil
	}

	return conn, nil
}
