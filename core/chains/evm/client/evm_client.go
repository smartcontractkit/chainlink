package client

import (
	"fmt"
	"math"
	"math/big"
	"net/url"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	commonclient "github.com/smartcontractkit/chainlink/v2/common/client"
	evmconfig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
)

func NewEvmClient(cfg evmconfig.NodePool, chainCfg commonclient.ChainConfig, clientErrors evmconfig.ClientErrors, lggr logger.Logger, chainID *big.Int, nodes []*toml.Node, chainType chaintype.ChainType) (Client, error) {
	var empty url.URL
	var primaries []commonclient.Node[*big.Int, *RPCClient]
	var sendonlys []commonclient.SendOnlyNode[*big.Int, *RPCClient]
	largePayloadRPCTimeout, defaultRPCTimeout := getRPCTimeouts(chainType)

	for i, node := range nodes {
		if i < math.MinInt32 || i > math.MaxInt32 {
			return nil, fmt.Errorf("integer overflow: cannot convert %d to int32", i)
		}
		int32i := int32(i)
		if node.SendOnly != nil && *node.SendOnly {
			rpc := NewRPCClient(cfg, lggr, empty, (*url.URL)(node.HTTPURL), *node.Name, int32i, chainID,
				commonclient.Secondary, largePayloadRPCTimeout, defaultRPCTimeout, chainType)
			sendonly := commonclient.NewSendOnlyNode(lggr, (url.URL)(*node.HTTPURL),
				*node.Name, chainID, rpc)
			sendonlys = append(sendonlys, sendonly)
		} else {
			rpc := NewRPCClient(cfg, lggr, (url.URL)(*node.WSURL), (*url.URL)(node.HTTPURL), *node.Name, int32i,
				chainID, commonclient.Primary, largePayloadRPCTimeout, defaultRPCTimeout, chainType)
			primaryNode := commonclient.NewNode(cfg, chainCfg,
				lggr, (url.URL)(*node.WSURL), (*url.URL)(node.HTTPURL), *node.Name, int32i, chainID, *node.Order,
				rpc, "EVM")
			primaries = append(primaries, primaryNode)
		}
	}

	return NewChainClient(lggr, cfg.SelectionMode(), cfg.LeaseDuration(),
		primaries, sendonlys, chainID, clientErrors, cfg.DeathDeclarationDelay(), chainType), nil
}

func getRPCTimeouts(chainType chaintype.ChainType) (largePayload, defaultTimeout time.Duration) {
	if chaintype.ChainHedera == chainType {
		return 30 * time.Second, commonclient.QueryTimeout
	}

	return commonclient.QueryTimeout, commonclient.QueryTimeout
}
