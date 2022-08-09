package reportingplugin

import (
	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	httypes "github.com/smartcontractkit/chainlink/core/chains/evm/headtracker/types"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/core/logger"
)

// factory implements types.ReportingPluginFactory interface and creates keepers reporting plugin.
type factory struct {
	logger          logger.Logger
	cfg             Config
	orm             ORM
	ethClient       evmclient.Client
	hb              httypes.HeadBroadcaster
	contractAddress ethkey.EIP55Address
}

// NewFactory is the constructor of factory
func NewFactory(logger logger.Logger, cfg Config, orm ORM, ethClient evmclient.Client, hb httypes.HeadBroadcaster, contractAddress ethkey.EIP55Address) types.ReportingPluginFactory {
	return &factory{
		logger:          logger,
		cfg:             cfg,
		orm:             orm,
		ethClient:       ethClient,
		hb:              hb,
		contractAddress: contractAddress,
	}
}

func (f *factory) NewReportingPlugin(rpc types.ReportingPluginConfig) (types.ReportingPlugin, types.ReportingPluginInfo, error) {
	p := NewPlugin(f.logger, f.cfg, f.orm, f.ethClient, f.hb, f.contractAddress)
	pi := types.ReportingPluginInfo{
		Name:          "OCR2Keeper",
		UniqueReports: false,
		Limits: types.ReportingPluginLimits{
			MaxQueryLength:       1000,
			MaxObservationLength: 1000,
			MaxReportLength:      1000,
		},
	}
	return p, pi, nil
}
