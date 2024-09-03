package chainlink

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/pelletier/go-toml/v2"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/adapters/relay"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	coretypes "github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos"
	coscfg "github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/config"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana"
	solcfg "github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
	pkgstarknet "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink"
	starkchain "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/chain"
	starkcfg "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	coreconfig "github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/config/env"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	corerelay "github.com/smartcontractkit/chainlink/v2/core/services/relay"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/dummy"
	evmrelay "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

type RelayerFactory struct {
	logger.Logger
	*plugins.LoopRegistry
	loop.GRPCOpts
	MercuryPool          wsrpc.Pool
	CapabilitiesRegistry coretypes.CapabilitiesRegistry
	HTTPClient           *http.Client
}

type DummyFactoryConfig struct {
	ChainID string
}

func (r *RelayerFactory) NewDummy(config DummyFactoryConfig) (loop.Relayer, error) {
	return dummy.NewRelayer(r.Logger, config.ChainID), nil
}

type EVMFactoryConfig struct {
	legacyevm.ChainOpts
	evmrelay.CSAETHKeystore
	coreconfig.MercuryTransmitter
}

func (r *RelayerFactory) NewEVM(ctx context.Context, config EVMFactoryConfig) (map[types.RelayID]evmrelay.LoopRelayAdapter, error) {
	// TODO impl EVM loop. For now always 'fallback' to an adapter and embedded chain

	relayers := make(map[types.RelayID]evmrelay.LoopRelayAdapter)

	lggr := r.Logger.Named("EVM")

	// override some common opts with the factory values. this seems weird... maybe other signatures should change, or this should take a different type...
	ccOpts := legacyevm.ChainRelayExtenderConfig{
		Logger:    lggr,
		KeyStore:  config.CSAETHKeystore.Eth(),
		ChainOpts: config.ChainOpts,
	}

	evmRelayExtenders, err := evmrelay.NewChainRelayerExtenders(ctx, ccOpts)
	if err != nil {
		return nil, err
	}
	legacyChains := evmrelay.NewLegacyChainsFromRelayerExtenders(evmRelayExtenders)
	for _, ext := range evmRelayExtenders.Slice() {
		relayID := types.RelayID{Network: corerelay.NetworkEVM, ChainID: ext.Chain().ID().String()}
		chain, err2 := legacyChains.Get(relayID.ChainID)
		if err2 != nil {
			return nil, err2
		}

		relayerOpts := evmrelay.RelayerOpts{
			DS:                   ccOpts.DS,
			CSAETHKeystore:       config.CSAETHKeystore,
			MercuryPool:          r.MercuryPool,
			TransmitterConfig:    config.MercuryTransmitter,
			CapabilitiesRegistry: r.CapabilitiesRegistry,
			HTTPClient:           r.HTTPClient,
		}
		relayer, err2 := evmrelay.NewRelayer(lggr.Named(relayID.ChainID), chain, relayerOpts)
		if err2 != nil {
			err = errors.Join(err, err2)
			continue
		}

		relayers[relayID] = evmrelay.NewLoopRelayServerAdapter(relayer, ext)
	}

	// always return err because it is accumulating individual errors
	return relayers, err
}

type SolanaFactoryConfig struct {
	Keystore keystore.Solana
	solcfg.TOMLConfigs
}

