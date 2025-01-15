package client

import (
	"math/big"
	"net/url"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-framework/multinode"

	evmconfig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
)

const QueryTimeout = 10 * time.Second

func NewEvmClient(cfg evmconfig.NodePool, chainCfg multinode.ChainConfig, clientErrors evmconfig.ClientErrors, lggr logger.Logger, chainID *big.Int, nodes []*toml.Node, chainType chaintype.ChainType) (Client, error) {
	var primaries []multinode.Node[*big.Int, *RPCClient]
	var sendonlys []multinode.SendOnlyNode[*big.Int, *RPCClient]
	largePayloadRPCTimeout, defaultRPCTimeout := getRPCTimeouts(chainType)

	for i, node := range nodes {
		if node.SendOnly != nil && *node.SendOnly {
			rpc := NewRPCClient(cfg, lggr, nil, node.HTTPURL.URL(), *node.Name, i, chainID,
				multinode.Secondary, largePayloadRPCTimeout, defaultRPCTimeout, chainType)
			sendonly := multinode.NewSendOnlyNode(lggr, (url.URL)(*node.HTTPURL),
				*node.Name, chainID, rpc)
			sendonlys = append(sendonlys, sendonly)
		} else {
			rpc := NewRPCClient(cfg, lggr, node.WSURL.URL(), node.HTTPURL.URL(), *node.Name, i,
				chainID, multinode.Primary, largePayloadRPCTimeout, defaultRPCTimeout, chainType)
			primaryNode := multinode.NewNode(cfg, chainCfg,
				lggr, node.WSURL.URL(), node.HTTPURL.URL(), *node.Name, i, chainID, *node.Order,
				rpc, "EVM")
			primaries = append(primaries, primaryNode)
		}
	}

	return NewChainClient(lggr, cfg.SelectionMode(), cfg.LeaseDuration(),
		primaries, sendonlys, chainID, clientErrors, cfg.DeathDeclarationDelay(), chainType), nil
}

func getRPCTimeouts(chainType chaintype.ChainType) (largePayload, defaultTimeout time.Duration) {
	if chaintype.ChainHedera == chainType {
		return 30 * time.Second, QueryTimeout
	}

	return QueryTimeout, QueryTimeout
}
