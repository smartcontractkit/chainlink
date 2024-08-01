package solana

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	relaytypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"

	"github.com/smartcontractkit/chainlink-solana/pkg/solana/client"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/txm"
)

var _ TxManager = (*txm.Txm)(nil)

type TxManager interface {
	Enqueue(accountID string, msg *solana.Transaction) error
}

var _ relaytypes.Relayer = &Relayer{} //nolint:staticcheck

type Relayer struct {
	lggr   logger.Logger
	chain  Chain
	stopCh services.StopChan
}

// Note: constructed in core
func NewRelayer(lggr logger.Logger, chain Chain, capabilitiesRegistry core.CapabilitiesRegistry) *Relayer {
	return &Relayer{
		lggr:   lggr,
		chain:  chain,
		stopCh: make(chan struct{}),
	}
}

func (r *Relayer) Name() string {
	return r.lggr.Name()
}

// Start starts the relayer respecting the given context.
func (r *Relayer) Start(context.Context) error {
	// No subservices started on relay start, but when the first job is started
	if r.chain == nil {
		return errors.New("Solana unavailable")
	}
	return nil
}

// Close will close all open subservices
func (r *Relayer) Close() error {
	close(r.stopCh)
	return nil
}

func (r *Relayer) Ready() error {
	return r.chain.Ready()
}

// Healthy only if all subservices are healthy
func (r *Relayer) Healthy() error { return nil }

func (r *Relayer) HealthReport() map[string]error {
	return map[string]error{r.Name(): r.Healthy()}
}

func (r *Relayer) NewMercuryProvider(rargs relaytypes.RelayArgs, pargs relaytypes.PluginArgs) (relaytypes.MercuryProvider, error) {
	return nil, errors.New("mercury is not supported for solana")
}

func (r *Relayer) NewLLOProvider(rargs relaytypes.RelayArgs, pargs relaytypes.PluginArgs) (relaytypes.LLOProvider, error) {
	return nil, errors.New("data streams is not supported for solana")
}

func (r *Relayer) NewCCIPCommitProvider(rargs relaytypes.RelayArgs, pargs relaytypes.PluginArgs) (relaytypes.CCIPCommitProvider, error) {
	return nil, errors.New("ccip.commit is not supported for solana")
}

func (r *Relayer) NewCCIPExecProvider(rargs relaytypes.RelayArgs, pargs relaytypes.PluginArgs) (relaytypes.CCIPExecProvider, error) {
	return nil, errors.New("ccip.exec is not supported for solana")
}

func (r *Relayer) NewConfigProvider(args relaytypes.RelayArgs) (relaytypes.ConfigProvider, error) {
	ctx, cancel := r.stopCh.NewCtx()
	defer cancel()
	configWatcher, err := newConfigProvider(ctx, r.lggr, r.chain, args)
	if err != nil {
		// Never return (*configProvider)(nil)
		return nil, err
	}
	return configWatcher, err
}

func (r *Relayer) NewChainWriter(_ context.Context, _ []byte) (relaytypes.ChainWriter, error) {
	return nil, errors.New("chain writer is not supported for solana")
}

func (r *Relayer) NewContractReader(_ []byte) (relaytypes.ContractReader, error) {
	return nil, errors.New("contract reader is not supported for solana")
}

func (r *Relayer) NewMedianProvider(rargs relaytypes.RelayArgs, pargs relaytypes.PluginArgs) (relaytypes.MedianProvider, error) {
	ctx, cancel := r.stopCh.NewCtx()
	defer cancel()
	lggr := logger.Named(r.lggr, "MedianProvider")
	configWatcher, err := newConfigProvider(ctx, lggr, r.chain, rargs)
	if err != nil {
		return nil, err
	}

	// parse transmitter account
	transmitterAccount, err := solana.PublicKeyFromBase58(pargs.TransmitterID)
	if err != nil {
		return nil, fmt.Errorf("error on 'solana.PublicKeyFromBase58' for 'spec.PluginArgs.TransmissionsID: %w", err)
	}

	// parse transmissions state account
	var relayConfig RelayConfig
	err = json.Unmarshal(rargs.RelayConfig, &relayConfig)
	if err != nil {
		return nil, err
	}
	transmissionsID, err := solana.PublicKeyFromBase58(relayConfig.TransmissionsID)
	if err != nil {
		return nil, fmt.Errorf("error on 'solana.PublicKeyFromBase58' for 'spec.RelayConfig.TransmissionsID: %w", err)
	}

	cfg := configWatcher.chain.Config()
	transmissionsCache := NewTransmissionsCache(transmissionsID, relayConfig.ChainID, cfg, configWatcher.reader, r.lggr)
	return &medianProvider{
		configProvider:     configWatcher,
		transmissionsCache: transmissionsCache,
		reportCodec:        ReportCodec{},
		contract: &MedianContract{
			stateCache:         configWatcher.stateCache,
			transmissionsCache: transmissionsCache,
		},
		transmitter: &Transmitter{
			stateID:            configWatcher.stateID,
			programID:          configWatcher.programID,
			storeProgramID:     configWatcher.storeProgramID,
			transmissionsID:    transmissionsID,
			transmissionSigner: transmitterAccount,
			reader:             configWatcher.reader,
			stateCache:         configWatcher.stateCache,
			lggr:               r.lggr,
			txManager:          configWatcher.chain.TxManager(),
		},
	}, nil
}

