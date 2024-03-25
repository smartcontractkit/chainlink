package pluginprovider

import (
	"context"
	"fmt"

	libocr "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/assert"
)

type staticOffchainConfigDigesterConfig struct {
	contractConfig     libocr.ContractConfig
	configDigest       libocr.ConfigDigest
	configDigestPrefix libocr.ConfigDigestPrefix
}

// staticOffchainConfigDigester is a static implementation of OffchainConfigDigesterEvaluator
type staticOffchainConfigDigester struct {
	staticOffchainConfigDigesterConfig
}

var _ libocr.OffchainConfigDigester = staticOffchainConfigDigester{}

func (s staticOffchainConfigDigester) ConfigDigest(ctx context.Context, config libocr.ContractConfig) (libocr.ConfigDigest, error) {
	if !assert.ObjectsAreEqual(s.contractConfig, config) {
		return libocr.ConfigDigest{}, fmt.Errorf("expected contract config %v but got %v", s.configDigest, config)
	}
	return s.configDigest, nil
}

func (s staticOffchainConfigDigester) ConfigDigestPrefix(ctx context.Context) (libocr.ConfigDigestPrefix, error) {
	return s.configDigestPrefix, nil
}

func (s staticOffchainConfigDigester) Evaluate(ctx context.Context, ocd libocr.OffchainConfigDigester) error {
	gotDigestPrefix, err := ocd.ConfigDigestPrefix(ctx)
	if err != nil {
		return fmt.Errorf("failed to get ConfigDigestPrefix: %w", err)
	}
	if gotDigestPrefix != s.configDigestPrefix {
		return fmt.Errorf("expected ConfigDigestPrefix %x but got %x", s.configDigestPrefix, gotDigestPrefix)
	}
	gotDigest, err := ocd.ConfigDigest(ctx, contractConfig)
	if err != nil {
		return fmt.Errorf("failed to get ConfigDigest: %w", err)
	}
	if gotDigest != s.configDigest {
		return fmt.Errorf("expected ConfigDigest %x but got %x", s.configDigest, gotDigest)
	}
	return nil
}
