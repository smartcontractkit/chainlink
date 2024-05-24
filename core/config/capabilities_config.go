package config

type Capabilities interface {
	Peering() P2P
	// NOTE: RegistrySyncer will need config with relay ID, chain ID and contract address when implemented
}
