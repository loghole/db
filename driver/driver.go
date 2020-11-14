package driver

import (
	"context"
	"database/sql/driver"
)

type Hooks []Hook

type Hook interface {
	Before(ctx context.Context, query string, args ...interface{}) (context.Context, error)
	After(ctx context.Context, query string, args ...interface{}) (context.Context, error)
}

type Driver struct {
	driver.Driver
	hooks Hooks
}
