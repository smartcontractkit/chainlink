package mocks

import "time"

type ChainConfig struct {
	IsFinalityTagEnabled            bool
	FinalityDepthVal                uint32
	NoNewHeadsThresholdVal          time.Duration
	FinalizedBlockOffsetVal         uint32
	NoNewFinalizedHeadsThresholdVal time.Duration
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

func (t ChainConfig) FinalizedBlockOffset() uint32 {
	return t.FinalizedBlockOffsetVal
}

func (t ChainConfig) NoNewFinalizedHeadsThreshold() time.Duration {
	return t.NoNewFinalizedHeadsThresholdVal
}
