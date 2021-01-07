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
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/libocr/gethwrappers/offchainaggregator"
	ocr "github.com/smartcontractkit/libocr/offchainreporting"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"
)

func RegisterJobType(
	db *gorm.DB,
	jobORM job.ORM,
	config *orm.Config,
	keyStore *KeyStore,
	jobSpawner job.Spawner,
	pipelineRunner pipeline.Runner,
	ethClient eth.Client,
	logBroadcaster log.Broadcaster,
	peerWrapper *SingletonPeerWrapper,
) {
	jobSpawner.RegisterDelegate(
		NewJobSpawnerDelegate(db, jobORM, config, keyStore, pipelineRunner, ethClient, logBroadcaster, peerWrapper),
	)
}

type jobSpawnerDelegate struct {
	db             *gorm.DB
	jobORM         job.ORM
	config         *orm.Config
	keyStore       *KeyStore
	pipelineRunner pipeline.Runner
	ethClient      eth.Client
	logBroadcaster log.Broadcaster
	peerWrapper    *SingletonPeerWrapper
}

func NewJobSpawnerDelegate(
	db *gorm.DB,
	jobORM job.ORM,
	config *orm.Config,
	keyStore *KeyStore,
	pipelineRunner pipeline.Runner,
	ethClient eth.Client,
	logBroadcaster log.Broadcaster,
	peerWrapper *SingletonPeerWrapper,
) *jobSpawnerDelegate {
	return &jobSpawnerDelegate{db, jobORM, config, keyStore, pipelineRunner, ethClient, logBroadcaster, peerWrapper}
}

func (d jobSpawnerDelegate) JobType() job.Type {
	return job.OffchainReporting
}

func (d jobSpawnerDelegate) ServicesForSpec(spec job.SpecDB) (services []job.Service, err error) {
	if spec.OffchainreportingOracleSpec == nil {
		return nil, errors.Errorf("offchainreporting.jobSpawnerDelegate expects an *job.OffchainreportingOracleSpec to be present, got %v", spec)
	}
	concreteSpec := spec.OffchainreportingOracleSpec

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
		concreteSpec.ID,
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
		"jobID", concreteSpec.ID))
	ocrLogger := NewLogger(loggerWith, d.config.OCRTraceLogging(), func(msg string) {
		d.jobORM.RecordError(context.Background(), spec.ID, msg)
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
			Datasource:                   dataSource{jobID: concreteSpec.ID, pipelineRunner: d.pipelineRunner},
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

// dataSource is an abstraction over the process of initiating a pipeline run
// and capturing the result.  Additionally, it converts the result to an
// ocrtypes.Observation (*big.Int), as expected by the offchain reporting library.
type dataSource struct {
	pipelineRunner pipeline.Runner
	jobID          int32
}

var _ ocrtypes.DataSource = (*dataSource)(nil)

func (ds dataSource) Observe(ctx context.Context) (ocrtypes.Observation, error) {
	runID, err := ds.pipelineRunner.CreateRun(ctx, ds.jobID, nil)
	if err != nil {
		return nil, err
	}

	err = ds.pipelineRunner.AwaitRun(ctx, runID)
	if err != nil {
		return nil, err
	}

	results, err := ds.pipelineRunner.ResultsForRun(ctx, runID)
	if err != nil {
		return nil, errors.Wrapf(err, "pipeline error")
	} else if len(results) != 1 {
		return nil, errors.Errorf("offchain reporting pipeline should have a single output (job spec ID: %v, pipeline run ID: %v)", ds.jobID, runID)
	}

	if results[0].Error != nil {
		return nil, results[0].Error
	}

	asDecimal, err := utils.ToDecimal(results[0].Value)
	if err != nil {
		return nil, err
	}
	return ocrtypes.Observation(asDecimal.BigInt()), nil
}
