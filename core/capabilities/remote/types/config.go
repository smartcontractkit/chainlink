package types

const (
	DefaultRegistrationRefreshMs = 30_000
	DefaultRegistrationExpiryMs  = 120_000
	DefaultMessageExpiryMs       = 120_000
)

// NOTE: consider splitting this config into values stored in Registry (KS-118)
// and values defined locally by Capability owners.
func (c *RemoteTriggerConfig) ApplyDefaults() {
	if c.RegistrationRefreshMs == 0 {
		c.RegistrationRefreshMs = DefaultRegistrationRefreshMs
	}
	if c.RegistrationExpiryMs == 0 {
		c.RegistrationExpiryMs = DefaultRegistrationExpiryMs
	}
	if c.MessageExpiryMs == 0 {
		c.MessageExpiryMs = DefaultMessageExpiryMs
	}
}
