package chainlink

import (
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
)

var _ config.OCR2 = (*ocr2Config)(nil)

type ocr2Config struct {
	c toml.OCR2
}

func (o *ocr2Config) Enabled() bool {
	return *o.c.Enabled
}

func (o *ocr2Config) ContractConfirmations() uint16 {
	return uint16(*o.c.ContractConfirmations)
}

func (o *ocr2Config) ContractTransmitterTransmitTimeout() time.Duration {
	return o.c.ContractTransmitterTransmitTimeout.Duration()
}

func (o *ocr2Config) BlockchainTimeout() time.Duration {
	return o.c.BlockchainTimeout.Duration()
}

func (o *ocr2Config) DatabaseTimeout() time.Duration {
	return o.c.DatabaseTimeout.Duration()
}

func (o *ocr2Config) ContractPollInterval() time.Duration {
	return o.c.ContractPollInterval.Duration()
}

func (o *ocr2Config) ContractSubscribeInterval() time.Duration {
	return o.c.ContractSubscribeInterval.Duration()
}

func (o *ocr2Config) KeyBundleID() (string, error) {
	b := o.c.KeyBundleID
	if *b == zeroSha256Hash {
		return "", nil
	}
	return b.String(), nil
}

func (o *ocr2Config) TraceLogging() bool {
	return *o.c.TraceLogging
}

func (o *ocr2Config) CaptureEATelemetry() bool {
	return *o.c.CaptureEATelemetry
}

func (o *ocr2Config) CaptureAutomationCustomTelemetry() bool {
	return *o.c.CaptureAutomationCustomTelemetry
}

func (o *ocr2Config) DefaultTransactionQueueDepth() uint32 {
	return *o.c.DefaultTransactionQueueDepth
}

func (o *ocr2Config) SimulateTransactions() bool {
	return *o.c.SimulateTransactions
}
