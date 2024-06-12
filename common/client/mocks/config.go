package mocks

import "time"

type ChainConfig struct {
	IsFinalityTagEnabled   bool
	FinalityDepthVal       uint32
	NoNewHeadsThresholdVal time.Duration
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
