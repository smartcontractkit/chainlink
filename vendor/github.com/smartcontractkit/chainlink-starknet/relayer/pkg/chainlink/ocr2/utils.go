package ocr2

import (
	"encoding/binary"
	"math/big"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/starknet"
)

func parseEpochAndRound(felt *big.Int) (epoch uint32, round uint8) {
	var epochAndRound [starknet.FeltLength]byte
	felt.FillBytes(epochAndRound[:])
	epoch = binary.BigEndian.Uint32(epochAndRound[starknet.FeltLength-5 : starknet.FeltLength-1])
	round = epochAndRound[starknet.FeltLength-1]
	return epoch, round
}

/* Testing utils - do not use (XXX) outside testing context */

func XXXMustBytesToConfigDigest(b []byte) types.ConfigDigest {
	c, err := types.BytesToConfigDigest(b)
	if err != nil {
		panic(err)
	}
	return c
}
