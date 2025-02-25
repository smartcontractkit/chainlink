package chainlink

import "github.com/smartcontractkit/chainlink/v2/core/config/toml"

type importedP2PKeyConfig struct {
	s toml.P2PKey
}

func (t *importedP2PKeyConfig) JSON() string {
	if t.s.JSON == nil {
		return ""
	}
	return string(*t.s.JSON)
}

func (t *importedP2PKeyConfig) Password() string {
	if t.s.Password == nil {
		return ""
	}
	return string(*t.s.Password)
}
