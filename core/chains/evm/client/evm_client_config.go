package client

import (
	"fmt"
	"math/big"
	"net/url"
	"time"

	"go.uber.org/multierr"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink/v2/common/config"
	evmconfig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
)

var _ evmconfig.NodePool = NodePool{}

type NodePool struct {
	selectionMode        string
	leaseDuration        time.Duration
	pollFailureThreshold uint32
	pollInterval         time.Duration
	syncThreshold        uint32
	nodeIsSyncingEnabled bool
}

func (c NodePool) PollFailureThreshold() uint32 {
	return c.pollFailureThreshold
}

func (c NodePool) PollInterval() time.Duration {
	return c.pollInterval
}

func (c NodePool) SelectionMode() string {
	return c.selectionMode
}

func (c NodePool) LeaseDuration() time.Duration {
	return c.leaseDuration
}

func (c NodePool) SyncThreshold() uint32 {
	return c.syncThreshold
}

func (c NodePool) NodeIsSyncingEnabled() bool {
	return c.nodeIsSyncingEnabled
}

type NodeConfig struct {
	Name     *string
	WSURL    *string
	HTTPURL  *string
	SendOnly *bool
	Order    *int32
}

// Build the configs needed to initialize the chain client
// Parameters should only be basic go types to make it accessible for external users
// Configs can be stored in a variety of ways
func NewClientConfigs(
	selectionMode string,
	leaseDuration time.Duration,
	noNewHeadsThreshold time.Duration,
	chainID *big.Int,
	chainType string,
	nodeCfgs []NodeConfig,
	pollFailureThreshold uint32,
	pollInterval time.Duration,
	syncThreshold uint32,
	nodeIsSyncingEnabled bool,
) (evmconfig.NodePool, []*toml.Node, config.ChainType, error) {
	nodes, err := parseNodeConfigs(nodeCfgs)
	if err != nil {
		return nil, nil, "", err
	}
	nodePool := NodePool{
		selectionMode:        selectionMode,
		leaseDuration:        leaseDuration,
		pollFailureThreshold: pollFailureThreshold,
		pollInterval:         pollInterval,
		syncThreshold:        syncThreshold,
		nodeIsSyncingEnabled: nodeIsSyncingEnabled,
	}
	return nodePool, nodes, config.ChainType(chainType), nil
}

func parseNodeConfigs(nodeCfgs []NodeConfig) ([]*toml.Node, error) {
	nodes := make([]*toml.Node, len(nodeCfgs))
	for _, nodeCfg := range nodeCfgs {
		wsUrl, err := commonconfig.ParseURL(*nodeCfg.WSURL)
		if err != nil {
			return nil, err
		}
		var httpUrl *commonconfig.URL
		httpUrl, err = commonconfig.ParseURL(*nodeCfg.HTTPURL)
		if err != nil {
			return nil, err
		}
		node := toml.Node{
			Name:     nodeCfg.Name,
			WSURL:    wsUrl,
			HTTPURL:  httpUrl,
			SendOnly: nodeCfg.SendOnly,
			Order:    nodeCfg.Order,
		}
		nodes = append(nodes, &node)
	}

	if err := validateNodeConfigs(nodes); err != nil {
		return nil, err
	}

	return nodes, nil
}

func validateNodeConfigs(nodes []*toml.Node) (err error) {
	names := commonconfig.UniqueStrings{}
	wsURLs := commonconfig.UniqueStrings{}
	httpURLs := commonconfig.UniqueStrings{}
	for i, node := range nodes {
		if nodeErr := node.ValidateConfig(); nodeErr != nil {
			err = multierr.Append(err, nodeErr)
		}
		if names.IsDupe(node.Name) {
			err = multierr.Append(err, commonconfig.NewErrDuplicate(fmt.Sprintf("Nodes.%d.Name", i), *node.Name))
		}
		u := (*url.URL)(node.WSURL)
		if wsURLs.IsDupeFmt(u) {
			err = multierr.Append(err, commonconfig.NewErrDuplicate(fmt.Sprintf("Nodes.%d.WSURL", i), u.String()))
		}
		u = (*url.URL)(node.HTTPURL)
		if httpURLs.IsDupeFmt(u) {
			err = multierr.Append(err, commonconfig.NewErrDuplicate(fmt.Sprintf("Nodes.%d.HTTPURL", i), u.String()))
		}
	}

	return err
}
