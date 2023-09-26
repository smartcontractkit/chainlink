package chainlink

import (
	"context"
	"errors"
	"fmt"

	"github.com/pelletier/go-toml/v2"

	"github.com/smartcontractkit/sqlx"

	pkgcosmos "github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos"
	"github.com/smartcontractkit/chainlink-relay/pkg/loop"
	pkgsolana "github.com/smartcontractkit/chainlink-solana/pkg/solana"
	pkgstarknet "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink"

	"github.com/smartcontractkit/chainlink/v2/core/chains/cosmos"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/chains/solana"
	"github.com/smartcontractkit/chainlink/v2/core/chains/starknet"
	"github.com/smartcontractkit/chainlink/v2/core/config/env"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	evmrelay "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

type RelayerFactory struct {
	logger.Logger
	*plugins.LoopRegistry
	loop.GRPCOpts
}

type EVMFactoryConfig struct {
	evm.ChainOpts
	evmrelay.CSAETHKeystore
}

func (r *RelayerFactory) NewEVM(ctx context.Context, config EVMFactoryConfig) (map[relay.ID]evmrelay.LoopRelayAdapter, error) {
	// TODO impl EVM loop. For now always 'fallback' to an adapter and embedded chain

	relayers := make(map[relay.ID]evmrelay.LoopRelayAdapter)

	// override some common opts with the factory values. this seems weird... maybe other signatures should change, or this should take a different type...
	ccOpts := evm.ChainRelayExtenderConfig{
		Logger:    r.Logger.Named("EVM"),
		KeyStore:  config.CSAETHKeystore.Eth(),
		ChainOpts: config.ChainOpts,
	}

	evmRelayExtenders, err := evmrelay.NewChainRelayerExtenders(ctx, ccOpts)
	if err != nil {
		return nil, err
	}
	legacyChains := evmrelay.NewLegacyChainsFromRelayerExtenders(evmRelayExtenders)
	for _, ext := range evmRelayExtenders.Slice() {
		relayID := relay.ID{Network: relay.EVM, ChainID: relay.ChainID(ext.Chain().ID().String())}
		chain, err2 := legacyChains.Get(relayID.ChainID)
		if err2 != nil {
			return nil, err2
		}

		relayerOpts := evmrelay.RelayerOpts{
			DB:               ccOpts.DB,
			QConfig:          ccOpts.AppConfig.Database(),
			CSAETHKeystore:   config.CSAETHKeystore,
			EventBroadcaster: ccOpts.EventBroadcaster,
		}
		relayer, err2 := evmrelay.NewRelayer(ccOpts.Logger, chain, relayerOpts)
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
	solana.SolanaConfigs
}

func (r *RelayerFactory) NewSolana(ks keystore.Solana, chainCfgs solana.SolanaConfigs) (map[relay.ID]loop.Relayer, error) {
	solanaRelayers := make(map[relay.ID]loop.Relayer)
	var (
		solLggr = r.Logger.Named("Solana")
		signer  = &keystore.SolanaSigner{Solana: ks}
	)

	unique := make(map[string]struct{})
	// create one relayer per chain id
	for _, chainCfg := range chainCfgs {

		relayId := relay.ID{Network: relay.Solana, ChainID: relay.ChainID(*chainCfg.ChainID)}
		_, alreadyExists := unique[relayId.Name()]
		if alreadyExists {
			return nil, fmt.Errorf("duplicate chain definitions for %s", relayId.Name())
		}
		unique[relayId.Name()] = struct{}{}

		// skip disabled chains from further processing
		if !chainCfg.IsEnabled() {
			solLggr.Warnw("Skipping disabled chain", "id", chainCfg.ChainID)
			continue
		}

		if cmdName := env.SolanaPluginCmd.Get(); cmdName != "" {

			// setup the solana relayer to be a LOOP
			cfgTOML, err := toml.Marshal(struct {
				Solana solana.SolanaConfig
			}{Solana: *chainCfg})

			if err != nil {
				return nil, fmt.Errorf("failed to marshal Solana configs: %w", err)
			}

			solCmdFn, err := plugins.NewCmdFactory(r.Register, plugins.CmdConfig{
				ID:  relayId.Name(),
				Cmd: cmdName,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to create Solana LOOP command: %w", err)
			}

			solanaRelayers[relayId] = loop.NewRelayerService(solLggr, r.GRPCOpts, solCmdFn, string(cfgTOML), signer)

		} else {
			// fallback to embedded chain
			opts := solana.ChainOpts{
				Logger:   solLggr,
				KeyStore: signer,
			}

			chain, err := solana.NewChain(chainCfg, opts)
			if err != nil {
				return nil, err
			}
			solanaRelayers[relayId] = relay.NewRelayerServerAdapter(pkgsolana.NewRelayer(solLggr, chain), chain)
		}
	}
	return solanaRelayers, nil
}

type StarkNetFactoryConfig struct {
	Keystore keystore.StarkNet
	starknet.StarknetConfigs
}

// TODO BCF-2606 consider consolidating the driving logic with that of NewSolana above via generics
// perhaps when we implement a Cosmos LOOP
func (r *RelayerFactory) NewStarkNet(ks keystore.StarkNet, chainCfgs starknet.StarknetConfigs) (map[relay.ID]loop.Relayer, error) {
	starknetRelayers := make(map[relay.ID]loop.Relayer)

	var (
		starkLggr = r.Logger.Named("StarkNet")
		loopKs    = &keystore.StarknetLooppSigner{StarkNet: ks}
	)

	unique := make(map[string]struct{})
	// create one relayer per chain id
	for _, chainCfg := range chainCfgs {
		relayId := relay.ID{Network: relay.StarkNet, ChainID: relay.ChainID(*chainCfg.ChainID)}
		_, alreadyExists := unique[relayId.Name()]
		if alreadyExists {
			return nil, fmt.Errorf("duplicate chain definitions for %s", relayId.Name())
		}
		unique[relayId.Name()] = struct{}{}

		// skip disabled chains from further processing
		if !chainCfg.IsEnabled() {
			starkLggr.Warnw("Skipping disabled chain", "id", chainCfg.ChainID)
			continue
		}

		if cmdName := env.StarknetPluginCmd.Get(); cmdName != "" {
			// setup the starknet relayer to be a LOOP
			cfgTOML, err := toml.Marshal(struct {
				Starknet starknet.StarknetConfig
			}{Starknet: *chainCfg})
			if err != nil {
				return nil, fmt.Errorf("failed to marshal StarkNet configs: %w", err)
			}

			starknetCmdFn, err := plugins.NewCmdFactory(r.Register, plugins.CmdConfig{
				ID:  relayId.Name(),
				Cmd: cmdName,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to create StarkNet LOOP command: %w", err)
			}
			// the starknet relayer service has a delicate keystore dependency. the value that is passed to NewRelayerService must
			// be compatible with instantiating a starknet transaction manager KeystoreAdapter within the LOOPp executable.
			starknetRelayers[relayId] = loop.NewRelayerService(starkLggr, r.GRPCOpts, starknetCmdFn, string(cfgTOML), loopKs)
		} else {
			// fallback to embedded chain
			opts := starknet.ChainOpts{
				Logger:   starkLggr,
				KeyStore: loopKs,
			}

			chain, err := starknet.NewChain(chainCfg, opts)
			if err != nil {
				return nil, err
			}

			starknetRelayers[relayId] = relay.NewRelayerServerAdapter(pkgstarknet.NewRelayer(starkLggr, chain), chain)
		}
	}
	return starknetRelayers, nil

}

type CosmosFactoryConfig struct {
	Keystore keystore.Cosmos
	cosmos.CosmosConfigs
	EventBroadcaster pg.EventBroadcaster
	*sqlx.DB
	pg.QConfig
}

func (c CosmosFactoryConfig) Validate() error {
	var err error
	if c.Keystore == nil {
		err = errors.Join(err, fmt.Errorf("nil Keystore"))
	}
	if len(c.CosmosConfigs) == 0 {
		err = errors.Join(err, fmt.Errorf("no CosmosConfigs provided"))
	}
	if c.EventBroadcaster == nil {
		err = errors.Join(err, fmt.Errorf("nil EventBroadcaster"))
	}
	if c.DB == nil {
		err = errors.Join(err, fmt.Errorf("nil DB"))
	}
	if c.QConfig == nil {
		err = errors.Join(err, fmt.Errorf("nil QConfig"))
	}

	if err != nil {
		err = fmt.Errorf("invalid CosmosFactoryConfig: %w", err)
	}
	return err
}

func (r *RelayerFactory) NewCosmos(ctx context.Context, config CosmosFactoryConfig) (map[relay.ID]cosmos.LoopRelayerChainer, error) {
	err := config.Validate()
	if err != nil {
		return nil, fmt.Errorf("cannot create Cosmos relayer: %w", err)
	}
	relayers := make(map[relay.ID]cosmos.LoopRelayerChainer)

	var (
		lggr   = r.Logger.Named("Cosmos")
		loopKs = &keystore.CosmosLoopKeystore{Cosmos: config.Keystore}
	)

	// create one relayer per chain id
	for _, chainCfg := range config.CosmosConfigs {
		relayId := relay.ID{Network: relay.Cosmos, ChainID: relay.ChainID(*chainCfg.ChainID)}

		opts := cosmos.ChainOpts{
			QueryConfig:      config.QConfig,
			Logger:           lggr.Named(relayId.ChainID),
			DB:               config.DB,
			KeyStore:         loopKs,
			EventBroadcaster: config.EventBroadcaster,
		}

		chain, err := cosmos.NewChain(chainCfg, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to load Cosmos chain %q: %w", relayId, err)
		}

		relayers[relayId] = cosmos.NewLoopRelayerChain(pkgcosmos.NewRelayer(lggr, chain), chain)

	}
	return relayers, nil

}
