package chainlink

import (
	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
)

type importedEthKeyConfig struct {
	s toml.EthKey
}

func (t *importedEthKeyConfig) JSON() string {
	if t.s.JSON == nil {
		return ""
	}
	return string(*t.s.JSON)
}

func (t *importedEthKeyConfig) ChainDetails() chain_selectors.ChainDetails {
	if t.s.ID == nil {
		return chain_selectors.ChainDetails{}
	}
	d, ok := chain_selectors.ChainByEvmChainID(uint64(*t.s.ID)) //nolint:gosec // disable G115
	if !ok {
		return chain_selectors.ChainDetails{}
	}
	return chain_selectors.ChainDetails{
		ChainSelector: d.Selector,
		ChainName:     d.Name,
	}
}

func (t *importedEthKeyConfig) Password() string {
	if t.s.Password == nil {
		return ""
	}
	return string(*t.s.Password)
}

type importedEthKeyConfigs struct {
	s toml.EthKeys
}

func (t *importedEthKeyConfigs) List() []config.ImportableEthKey {
	res := make([]config.ImportableEthKey, len(t.s.Keys))

	if len(t.s.Keys) == 0 {
		return res
	}

	for i, v := range t.s.Keys {
		res[i] = &importedEthKeyConfig{s: *v}
	}
	return res
}
