package orm

import "github.com/smartcontractkit/chainlink/core/logger"

type ormLogWrapper struct {
	*logger.Logger
}

func (l ormLogWrapper) Print(args ...interface{}) {
	logger.Info("", args)
}
