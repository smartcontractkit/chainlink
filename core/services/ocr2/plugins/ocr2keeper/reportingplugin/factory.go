package reportingplugin

import (
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/chains/evm/gas"
	httypes "github.com/smartcontractkit/chainlink/core/chains/evm/headtracker/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

const (
	reportingPluginName   = "OCR2Keeper"
	generateUniqueReports = false
)

// FactoryOptions contains required options to create a reporting plugin factory
type FactoryOptions struct {
	Logger          logger.Logger
	JobID           int32
	ChainID         int64
	Cfg             Config
	ORM             ORM
	EthClient       evmclient.Client
	HeadBroadcaster httypes.HeadBroadcaster
	ContractAddress string
	PipelineRunner  pipeline.Runner
	GasEstimator    gas.Estimator
	PluginLimits    types.ReportingPluginLimits
}

// factory implements types.ReportingPluginFactory interface and creates keepers reporting plugin.
type factory struct {
	logger          logger.Logger
	jobID           int32
	chainID         int64
	cfg             Config
	orm             ORM
	ethClient       evmclient.Client
	hb              httypes.HeadBroadcaster
	contractAddress string
	pr              pipeline.Runner
	gasEstimator    gas.Estimator
	pluginLimits    types.ReportingPluginLimits
}

// NewFactory is the constructor of factory
func NewFactory(opts FactoryOptions) types.ReportingPluginFactory {
	return &factory{
		logger:          opts.Logger,
		jobID:           opts.JobID,
		chainID:         opts.ChainID,
		cfg:             opts.Cfg,
		orm:             opts.ORM,
		ethClient:       opts.EthClient,
		hb:              opts.HeadBroadcaster,
		contractAddress: opts.ContractAddress,
		pr:              opts.PipelineRunner,
		gasEstimator:    opts.GasEstimator,
		pluginLimits:    opts.PluginLimits,
	}
}

func (f *factory) NewReportingPlugin(rpc types.ReportingPluginConfig) (types.ReportingPlugin, types.ReportingPluginInfo, error) {
	p := NewPlugin(f.logger, f.jobID, f.chainID, f.cfg, f.orm, f.ethClient, f.hb, f.contractAddress, f.pr, f.gasEstimator)
	return p, types.ReportingPluginInfo{
		Name:          reportingPluginName,
		UniqueReports: generateUniqueReports,
		Limits:        f.pluginLimits,
	}, nil
}
