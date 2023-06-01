package config

import (
	"net/url"
	"time"
)

type TelemetryIngress interface {
	TelemetryIngressLogging() bool
	TelemetryIngressUniConn() bool
	TelemetryIngressServerPubKey() string
	TelemetryIngressURL() *url.URL
	TelemetryIngressBufferSize() uint
	TelemetryIngressMaxBatchSize() uint
	TelemetryIngressSendInterval() time.Duration
	TelemetryIngressSendTimeout() time.Duration
	TelemetryIngressUseBatchSend() bool
}