func (r *RelayerFactory) NewSolana(ks keystore.Solana, chainCfgs solcfg.TOMLConfigs) (map[types.RelayID]loop.Relayer, error) {
	solanaRelayers := make(map[types.RelayID]loop.Relayer)
	var (
		solLggr = r.Logger.Named("Solana")
		signer  = &keystore.SolanaSigner{Solana: ks}
	)

	unique := make(map[string]struct{})
	// create one relayer per chain id
	for _, chainCfg := range chainCfgs {
		relayID := types.RelayID{Network: corerelay.NetworkSolana, ChainID: *chainCfg.ChainID}
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

			solanaRelayers[relayID] = loop.NewRelayerService(lggr, r.GRPCOpts, solCmdFn, string(cfgTOML), signer, r.CapabilitiesRegistry)
		} else {
			// fallback to embedded chain
			opts := solana.ChainOpts{
				Logger:   lggr,
				KeyStore: signer,
			}

			chain, err := solana.NewChain(chainCfg, opts)
			if err != nil {
				return nil, err
			}
			solanaRelayers[relayID] = relay.NewServerAdapter(solana.NewRelayer(lggr, chain, r.CapabilitiesRegistry), chain)
		}
	}
	return solanaRelayers, nil
}

type StarkNetFactoryConfig struct {
	Keystore keystore.StarkNet
	starkcfg.TOMLConfigs
}

