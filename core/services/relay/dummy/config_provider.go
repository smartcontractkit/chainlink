package dummy

import (
	"context"

	"github.com/ethereum/go-ethereum/common/hexutil"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

// Stub ConfigTracker that uses static config

// interim struct used for unmarshalling from relay config
type ConfigTrackerCfg struct {
	// OCR Config
	ConfigDigest          hexutil.Bytes
	ConfigCount           uint64
	Signers               []hexutil.Bytes
	Transmitters          []string
	F                     uint8
	OnchainConfig         hexutil.Bytes
	OffchainConfigVersion uint64
	OffchainConfig        hexutil.Bytes

	// Tracker config
	ChangedInBlock uint64
	BlockHeight    uint64
}

func (cfg ConfigTrackerCfg) ToContractConfig() (ocrtypes.ContractConfig, error) {
	cd, err := ocrtypes.BytesToConfigDigest(cfg.ConfigDigest)
	if err != nil {
		return ocrtypes.ContractConfig{}, err
	}
	signers := make([]ocrtypes.OnchainPublicKey, len(cfg.Signers))
	for i, s := range cfg.Signers {
		signers[i] = ocrtypes.OnchainPublicKey(s)
	}
	transmitters := make([]ocrtypes.Account, len(cfg.Transmitters))
	for i, t := range cfg.Transmitters {
		transmitters[i] = ocrtypes.Account(t)
	}
	return ocrtypes.ContractConfig{
		ConfigDigest:          cd,
		ConfigCount:           cfg.ConfigCount,
		Signers:               signers,
		Transmitters:          transmitters,
		F:                     cfg.F,
		OnchainConfig:         cfg.OnchainConfig,
		OffchainConfigVersion: cfg.OffchainConfigVersion,
		OffchainConfig:        cfg.OffchainConfig,
	}, nil
}

type configProvider struct {
	lggr logger.Logger

	digester ocrtypes.OffchainConfigDigester
	tracker  ocrtypes.ContractConfigTracker
}

func NewConfigProvider(lggr logger.Logger, cfg RelayConfig) (types.ConfigProvider, error) {
	cp := &configProvider{lggr: lggr.Named("DummyConfigProvider").Named(cfg.ConfigTracker.ConfigDigest.String())}

	{
		contractConfig, err := cfg.ConfigTracker.ToContractConfig()
		if err != nil {
			return nil, err
		}

		cp.digester, err = NewOffchainConfigDigester(contractConfig.ConfigDigest)
		if err != nil {
			return nil, err
		}
	}
	var err error
	cp.tracker, err = NewContractConfigTracker(cp.lggr, cfg.ConfigTracker)
	if err != nil {
		return nil, err
	}
	return cp, nil
}

func (cp *configProvider) OffchainConfigDigester() ocrtypes.OffchainConfigDigester {
	return cp.digester
}
func (cp *configProvider) ContractConfigTracker() ocrtypes.ContractConfigTracker { return cp.tracker }
func (cp *configProvider) Name() string                                          { return cp.lggr.Name() }
func (*configProvider) Start(context.Context) error                              { return nil }
func (*configProvider) Close() error                                             { return nil }
func (*configProvider) Ready() error                                             { return nil }
func (cp *configProvider) HealthReport() map[string]error {
	return map[string]error{cp.lggr.Name(): nil}
}