func (r *Relayer) NewFunctionsProvider(rargs relaytypes.RelayArgs, pargs relaytypes.PluginArgs) (relaytypes.FunctionsProvider, error) {
	return nil, errors.New("functions are not supported for solana")
}

func (r *Relayer) NewAutomationProvider(rargs relaytypes.RelayArgs, pargs relaytypes.PluginArgs) (relaytypes.AutomationProvider, error) {
	return nil, errors.New("automation is not supported for solana")
}

var _ relaytypes.ConfigProvider = &configProvider{}

type configProvider struct {
	services.StateMachine
	chainID                            string
	programID, storeProgramID, stateID solana.PublicKey
	stateCache                         *StateCache
	offchainConfigDigester             types.OffchainConfigDigester
	configTracker                      types.ContractConfigTracker
	chain                              Chain
	reader                             client.Reader
}

func newConfigProvider(ctx context.Context, lggr logger.Logger, chain Chain, args relaytypes.RelayArgs) (*configProvider, error) {
	lggr = logger.Named(lggr, "ConfigProvider")
	var relayConfig RelayConfig
	err := json.Unmarshal(args.RelayConfig, &relayConfig)
	if err != nil {
		return nil, err
	}
	stateID, err := solana.PublicKeyFromBase58(args.ContractID)
	if err != nil {
		return nil, fmt.Errorf("error on 'solana.PublicKeyFromBase58' for 'spec.ContractID: %w", err)
	}
	programID, err := solana.PublicKeyFromBase58(relayConfig.OCR2ProgramID)
	if err != nil {
		return nil, fmt.Errorf("error on 'solana.PublicKeyFromBase58' for 'spec.RelayConfig.OCR2ProgramID: %w", err)
	}
	storeProgramID, err := solana.PublicKeyFromBase58(relayConfig.StoreProgramID)
	if err != nil {
		return nil, fmt.Errorf("error on 'solana.PublicKeyFromBase58' for 'spec.RelayConfig.StateID: %w", err)
	}
	offchainConfigDigester := OffchainConfigDigester{
		ProgramID: programID,
		StateID:   stateID,
	}

	reader, err := chain.Reader()
	if err != nil {
		return nil, fmt.Errorf("error in NewMedianProvider.chain.Reader: %w", err)
	}
	stateCache := NewStateCache(stateID, relayConfig.ChainID, chain.Config(), reader, lggr)
	return &configProvider{
		chainID:                relayConfig.ChainID,
		stateID:                stateID,
		programID:              programID,
		storeProgramID:         storeProgramID,
		stateCache:             stateCache,
		offchainConfigDigester: offchainConfigDigester,
		configTracker:          &ConfigTracker{stateCache: stateCache, reader: reader},
		chain:                  chain,
		reader:                 reader,
	}, nil
}

func (c *configProvider) Name() string {
	return c.stateCache.Name()
}

func (c *configProvider) Start(ctx context.Context) error {
	return c.StartOnce("SolanaConfigProvider", func() error {
		return c.stateCache.Start(ctx)
	})
}

func (c *configProvider) Close() error {
	return c.StopOnce("SolanaConfigProvider", func() error {
		return c.stateCache.Close()
	})
}

func (c *configProvider) HealthReport() map[string]error {
	return map[string]error{c.Name(): c.Healthy()}
}

func (c *configProvider) OffchainConfigDigester() types.OffchainConfigDigester {
	return c.offchainConfigDigester
}

func (c *configProvider) ContractConfigTracker() types.ContractConfigTracker {
	return c.configTracker
}

var _ relaytypes.MedianProvider = &medianProvider{}

type medianProvider struct {
	*configProvider
	transmissionsCache *TransmissionsCache
	reportCodec        median.ReportCodec
	contract           median.MedianContract
	transmitter        types.ContractTransmitter
}

func (p *medianProvider) Name() string {
	return p.stateCache.Name()
}

// start both cache services
func (p *medianProvider) Start(ctx context.Context) error {
	return p.StartOnce("SolanaMedianProvider", func() error {
		if err := p.configProvider.stateCache.Start(ctx); err != nil {
			return err
		}
		return p.transmissionsCache.Start(ctx)
	})
}

// close both cache services
func (p *medianProvider) Close() error {
	return p.StopOnce("SolanaMedianProvider", func() error {
		if err := p.configProvider.stateCache.Close(); err != nil {
			return err
		}
		return p.transmissionsCache.Close()
	})
}

func (p *medianProvider) ContractTransmitter() types.ContractTransmitter {
	return p.transmitter
}

func (p *medianProvider) ReportCodec() median.ReportCodec {
	return p.reportCodec
}

func (p *medianProvider) MedianContract() median.MedianContract {
	return p.contract
}

func (p *medianProvider) OnchainConfigCodec() median.OnchainConfigCodec {
	return median.StandardOnchainConfigCodec{}
}

func (p *medianProvider) ChainReader() relaytypes.ContractReader {
	return nil
}

func (p *medianProvider) Codec() relaytypes.Codec {
	return nil
}

func (r *Relayer) NewPluginProvider(rargs relaytypes.RelayArgs, pargs relaytypes.PluginArgs) (relaytypes.PluginProvider, error) {
	return nil, errors.New("plugin provider is not supported for solana")
}

func (r *Relayer) NewOCR3CapabilityProvider(rargs relaytypes.RelayArgs, pargs relaytypes.PluginArgs) (relaytypes.OCR3CapabilityProvider, error) {
	return nil, errors.New("ocr3 capability provider is not supported for solana")
}
