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

type telemetryIngressEndpointConfig struct {
	c toml.TelemetryIngressEndpoint
}

func (t *telemetryIngressConfig) Logging() bool {
	return *t.c.Logging
}

func (t *telemetryIngressConfig) UniConn() bool {
	return *t.c.UniConn
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

// Deprecated: Use TelemetryIngressEndpoint.ServerPubKey instead, if this field is set it will trigger an error, only used to warn NOPs of change
func (t *telemetryIngressConfig) ServerPubKey() string {
	return t.ServerPubKey()
}

// Deprecated: Use TelemetryIngressEndpoint.URL instead, if this field is set it will trigger an error, only used to warn NOPs of change
func (t *telemetryIngressConfig) URL() *url.URL {
	return t.URL()
}

func (t *telemetryIngressConfig) Endpoints() []config.TelemetryIngressEndpoint {
	var endpoints []config.TelemetryIngressEndpoint
	for _, e := range t.c.Endpoints {
		endpoints = append(endpoints, &telemetryIngressEndpointConfig{
			c: e,
		})
	}
	return endpoints
}

func (t *telemetryIngressEndpointConfig) Network() string {
	return *t.c.Network
}

func (t *telemetryIngressEndpointConfig) ChainID() string {
	return *t.c.ChainID
}

func (t *telemetryIngressEndpointConfig) URL() *url.URL {
	if t.c.URL.IsZero() {
		return nil
	}
	return t.c.URL.URL()
}

func (t *telemetryIngressEndpointConfig) ServerPubKey() string {
	return *t.c.ServerPubKey
}
