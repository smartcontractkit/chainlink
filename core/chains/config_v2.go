package chains

import (
	"github.com/smartcontractkit/chainlink-relay/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
)

type configsV2AsV1[N Node] struct {
	*configChains
	*configNodes[N]
}

type ConfigsV2[N Node] interface {
	chainConfigsV2
	NodeConfigs[N]
}

// NewConfigs returns a [Configs] backed by [ConfigsV2].
func NewConfigs[N Node](cfgs ConfigsV2[N]) Configs[N] {
	return configsV2AsV1[N]{
		newConfigChains(cfgs),
		newConfigNodes[N](cfgs),
	}
}

// configChains is a generic, immutable Configs for chains.
type configChains struct {
	v2 chainConfigsV2
}

type chainConfigsV2 interface {
	Chains(ids ...relay.ChainID) ([]types.ChainStatus, error)
}

// newConfigChains returns a chains backed by chains.
func newConfigChains(d chainConfigsV2) *configChains {
	return &configChains{v2: d}
}

func (o *configChains) Chains(offset, limit int, ids ...relay.ChainID) (chains []types.ChainStatus, count int, err error) {
	chains, err = o.v2.Chains(ids...)
	if err != nil {
		return
	}
	count = len(chains)
	if offset < len(chains) {
		chains = chains[offset:]
	} else {
		chains = nil
	}
	if limit > 0 && len(chains) > limit {
		chains = chains[:limit]
	}
	return
}

// configNodes is a generic Configs for nodes.
type configNodes[N Node] struct {
	NodeConfigs[N]
}

func newConfigNodes[N Node](d NodeConfigs[N]) *configNodes[N] {
	return &configNodes[N]{d}
}
