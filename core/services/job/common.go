package job

import (
	"net/url"
	"time"
)

//go:generate mockery --name Service --output ./mocks/ --case=underscore

type Type string

func (t Type) String() string {
	return string(t)
}

type Service interface {
	Start() error
	Close() error
}

type Config interface {
	DatabaseMaximumTxDuration() time.Duration
	DatabaseURL() url.URL
	TriggerFallbackDBPollInterval() time.Duration
	JobPipelineParallelism() uint8
}
