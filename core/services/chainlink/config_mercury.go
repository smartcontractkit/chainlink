package chainlink

import (
	"time"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
)

var _ config.MercuryCache = (*mercuryCacheConfig)(nil)

type mercuryCacheConfig struct {
	c toml.MercuryCache
}

func (m *mercuryCacheConfig) LatestReportTTL() time.Duration {
	return m.c.LatestReportTTL.Duration()
}
func (m *mercuryCacheConfig) MaxStaleAge() time.Duration {
	return m.c.MaxStaleAge.Duration()
}
func (m *mercuryCacheConfig) LatestReportDeadline() time.Duration {
	return m.c.LatestReportDeadline.Duration()
}

var _ config.MercuryTLS = (*mercuryTLSConfig)(nil)

type mercuryTLSConfig struct {
	c toml.MercuryTLS
}

func (m *mercuryTLSConfig) CertFile() string {
	return *m.c.CertFile
}

var _ config.MercuryTransmitter = (*mercuryTransmitterConfig)(nil)

type mercuryTransmitterConfig struct {
	c toml.MercuryTransmitter
}

func (m *mercuryTransmitterConfig) Protocol() config.MercuryTransmitterProtocol {
	return *m.c.Protocol
}

func (m *mercuryTransmitterConfig) TransmitQueueMaxSize() uint32 {
	return *m.c.TransmitQueueMaxSize
}

func (m *mercuryTransmitterConfig) TransmitTimeout() commonconfig.Duration {
	return *m.c.TransmitTimeout
}

func (m *mercuryTransmitterConfig) TransmitConcurrency() uint32 {
	return *m.c.TransmitConcurrency
}

func (m *mercuryTransmitterConfig) ReaperFrequency() commonconfig.Duration {
	return *m.c.ReaperFrequency
}

func (m *mercuryTransmitterConfig) ReaperMaxAge() commonconfig.Duration {
	return *m.c.ReaperMaxAge
}

type mercuryConfig struct {
	c toml.Mercury
	s toml.MercurySecrets
}

func (m *mercuryConfig) Credentials(credName string) *types.MercuryCredentials {
	if mc, ok := m.s.Credentials[credName]; ok {
		c := &types.MercuryCredentials{
			URL:      mc.URL.URL().String(),
			Password: string(*mc.Password),
			Username: string(*mc.Username),
		}
		if mc.LegacyURL != nil && mc.LegacyURL.URL() != nil {
			c.LegacyURL = mc.LegacyURL.URL().String()
		}
		return c
	}
	return nil
}

func (m *mercuryConfig) Cache() config.MercuryCache {
	return &mercuryCacheConfig{c: m.c.Cache}
}

func (m *mercuryConfig) TLS() config.MercuryTLS {
	return &mercuryTLSConfig{c: m.c.TLS}
}

func (m *mercuryConfig) Transmitter() config.MercuryTransmitter {
	return &mercuryTransmitterConfig{c: m.c.Transmitter}
}

func (m *mercuryConfig) VerboseLogging() bool {
	return *m.c.VerboseLogging
}
