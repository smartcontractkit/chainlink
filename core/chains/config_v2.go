package chains

import (
	"context"

	"github.com/smartcontractkit/chainlink-relay/pkg/types"
)

type configsV2AsV1[N Node] struct {
	*configChains
	*configNodes[N]
}

type ConfigsV2[N Node] interface {
	//chainConfigsV2
	ChainStatuser
	nodeConfigsV2[N]
}

// NewConfigs returns a [Configs] backed by [ConfigsV2].
func NewConfigs[N Node](cfgs ConfigsV2[N]) Statuser[N] {
	return configsV2AsV1[N]{
		newConfigChains(cfgs),
		newConfigNodes[N](cfgs),
	}
}

// configChains is a generic, immutable Configs for chains.
type configChains struct {
	x ChainStatuser
}

type chainConfigsV2 interface {
	ChainStatus() (types.ChainStatus, error)
	//Chains(ids ...string) ([]types.ChainStatus, error)

}

// newConfigChains returns a chains backed by chains.
func newConfigChains(d ChainStatuser) *configChains {
	return &configChains{x: d}
}

func (c *configChains) ChainStatus() (types.ChainStatus, error) {
	return c.x.ChainStatus(context.Background())
}

type nodeConfigsV2[N Node] interface {
	NodeConfigs[N]
	/*
		Nodes(names ...string) (nodes []N, err error)

		NodeStatus(name string) (types.NodeStatus, error)
		NodeStatusesPaged(offset, limit int) (nodes []types.NodeStatus, count int, err error)

		Node(name string) (N, error)
		NodeStatuses(name ...string) ([]types.NodeStatus, error)
	*/
}

// configNodes is a generic Configs for nodes.
type configNodes[N Node] struct {
	nodeConfigsV2[N]
}

func newConfigNodes[N Node](d nodeConfigsV2[N]) *configNodes[N] {
	return &configNodes[N]{d}
}
