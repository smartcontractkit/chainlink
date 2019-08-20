package orm

import (
	"github.com/smartcontractkit/chainlink/core/logger"
)

type ormLogWrapper struct {
	*logger.Logger
}

func (l ormLogWrapper) Print(args ...interface{}) {
	switch args[0] {
	case "error":
		logger.Error(args[2])
	case "sql":
		logger.Debugw(args[3].(string), "time", args[2], "rows_affected", args[5])
	default:
		logger.Info("", args[2:])
	}
}
