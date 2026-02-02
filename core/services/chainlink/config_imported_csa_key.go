package chainlink

import "github.com/smartcontractkit/chainlink/v2/core/config/toml"

type importedCSAKeyConfig struct {
	s toml.CSAKey
}

func (t *importedCSAKeyConfig) JSON() string {
	if t.s.JSON == nil {
		return ""
	}
	return string(*t.s.JSON)
}

func (t *importedCSAKeyConfig) Password() string {
	if t.s.Password == nil {
		return ""
	}
	return string(*t.s.Password)
}
