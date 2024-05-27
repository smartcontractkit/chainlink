package config

type RemoteRegistry interface {
	ContractAddress() string
	Network() string
	ChainID() string
}

type Capabilities interface {
	// This enables Capability related features, such as DON-to-DON
	// communication and CapabilityRegistry syncing.
	Enabled() bool
	Peering() P2P
	RemoteRegistry() RemoteRegistry
}
