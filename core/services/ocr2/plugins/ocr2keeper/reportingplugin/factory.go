package reportingplugin

import (
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/chains/evm/gas"
	httypes "github.com/smartcontractkit/chainlink/core/chains/evm/headtracker/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

// factory implements types.ReportingPluginFactory interface and creates keepers reporting plugin.
type factory struct {
	logger          logger.Logger
	jobID           int32
	chainID         string
	cfg             Config
	orm             ORM
	ethClient       evmclient.Client
	hb              httypes.HeadBroadcaster
	contractAddress ethkey.EIP55Address
	pr              pipeline.Runner
	gasEstimator    gas.Estimator
}

// NewFactory is the constructor of factory
func NewFactory(
	logger logger.Logger,
	jobID int32,
	chainID string,
	cfg Config,
	orm ORM,
	ethClient evmclient.Client,
	hb httypes.HeadBroadcaster,
	contractAddress ethkey.EIP55Address,
	pr pipeline.Runner,
	gasEstimator gas.Estimator,
) types.ReportingPluginFactory {
	return &factory{
		logger:          logger,
		jobID:           jobID,
		chainID:         chainID,
		cfg:             cfg,
		orm:             orm,
		ethClient:       ethClient,
		hb:              hb,
		contractAddress: contractAddress,
		pr:              pr,
		gasEstimator:    gasEstimator,
	}
}

func (f *factory) NewReportingPlugin(rpc types.ReportingPluginConfig) (types.ReportingPlugin, types.ReportingPluginInfo, error) {
	p := NewPlugin(f.logger, f.jobID, f.chainID, f.cfg, f.orm, f.ethClient, f.hb, f.contractAddress, f.pr, f.gasEstimator)
	pi := types.ReportingPluginInfo{
		Name:          "OCR2Keeper",
		UniqueReports: false,
		Limits: types.ReportingPluginLimits{
			MaxQueryLength:       1000, // TODO: Configure
			MaxObservationLength: 1000, // TODO: Configure
			MaxReportLength:      1000, // TODO: Configure
		},
	}
	return p, pi, nil
}
