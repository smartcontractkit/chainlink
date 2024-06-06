package ccipcapability

import (
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

var (
	_ ocrtypes.OffchainConfigDigester = (*offchainConfigDigester)(nil)
)

type offchainConfigDigester struct{}

// Compute ConfigDigest for the given ContractConfig. The first two bytes of the
// ConfigDigest must be the big-endian encoding of ConfigDigestPrefix!
func (o *offchainConfigDigester) ConfigDigest(ocrtypes.ContractConfig) (ocrtypes.ConfigDigest, error) {
	// TODO: implement.
	// Config digest calculation seems to involve things that are not in the ContractConfig.
	panic("unimplemented")
}

// This should return the same constant value on every invocation
func (o *offchainConfigDigester) ConfigDigestPrefix() (ocrtypes.ConfigDigestPrefix, error) {
	return ocrtypes.ConfigDigestPrefixEVM, nil
}
