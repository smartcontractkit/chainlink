package client

import (
	"math/big"
	"net/url"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	commonclient "github.com/smartcontractkit/chainlink/v2/common/client"
	evmconfig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

func NewEvmClient(cfg evmconfig.NodePool, chainCfg commonclient.ChainConfig, clientErrors evmconfig.ClientErrors, lggr logger.Logger, chainID *big.Int, nodes []*toml.Node) Client {
	var empty url.URL
	var primaries []commonclient.Node[*big.Int, *evmtypes.Head, RPCClient]
	var sendonlys []commonclient.SendOnlyNode[*big.Int, RPCClient]
	for i, node := range nodes {
		if node.SendOnly != nil && *node.SendOnly {
			rpc := NewRPCClient(lggr, empty, (*url.URL)(node.HTTPURL), *node.Name, int32(i), chainID,
				commonclient.Secondary)
			sendonly := commonclient.NewSendOnlyNode(lggr, (url.URL)(*node.HTTPURL),
				*node.Name, chainID, rpc)
			sendonlys = append(sendonlys, sendonly)
		} else {
			rpc := NewRPCClient(lggr, (url.URL)(*node.WSURL), (*url.URL)(node.HTTPURL), *node.Name, int32(i),
				chainID, commonclient.Primary)
			primaryNode := commonclient.NewNode(cfg, chainCfg,
				lggr, (url.URL)(*node.WSURL), (*url.URL)(node.HTTPURL), *node.Name, int32(i), chainID, *node.Order,
				rpc, "EVM")
			primaries = append(primaries, primaryNode)
		}
	}

	return NewChainClient(lggr, cfg.SelectionMode(), cfg.LeaseDuration(), chainCfg.NodeNoNewHeadsThreshold(),
		primaries, sendonlys, chainID, chainCfg.ChainType(), clientErrors)
}
