package v1_6

// ConnectionConfig defines how a chain should connect with other chains.
type ConnectionConfig struct {
	// RMNVerificationDisabled is true if we do not want the RMN to bless messages FROM this chain.
	RMNVerificationDisabled bool
	// AllowListEnabled is true if we want an allowlist to dictate who can send messages TO this chain.
	AllowListEnabled bool
}
