package configtest

import (
	"time"

	"github.com/smartcontractkit/chainlink/core/config"
)

var _ config.OCR2Config = &TestGeneralConfig{}

// OCR2DatabaseTimeout returns the overridden value, if one exists.
func (c *TestGeneralConfig) OCR2DatabaseTimeout() time.Duration {
	if c.Overrides.OCR2DatabaseTimeout != nil {
		return *c.Overrides.OCR2DatabaseTimeout
	}
	return c.GeneralConfig.OCR2DatabaseTimeout()
}
