package cosmwasm

import (
	"encoding/hex"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

// HexToByteArray is a wrapper for hex.DecodeString
func HexToByteArray(s string, b *[]byte) (err error) {
	*b, err = hex.DecodeString(s)
	return err
}

// HexToConfigDigest converts a hex string to ConfigDigest
func HexToConfigDigest(s string, digest *types.ConfigDigest) (err error) {
	// parse byte array encoded as hex string
	var byteArr []byte
	if err = HexToByteArray(s, &byteArr); err != nil {
		return
	}

	*digest, err = types.BytesToConfigDigest(byteArr)
	return
}
