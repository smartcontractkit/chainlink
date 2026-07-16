package pipeline

import (
	"math/big"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/vrfkey"
)

// VRFKeyStore is the keystore surface required by vrfv2 and vrfv2plus pipeline tasks.
type VRFKeyStore interface {
	GenerateProof(id string, seed *big.Int) (vrfkey.Proof, error)
}
