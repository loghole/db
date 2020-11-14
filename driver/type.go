package driver

import (
	"context"
	"database/sql/driver"
	"errors"
	"log"

	"github.com/loghole/db/hooks"
)

type EQ struct {
	*Conn
	*ExecerContext
	*QueryerContext
}

type EQSR struct {
	*Conn
	*ExecerContext
	*QueryerContext
	*SessionResetter
}

type Conn struct {
	Conn  driver.Conn
	hooks hooks.Hooks
}

func (conn *Conn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	log.Println("prepeared context", query)

	var (
		stmt driver.Stmt
		err  error
	)

	if c, ok := conn.Conn.(driver.ConnPrepareContext); ok {
		stmt, err = c.PrepareContext(ctx, query)
	} else {
		stmt, err = conn.Prepare(query)
	}

	if err != nil {
		return nil, err
	}

	return &Stmt{stmt, conn.hooks, query}, nil
}

func (conn *QueryerContext) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	log.Println("conn", "QueryContext", query)

	for _, h := range conn.hooks {
		h.Before(ctx, query)
	}

	results, err := conn.queryContext(ctx, query, args)
	if err != nil {
		return results, err
	}

	for _, h := range conn.hooks {
		h.After(ctx, query)
	}

	return results, err
}

func (conn *QueryerContext) queryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	switch c := conn.Conn.Conn.(type) {
	case driver.QueryerContext:
		return c.QueryContext(ctx, query, args)
	case driver.Queryer:
		dargs, err := namedValueToValue(args)
		if err != nil {
			return nil, err
		}
		return c.Query(query, dargs)
	default:
		// This should not happen
		return nil, errors.New("QueryerContext created for a non Queryer driver.Conn")
	}
}

func (conn *Conn) Prepare(query string) (driver.Stmt, error) { return conn.Conn.Prepare(query) }
func (conn *Conn) Close() error                              { return conn.Conn.Close() }
func (conn *Conn) Begin() (driver.Tx, error)                 { return conn.Conn.Begin() }

type ExecerContext struct {
	*Conn
}

func (conn *ExecerContext) execContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	switch c := conn.Conn.Conn.(type) {
	case driver.ExecerContext:
		return c.ExecContext(ctx, query, args)
	case driver.Execer:
		dargs, err := namedValueToValue(args)
		if err != nil {
			return nil, err
		}
		return c.Exec(query, dargs)
	default:
		// This should not happen
		return nil, errors.New("ExecerContext created for a non Execer driver.Conn")
	}
}

func (conn *ExecerContext) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	log.Println("execer")

	for _, h := range conn.hooks {
		h.Before(ctx, query)
	}

	results, err := conn.execContext(ctx, query, args)
	if err != nil {
		return results, err
	}

	for _, h := range conn.hooks {
		h.After(ctx, query)
	}

	return results, err
}

type QueryerContext struct {
	*Conn
}

type SessionResetter struct {
	*Conn
}

// namedValueToValue copied from database/sql
func namedValueToValue(named []driver.NamedValue) ([]driver.Value, error) {
	dargs := make([]driver.Value, len(named))
	for n, param := range named {
		if len(param.Name) > 0 {
			return nil, errors.New("sql: driver does not support the use of Named Parameters")
		}
		dargs[n] = param.Value
	}
	return dargs, nil
}
