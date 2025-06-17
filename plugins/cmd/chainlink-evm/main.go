package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-plugin"
	"github.com/pelletier/go-toml/v2"
	"github.com/prometheus/client_golang/prometheus"

	clhttp "github.com/smartcontractkit/chainlink-common/pkg/http"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox"
	"github.com/smartcontractkit/chainlink-evm/pkg/chains/legacyevm"
	evmcfg "github.com/smartcontractkit/chainlink-evm/pkg/config/toml"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo/retirement"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/cache"
)

func main() {
	s := loop.MustNewStartedServer("PluginEVM")
	defer s.Stop()

	p := &pluginRelayer{EnvConfig: s.EnvConfig, Plugin: loop.Plugin{Logger: s.Logger}, DataSource: s.DataSource}
	defer s.Logger.ErrorIfFn(p.Close, "Failed to close")

	s.MustRegister(p)

	stopCh := make(chan struct{})
	defer close(stopCh)

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: loop.PluginRelayerHandshakeConfig(),
		Plugins: map[string]plugin.Plugin{
			loop.PluginRelayerName: &loop.GRPCPluginRelayer{
				PluginServer: p,
				BrokerConfig: loop.BrokerConfig{
					StopCh:   stopCh,
					Logger:   s.Logger,
					GRPCOpts: s.GRPCOpts,
				},
			},
		},
		GRPCServer: s.GRPCOpts.NewServer,
	})
}

type pluginRelayer struct {
	loop.EnvConfig
	loop.Plugin
	sqlutil.DataSource
}

func (c *pluginRelayer) NewRelayer(ctx context.Context, configTOML string, keystore, csaKeystore core.Keystore, capRegistry core.CapabilitiesRegistry) (loop.Relayer, error) {
	d := toml.NewDecoder(strings.NewReader(configTOML))
	d.DisallowUnknownFields()
	var cfg struct {
		EVM evmcfg.EVMConfig
	}

	if err := d.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config toml: %w:\n\t%s", err, configTOML)
	}

	// TODO validate?

	evmKeystore := keys.NewChainStore(keystore, cfg.EVM.ChainID.ToInt())

	mailMon := mailbox.NewMonitor(c.AppID, logger.Named(c.Logger, "Mailbox"))
	c.SubService(mailMon)

	chain, err := legacyevm.NewTOMLChain(&cfg.EVM, legacyevm.ChainRelayOpts{
		Logger:   c.Logger,
		KeyStore: evmKeystore,
		ChainOpts: legacyevm.ChainOpts{
			ChainConfigs: evmcfg.EVMConfigs{&cfg.EVM},
			DatabaseConfig: &DatabaseConfig{
				defaultQueryTimeout: c.DatabaseQueryTimeout,
				logSQL:              c.DatabaseLogSQL,
			},
			FeatureConfig: &FeatureConfig{
				logPoller: c.FeatureLogPoller,
			},
			ListenerConfig: &ListenerConfig{
				fallbackPollInterval: c.DatabaseListenerFallbackPollInterval,
			},
			MailMon: mailMon,
			DS:      c.DataSource,
		},
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create chain: %w", err)
	}

	ra, err := evm.NewRelayer(c.Logger, chain, evm.RelayerOpts{
		DS:          c.DataSource,
		Registerer:  prometheus.DefaultRegisterer,
		EVMKeystore: evmKeystore,
		CSAKeystore: csaKeystore,
		MercuryPool: wsrpc.NewPool(c.Logger, cache.Config{
			LatestReportTTL:      c.MercuryCacheLatestReportTTL,
			MaxStaleAge:          c.MercuryCacheMaxStaleAge,
			LatestReportDeadline: c.MercuryCacheLatestReportDeadline,
		}),
		RetirementReportCache: retirement.NewRetirementReportCache(c.Logger, c.DataSource),
		MercuryConfig: &MercuryConfig{
			transmitter: &Transmitter{
				protocol:             config.MercuryTransmitterProtocol(c.MercuryTransmitterProtocol),
				transmitQueueMaxSize: c.MercuryTransmitterTransmitQueueMaxSize,
				transmitTimeout:      c.MercuryTransmitterTransmitTimeout,
				transmitConcurrency:  c.MercuryTransmitterTransmitConcurrency,
				reaperFrequency:      c.MercuryTransmitterReaperFrequency,
				reaperMaxAge:         c.MercuryTransmitterReaperMaxAge,
			},
			verboseLogging: c.MercuryVerboseLogging,
		},
		CapabilitiesRegistry: capRegistry,
		HTTPClient:           clhttp.NewUnrestrictedClient(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create relayer: %w", err)
	}

	c.SubService(ra)

	return ra, nil
}

type DatabaseConfig struct {
	defaultQueryTimeout time.Duration
	logSQL              bool
}

func (d *DatabaseConfig) DefaultQueryTimeout() time.Duration {
	return d.defaultQueryTimeout
}

func (d *DatabaseConfig) LogSQL() bool {
	return d.logSQL
}

type FeatureConfig struct {
	logPoller bool
}

func (f *FeatureConfig) LogPoller() bool {
	return f.logPoller
}

type ListenerConfig struct {
	fallbackPollInterval time.Duration
}

func (l *ListenerConfig) FallbackPollInterval() time.Duration {
	return l.fallbackPollInterval
}

type MercuryConfig struct {
	transmitter    *Transmitter
	verboseLogging bool
}

func (m *MercuryConfig) Transmitter() config.MercuryTransmitter {
	return m.transmitter
}

func (m *MercuryConfig) VerboseLogging() bool {
	return m.verboseLogging
}

type Transmitter struct {
	protocol             config.MercuryTransmitterProtocol
	transmitQueueMaxSize uint32
	transmitTimeout      time.Duration
	transmitConcurrency  uint32
	reaperFrequency      time.Duration
	reaperMaxAge         time.Duration
}

func (t *Transmitter) Protocol() config.MercuryTransmitterProtocol {
	return t.protocol
}

func (t *Transmitter) TransmitQueueMaxSize() uint32 {
	return t.transmitQueueMaxSize
}

func (t *Transmitter) TransmitTimeout() time.Duration {
	return t.transmitTimeout
}

func (t *Transmitter) TransmitConcurrency() uint32 {
	return t.transmitConcurrency
}

func (t *Transmitter) ReaperFrequency() time.Duration {
	return t.reaperFrequency
}

func (t *Transmitter) ReaperMaxAge() time.Duration {
	return t.reaperMaxAge
}
