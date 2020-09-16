package internal

import (
	"context"
	"database/sql/driver"

	"github.com/loghole/db/wrapper"
)

type Tx struct {
	driver.Tx
	wrapper.Wrapper
	ctx   context.Context
	txCtx context.Context
}

func NewTx(tx driver.Tx, wrapper wrapper.Wrapper, ctx, txCtx context.Context) *Tx {
	return &Tx{
		Tx:      tx,
		Wrapper: wrapper,
		ctx:     ctx,
		txCtx:   txCtx,
	}
}

func (t *Tx) Commit() (err error) {
	t.ctx = t.Wrapper.BeforeQuery(t.ctx, ActionCommit)

	err = t.Tx.Commit()

	t.Wrapper.AfterQuery(t.ctx, err)

	t.Wrapper.AfterQuery(t.txCtx, err)

	return
}

func (t *Tx) Rollback() (err error) {
	t.ctx = t.Wrapper.BeforeQuery(t.ctx, ActionRollback)

	err = t.Tx.Rollback()

	t.Wrapper.AfterQuery(t.ctx, err)

	t.Wrapper.AfterQuery(t.txCtx, err)

	return
}
