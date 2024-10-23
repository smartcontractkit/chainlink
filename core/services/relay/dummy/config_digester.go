package dummy

import (
	"context"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

// Stub config digester that uses a static config digest

type configDigester struct {
	configDigest ocrtypes.ConfigDigest
}

func NewOffchainConfigDigester(cd ocrtypes.ConfigDigest) (ocrtypes.OffchainConfigDigester, error) {
	return &configDigester{cd}, nil
}

// Compute ConfigDigest for the given ContractConfig. The first two bytes of the
// ConfigDigest must be the big-endian encoding of ConfigDigestPrefix!
func (cd *configDigester) ConfigDigest(context.Context, ocrtypes.ContractConfig) (ocrtypes.ConfigDigest, error) {
	return cd.configDigest, nil
}

// This should return the same constant value on every invocation
func (cd *configDigester) ConfigDigestPrefix(ctx context.Context) (ocrtypes.ConfigDigestPrefix, error) {
	return ocrtypes.ConfigDigestPrefixFromConfigDigest(cd.configDigest), nil
}
