package txm

import "time"

//go:generate mockery --name Config --output ./mocks/ --case=underscore --filename config.go

// txm config
type Config interface {
	TxTimeout() time.Duration
	TxSendFrequency() time.Duration
	TxMaxBatchSize() int
}
