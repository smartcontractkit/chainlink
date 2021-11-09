package loader

import (
	"context"
	"errors"

	"github.com/graph-gophers/dataloader"

	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
)

// GetChainByID fetches the chain by it's id.
func GetChainByID(ctx context.Context, id string) (*types.Chain, error) {
	ldr := For(ctx)

	thunk := ldr.ChainsByIDLoader.Load(ctx, dataloader.StringKey(id))
	result, err := thunk()
	if err != nil {
		return nil, err
	}

	chain, ok := result.(types.Chain)
	if !ok {
		return nil, errors.New("invalid type")
	}

	return &chain, nil
}

// GetNodesByChainID fetches the nodes for a chain.
func GetNodesByChainID(ctx context.Context, id string) ([]types.Node, error) {
	ldr := For(ctx)

	thunk := ldr.NodesByChainIDLoader.Load(ctx, dataloader.StringKey(id))
	result, err := thunk()
	if err != nil {
		return nil, err
	}

	nodes, ok := result.([]types.Node)
	if !ok {
		return nil, errors.New("invalid type")
	}

	return nodes, nil
}
