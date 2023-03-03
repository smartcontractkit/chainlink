package ocr2

import (
	"encoding/binary"
	"math/big"

	junotypes "github.com/NethermindEth/juno/pkg/types"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

func parseEpochAndRound(felt *big.Int) (epoch uint32, round uint8) {
	var epochAndRound [junotypes.FeltLength]byte
	felt.FillBytes(epochAndRound[:])
	epoch = binary.BigEndian.Uint32(epochAndRound[junotypes.FeltLength-5 : junotypes.FeltLength-1])
	round = epochAndRound[junotypes.FeltLength-1]
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
