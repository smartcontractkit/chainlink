package client

import (
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"

	commonclient "github.com/smartcontractkit/chainlink/v2/common/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

// Type aliases for EVM
type (
	EvmNode         = commonclient.Node[*big.Int, *evmtypes.Head, RPCClient]
	EvmSendOnlyNode = commonclient.SendOnlyNode[*big.Int, RPCClient]
	EvmMultiNode    = commonclient.MultiNode[*big.Int, evmtypes.Nonce, common.Address, common.Hash, *types.Transaction, common.Hash, types.Log, ethereum.FilterQuery, *evmtypes.Receipt, *assets.Wei, *evmtypes.Head, RPCClient, rpc.BatchElem]
)
