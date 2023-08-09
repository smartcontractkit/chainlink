package client

import (
	"context"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type erroringNode[
	CHAIN_ID types.ID,
	BLOCK_HASH types.Hashable,
	HEAD types.Head[BLOCK_HASH],
	SUB types.Subscription,
	RPC_CLIENT NodeClientAPI[CHAIN_ID, BLOCK_HASH, HEAD, SUB],
] struct {
	errMsg string
}

func (e *erroringNode[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) ChainID() (chainID CHAIN_ID, err error) {
	return chainID, errors.New(e.errMsg)
}

func (e *erroringNode[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) Start(ctx context.Context) error {
	return errors.New(e.errMsg)
}

func (e *erroringNode[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) Close() error {
	return nil
}

func (e *erroringNode[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) String() string {
	return "<erroring node>"
}

func (e *erroringNode[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) State() NodeState {
	return NodeStateUnreachable
}

func (e *erroringNode[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) StateAndLatest() (NodeState, int64, *utils.Big) {
	return NodeStateUnreachable, -1, nil
}

func (e *erroringNode[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) Order() int32 {
	return 100
}

func (e *erroringNode[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) Name() string {
	return ""
}
func (e *erroringNode[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) NodeStates() map[int32]string {
	return nil
}

func (e *erroringNode[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) RPCClient() (rpcClient RPC_CLIENT) {
	return rpcClient
}
