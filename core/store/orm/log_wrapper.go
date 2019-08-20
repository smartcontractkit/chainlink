package orm

import (
	"github.com/smartcontractkit/chainlink/core/logger"
	"go.uber.org/zap"
)

type ormLogWrapper struct {
	*zap.SugaredLogger
}

func newOrmLogWrapper(logger *logger.Logger) *ormLogWrapper {
	newLogger := logger.
		SugaredLogger.
		Desugar().
		WithOptions(zap.AddCallerSkip(2)).
		Sugar()
	return &ormLogWrapper{newLogger}
}

func (l ormLogWrapper) Print(args ...interface{}) {
	switch args[0] {
	case "error":
		logger.Error(args[2])
	case "log":
		logger.Warn(args[2])
	case "sql":
		logger.Debugw(args[3].(string), "time", args[2], "rows_affected", args[5])
	default:
		logger.Info(args...)
	}
}
