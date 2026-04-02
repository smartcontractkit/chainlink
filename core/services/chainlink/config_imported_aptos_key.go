package chainlink

import (
	"strconv"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
)

// importedAptosKeyConfig adapts a single chain-aware Aptos imported key entry
// to the shared ImportableChainKey interface used during keystore startup.
type importedAptosKeyConfig struct {
	s toml.AptosKey
}

func (t *importedAptosKeyConfig) JSON() string {
	if t.s.JSON == nil {
		return ""
	}
	return string(*t.s.JSON)
}

func (t *importedAptosKeyConfig) ChainDetails() chain_selectors.ChainDetails {
	if t.s.ID == nil {
		return chain_selectors.ChainDetails{}
	}
	details, err := chain_selectors.GetChainDetailsByChainIDAndFamily(strconv.FormatUint(*t.s.ID, 10), chain_selectors.FamilyAptos)
	if err != nil {
		return chain_selectors.ChainDetails{}
	}
	return details
}

func (t *importedAptosKeyConfig) Password() string {
	if t.s.Password == nil {
		return ""
	}
	return string(*t.s.Password)
}

type importedAptosKeyConfigs struct {
	s toml.AptosKeys
}

func (t *importedAptosKeyConfigs) List() []config.ImportableChainKey {
	res := make([]config.ImportableChainKey, len(t.s.Keys))

	if len(t.s.Keys) == 0 {
		return res
	}

	for i, v := range t.s.Keys {
		res[i] = &importedAptosKeyConfig{s: *v}
	}
	return res
}
