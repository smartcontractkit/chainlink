package config

import (
	"math"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/types"
)

type PublicConfig struct {
	DeltaProgress    time.Duration
	DeltaResend      time.Duration
	DeltaRound       time.Duration
	DeltaGrace       time.Duration
	DeltaC           time.Duration
	Alpha            float64
	DeltaStage       time.Duration
	RMax             uint8
	S                []int
	OracleIdentities []OracleIdentity

	F            int
	ConfigDigest types.ConfigDigest
}

type OracleIdentity struct {
	PeerID                string
	OffchainPublicKey     types.OffchainPublicKey
	OnChainSigningAddress types.OnChainSigningAddress
	TransmitAddress       common.Address
}

func (c *PublicConfig) N() int {
	return len(c.OracleIdentities)
}

func (c *PublicConfig) CheckParameterBounds() error {
	if c.F < 0 || c.F > math.MaxUint8 {
		return errors.Errorf("number of potentially faulty oracles must fit in 8 bits.")
	}
	return nil
}

func PublicConfigFromContractConfig(change types.ContractConfig) (PublicConfig, error) {
	pubcon, _, err := publicConfigFromContractConfig(change)
	return pubcon, err
}

func publicConfigFromContractConfig(change types.ContractConfig) (PublicConfig, SharedSecretEncryptions, error) {
	oc, err := decodeContractSetConfigEncodedComponents(change.Encoded)
	if err != nil {
		return PublicConfig{}, SharedSecretEncryptions{}, err
	}

	identities := []OracleIdentity{}
	for i := range change.Signers {
		identities = append(identities, OracleIdentity{
			oc.PeerIDs[i],
			oc.OffchainPublicKeys[i],
			types.OnChainSigningAddress(change.Signers[i]),
			change.Transmitters[i],
		})
	}

	return PublicConfig{
		oc.DeltaProgress,
		oc.DeltaResend,
		oc.DeltaRound,
		oc.DeltaGrace,
		oc.DeltaC,
		oc.Alpha,
		oc.DeltaStage,
		oc.RMax,
		oc.S,
		identities,
		int(change.Threshold),
		change.ConfigDigest,
	}, oc.SharedSecretEncryptions, nil
}
