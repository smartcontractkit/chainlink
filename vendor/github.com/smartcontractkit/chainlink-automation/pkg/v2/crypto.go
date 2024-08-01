package ocr2keepers

import (
	"encoding/binary"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"golang.org/x/crypto/sha3"
)

// Generates a randomness source derived from the report timestamp (config, epoch, round) so
// that it's the same across the network for the same round
func getRandomKeySource(rt types.ReportTimestamp) [16]byte {
	// similar key building as libocr transmit selector
	hash := sha3.NewLegacyKeccak256()
	hash.Write(rt.ConfigDigest[:])
	temp := make([]byte, 8)
	binary.LittleEndian.PutUint64(temp, uint64(rt.Epoch))
	hash.Write(temp)
	binary.LittleEndian.PutUint64(temp, uint64(rt.Round))
	hash.Write(temp)

	var keyRandSource [16]byte
	copy(keyRandSource[:], hash.Sum(nil))
	return keyRandSource
}
