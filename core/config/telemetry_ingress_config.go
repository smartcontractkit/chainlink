package config

import (
	"net/url"
	"time"
)

type TelemetryIngress interface {
	Logging() bool
	UniConn() bool
	ServerPubKey() string
	URL() *url.URL
	BufferSize() uint
	MaxBatchSize() uint
	SendInterval() time.Duration
	SendTimeout() time.Duration
	UseBatchSend() bool
}
