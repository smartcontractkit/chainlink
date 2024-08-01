package ocr2

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	relaytypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"

	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/ocr2/medianreport"
	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/txm"
	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/starknet"
)

var _ relaytypes.ConfigProvider = (*configProvider)(nil)

type configProvider struct {
	utils.StartStopOnce

	reader        Reader
	contractCache *contractCache
	digester      types.OffchainConfigDigester

	lggr logger.Logger
}

func NewConfigProvider(chainID string, contractAddress string, basereader starknet.Reader, cfg Config, lggr logger.Logger) (*configProvider, error) {
	lggr = logger.Named(lggr, "ConfigProvider")
	chainReader, err := NewClient(basereader, lggr)
	if err != nil {
		return nil, fmt.Errorf("err in NewConfigProvider.NewClient: %w", err)
	}

	reader := NewContractReader(contractAddress, chainReader, lggr)
	cache := NewContractCache(cfg, reader, lggr)
	digester := NewOffchainConfigDigester(chainID, contractAddress)

	return &configProvider{
		reader:        reader,
		contractCache: cache,
		digester:      digester,
		lggr:          lggr,
	}, nil
}

func (p *configProvider) Name() string {
	return p.lggr.Name()
}

func (p *configProvider) Start(context.Context) error {
	return p.StartOnce("ConfigProvider", func() error {
		p.lggr.Debugf("Config provider starting")
		return p.contractCache.Start()
	})
}

func (p *configProvider) Close() error {
	return p.StopOnce("ConfigProvider", func() error {
		p.lggr.Debugf("Config provider stopping")
		return p.contractCache.Close()
	})
}

func (p *configProvider) HealthReport() map[string]error {
	return map[string]error{p.Name(): p.Healthy()}
}

func (p *configProvider) ContractConfigTracker() types.ContractConfigTracker {
	return p.contractCache
}

func (p *configProvider) OffchainConfigDigester() types.OffchainConfigDigester {
	return p.digester
}

var _ relaytypes.MedianProvider = (*medianProvider)(nil)

type medianProvider struct {
	*configProvider
	transmitter        types.ContractTransmitter
	transmissionsCache *transmissionsCache
	reportCodec        median.ReportCodec
}

func NewMedianProvider(chainID string, contractAddress string, senderAddress string, accountAddress string, basereader starknet.Reader, cfg Config, txm txm.TxManager, lggr logger.Logger) (*medianProvider, error) {
	lggr = logger.Named(lggr, "MedianProvider")
	configProvider, err := NewConfigProvider(chainID, contractAddress, basereader, cfg, lggr)
	if err != nil {
		return nil, fmt.Errorf("error in NewMedianProvider.NewConfigProvider: %w", err)
	}

	cache := NewTransmissionsCache(cfg, configProvider.reader, lggr)
	transmitter := NewContractTransmitter(cache, contractAddress, senderAddress, accountAddress, txm)

	return &medianProvider{
		configProvider:     configProvider,
		transmitter:        transmitter,
		transmissionsCache: cache,
		reportCodec:        medianreport.ReportCodec{},
	}, nil
}

func (p *medianProvider) Name() string {
	return p.lggr.Name()
}

func (p *medianProvider) Start(context.Context) error {
	return p.StartOnce("MedianProvider", func() error {
		p.lggr.Debugf("Median provider starting")
		// starting both cache services here
		// todo: find a better way
		if err := p.configProvider.contractCache.Start(); err != nil {
			return fmt.Errorf("couldn't start contractCache: %w", err)
		}
		return p.transmissionsCache.Start()
	})
}

func (p *medianProvider) Close() error {
	return p.StopOnce("MedianProvider", func() error {
		p.lggr.Debugf("Median provider stopping")
		// stopping both cache services here
		// todo: find a better way
		if err := p.configProvider.contractCache.Close(); err != nil {
			return fmt.Errorf("coulnd't stop contractCache: %w", err)
		}
		return p.transmissionsCache.Close()
	})
}

func (p *medianProvider) HealthReport() map[string]error {
	return map[string]error{p.Name(): p.Healthy()}
}

func (p *medianProvider) ContractTransmitter() types.ContractTransmitter {
	return p.transmitter
}

func (p *medianProvider) ReportCodec() median.ReportCodec {
	return p.reportCodec
}

func (p *medianProvider) MedianContract() median.MedianContract {
	return p.transmissionsCache
}

func (p *medianProvider) OnchainConfigCodec() median.OnchainConfigCodec {
	return medianreport.OnchainConfigCodec{}
}

func (p *medianProvider) ChainReader() relaytypes.ContractReader {
	return nil
}

func (p *medianProvider) Codec() relaytypes.Codec {
	return nil
}
