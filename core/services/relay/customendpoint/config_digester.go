package customendpoint

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

const ConfigDigestPrefixCustomEndpoint types.ConfigDigestPrefix = 4

var _ types.OffchainConfigDigester = (*OffchainConfigDigester)(nil)

type OffchainConfigDigester struct {
	// This uniquely identifies a custom endpoint class. For example, dydx.
	EndpointName string
	// Endpoint class specific target. Example, staging/prod
	EndpointTarget string
	// Uniquely identifies the type of data being uploaded to the endpoint
	// For example, ETHUSD represents ETH price in USD.
	PayloadType string
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

	if _, err := buf.Write([]byte(d.EndpointName)); err != nil {
		return digest, err
	}
	if _, err := buf.Write([]byte(d.EndpointTarget)); err != nil {
		return digest, err
	}
	if _, err := buf.Write([]byte(d.PayloadType)); err != nil {
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
	return ConfigDigestPrefixCustomEndpoint
}
