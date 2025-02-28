package config

import chain_selectors "github.com/smartcontractkit/chain-selectors"

type ImportableEthKey interface {
	// ChainDetails returns the chain details for the key.
	ChainDetails() chain_selectors.ChainDetails
	ImportableKey
}

type ImportableKey interface {
	// JSON must be a valid JSON string conforming to the
	// particular key format.
	JSON() string
	// Password is the password used to encrypt the key.
	Password() string
}

type ImportableEthKeyLister interface {
	List() []ImportableEthKey
}
