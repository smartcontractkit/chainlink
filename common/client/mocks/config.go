package mocks

import "time"

type NodeConfig struct {
	PollFailureThresholdVal       uint32
	PollIntervalVal               time.Duration
	SelectionModeVal              string
	SyncThresholdVal              uint32
	NodeIsSyncingEnabledVal       bool
	FinalizedBlockPollIntervalVal time.Duration
}

func (n NodeConfig) PollFailureThreshold() uint32 {
	return n.PollFailureThresholdVal
}

func (n NodeConfig) PollInterval() time.Duration {
	return n.PollIntervalVal
}

func (n NodeConfig) SelectionMode() string {
	return n.SelectionModeVal
}

func (n NodeConfig) SyncThreshold() uint32 {
	return n.SyncThresholdVal
}

func (n NodeConfig) NodeIsSyncingEnabled() bool {
	return n.NodeIsSyncingEnabledVal
}

func (n NodeConfig) FinalizedBlockPollInterval() time.Duration {
	return n.FinalizedBlockPollIntervalVal
}

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
