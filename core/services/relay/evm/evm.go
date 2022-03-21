package evm

import (
	"context"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	"github.com/smartcontractkit/libocr/offchainreporting2/chains/evmutil"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median/evmreportcodec"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/sqlx"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	txm "github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/ocrcommon"
	types2 "github.com/smartcontractkit/chainlink/core/services/relay/types"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type Relayer struct {
	db       *sqlx.DB
	chainSet evm.ChainSet
	lggr     logger.Logger
}

func NewRelayer(db *sqlx.DB, chainSet evm.ChainSet, lggr logger.Logger) *Relayer {
	return &Relayer{
		db:       db,
		chainSet: chainSet,
		lggr:     lggr.Named("Relayer"),
	}
}

// Start does noop: no subservices started on relay start, but when the first job is started
func (r *Relayer) Start(context.Context) error {
	return nil
}

// Close does noop: no persistent subservices to close on relay close
func (r *Relayer) Close() error {
	return nil
}

// Ready does noop: always ready
func (r *Relayer) Ready() error {
	return nil
}

// Healthy does noop: always healthy
func (r *Relayer) Healthy() error {
	return nil
}

// NewOCR2Provider provides all evm specific implementations of OCR2 components
// including components generic across all plugins and ones specific to plugins.
func (r *Relayer) NewOCR2Provider(externalJobID uuid.UUID, s interface{}) (types2.OCR2ProviderCtx, error) {
	// Expect trusted input
	spec := s.(OCR2Spec)
	chain, err := r.chainSet.Get(spec.ChainID)
	if err != nil {
		return nil, err
	}
	if !common.IsHexAddress(spec.ContractID) {
		return nil, errors.Errorf("invalid contractID, expected hex address")
	}
	contractAddress := common.HexToAddress(spec.ContractID)
	contractABI, err := abi.JSON(strings.NewReader(ocr2aggregator.OCR2AggregatorABI))
	if err != nil {
		return nil, errors.Wrap(err, "could not get contract ABI JSON")
	}
	configTracker := NewConfigTracker(r.lggr, contractABI,
		chain.Client(),
		contractAddress,
		chain.Config().ChainType(),
		chain.HeadBroadcaster())

	offchainConfigDigester := evmutil.EVMOffchainConfigDigester{
		ChainID:         chain.Config().ChainID().Uint64(),
		ContractAddress: contractAddress,
	}

	if spec.IsBootstrap {
		// Return early if bootstrap node (doesn't require the full OCR2 provider)
		return &ocr2Provider{
			tracker:                configTracker,
			offchainConfigDigester: offchainConfigDigester,
			plugin:                 spec.Plugin,
		}, nil
	}

	if !spec.TransmitterID.Valid {
		return nil, errors.New("transmitterID is required for non-bootstrap jobs")
	}
	transmitterAddress := common.HexToAddress(spec.TransmitterID.String)
	strategy := txm.NewQueueingTxStrategy(externalJobID, chain.Config().OCRDefaultTransactionQueueDepth())

	contractTransmitter := NewOCRContractTransmitter(
		contractAddress,
		chain.Client(),
		contractABI,
		ocrcommon.NewTransmitter(chain.TxManager(), transmitterAddress, chain.Config().EvmGasLimitDefault(), strategy, txm.TransmitCheckerSpec{}),
		r.lggr,
	)

	medianContract, err := newMedianContract(contractAddress, chain, spec.ID, r.db, r.lggr)
	if err != nil {
		return nil, errors.Wrap(err, "error during median contract setup")
	}

	reportCodec := evmreportcodec.ReportCodec{}

	return &ocr2Provider{
		tracker:                configTracker,
		offchainConfigDigester: offchainConfigDigester,
		reportCodec:            reportCodec,
		contractTransmitter:    contractTransmitter,
		medianContract:         medianContract,
	}, nil
}

type RelayConfig struct {
	ChainID *utils.Big `json:"chainID"`
}

type OCR2Spec struct {
	ID            int32
	ContractID    string
	TransmitterID null.String // Will be null for bootstrap jobs
	IsBootstrap   bool
	ChainID       *big.Int
	Plugin        job.OCR2PluginType
}

var _ services.ServiceCtx = (*ocr2Provider)(nil)

type ocr2Provider struct {
	tracker                *ConfigTracker
	offchainConfigDigester types.OffchainConfigDigester
	contractTransmitter    *ContractTransmitter
	plugin                 job.OCR2PluginType
	// Median specific
	reportCodec    median.ReportCodec
	medianContract *medianContract
}

// Start an ethereum ocr2 provider will start the contract tracker.
func (p ocr2Provider) Start(context.Context) error {
	err := p.tracker.Start()
	if err != nil {
		return err
	}
	// TODO (https://app.shortcut.com/chainlinklabs/story/32017/plugin-specific-relay-interfaces):
	// We need to break up ocr2Provider into more granular components
	// per plugin (would require changes in solana/terra relay repos)
	if p.plugin == job.Median && p.medianContract != nil {
		return p.medianContract.Start()
	}
	return nil
}

// Close an ethereum ocr2 provider will close the contract tracker.
func (p ocr2Provider) Close() error {
	err := p.tracker.Close()
	if err != nil {
		return err
	}
	if p.plugin == job.Median && p.medianContract != nil {
		return p.medianContract.Close()
	}
	return nil
}

// Ready always returns ready.
func (p ocr2Provider) Ready() error {
	return nil
}

// Healthy always returns healthy.
func (p ocr2Provider) Healthy() error {
	return nil
}

func (p ocr2Provider) ContractTransmitter() types.ContractTransmitter {
	return p.contractTransmitter
}

func (p ocr2Provider) ContractConfigTracker() types.ContractConfigTracker {
	return p.tracker
}

func (p ocr2Provider) OffchainConfigDigester() types.OffchainConfigDigester {
	return p.offchainConfigDigester
}

func (p ocr2Provider) ReportCodec() median.ReportCodec {
	return p.reportCodec
}

func (p ocr2Provider) MedianContract() median.MedianContract {
	return p.medianContract
}
