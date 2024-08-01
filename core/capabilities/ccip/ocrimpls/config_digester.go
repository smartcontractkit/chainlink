package ocrimpls

import "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

type configDigester struct {
	d types.ConfigDigest
}

func NewConfigDigester(d types.ConfigDigest) *configDigester {
	return &configDigester{d: d}
}

// ConfigDigest implements types.OffchainConfigDigester.
func (c *configDigester) ConfigDigest(types.ContractConfig) (types.ConfigDigest, error) {
	return c.d, nil
}

// ConfigDigestPrefix implements types.OffchainConfigDigester.
func (c *configDigester) ConfigDigestPrefix() (types.ConfigDigestPrefix, error) {
	return types.ConfigDigestPrefixCCIPMultiRole, nil
}

var _ types.OffchainConfigDigester = (*configDigester)(nil)
