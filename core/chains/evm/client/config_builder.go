package client

import (
	"fmt"
	"net/url"
	"time"

	"go.uber.org/multierr"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink/v2/common/config"
	evmconfig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
)

type nodeConfig struct {
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
	selectionMode *string,
	leaseDuration time.Duration,
	chainType string,
	nodeCfgs []nodeConfig,
	pollFailureThreshold *uint32,
	pollInterval time.Duration,
	syncThreshold *uint32,
	nodeIsSyncingEnabled *bool,
) (evmconfig.NodePool, []*toml.Node, config.ChainType, error) {
	nodes, err := parseNodeConfigs(nodeCfgs)
	if err != nil {
		return nil, nil, "", err
	}
	nodePool := toml.NodePool{
		SelectionMode:        selectionMode,
		LeaseDuration:        commonconfig.MustNewDuration(leaseDuration),
		PollFailureThreshold: pollFailureThreshold,
		PollInterval:         commonconfig.MustNewDuration(pollInterval),
		SyncThreshold:        syncThreshold,
		NodeIsSyncingEnabled: nodeIsSyncingEnabled,
	}
	nodePoolCfg := &evmconfig.NodePoolConfig{C: nodePool}
	return nodePoolCfg, nodes, config.ChainType(chainType), nil
}

func parseNodeConfigs(nodeCfgs []nodeConfig) ([]*toml.Node, error) {
	nodes := make([]*toml.Node, len(nodeCfgs))
	for i, nodeCfg := range nodeCfgs {
		if nodeCfg.WSURL == nil || nodeCfg.HTTPURL == nil {
			return nil, fmt.Errorf("node config [%d]: missing WS or HTTP URL", i)
		}
		wsUrl := commonconfig.MustParseURL(*nodeCfg.WSURL)
		httpUrl := commonconfig.MustParseURL(*nodeCfg.HTTPURL)
		node := &toml.Node{
			Name:     nodeCfg.Name,
			WSURL:    wsUrl,
			HTTPURL:  httpUrl,
			SendOnly: nodeCfg.SendOnly,
			Order:    nodeCfg.Order,
		}
		nodes[i] = node
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
