package config

import (
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
)

type ocrConfig struct {
	c toml.OCR
}

func (o *ocrConfig) ContractConfirmations() uint16 {
	return *o.c.ContractConfirmations
}

func (o *ocrConfig) ContractTransmitterTransmitTimeout() time.Duration {
	return o.c.ContractTransmitterTransmitTimeout.Duration()
}

func (o *ocrConfig) ObservationGracePeriod() time.Duration {
	return o.c.ObservationGracePeriod.Duration()
}

func (o *ocrConfig) DatabaseTimeout() time.Duration {
	return o.c.DatabaseTimeout.Duration()
}

func (o *ocrConfig) DeltaCOverride() time.Duration {
	return o.c.DeltaCOverride.Duration()
}

func (o *ocrConfig) DeltaCJitterOverride() time.Duration {
	return o.c.DeltaCJitterOverride.Duration()
}
