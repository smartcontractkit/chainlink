package chainlink

import (
	"math"

	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
)

type looppConfig struct {
	l toml.LOOPP
}

func (l *looppConfig) GRPCServerMaxRecvMsgSizeBytes() int {
	if *l.l.GRPCServerMaxRecvMsgSize > math.MaxInt {
		return math.MaxInt
	}
	return int(*l.l.GRPCServerMaxRecvMsgSize) //nolint:gosec // G115 false positive
}
