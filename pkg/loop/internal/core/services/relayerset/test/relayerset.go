package test

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
)

type RelayerSet struct {
}

func (s RelayerSet) Get(ctx context.Context, relayID types.RelayID) (core.Relayer, error) {
	//TODO implement me
	panic("implement me")
}

func (s RelayerSet) List(ctx context.Context, relayIDs ...types.RelayID) (map[types.RelayID]core.Relayer, error) {
	//TODO implement me
	panic("implement me")
}
