package config

import "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

type OracleIdentity struct {
	OffchainPublicKey types.OffchainPublicKey
	OnchainPublicKey  types.OnchainPublicKey
	PeerID            string
	TransmitAccount   types.Account
}
