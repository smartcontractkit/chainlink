package evm

import (
	"math/big"
	"strings"

	"github.com/smartcontractkit/chainlink/core/services/job"

	types2 "github.com/smartcontractkit/chainlink/core/services/relay/types"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	txm "github.com/smartcontractkit/chainlink/core/chains/evm/bulletprooftxmanager"
	offchain_aggregator_wrapper "github.com/smartcontractkit/chainlink/core/internal/gethwrappers2/generated/offchainaggregator"
	"github.com/smartcontractkit/chainlink/core/services/ocrcommon"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	"github.com/smartcontractkit/libocr/offchainreporting2/chains/evmutil"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median/evmreportcodec"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/sqlx"
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
func (r *Relayer) Start() error {
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

// NewOCR2Provider creates a per-job set of OCR2 dependencies.
func (r *Relayer) NewOCR2Provider(externalJobID uuid.UUID, s interface{}, contractReady chan<- struct{}) (types2.OCR2Provider, error) {
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

	contract, err := offchain_aggregator_wrapper.NewOffchainAggregator(contractAddress, chain.Client())
	if err != nil {
		return nil, errors.Wrap(err, "could not instantiate NewOffchainAggregator")
	}

	contractFilterer, err := ocr2aggregator.NewOCR2AggregatorFilterer(contractAddress, chain.Client())
	if err != nil {
		return nil, errors.Wrap(err, "could not instantiate NewOffchainAggregatorFilterer")
	}

	contractCaller, err := ocr2aggregator.NewOCR2AggregatorCaller(contractAddress, chain.Client())
	if err != nil {
		return nil, errors.Wrap(err, "could not instantiate NewOffchainAggregatorCaller")
	}

	ocrDB := NewContractDB(r.db.DB, spec.ID, r.lggr)

	tracker := NewOCRContractTracker(
		contract,
		contractFilterer,
		contractCaller,
		chain.Client(),
		chain.LogBroadcaster(),
		spec.ID,
		r.lggr,
		r.db,
		ocrDB,
		chain.Config(),
		chain.HeadBroadcaster(),
		contractReady,
	)

	offchainConfigDigester := evmutil.EVMOffchainConfigDigester{
		ChainID:         chain.Config().ChainID().Uint64(),
		ContractAddress: contractAddress,
	}

	if spec.IsBootstrap {
		// Return early if bootstrap node (doesn't require the full OCR2 provider)
		return &ocr2Provider{
			tracker:                tracker,
			offchainConfigDigester: offchainConfigDigester,
		}, nil
	}

	reportCodec := evmreportcodec.ReportCodec{}

	contractABI, err := abi.JSON(strings.NewReader(ocr2aggregator.OCR2AggregatorABI))
	if err != nil {
		return nil, errors.Wrap(err, "could not get contract ABI JSON")
	}

	if !spec.TransmitterID.Valid {
		return nil, errors.New("transmitterID is required for non-bootstrap jobs")
	}
	transmitterAddress := common.HexToAddress(spec.TransmitterID.String)
	strategy := txm.NewQueueingTxStrategy(externalJobID, chain.Config().OCRDefaultTransactionQueueDepth())

	contractTransmitter := NewOCRContractTransmitter(
		contract.Address(),
		contractCaller,
		contractABI,
		ocrcommon.NewTransmitter(chain.TxManager(), transmitterAddress, chain.Config().EvmGasLimitDefault(), strategy, txm.TransmitCheckerSpec{}),
		tracker,
		r.lggr,
	)

	return &ocr2Provider{
		tracker:                tracker,
		offchainConfigDigester: offchainConfigDigester,
		reportCodec:            reportCodec,
		contractTransmitter:    contractTransmitter,
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
}

var _ job.Service = (*ocr2Provider)(nil)

type ocr2Provider struct {
	contractReady          chan struct{}
	tracker                *ContractTracker
	offchainConfigDigester types.OffchainConfigDigester
	reportCodec            median.ReportCodec
	contractTransmitter    *ContractTransmitter
}

// On start, an ethereum ocr2 provider will start the contract tracker.
func (p ocr2Provider) Start() error {
	return p.tracker.Start()
}

// On close, an ethereum ocr2 provider will close the contract tracker.
func (p ocr2Provider) Close() error {
	return p.tracker.Close()
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
	return p.contractTransmitter
}
