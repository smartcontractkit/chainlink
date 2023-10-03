package config

import (
	"net/url"
	"time"
)

//go:generate mockery --quiet --name TelemetryIngress --output ./mocks/ --case=underscore --filename telemetry_ingress.go

type TelemetryIngress interface {
	Logging() bool
	UniConn() bool
	BufferSize() uint
	MaxBatchSize() uint
	SendInterval() time.Duration
	SendTimeout() time.Duration
	UseBatchSend() bool
	Endpoints() []TelemetryIngressEndpoint

	ServerPubKey() string // Deprecated: Use TelemetryIngressEndpoint.ServerPubKey instead, if this field is set it will trigger an error, only used to warn NOPs of change
	URL() *url.URL        // Deprecated: Use TelemetryIngressEndpoint.URL instead, if this field is set it will trigger an error, only used to warn NOPs of change
}

//go:generate mockery --quiet --name TelemetryIngressEndpoint --output ./mocks/ --case=underscore --filename telemetry_ingress_endpoint.go
type TelemetryIngressEndpoint interface {
	Network() string
	ChainID() string
	ServerPubKey() string
	URL() *url.URL
}
