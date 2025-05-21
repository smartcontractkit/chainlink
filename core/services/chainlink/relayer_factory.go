package chainlink

import (
	"cmp"
	"errors"
	"fmt"
	"net/http"

	"github.com/pelletier/go-toml/v2"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	coretypes "github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana"
	solcfg "github.com/smartcontractkit/chainlink-solana/pkg/solana/config"

	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	coreconfig "github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/config/env"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo/retirement"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/dummy"
	evmrelay "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

type RelayerFactory struct {
	logger.Logger
	*plugins.LoopRegistry
	loop.GRPCOpts
	Registerer            prometheus.Registerer
	MercuryPool           wsrpc.Pool
	CapabilitiesRegistry  coretypes.CapabilitiesRegistry
	HTTPClient            *http.Client
	RetirementReportCache retirement.RetirementReportCache
}

type DummyFactoryConfig struct {
	ChainID string
}

func (r *RelayerFactory) NewDummy(config DummyFactoryConfig) (loop.Relayer, error) {
	return dummy.NewRelayer(r.Logger, config.ChainID), nil
}

type EVMFactoryConfig struct {
	legacyevm.ChainOpts
	EthKeystore   keystore.Eth
	CSAKeystore   coretypes.Keystore
	MercuryConfig coreconfig.Mercury
}

func (r *RelayerFactory) NewEVM(config EVMFactoryConfig) (map[types.RelayID]evmrelay.LOOPRelayAdapter, error) {
	// TODO impl EVM loop. For now always 'fallback' to an adapter and embedded chain
	relayers := make(map[types.RelayID]evmrelay.LOOPRelayAdapter)
	lggr := r.Logger.Named("EVM")

	legacyChains, err := evmrelay.NewLegacyChains(lggr, config.EthKeystore, config.ChainOpts)
	if err != nil {
		return nil, err
	}
	newChainStore := config.GenChainStore
	if newChainStore == nil {
		newChainStore = keys.NewChainStore
	}
	for _, chain := range legacyChains {
		relayID := types.RelayID{Network: relay.NetworkEVM, ChainID: chain.ID().String()}
		chain := chain

		relayerOpts := evmrelay.RelayerOpts{
			DS:                    config.DS,
			Registerer:            r.Registerer,
			EVMKeystore:           newChainStore(keystore.NewEthSigner(config.EthKeystore, chain.ID()), chain.ID()),
			CSAKeystore:           config.CSAKeystore,
			MercuryPool:           r.MercuryPool,
			MercuryConfig:         config.MercuryConfig,
			CapabilitiesRegistry:  r.CapabilitiesRegistry,
			HTTPClient:            r.HTTPClient,
			RetirementReportCache: r.RetirementReportCache,
		}
		relayer, err2 := evmrelay.NewRelayer(lggr.Named(relayID.ChainID), chain, relayerOpts)
		if err2 != nil {
			err = errors.Join(err, err2)
			continue
		}

		relayers[relayID] = evmrelay.NewLOOPRelayAdapter(relayer)
	}

	// always return err because it is accumulating individual errors
	return relayers, err
}

type SolanaFactoryConfig struct {
	solcfg.TOMLConfigs
	DS sqlutil.DataSource
}

func (r *RelayerFactory) NewSolana(ks coretypes.Keystore, config SolanaFactoryConfig) (map[types.RelayID]loop.Relayer, error) {
	chainCfgs, ds := config.TOMLConfigs, config.DS
	solanaRelayers := make(map[types.RelayID]loop.Relayer)
	var solLggr = r.Logger.Named("Solana")

	unique := make(map[string]struct{})
	// create one relayer per chain id
	for _, chainCfg := range chainCfgs {
		relayID := types.RelayID{Network: relay.NetworkSolana, ChainID: *chainCfg.ChainID}
		_, alreadyExists := unique[relayID.Name()]
		if alreadyExists {
			return nil, fmt.Errorf("duplicate chain definitions for %s", relayID.Name())
		}
		unique[relayID.Name()] = struct{}{}

		// skip disabled chains from further processing
		if !chainCfg.IsEnabled() {
			solLggr.Warnw("Skipping disabled chain", "id", chainCfg.ChainID)
			continue
		}

		lggr := solLggr.Named(relayID.ChainID)

		if cmdName := env.SolanaPlugin.Cmd.Get(); cmdName != "" {
			// setup the solana relayer to be a LOOP
			cfgTOML, err := toml.Marshal(struct {
				Solana solcfg.TOMLConfig
			}{Solana: *chainCfg})
			if err != nil {
				return nil, fmt.Errorf("failed to marshal Solana configs: %w", err)
			}
			envVars, err := plugins.ParseEnvFile(env.SolanaPlugin.Env.Get())
			if err != nil {
				return nil, fmt.Errorf("failed to parse Solana env file: %w", err)
			}
			solCmdFn, err := plugins.NewCmdFactory(r.Register, plugins.CmdConfig{
				ID:  relayID.Name(),
				Cmd: cmdName,
				Env: envVars,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to create Solana LOOP command: %w", err)
			}

			solanaRelayers[relayID] = loop.NewRelayerService(lggr, r.GRPCOpts, solCmdFn, string(cfgTOML), ks, r.CapabilitiesRegistry)
		} else {
			// fallback to embedded chain
			opts := solana.ChainOpts{
				Logger:   lggr,
				KeyStore: ks,
				DS:       ds,
			}

			chain, err := solana.NewChain(chainCfg, opts)
			if err != nil {
				return nil, err
			}
			solanaRelayers[relayID] = relay.NewServerAdapter(solana.NewRelayer(lggr, chain, r.CapabilitiesRegistry))
		}
	}
	return solanaRelayers, nil
}

func (r *RelayerFactory) NewStarkNet(ks coretypes.Keystore, chainCfgs RawConfigs) (map[types.RelayID]loop.Relayer, error) {
	return r.NewLOOPRelayer("StarkNet", relay.NetworkStarkNet, env.StarknetPlugin, ks, chainCfgs)
}

type CosmosFactoryConfig struct {
	Keystore    keystore.Cosmos
	TOMLConfigs RawConfigs
}

func (c CosmosFactoryConfig) Validate() error {
	var err error
	if c.Keystore == nil {
		err = errors.Join(err, errors.New("nil Keystore"))
	}
	if len(c.TOMLConfigs) == 0 {
		err = errors.Join(err, errors.New("no CosmosConfigs provided"))
	}

	if err != nil {
		err = fmt.Errorf("invalid CosmosFactoryConfig: %w", err)
	}
	return err
}

func (r *RelayerFactory) NewCosmos(ks coretypes.Keystore, chainCfgs RawConfigs) (map[types.RelayID]loop.Relayer, error) {
	return r.NewLOOPRelayer("Cosmos", relay.NetworkCosmos, env.CosmosPlugin, ks, chainCfgs)
}

func (r *RelayerFactory) NewAptos(ks coretypes.Keystore, chainCfgs RawConfigs) (map[types.RelayID]loop.Relayer, error) {
	return r.NewLOOPRelayer("Aptos", relay.NetworkAptos, env.AptosPlugin, ks, chainCfgs)
}

func (r *RelayerFactory) NewLOOPRelayer(name string, network string, plugin env.Plugin, ks coretypes.Keystore, chainCfgs RawConfigs) (map[types.RelayID]loop.Relayer, error) {
	relayers := make(map[types.RelayID]loop.Relayer)
	lggr := r.Logger.Named(name)

	cmdName := cmp.Or(plugin.Cmd.Get(), plugin.CmdDefault)
	if cmdName == "" {
		return nil, fmt.Errorf("plugin command not defined: %s", plugin.Cmd)
	}
	envFile := plugin.Env.Get()
	envVars, err := plugins.ParseEnvFile(envFile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse env file %s: %w", envFile, err)
	}

	unique := make(map[string]struct{})
	// create one relayer per chain id
	for _, chainCfg := range chainCfgs {
		relayID := types.RelayID{Network: network, ChainID: chainCfg.ChainID()}
		if _, alreadyExists := unique[relayID.Name()]; alreadyExists {
			return nil, fmt.Errorf("duplicate chain definitions for %s", relayID.Name())
		}
		unique[relayID.Name()] = struct{}{}

		// skip disabled chains from further processing
		if !chainCfg.IsEnabled() {
			lggr.Warnw("Skipping disabled chain", "id", relayID.ChainID)
			continue
		}

		cmdFn, err := plugins.NewCmdFactory(r.Register, plugins.CmdConfig{
			ID: relayID.Name(), Cmd: cmdName, Env: envVars,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create LOOP command: %w", err)
		}
		cfgTOML, err := toml.Marshal(chainCfg)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal configs: %w", err)
		}
		// the relayer service has a delicate keystore dependency. the value that is passed to NewRelayerService must
		// be compatible with instantiating a starknet transaction manager KeystoreAdapter within the LOOPp executable.
		relayers[relayID] = loop.NewRelayerService(lggr.Named(relayID.ChainID), r.GRPCOpts, cmdFn, string(cfgTOML), ks, r.CapabilitiesRegistry)
	}
	return relayers, nil
}

func (r *RelayerFactory) NewTron(ks coretypes.Keystore, chainCfgs RawConfigs) (map[types.RelayID]loop.Relayer, error) {
	return r.NewLOOPRelayer("Tron", relay.NetworkTron, env.TronPlugin, ks, chainCfgs)
}
