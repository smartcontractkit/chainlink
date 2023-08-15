package evm

import (
	"github.com/smartcontractkit/chainlink-relay/pkg/loop"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
)

//go:generate mockery --quiet --name LoopRelayAdapter --output ./mocks/ --case=underscore
type LoopRelayAdapter interface {
	loop.Relayer
	Chain() evm.Chain
	Default() bool
}
type LoopRelayer struct {
	loop.Relayer
	ext EVMChainRelayerExtender
}

var _ loop.Relayer = &LoopRelayer{}

func NewLoopRelayAdapter(r *Relayer, cs EVMChainRelayerExtender) *LoopRelayer {
	ra := relay.NewRelayerAdapter(r, cs)
	return &LoopRelayer{
		Relayer: ra,
		ext:     cs,
	}
}

func (la *LoopRelayer) Chain() evm.Chain {
	return la.ext.Chain()
}

func (la *LoopRelayer) Default() bool {
	return la.ext.Default()
}
