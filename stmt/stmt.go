package stmt

import (
	"database/sql/driver"

	"github.com/loghole/db/hooks"
)

type Stmt struct {
	Stmt  driver.Stmt
	hooks hooks.Hooks
	query string
}
