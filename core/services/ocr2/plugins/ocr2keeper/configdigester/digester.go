package configdigester

import "github.com/smartcontractkit/libocr/offchainreporting2/types"

type configDigester struct {
	types.OffchainConfigDigester
	instance uint8
}

func New(base types.OffchainConfigDigester, instance uint8) types.OffchainConfigDigester {
	return &configDigester{
		OffchainConfigDigester: base,
		instance:               instance,
	}
}

// Compute ConfigDigest for the given ContractConfig. The first two bytes of the
// ConfigDigest must be the big-endian encoding of ConfigDigestPrefix!
func (c *configDigester) ConfigDigest(cfg types.ContractConfig) (types.ConfigDigest, error) {
	digest, err := c.OffchainConfigDigester.ConfigDigest(cfg)
	if err != nil {
		return digest, err
	}

	return c.thresholdDigitalDigest(digest), nil
}

// This should return the same constant value on every invocation
func (c *configDigester) ConfigDigestPrefix() types.ConfigDigestPrefix {
	return c.OffchainConfigDigester.ConfigDigestPrefix()
}

func (c *configDigester) thresholdDigitalDigest(root types.ConfigDigest) types.ConfigDigest {
	return root
	var thresholdBytes types.ConfigDigest
	for i, b := range root[:2] {
		thresholdBytes[i] = b
	}
	for i, b := range root[2:] {
		thresholdBytes[i+2] = b ^ c.instance
	}
	return thresholdBytes
}
