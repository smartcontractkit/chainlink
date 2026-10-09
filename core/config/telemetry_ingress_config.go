package config

import (
	"net/url"
	"time"
)

type TelemetryIngress interface {
	Logging() bool
	UniConn() bool
	BufferSize() uint
	MaxBatchSize() uint
	SendInterval() time.Duration
	SendTimeout() time.Duration
	UseBatchSend() bool
	Endpoints() []TelemetryIngressEndpoint
	ChipIngressEnabled() bool
	ChipIngressEndpoint() string
	ChipIngressInsecureConnection() bool
	ChipIngressBufferSize() uint
	ChipIngressMaxBatchSize() uint
	ChipIngressMaxConcurrentSends() int
	ChipIngressDrainTimeout() time.Duration
}

type TelemetryIngressEndpoint interface {
	Network() string
	ChainID() string
	ServerPubKey() string
	URL() *url.URL
}
