package connection

import (
	"context"
	"database/sql/driver"

	"github.com/loghole/db/hooks"
)

// Conn is deprecated
type Conn struct {
	Conn  driver.Conn
	hooks hooks.Hooks
}

func (conn *Conn) Prepare(query string) (driver.Stmt, error) {
	panic("not implemented")
}

func (conn *Conn) Close() error {
	panic("not implemented")
}

func (conn *Conn) Begin() (driver.Tx, error) {
	panic("not implemented")
}

type ConnPrepareContext struct {
	Conn  driver.ConnPrepareContext
	hooks hooks.Hooks
}

func (conn *ConnPrepareContext) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	panic("not implemented")
}

type ConnBeginTx struct {
	Conn  driver.ConnBeginTx
	hooks hooks.Hooks
}

func (conn *ConnPrepareContext) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	panic("not implemented")
}
