package chainlink

import (
	"time"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
)

var _ config.OCR = (*ocrConfig)(nil)

type ocrConfig struct {
	c toml.OCR
}

func (o *ocrConfig) Enabled() bool {
	return *o.c.Enabled
}

func (o *ocrConfig) BlockchainTimeout() time.Duration {
	return o.c.BlockchainTimeout.Duration()
}

func (o *ocrConfig) ContractPollInterval() time.Duration {
	return o.c.ContractPollInterval.Duration()
}

func (o *ocrConfig) ContractSubscribeInterval() time.Duration {
	return o.c.ContractSubscribeInterval.Duration()
}

func (o *ocrConfig) KeyBundleID() (string, error) {
	b := o.c.KeyBundleID
	if *b == zeroSha256Hash {
		return "", nil
	}
	return b.String(), nil
}

func (o *ocrConfig) ObservationTimeout() time.Duration {
	return o.c.ObservationTimeout.Duration()
}

func (o *ocrConfig) SimulateTransactions() bool {
	return *o.c.SimulateTransactions
}

func (o *ocrConfig) TransmitterAddress() (types.EIP55Address, error) {
	a := *o.c.TransmitterAddress
	if a.IsZero() {
		return a, errors.Wrap(config.ErrEnvUnset, "OCR.TransmitterAddress is not set")
	}
	return a, nil
}

func (o *ocrConfig) TraceLogging() bool {
	return *o.c.TraceLogging
}

func (o *ocrConfig) DefaultTransactionQueueDepth() uint32 {
	return *o.c.DefaultTransactionQueueDepth
}

func (o *ocrConfig) CaptureEATelemetry() bool {
	return *o.c.CaptureEATelemetry
}
