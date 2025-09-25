package chainlink

import "github.com/smartcontractkit/chainlink/v2/core/config/toml"

type importedDKGRecipientKeyConfig struct {
	s toml.DKGRecipientKey
}

func (t *importedDKGRecipientKeyConfig) JSON() string {
	if t.s.JSON == nil {
		return ""
	}
	return string(*t.s.JSON)
}

func (t *importedDKGRecipientKeyConfig) Password() string {
	if t.s.Password == nil {
		return ""
	}
	return string(*t.s.Password)
}
