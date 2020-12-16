package hooks

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/loghole/dbhook"

	"github.com/loghole/db/internal/dbsqlx"
)

type ReconnectHook struct {
	db     *sqlx.DB
	config *Config
}

var ErrCanRetry = errors.New("connection reconnect")

func NewReconnectHook(db *sqlx.DB, config *Config) *ReconnectHook {
	return &ReconnectHook{
		db:     db,
		config: config,
	}
}

func (rh *ReconnectHook) Error(ctx context.Context, input *dbhook.HookInput) (context.Context, error) {
	if input.Error != nil && isReconnectError(input.Error) {
		tmpDB, err := dbsqlx.NewSQLx(rh.config.DriverName, rh.config.DataSourceName)
		if err != nil {
			return ctx, fmt.Errorf("reconnect error: %w", err)
		}

		*rh.db = *tmpDB

		return ctx, fmt.Errorf("%w: %s", ErrCanRetry, input.Error.Error()) // nolint:errorlint // need wrap ErrCanRetry
	}

	return ctx, input.Error
}

func isReconnectError(err error) bool {
	msg := err.Error()

	return strings.Contains(msg, "broken pipe") ||
		strings.Contains(msg, "bad connection")
}