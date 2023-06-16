package v2

import "time"

type ocrConfig struct {
	c OCR
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
