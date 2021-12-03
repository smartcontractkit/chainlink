package ethereum

import (
	"encoding/json"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	offchain_aggregator_wrapper "github.com/smartcontractkit/chainlink/core/internal/gethwrappers2/generated/offchainaggregator"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/service"
	txm "github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/core/services/ocrcommon"
	ocr2 "github.com/smartcontractkit/chainlink/core/services/offchainreporting2"
	"github.com/smartcontractkit/chainlink/core/services/relay"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	"github.com/smartcontractkit/libocr/offchainreporting2/chains/evmutil"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median/evmreportcodec"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/sqlx"
)

var _ service.Service = (*relayer)(nil)
var _ relay.Relayer = (*relayer)(nil)

type relayer struct {
	db       *sqlx.DB
	keystore keystore.Master
	chainSet evm.ChainSet
	lggr     logger.Logger
}

func NewRelayer(config relay.Config) *relayer {
	return &relayer{
		db:       config.Db,
		keystore: config.Keystore,
		chainSet: config.ChainSet,
		lggr:     config.Lggr,
	}
}

// No subservices started on relay start, but when the first job is started
func (r relayer) Start() error {
	return nil
}

// No peristent subservices to close on relay close
func (r relayer) Close() error {
	return nil
}

// Always ready
func (r relayer) Ready() error {
	return nil
}

// Always healthy
func (r relayer) Healthy() error {
	return nil
}

type Config struct {
	ChainID utils.Big `json:"chainID"`
}

type OCR2Spec struct {
	ID             int32
	ContractID     null.String
	OCRKeyBundleID null.String
	TransmitterID  null.String
	IsBootstrap    bool
	RelayConfig    models.JSON
}

func (r relayer) NewOCR2Provider(externalJobID uuid.UUID, s interface{}) (relay.OCR2Provider, error) {
	spec, ok := s.(OCR2Spec)
	if !ok {
		return nil, errors.New("unsuccessful cast to 'ethereum.OCR2Spec'")
	}
	var c Config
	err := json.Unmarshal(spec.RelayConfig.Bytes(), &c)
	if err != nil {
		return nil, err
	}

	chain, err := r.chainSet.Get(c.ChainID.ToInt())
	if err != nil {
		return nil, err
	}
	// TODO: more validation
	contractAddress := common.HexToAddress(spec.ContractID.String)

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

	ocrDB := ocr2.NewDB(r.db.DB, spec.ID, r.lggr)

	tracker := ocr2.NewOCRContractTracker(
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
			config:                 chain.Config(),
		}, nil
	}

	reportCodec := evmreportcodec.ReportCodec{}

	contractABI, err := abi.JSON(strings.NewReader(ocr2aggregator.OCR2AggregatorABI))
	if err != nil {
		return nil, errors.Wrap(err, "could not get contract ABI JSON")
	}

	if !spec.TransmitterID.Valid {
		return nil, errors.New("transmitter address is required")
	}

	transmitterAddress := common.HexToAddress(spec.TransmitterID.String)
	strategy := txm.NewQueueingTxStrategy(externalJobID, chain.Config().OCRDefaultTransactionQueueDepth(), false)

	contractTransmitter := ocr2.NewOCRContractTransmitter(
		contract.Address(),
		contractCaller,
		contractABI,
		ocrcommon.NewTransmitter(chain.TxManager(), transmitterAddress, chain.Config().EvmGasLimitDefault(), strategy),
		chain.LogBroadcaster(),
		tracker,
		r.lggr,
	)

	// Fetch the specified OCR2 key bundle
	var kbID string
	if spec.OCRKeyBundleID.Valid {
		kbID = spec.OCRKeyBundleID.String
	} else if kbID, err = chain.Config().OCR2KeyBundleID(); err != nil {
		return nil, err
	}

	kb, err := r.keystore.OCR2().Get(kbID)
	if err != nil {
		return nil, err
	}

	return &ocr2Provider{
		tracker:                tracker,
		offchainConfigDigester: offchainConfigDigester,
		reportCodec:            reportCodec,
		contractTransmitter:    contractTransmitter,
		keyBundle:              kb,
		config:                 chain.Config(),
	}, nil
}

var _ service.Service = (*ocr2Provider)(nil)

type ocr2Provider struct {
	tracker                *ocr2.OCRContractTracker
	offchainConfigDigester types.OffchainConfigDigester
	reportCodec            median.ReportCodec
	contractTransmitter    *ocr2.OCRContractTransmitter
	keyBundle              ocr2key.KeyBundle
	config                 config.OCR2Config
}

// On start, an ethereum ocr2 provider will start the contract tracker.
func (p ocr2Provider) Start() error {
	return p.tracker.Start()
}

// On close, an ethereum ocr2 provider will close the contract tracker.
func (p ocr2Provider) Close() error {
	return p.tracker.Close()
}

// An ethereum ocr2 provider is ready if the contract tracker is ready.
func (p ocr2Provider) Ready() error {
	return p.tracker.Ready()
}

// An ethereum ocr2 provider is healthy if the contract tracker is healthy.
func (p ocr2Provider) Healthy() error {
	return p.tracker.Healthy()
}

func (p ocr2Provider) OffchainKeyring() types.OffchainKeyring {
	return &p.keyBundle.OffchainKeyring
}

func (p ocr2Provider) OnchainKeyring() types.OnchainKeyring {
	return &p.keyBundle.OnchainKeyring
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

func (p ocr2Provider) OCRConfig() config.OCR2Config {
	return p.config
}
