package dydx

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

const ConfigDigestPrefixDydx types.ConfigDigestPrefix = 4

var _ types.OffchainConfigDigester = (*OffchainConfigDigester)(nil)

type OffchainConfigDigester struct {
	// endpoint type can be staging or prod
	endpointType string
}

func (d OffchainConfigDigester) ConfigDigest(cfg types.ContractConfig) (types.ConfigDigest, error) {
	return d.configDigest()
}

// The digest is unique per OffchainConfigDigester.endpointType value. This ensures
// protocol instances for staging vs prod are distinct, and we have separate
// monitoring for each.
func (d OffchainConfigDigester) configDigest() (types.ConfigDigest, error) {
	digest := types.ConfigDigest{}
	buf := sha256.New()

	if err := binary.Write(buf, binary.BigEndian, uint8(len(d.endpointType))); err != nil {
		return digest, err
	}

	rawHash := buf.Sum(nil)
	if n := copy(digest[:], rawHash[:]); n != len(digest) {
		return digest, fmt.Errorf("incorrect hash size %d, expected %d", n, len(digest))
	}

	binary.BigEndian.PutUint16(digest[0:2], uint16(d.ConfigDigestPrefix()))

	return digest, nil
}

// This should return the same constant value on every invocation
func (OffchainConfigDigester) ConfigDigestPrefix() types.ConfigDigestPrefix {
	return ConfigDigestPrefixDydx
}
