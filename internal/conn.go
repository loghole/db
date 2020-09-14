package internal

import (
	"context"
	"database/sql/driver"
	"log"
)

type ConnBeginTx struct {
	driverConnBeginTx
	Wrapper
}

func NewConnBeginTx(conn driverConnBeginTx) *ConnBeginTx {
	return &ConnBeginTx{driverConnBeginTx: conn}
}

func (c *ConnBeginTx) BeginTx(ctx context.Context, opts driver.TxOptions) (tx driver.Tx, err error) {
	log.Println("BeginTx")

	tx, err = c.driverConnBeginTx.BeginTx(ctx, opts)

	return &Tx{Tx: tx}, nil
}

type ConnNamedValue struct {
	*ConnBeginTx
	conn driverConnNamedValue
}

func NewConnNamedValue(conn driverConnNamedValue) *ConnNamedValue {
	return &ConnNamedValue{
		ConnBeginTx: NewConnBeginTx(conn),
		conn:        conn,
	}
}

func (c *ConnNamedValue) CheckNamedValue(value *driver.NamedValue) error {
	return c.conn.CheckNamedValue(value)
}

type ConnQueryExec struct {
	*ConnBeginTx
	conn driverConnQueryExec
}

func NewConnQueryExec(conn driverConnQueryExec) *ConnQueryExec {
	return &ConnQueryExec{
		ConnBeginTx: NewConnBeginTx(conn),
		conn:        conn,
	}
}

func (c *ConnQueryExec) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	log.Println("QueryContext")
	return c.conn.QueryContext(ctx, query, args)
}

func (c *ConnQueryExec) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	log.Println("ExecContext")
	return c.conn.ExecContext(ctx, query, args)
}

type ConnQueryExecAndNamedValue struct {
	*ConnQueryExec
	conn driverConnQueryExecAndNamedValue
}

func NewConnQueryExecAndNamedValue(conn driverConnQueryExecAndNamedValue) *ConnQueryExecAndNamedValue {
	return &ConnQueryExecAndNamedValue{
		ConnQueryExec: NewConnQueryExec(conn),
		conn:          conn,
	}
}

func (c *ConnQueryExecAndNamedValue) CheckNamedValue(value *driver.NamedValue) error {
	log.Println("CheckNamedValue")
	return c.conn.CheckNamedValue(value)
}
