package chainlink

import (
	"net/url"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
)

var _ config.TelemetryIngress = (*telemetryIngressConfig)(nil)

type telemetryIngressConfig struct {
	c toml.TelemetryIngress
}

func (t *telemetryIngressConfig) Logging() bool {
	return *t.c.Logging
}

func (t *telemetryIngressConfig) UniConn() bool {
	return *t.c.UniConn
}

func (t *telemetryIngressConfig) ServerPubKey() string {
	return *t.c.ServerPubKey
}

func (t *telemetryIngressConfig) URL() *url.URL {
	if t.c.URL.IsZero() {
		return nil
	}
	return t.c.URL.URL()
}

func (t *telemetryIngressConfig) BufferSize() uint {
	return uint(*t.c.BufferSize)
}

func (t *telemetryIngressConfig) MaxBatchSize() uint {
	return uint(*t.c.MaxBatchSize)
}

func (t *telemetryIngressConfig) SendInterval() time.Duration {
	return t.c.SendInterval.Duration()
}

func (t *telemetryIngressConfig) SendTimeout() time.Duration {
	return t.c.SendTimeout.Duration()
}

func (t *telemetryIngressConfig) UseBatchSend() bool {
	return *t.c.UseBatchSend
}
