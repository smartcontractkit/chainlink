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
	RPC NodeClient[CHAIN_ID, HEAD],
] struct {
	errMsg string
}

func (e *erroringNode[CHAIN_ID, HEAD, RPC]) ConfiguredChainID() (chainID CHAIN_ID) {
	return chainID
}

func (e *erroringNode[CHAIN_ID, HEAD, RPC]) Start(ctx context.Context) error {
	return errors.New(e.errMsg)
}

func (e *erroringNode[CHAIN_ID, HEAD, RPC]) Close() error {
	return nil
}

func (e *erroringNode[CHAIN_ID, HEAD, RPC]) String() string {
	return "<erroring node>"
}

func (e *erroringNode[CHAIN_ID, HEAD, RPC]) State() NodeState {
	return nodeStateUnreachable
}

func (e *erroringNode[CHAIN_ID, HEAD, RPC]) StateAndLatest() (NodeState, int64, *utils.Big) {
	return nodeStateUnreachable, -1, nil
}

func (e *erroringNode[CHAIN_ID, HEAD, RPC]) Order() int32 {
	return 100
}

func (e *erroringNode[CHAIN_ID, HEAD, RPC]) Name() string {
	return ""
}
func (e *erroringNode[CHAIN_ID, HEAD, RPC]) NodeStates() map[int32]string {
	return nil
}

func (e *erroringNode[CHAIN_ID, HEAD, RPC]) RPC() (rpc RPC) {
	return rpc
}
