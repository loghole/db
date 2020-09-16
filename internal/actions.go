package internal

import (
	"github.com/loghole/db/wrapper"
)

const (
	ActionTx       wrapper.Action = "TX"
	ActionBegin    wrapper.Action = "BEGIN"
	ActionQuery    wrapper.Action = "QUERY"
	ActionExec     wrapper.Action = "EXEC"
	ActionCommit   wrapper.Action = "COMMIT"
	ActionRollback wrapper.Action = "ROLLBACK"
)

