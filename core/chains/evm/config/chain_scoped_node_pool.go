package config

import (
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
)

type NodePoolConfig struct {
	C toml.NodePool
}

func (n *NodePoolConfig) PollFailureThreshold() uint32 {
	return *n.C.PollFailureThreshold
}

func (n *NodePoolConfig) PollInterval() time.Duration {
	return n.C.PollInterval.Duration()
}

func (n *NodePoolConfig) SelectionMode() string {
	return *n.C.SelectionMode
}

func (n *NodePoolConfig) SyncThreshold() uint32 {
	return *n.C.SyncThreshold
}

func (n *NodePoolConfig) LeaseDuration() time.Duration {
	return n.C.LeaseDuration.Duration()
}

func (n *NodePoolConfig) NodeIsSyncingEnabled() bool {
	return *n.C.NodeIsSyncingEnabled
}

func (n *NodePoolConfig) FinalizedBlockPollInterval() time.Duration {
	return n.C.FinalizedBlockPollInterval.Duration()
}

func (n *NodePoolConfig) NewHeadsPollInterval() time.Duration {
	return n.C.NewHeadsPollInterval.Duration()
}

func (n *NodePoolConfig) Errors() ClientErrors { return &clientErrorsConfig{c: n.C.Errors} }

func (n *NodePoolConfig) EnforceRepeatableRead() bool {
	return *n.C.EnforceRepeatableRead
}

func (n *NodePoolConfig) DeathDeclarationDelay() time.Duration {
	return n.C.DeathDeclarationDelay.Duration()
}
