package chain_client

import (
	"math/big"
	"net/url"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commonclient "github.com/smartcontractkit/chainlink/v2/common/client"
	commonconfig "github.com/smartcontractkit/chainlink/v2/common/config"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

// Wrapper around the client package NewChainClient
// Allows the chain client to be more accessible to external users that may not have the know how to properly initialize the different components
// Configs should only be common go types
func NewChainClient(lggrName string, cfg ChainClientConfig) (client.Client, error) {
	lggr, err := logger.New()
	if err != nil {
		return nil, err
	}
	lggr = logger.Named(lggr, lggrName)
	var empty url.URL
	var primaries []commonclient.Node[*big.Int, *evmtypes.Head, client.RPCClient]
	var sendonlys []commonclient.SendOnlyNode[*big.Int, client.RPCClient]
	for i, node := range cfg.nodes {
		if node.SendOnly != nil && *node.SendOnly {
			rpc := client.NewRPCClient(lggr, empty, (*url.URL)(node.HTTPURL), *node.Name, int32(i), cfg.chainID,
				commonclient.Secondary)
			sendonly := commonclient.NewSendOnlyNode(lggr, (url.URL)(*node.HTTPURL),
				*node.Name, cfg.chainID, rpc)
			sendonlys = append(sendonlys, sendonly)
		} else {
			rpc := client.NewRPCClient(lggr, (url.URL)(*node.WSURL), (*url.URL)(node.HTTPURL), *node.Name, int32(i),
				cfg.chainID, commonclient.Primary)
			primaryNode := commonclient.NewNode(cfg, cfg.noNewHeadsThreshold,
				lggr, (url.URL)(*node.WSURL), (*url.URL)(node.HTTPURL), *node.Name, int32(i), cfg.chainID, *node.Order,
				rpc, "EVM")
			primaries = append(primaries, primaryNode)
		}
	}
	return client.NewChainClient(lggr, cfg.SelectionMode(), cfg.leaseDuration, cfg.noNewHeadsThreshold, primaries, sendonlys, cfg.chainID, commonconfig.ChainType(cfg.chainType)), nil
}
