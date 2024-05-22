package config

type Capabilities interface {
	// This enables Capability related features, such as DON-to-DON
	// communication and CapabilityRegistry syncing.
	Enabled() bool
	Peering() P2P
	// NOTE: RegistrySyncer will need config with relay ID, chain ID and contract address when implemented
	RelayerID() string
	OnchainCapabilityRegistryAddress() string
}
