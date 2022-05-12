package evm

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	"github.com/smartcontractkit/libocr/offchainreporting2/chains/evmutil"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median/evmreportcodec"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	txm "github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/core/logger"
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

func (r *Relayer) NewMedianProvider(args types2.OCR2Args) (types2.MedianProvider, error) {
	relayConfigBytes, err := json.Marshal(args.RelayConfig)
	if err != nil {
		return nil, err
	}
	var relayConfig RelayConfig
	err = json.Unmarshal(relayConfigBytes, &relayConfig)
	if err != nil {
		return nil, err
	}
	chain, err := r.chainSet.Get(relayConfig.ChainID.ToInt())
	if err != nil {
		return nil, err
	}
	if !common.IsHexAddress(args.ContractID) {
		return nil, errors.Errorf("invalid contractID, expected hex address")
	}
	contractAddress := common.HexToAddress(args.ContractID)
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

	if args.IsBootstrap {
		// Return early if bootstrap node (doesn't require the full OCR2 provider)
		return &medianProvider{
			tracker:                configTracker,
			offchainConfigDigester: offchainConfigDigester,
		}, nil
	}

	if !args.TransmitterID.Valid {
		return nil, errors.New("transmitterID is required for non-bootstrap jobs")
	}
	transmitterAddress := common.HexToAddress(args.TransmitterID.String)
	strategy := txm.NewQueueingTxStrategy(args.ExternalJobID, chain.Config().OCRDefaultTransactionQueueDepth())

	contractTransmitter := NewOCRContractTransmitter(
		contractAddress,
		chain.Client(),
		contractABI,
		ocrcommon.NewTransmitter(chain.TxManager(), transmitterAddress, chain.Config().EvmGasLimitDefault(), strategy, txm.TransmitCheckerSpec{}),
		r.lggr,
	)

	medianContract, err := newMedianContract(contractAddress, chain, args.JobID, r.db, r.lggr)
	if err != nil {
		return nil, errors.Wrap(err, "error during median contract setup")
	}

	reportCodec := evmreportcodec.ReportCodec{}

	return &medianProvider{
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

var _ types2.MedianProvider = (*medianProvider)(nil)

type medianProvider struct {
	tracker                *ConfigTracker
	offchainConfigDigester types.OffchainConfigDigester
	contractTransmitter    *ContractTransmitter
	reportCodec            median.ReportCodec
	medianContract         *medianContract
}

// Start an ethereum ocr2 provider will start the contract tracker.
func (p medianProvider) Start(context.Context) error {
	err := p.tracker.Start()
	if err != nil {
		return err
	}
	// Bootstrap does not need a median contract.
	if p.medianContract != nil {
		return p.medianContract.Start()
	}
	return nil
}

// Close an ethereum ocr2 provider will close the contract tracker.
func (p medianProvider) Close() error {
	err := p.tracker.Close()
	if err != nil {
		return err
	}
	if p.medianContract != nil {
		return p.medianContract.Close()
	}
	return nil
}

// Ready always returns ready.
func (p medianProvider) Ready() error {
	return nil
}

// Healthy always returns healthy.
func (p medianProvider) Healthy() error {
	return nil
}

func (p medianProvider) ContractTransmitter() types.ContractTransmitter {
	return p.contractTransmitter
}

func (p medianProvider) ContractConfigTracker() types.ContractConfigTracker {
	return p.tracker
}

func (p medianProvider) OffchainConfigDigester() types.OffchainConfigDigester {
	return p.offchainConfigDigester
}

func (p medianProvider) ReportCodec() median.ReportCodec {
	return p.reportCodec
}

func (p medianProvider) MedianContract() median.MedianContract {
	return p.medianContract
}
