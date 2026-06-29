package chainlink

import (
	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
)

// importedStellarKeyConfig adapts a single chain-aware Stellar imported key entry
// to the shared ImportableChainKey interface used during keystore startup.
type importedStellarKeyConfig struct {
	s toml.StellarKey
}

func (t *importedStellarKeyConfig) JSON() string {
	if t.s.JSON == nil {
		return ""
	}
	return string(*t.s.JSON)
}

func (t *importedStellarKeyConfig) ChainDetails() chain_selectors.ChainDetails {
	if t.s.ID == nil {
		return chain_selectors.ChainDetails{}
	}
	details, err := chain_selectors.GetChainDetailsByChainIDAndFamily(*t.s.ID, chain_selectors.FamilyStellar)
	if err != nil {
		return chain_selectors.ChainDetails{}
	}
	return details
}

func (t *importedStellarKeyConfig) Password() string {
	if t.s.Password == nil {
		return ""
	}
	return string(*t.s.Password)
}

type importedStellarKeyConfigs struct {
	s toml.StellarKeys
}

func (t *importedStellarKeyConfigs) List() []config.ImportableChainKey {
	res := make([]config.ImportableChainKey, len(t.s.Keys))

	if len(t.s.Keys) == 0 {
		return res
	}

	for i, v := range t.s.Keys {
		res[i] = &importedStellarKeyConfig{s: *v}
	}
	return res
}
