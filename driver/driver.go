package driver

import (
	"database/sql/driver"

	"github.com/loghole/db/hooks"
)

func Wrap(drv driver.Driver, hks ...hooks.Hook) driver.Driver {
	return &Driver{Driver: drv, hooks: hks}
}

type Driver struct {
	driver.Driver
	hooks hooks.Hooks
}

// Open opens a connection
func (drv *Driver) Open(name string) (driver.Conn, error) {
	conn, err := drv.Driver.Open(name)
	if err != nil {
		return nil, err
	}

	wrapped := &Conn{Conn: conn, hooks: drv.hooks}

	switch {
	case isEQSR(conn):
		return &EQSR{
			wrapped,
			&ExecerContext{wrapped},
			&QueryerContext{wrapped},
			&SessionResetter{wrapped},
		}, nil
	case isEQ(conn):
		return &EQ{
			wrapped,
			&ExecerContext{wrapped},
			&QueryerContext{wrapped},
		}, nil
	case isE(conn):
		return &ExecerContext{wrapped}, nil
	case isQ(conn):
		return &QueryerContext{wrapped}, nil
	}

	return wrapped, nil
}

func isEQ(conn driver.Conn) bool {
	return isE(conn) && isQ(conn)
}

func isEQSR(conn driver.Conn) bool {
	return isE(conn) && isQ(conn) && isSR(conn)
}

func isE(conn driver.Conn) bool {
	_, ok := conn.(driver.ExecerContext)

	return ok
}

func isQ(conn driver.Conn) bool {
	_, ok := conn.(driver.QueryerContext)

	return ok
}

func isSR(conn driver.Conn) bool {
	_, ok := conn.(driver.SessionResetter)

	return ok
}
