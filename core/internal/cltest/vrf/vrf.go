package cltest

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/services/vrf"
)

// SeedData returns the request data needed to construct/validate a VRF proof,
// modulo the key.
func SeedData(t *testing.T, preseed *big.Int, blockHash common.Hash,
	blockNum int) vrf.PreSeedData {
	seedAsSeed, err := vrf.BigToSeed(big.NewInt(0x10))
	require.NoError(t, err, "seed %x out of range", 0x10)
	return vrf.PreSeedData{
		PreSeed:   seedAsSeed,
		BlockNum:  uint64(blockNum),
		BlockHash: blockHash,
	}
}
