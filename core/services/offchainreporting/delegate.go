package offchainreporting

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/log"
	"github.com/smartcontractkit/chainlink/core/services/synchronization"
	"github.com/smartcontractkit/chainlink/core/services/telemetry"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/libocr/gethwrappers/offchainaggregator"
	ocr "github.com/smartcontractkit/libocr/offchainreporting"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"
)

type Delegate struct {
	db             *gorm.DB
	jobORM         job.ORM
	config         *orm.Config
	keyStore       *KeyStore
	pipelineRunner pipeline.Runner
	ethClient      eth.Client
	logBroadcaster log.Broadcaster
	peerWrapper    *SingletonPeerWrapper
}

func NewDelegate(
	db *gorm.DB,
	jobORM job.ORM,
	config *orm.Config,
	keyStore *KeyStore,
	pipelineRunner pipeline.Runner,
	ethClient eth.Client,
	logBroadcaster log.Broadcaster,
	peerWrapper *SingletonPeerWrapper,
) *Delegate {
	return &Delegate{db, jobORM, config, keyStore, pipelineRunner, ethClient, logBroadcaster, peerWrapper}
}

func (d Delegate) JobType() job.Type {
	return job.OffchainReporting
}

func (d Delegate) ServicesForSpec(jobSpec job.SpecDB) (services []job.Service, err error) {
	if jobSpec.OffchainreportingOracleSpec == nil {
		return nil, errors.Errorf("offchainreporting.Delegate expects an *job.OffchainreportingOracleSpec to be present, got %v", jobSpec)
	}
	concreteSpec := jobSpec.OffchainreportingOracleSpec

	contractFilterer, err := offchainaggregator.NewOffchainAggregatorFilterer(concreteSpec.ContractAddress.Address(), d.ethClient)
	if err != nil {
		return nil, errors.Wrap(err, "could not instantiate NewOffchainAggregatorFilterer")
	}

	contractCaller, err := offchainaggregator.NewOffchainAggregatorCaller(concreteSpec.ContractAddress.Address(), d.ethClient)
	if err != nil {
		return nil, errors.Wrap(err, "could not instantiate NewOffchainAggregatorCaller")
	}

	ocrContract, err := NewOCRContractConfigTracker(
		concreteSpec.ContractAddress.Address(),
		contractFilterer,
		contractCaller,
		d.ethClient,
		d.logBroadcaster,
		jobSpec.ID,
		*logger.Default,
	)
	if err != nil {
		return nil, errors.Wrap(err, "error calling NewOCRContract")
	}

	peerID, err := d.config.P2PPeerID(concreteSpec.P2PPeerID)
	if err != nil {
		return nil, err
	}
	peerWrapper := d.peerWrapper
	if peerWrapper == nil {
		return nil, errors.New("cannot setup OCR job service, libp2p peer was missing")
	} else if !peerWrapper.IsStarted() {
		return nil, errors.New("peerWrapper is not started. OCR jobs require a started and running peer. Did you forget to specify P2P_LISTEN_PORT?")
	} else if peerWrapper.PeerID != peerID {
		return nil, errors.Errorf("given peer with ID '%s' does not match OCR configured peer with ID: %s", peerWrapper.PeerID.String(), peerID.String())
	}
	bootstrapPeers, err := d.config.P2PBootstrapPeers(concreteSpec.P2PBootstrapPeers)
	if err != nil {
		return nil, err
	}

	loggerWith := logger.CreateLogger(logger.Default.With(
		"contractAddress", concreteSpec.ContractAddress,
		"jobID", jobSpec.ID))
	ocrLogger := NewLogger(loggerWith, d.config.OCRTraceLogging(), func(msg string) {
		d.jobORM.RecordError(context.Background(), jobSpec.ID, msg)
	})

	var endpointURL *url.URL
	if me := d.config.OCRMonitoringEndpoint(concreteSpec.MonitoringEndpoint); me != "" {
		endpointURL, err = url.Parse(me)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid monitoring url: %s", me)
		}
	} else {
		endpointURL = d.config.ExplorerURL()
	}

	var monitoringEndpoint ocrtypes.MonitoringEndpoint
	if endpointURL != nil {
		client := synchronization.NewExplorerClient(endpointURL, d.config.ExplorerAccessKey(), d.config.ExplorerSecret())
		monitoringEndpoint = telemetry.NewAgent(client)
		services = append(services, client)
	} else {
		monitoringEndpoint = ocrtypes.MonitoringEndpoint(nil)
	}

	lc := ocrtypes.LocalConfig{
		BlockchainTimeout:                      d.config.OCRBlockchainTimeout(time.Duration(concreteSpec.BlockchainTimeout)),
		ContractConfigConfirmations:            d.config.OCRContractConfirmations(concreteSpec.ContractConfigConfirmations),
		ContractConfigTrackerPollInterval:      d.config.OCRContractPollInterval(time.Duration(concreteSpec.ContractConfigTrackerPollInterval)),
		ContractConfigTrackerSubscribeInterval: d.config.OCRContractSubscribeInterval(time.Duration(concreteSpec.ContractConfigTrackerSubscribeInterval)),
		ContractTransmitterTransmitTimeout:     d.config.OCRContractTransmitterTransmitTimeout(),
		DatabaseTimeout:                        d.config.OCRDatabaseTimeout(),
		DataSourceTimeout:                      d.config.OCRObservationTimeout(time.Duration(concreteSpec.ObservationTimeout)),
	}
	if d.config.Dev() {
		// Skips config validation so we can use any config parameters we want.
		// For example to lower contractConfigTrackerPollInterval to speed up tests.
		lc.DevelopmentMode = ocrtypes.EnableDangerousDevelopmentMode
	}
	if err := ocr.SanityCheckLocalConfig(lc); err != nil {
		return nil, err
	}
	logger.Info(fmt.Sprintf("OCR job using local config %+v", lc))

	if concreteSpec.IsBootstrapPeer {
		bootstrapper, err := ocr.NewBootstrapNode(ocr.BootstrapNodeArgs{
			BootstrapperFactory:   peerWrapper.Peer,
			Bootstrappers:         bootstrapPeers,
			ContractConfigTracker: ocrContract,
			Database:              NewDB(d.db.DB(), concreteSpec.ID),
			LocalConfig:           lc,
			Logger:                ocrLogger,
			MonitoringEndpoint:    monitoringEndpoint,
		})
		if err != nil {
			return nil, errors.Wrap(err, "error calling NewBootstrapNode")
		}
		services = append(services, bootstrapper)
	} else {
		if len(bootstrapPeers) < 1 {
			return nil, errors.New("need at least one bootstrap peer")
		}
		kb, err := d.config.OCRKeyBundleID(concreteSpec.EncryptedOCRKeyBundleID)
		if err != nil {
			return nil, err
		}
		ocrkey, exists := d.keyStore.DecryptedOCRKey(kb)
		if !exists {
			return nil, errors.Errorf("OCR key '%v' does not exist", concreteSpec.EncryptedOCRKeyBundleID)
		}
		contractABI, err := abi.JSON(strings.NewReader(offchainaggregator.OffchainAggregatorABI))
		if err != nil {
			return nil, errors.Wrap(err, "could not get contract ABI JSON")
		}

		ta, err := d.config.OCRTransmitterAddress(concreteSpec.TransmitterAddress)
		if err != nil {
			return nil, err
		}
		contractTransmitter := NewOCRContractTransmitter(concreteSpec.ContractAddress.Address(), contractCaller, contractABI,
			NewTransmitter(d.db.DB(), ta.Address(), d.config.EthGasLimitDefault()))

		oracle, err := ocr.NewOracle(ocr.OracleArgs{
			Database:                     NewDB(d.db.DB(), concreteSpec.ID),
			Datasource:                   dataSource{jobID: jobSpec.ID, pipelineRunner: d.pipelineRunner},
			LocalConfig:                  lc,
			ContractTransmitter:          contractTransmitter,
			ContractConfigTracker:        ocrContract,
			PrivateKeys:                  &ocrkey,
			BinaryNetworkEndpointFactory: peerWrapper.Peer,
			MonitoringEndpoint:           monitoringEndpoint,
			Logger:                       ocrLogger,
			Bootstrappers:                bootstrapPeers,
		})
		if err != nil {
			return nil, errors.Wrap(err, "error calling NewOracle")
		}
		services = append(services, oracle)
	}

	return services, nil
}