// TODO BCF-2606 consider consolidating the driving logic with that of NewSolana above via generics
// perhaps when we implement a Cosmos LOOP
func (r *RelayerFactory) NewStarkNet(ks keystore.StarkNet, chainCfgs starkcfg.TOMLConfigs) (map[types.RelayID]loop.Relayer, error) {
	starknetRelayers := make(map[types.RelayID]loop.Relayer)

	var (
		starkLggr = r.Logger.Named("StarkNet")
		loopKs    = &keystore.StarknetLooppSigner{StarkNet: ks}
	)

	unique := make(map[string]struct{})
	// create one relayer per chain id
	for _, chainCfg := range chainCfgs {
		relayID := types.RelayID{Network: corerelay.NetworkStarkNet, ChainID: *chainCfg.ChainID}
		_, alreadyExists := unique[relayID.Name()]
		if alreadyExists {
			return nil, fmt.Errorf("duplicate chain definitions for %s", relayID.Name())
		}
		unique[relayID.Name()] = struct{}{}

		// skip disabled chains from further processing
		if !chainCfg.IsEnabled() {
			starkLggr.Warnw("Skipping disabled chain", "id", chainCfg.ChainID)
			continue
		}

		lggr := starkLggr.Named(relayID.ChainID)

		if cmdName := env.StarknetPlugin.Cmd.Get(); cmdName != "" {
			// setup the starknet relayer to be a LOOP
			cfgTOML, err := toml.Marshal(struct {
				Starknet starkcfg.TOMLConfig
			}{Starknet: *chainCfg})
			if err != nil {
				return nil, fmt.Errorf("failed to marshal StarkNet configs: %w", err)
			}

			envVars, err := plugins.ParseEnvFile(env.StarknetPlugin.Env.Get())
			if err != nil {
				return nil, fmt.Errorf("failed to parse Starknet env file: %w", err)
			}
			starknetCmdFn, err := plugins.NewCmdFactory(r.Register, plugins.CmdConfig{
				ID:  relayID.Name(),
				Cmd: cmdName,
				Env: envVars,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to create StarkNet LOOP command: %w", err)
			}
			// the starknet relayer service has a delicate keystore dependency. the value that is passed to NewRelayerService must
			// be compatible with instantiating a starknet transaction manager KeystoreAdapter within the LOOPp executable.
			starknetRelayers[relayID] = loop.NewRelayerService(lggr, r.GRPCOpts, starknetCmdFn, string(cfgTOML), loopKs, r.CapabilitiesRegistry)
		} else {
			// fallback to embedded chain
			opts := starkchain.ChainOpts{
				Logger:   lggr,
				KeyStore: loopKs,
			}

			chain, err := starkchain.NewChain(chainCfg, opts)
			if err != nil {
				return nil, err
			}

			starknetRelayers[relayID] = relay.NewServerAdapter(pkgstarknet.NewRelayer(lggr, chain, r.CapabilitiesRegistry), chain)
		}
	}
	return starknetRelayers, nil
}

type CosmosFactoryConfig struct {
	Keystore keystore.Cosmos
	coscfg.TOMLConfigs
	DS sqlutil.DataSource
}

func (c CosmosFactoryConfig) Validate() error {
	var err error
	if c.Keystore == nil {
		err = errors.Join(err, fmt.Errorf("nil Keystore"))
	}
	if len(c.TOMLConfigs) == 0 {
		err = errors.Join(err, fmt.Errorf("no CosmosConfigs provided"))
	}
	if c.DS == nil {
		err = errors.Join(err, fmt.Errorf("nil DataStore"))
	}

	if err != nil {
		err = fmt.Errorf("invalid CosmosFactoryConfig: %w", err)
	}
	return err
}

func (r *RelayerFactory) NewCosmos(config CosmosFactoryConfig) (map[types.RelayID]CosmosLoopRelayerChainer, error) {
	err := config.Validate()
	if err != nil {
		return nil, fmt.Errorf("cannot create Cosmos relayer: %w", err)
	}
	relayers := make(map[types.RelayID]CosmosLoopRelayerChainer)

	var (
		cosmosLggr = r.Logger.Named("Cosmos")
		loopKs     = &keystore.CosmosLoopKeystore{Cosmos: config.Keystore}
	)

	// create one relayer per chain id
	for _, chainCfg := range config.TOMLConfigs {
		relayID := types.RelayID{Network: corerelay.NetworkCosmos, ChainID: *chainCfg.ChainID}

		lggr := cosmosLggr.Named(relayID.ChainID)

		opts := cosmos.ChainOpts{
			Logger:   lggr,
			DS:       config.DS,
			KeyStore: loopKs,
		}

		chain, err := cosmos.NewChain(chainCfg, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to load Cosmos chain %q: %w", relayID, err)
		}

		relayers[relayID] = NewCosmosLoopRelayerChain(cosmos.NewRelayer(lggr, chain), chain)
	}
	return relayers, nil
}

type AptosFactoryConfig struct {
	Keystore    keystore.Aptos
	TOMLConfigs RawConfigs
}

func (r *RelayerFactory) NewAptos(ks keystore.Aptos, chainCfgs RawConfigs) (map[types.RelayID]loop.Relayer, error) {
	plugin := env.NewPlugin("aptos")
	loopKs := &keystore.AptosLooppSigner{Aptos: ks}
	return r.NewLOOPRelayer("Aptos", corerelay.NetworkAptos, plugin, loopKs, chainCfgs)
}

func (r *RelayerFactory) NewLOOPRelayer(name string, network string, plugin env.Plugin, ks coretypes.Keystore, chainCfgs RawConfigs) (map[types.RelayID]loop.Relayer, error) {
	relayers := make(map[types.RelayID]loop.Relayer)
	lggr := r.Logger.Named(name)

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

		lggr2 := lggr.Named(relayID.ChainID)

		cmdName := plugin.Cmd.Get()
		if cmdName == "" {
			return nil, fmt.Errorf("plugin not defined: %s", "")
		}

		// setup the relayer as a LOOP
		cfgTOML, err := toml.Marshal(chainCfg)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal configs: %w", err)
		}

		envVars, err := plugins.ParseEnvFile(plugin.Env.Get())
		if err != nil {
			return nil, fmt.Errorf("failed to parse env file: %w", err)
		}
		cmdFn, err := plugins.NewCmdFactory(r.Register, plugins.CmdConfig{
			ID:  relayID.Name(),
			Cmd: cmdName,
			Env: envVars,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create LOOP command: %w", err)
		}
		// the relayer service has a delicate keystore dependency. the value that is passed to NewRelayerService must
		// be compatible with instantiating a starknet transaction manager KeystoreAdapter within the LOOPp executable.
		relayers[relayID] = loop.NewRelayerService(lggr2, r.GRPCOpts, cmdFn, string(cfgTOML), ks, r.CapabilitiesRegistry)
	}
	return relayers, nil
}
