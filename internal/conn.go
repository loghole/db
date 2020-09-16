package internal

import (
	"context"
	"database/sql/driver"
	"log"

	"github.com/loghole/db/wrapper"
)

type ConnBeginTx struct {
	driverConnBeginTx
	wrapper.Wrapper
}

func NewConnBeginTx(
	conn driverConnBeginTx,
	wrapper wrapper.Wrapper,
) *ConnBeginTx {
	return &ConnBeginTx{driverConnBeginTx: conn, Wrapper: wrapper}
}

func (c *ConnBeginTx) PrepareContext(ctx context.Context, query string) (stmt driver.Stmt, err error) {
	stmt, err = c.prepare(ctx, query)

	return NewStmt(stmt, c.Wrapper), nil
}

func (c *ConnBeginTx) BeginTx(ctx context.Context, opts driver.TxOptions) (tx driver.Tx, err error) {
	var (
		txCtx    = c.Wrapper.BeforeQuery(ctx, ActionTx)
		beginCtx = c.Wrapper.BeforeQuery(ctx, ActionBegin)
	)

	tx, err = c.driverConnBeginTx.BeginTx(ctx, opts)

	c.Wrapper.AfterQuery(beginCtx, err)

	return NewTx(tx, c.Wrapper, ctx, txCtx), err
}

func (c *ConnBeginTx) prepare(ctx context.Context, query string) (driver.Stmt, error) {
	if conn, ok := c.driverConnBeginTx.(driver.ConnPrepareContext); ok {
		return conn.PrepareContext(ctx, query)
	}

	stmt, err := c.Prepare(query)
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		if err := stmt.Close(); err != nil {
			return nil, err
		}

		return nil, ctx.Err()
	default:
		return stmt, err
	}
}

type ConnNamedValue struct {
	*ConnBeginTx
	conn driverConnNamedValue
}

func NewConnNamedValue(
	conn driverConnNamedValue,
	wrapper wrapper.Wrapper,
) *ConnNamedValue {
	return &ConnNamedValue{
		ConnBeginTx: NewConnBeginTx(conn, wrapper),
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

func NewConnQueryExec(
	conn driverConnQueryExec,
	wrapper wrapper.Wrapper,
) *ConnQueryExec {
	return &ConnQueryExec{
		ConnBeginTx: NewConnBeginTx(conn, wrapper),
		conn:        conn,
	}
}

func (c *ConnQueryExec) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (rows driver.Rows, err error) {
	ctx = c.Wrapper.BeforeQuery(ctx, ActionQuery)

	rows, err = c.conn.QueryContext(ctx, query, args)

	c.Wrapper.AfterQuery(ctx, err)

	return rows, err
}

func (c *ConnQueryExec) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (result driver.Result, err error) {
	ctx = c.Wrapper.BeforeQuery(ctx, ActionExec)

	result, err = c.conn.ExecContext(ctx, query, args)

	c.Wrapper.AfterQuery(ctx, err)

	return result, err
}

type ConnQueryExecAndNamedValue struct {
	*ConnQueryExec
	conn driverConnQueryExecAndNamedValue
}

func NewConnQueryExecAndNamedValue(
	conn driverConnQueryExecAndNamedValue,
	wrapper wrapper.Wrapper,
) *ConnQueryExecAndNamedValue {
	return &ConnQueryExecAndNamedValue{
		ConnQueryExec: NewConnQueryExec(conn, wrapper),
		conn:          conn,
	}
}

func (c *ConnQueryExecAndNamedValue) CheckNamedValue(value *driver.NamedValue) error {
	log.Println("CheckNamedValue")
	return c.conn.CheckNamedValue(value)
}
