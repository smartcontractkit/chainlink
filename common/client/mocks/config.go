package mocks

import (
	"time"

	commonconfig "github.com/smartcontractkit/chainlink/v2/common/config"
)

type ChainConfig struct {
	IsFinalityTagEnabled   bool
	FinalityDepthVal       uint32
	NoNewHeadsThresholdVal time.Duration
	ChainTypeVal           commonconfig.ChainType
}

func (t ChainConfig) ChainType() commonconfig.ChainType {
	return t.ChainTypeVal
}

func (t ChainConfig) NodeNoNewHeadsThreshold() time.Duration {
	return t.NoNewHeadsThresholdVal
}

func (t ChainConfig) FinalityDepth() uint32 {
	return t.FinalityDepthVal
}

func (t ChainConfig) FinalityTagEnabled() bool {
	return t.IsFinalityTagEnabled
}
