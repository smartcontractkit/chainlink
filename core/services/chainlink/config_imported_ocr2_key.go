package chainlink

import "github.com/smartcontractkit/chainlink/v2/core/config/toml"

type importedOCR2KeyConfig struct {
	s toml.OCR2Key
}

func (t *importedOCR2KeyConfig) JSON() string {
	if t.s.JSON == nil {
		return ""
	}
	return string(*t.s.JSON)
}

func (t *importedOCR2KeyConfig) Password() string {
	if t.s.Password == nil {
		return ""
	}
	return string(*t.s.Password)
}
