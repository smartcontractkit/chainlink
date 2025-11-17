package chainlink

import (
	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
)

type importedSolKeyConfig struct {
	s toml.SolKey
}

func (t *importedSolKeyConfig) JSON() string {
	if t.s.JSON == nil {
		return ""
	}
	return string(*t.s.JSON)
}

func (t *importedSolKeyConfig) ChainDetails() chain_selectors.ChainDetails {
	if t.s.ID == nil {
		return chain_selectors.ChainDetails{}
	}
	d, ok := chain_selectors.SolanaChainIdToChainSelector()[*t.s.ID]
	if !ok {
		return chain_selectors.ChainDetails{}
	}
	return chain_selectors.ChainDetails{
		ChainSelector: d,
		ChainName:     "solana",
	}
}

func (t *importedSolKeyConfig) Password() string {
	if t.s.Password == nil {
		return ""
	}
	return string(*t.s.Password)
}

type importedSolKeyConfigs struct {
	s toml.SolKeys
}

func (t *importedSolKeyConfigs) List() []config.ImportableChainKey {
	res := make([]config.ImportableChainKey, len(t.s.Keys))

	if len(t.s.Keys) == 0 {
		return res
	}

	for i, v := range t.s.Keys {
		res[i] = &importedSolKeyConfig{s: *v}
	}
	return res
}
