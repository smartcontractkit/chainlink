package vrf

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/utils"
)

// Seed represents a VRF seed as a serialized uint256
type Seed [32]byte

// BigToSeed returns seed x represented as a Seed, or an error if x is too big
func BigToSeed(x *big.Int) (Seed, error) {
	seed, err := utils.Uint256ToBytes(x)
	if err != nil {
		return Seed{}, err
	}
	return Seed(common.BytesToHash(seed)), nil
}

// Big returns the uint256 seed represented by s
func (s *Seed) Big() *big.Int {
	return common.Hash(*s).Big()
}

// BytesToSeed returns the Seed corresponding to b, or an error if b is too long
func BytesToSeed(b []byte) (*Seed, error) {
	if len(b) > 32 {
		return nil, errors.Errorf("Seed representation can be at most 32 bytes, "+
			"got %d", len(b))
	}
	seed := Seed(common.BytesToHash(b))
	return &seed, nil
}

// PreSeedData contains the data the VRF provider needs to compute the final VRF
// output and marshal the proof for transmission to the VRFCoordinator contract.
type PreSeedData struct {
	PreSeed   Seed        // Seed to be mixed with hash of containing block
	BlockHash common.Hash // Hash of block containing VRF request
	BlockNum  uint64      // Cardinal number of block containing VRF request
}

// FinalSeed is the seed which is actually passed to the VRF proof generator,
// given the preseed and the hash of the block in which the VRFCoordinator
// emitted the log for the request this is responding to.
func FinalSeed(s PreSeedData) (finalSeed *big.Int) {
	seedHashMsg := append(s.PreSeed[:], s.BlockHash.Bytes()...)
	return utils.MustHash(string(seedHashMsg)).Big()
}
