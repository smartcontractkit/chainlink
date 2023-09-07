package client

import (
	"context"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type erroringNode[
	CHAIN_ID types.ID,
	HEAD Head,
	RPC_CLIENT NodeClient[CHAIN_ID, HEAD],
] struct {
	errMsg string
}

func (e *erroringNode[CHAIN_ID, HEAD, RPC_CLIENT]) ConfiguredChainID() (chainID CHAIN_ID) {
	return chainID
}

func (e *erroringNode[CHAIN_ID, HEAD, RPC_CLIENT]) Start(ctx context.Context) error {
	return errors.New(e.errMsg)
}

func (e *erroringNode[CHAIN_ID, HEAD, RPC_CLIENT]) Close() error {
	return nil
}

func (e *erroringNode[CHAIN_ID, HEAD, RPC_CLIENT]) String() string {
	return "<erroring node>"
}

func (e *erroringNode[CHAIN_ID, HEAD, RPC_CLIENT]) State() NodeState {
	return nodeStateUnreachable
}

func (e *erroringNode[CHAIN_ID, HEAD, RPC_CLIENT]) StateAndLatest() (NodeState, int64, *utils.Big) {
	return nodeStateUnreachable, -1, nil
}

func (e *erroringNode[CHAIN_ID, HEAD, RPC_CLIENT]) Order() int32 {
	return 100
}

func (e *erroringNode[CHAIN_ID, HEAD, RPC_CLIENT]) Name() string {
	return ""
}
func (e *erroringNode[CHAIN_ID, HEAD, RPC_CLIENT]) NodeStates() map[int32]string {
	return nil
}

func (e *erroringNode[CHAIN_ID, HEAD, RPC_CLIENT]) RPCClient() (rpcClient RPC_CLIENT) {
	return rpcClient
}
