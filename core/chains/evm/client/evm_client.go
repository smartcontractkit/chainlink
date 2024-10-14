package client

import (
	"math/big"
	"net/url"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	commonclient "github.com/smartcontractkit/chainlink/v2/common/client"
	evmconfig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

func NewEvmClient(cfg evmconfig.NodePool, chainCfg commonclient.ChainConfig, clientErrors evmconfig.ClientErrors, lggr logger.Logger, chainID *big.Int, nodes []*toml.Node, chainType chaintype.ChainType) Client {
	var primaries []commonclient.Node[*big.Int, *evmtypes.Head, RPCClient]
	var sendonlys []commonclient.SendOnlyNode[*big.Int, RPCClient]
	largePayloadRPCTimeout, defaultRPCTimeout := getRPCTimeouts(chainType)
	for i, node := range nodes {
		var ws *url.URL
		if node.WSURL != nil {
			ws = (*url.URL)(node.WSURL)
		}
		if node.SendOnly != nil && *node.SendOnly {
			rpc := NewRPCClient(lggr, nil, (*url.URL)(node.HTTPURL), *node.Name, i, chainID,
				commonclient.Secondary, cfg.FinalizedBlockPollInterval(), cfg.NewHeadsPollInterval(), largePayloadRPCTimeout, defaultRPCTimeout, chainType)
			sendonly := commonclient.NewSendOnlyNode(lggr, (url.URL)(*node.HTTPURL),
				*node.Name, chainID, rpc)
			sendonlys = append(sendonlys, sendonly)
		} else {
			rpc := NewRPCClient(lggr, ws, (*url.URL)(node.HTTPURL), *node.Name, i,
				chainID, commonclient.Primary, cfg.FinalizedBlockPollInterval(), cfg.NewHeadsPollInterval(), largePayloadRPCTimeout, defaultRPCTimeout, chainType)
			primaryNode := commonclient.NewNode(cfg, chainCfg,
				lggr, ws, (*url.URL)(node.HTTPURL), *node.Name, i, chainID, *node.Order,
				rpc, "EVM")
			primaries = append(primaries, primaryNode)
		}
	}

	return NewChainClient(lggr, cfg.SelectionMode(), cfg.LeaseDuration(), chainCfg.NodeNoNewHeadsThreshold(),
		primaries, sendonlys, chainID, chainType, clientErrors, cfg.DeathDeclarationDelay())
}

func getRPCTimeouts(chainType chaintype.ChainType) (largePayload, defaultTimeout time.Duration) {
	if chaintype.ChainHedera == chainType {
		return 30 * time.Second, commonclient.QueryTimeout
	}

	return commonclient.QueryTimeout, commonclient.QueryTimeout
}
