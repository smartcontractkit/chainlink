package chainlink

import "github.com/smartcontractkit/chainlink/v2/core/config/toml"

// Aptos key import config is node-scoped rather than chain-scoped, so unlike
// EVM/Solana imported keys it intentionally only implements ImportableKey.
type importedAptosKeyConfig struct {
	s toml.AptosKey
}

func (t *importedAptosKeyConfig) JSON() string {
	if t.s.JSON == nil {
		return ""
	}
	return string(*t.s.JSON)
}

func (t *importedAptosKeyConfig) Password() string {
	if t.s.Password == nil {
		return ""
	}
	return string(*t.s.Password)
}
