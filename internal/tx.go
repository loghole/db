package internal

import (
	"context"
	"database/sql/driver"
)

type Tx struct {
	driver.Tx
	Wrapper
	ctx   context.Context
	txCtx context.Context
}

func NewTx(tx driver.Tx, wrapper Wrapper, ctx, txCtx context.Context) *Tx {
	return &Tx{
		Tx:      tx,
		Wrapper: wrapper,
		ctx:     ctx,
		txCtx:   txCtx,
	}
}

func (t *Tx) Commit() (err error) {
	t.ctx = t.Wrapper.BeforeQuery(t.ctx, "commit")

	err = t.Tx.Commit()

	t.Wrapper.AfterQuery(t.ctx, err)

	t.Wrapper.AfterQuery(t.txCtx, err)

	return
}

func (t *Tx) Rollback() (err error) {
	t.ctx = t.Wrapper.BeforeQuery(t.ctx, "rollback")

	err = t.Tx.Rollback()

	t.Wrapper.AfterQuery(t.ctx, err)

	t.Wrapper.AfterQuery(t.txCtx, err)

	return
}
