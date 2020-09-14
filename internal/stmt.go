package internal

import (
	"context"
	"database/sql/driver"
)

type Stmt struct {
	driver.Stmt
	Wrapper
}

func NewStmt(stmt driver.Stmt, wrapper Wrapper) *Stmt {
	return &Stmt{
		Stmt:    stmt,
		Wrapper: wrapper,
	}
}

func (s *Stmt) ExecContext(ctx context.Context, args []driver.NamedValue) (result driver.Result, err error) {
	ctx = s.Wrapper.BeforeQuery(ctx, "exec")

	result, err = s.exec(ctx, args)

	s.Wrapper.AfterQuery(ctx, err)

	return result, err
}

func (s *Stmt) QueryContext(ctx context.Context, args []driver.NamedValue) (rows driver.Rows, err error) {
	ctx = s.Wrapper.BeforeQuery(ctx, "query")

	rows, err = s.query(ctx, args)

	s.Wrapper.AfterQuery(ctx, err)

	return rows, err
}

func (s *Stmt) exec(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	if s, ok := s.Stmt.(driver.StmtExecContext); ok {
		return s.ExecContext(ctx, args)
	}

	values := make([]driver.Value, len(args))

	for _, arg := range args {
		values[arg.Ordinal-1] = arg.Value
	}

	return s.Exec(values)
}

func (s *Stmt) query(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	if s, ok := s.Stmt.(driver.StmtQueryContext); ok {
		return s.QueryContext(ctx, args)
	}

	values := make([]driver.Value, len(args))

	for _, arg := range args {
		values[arg.Ordinal-1] = arg.Value
	}

	return s.Query(values)
}
