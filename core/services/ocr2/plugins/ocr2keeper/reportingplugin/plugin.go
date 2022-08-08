package reportingplugin

import (
	"context"
	"math/big"
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	httypes "github.com/smartcontractkit/chainlink/core/chains/evm/headtracker/types"
	"github.com/smartcontractkit/chainlink/core/logger"
)

type Config interface {
	EvmEIP1559DynamicFees() bool
	KeySpecificMaxGasPriceWei(addr gethcommon.Address) *big.Int
	KeeperDefaultTransactionQueueDepth() uint32
	KeeperGasPriceBufferPercent() uint32
	KeeperGasTipCapBufferPercent() uint32
	KeeperBaseFeeBufferPercent() uint32
	KeeperMaximumGracePeriod() int64
	KeeperRegistryCheckGasOverhead() uint64
	KeeperRegistryPerformGasOverhead() uint64
	KeeperRegistrySyncInterval() time.Duration
	KeeperRegistrySyncUpkeepQueueSize() uint32
	KeeperCheckUpkeepGasPriceFeatureEnabled() bool
	KeeperTurnLookBack() int64
	KeeperTurnFlagEnabled() bool
	LogSQL() bool
}

// plugin implements types.ReportingPlugin interface with the keepers-specific logic.
type plugin struct {
	logger    logger.Logger
	headsMngr *headsMngr
	// TODO: Keepers ORM
}

// NewPlugin is the constructor of plugin
func NewPlugin(logger logger.Logger, headBroadcaster httypes.HeadBroadcaster) types.ReportingPlugin {
	hm := newHeadsMngr(logger, headBroadcaster)
	hm.start()

	return &plugin{
		logger:    logger,
		headsMngr: hm,
	}
}

func (p *plugin) Query(context.Context, types.ReportTimestamp) (types.Query, error) {
	currentHead := p.headsMngr.getCurrentHead()

	p.logger.Info("Query()", currentHead)
	return []byte("Query()"), nil
}

func (p *plugin) Observation(_ context.Context, _ types.ReportTimestamp, q types.Query) (types.Observation, error) {
	p.logger.Info("Observation()", string(q))
	return []byte("Observation()"), nil
}

func (p *plugin) Report(_ context.Context, _ types.ReportTimestamp, q types.Query, _ []types.AttributedObservation) (bool, types.Report, error) {
	p.logger.Info("Report()", string(q))
	return true, []byte("Report()"), nil
}

func (p *plugin) ShouldAcceptFinalizedReport(_ context.Context, _ types.ReportTimestamp, r types.Report) (bool, error) {
	p.logger.Info("ShouldAcceptFinalizedReport()", string(r))
	return true, nil
}

func (p *plugin) ShouldTransmitAcceptedReport(_ context.Context, _ types.ReportTimestamp, r types.Report) (bool, error) {
	p.logger.Info("ShouldTransmitAcceptedReport()", string(r))
	return true, nil
}

func (p *plugin) Close() error {
	p.headsMngr.stop()
	return nil
}
