package job

import (
	"net/url"
	"time"
)

//go:generate mockery --name Service --output ./mocks/ --case=underscore

type Service interface {
	Start() error
	Close() error
}

type Config interface {
	DatabaseMaximumTxDuration() time.Duration
	DatabaseURL() url.URL
	OCRBlockchainTimeout() time.Duration
	OCRContractConfirmations() uint16
	OCRContractPollInterval() time.Duration
	OCRContractSubscribeInterval() time.Duration
	OCRObservationTimeout() time.Duration
	TriggerFallbackDBPollInterval() time.Duration
}
