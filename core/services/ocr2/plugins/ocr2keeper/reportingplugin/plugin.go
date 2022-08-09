package reportingplugin

import (
	"context"
	"fmt"
	"math/big"
	"time"

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	httypes "github.com/smartcontractkit/chainlink/core/chains/evm/headtracker/types"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keeper"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
)

type ORM interface {
	RegistryByContractAddress(registryAddress ethkey.EIP55Address) (keeper.Registry, error)
	NewEligibleUpkeepsForRegistry(registryAddress ethkey.EIP55Address, blockNumber int64, gracePeriod int64, binaryHash string) (upkeeps []keeper.UpkeepRegistration, err error)
	EligibleUpkeepsForRegistry(registryAddress ethkey.EIP55Address, blockNumber, gracePeriod int64) (upkeeps []keeper.UpkeepRegistration, err error)
}

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
	logger          logger.Logger
	cfg             Config
	orm             ORM
	ethClient       evmclient.Client
	headsMngr       *headsMngr
	contractAddress ethkey.EIP55Address
	// TODO: Keepers ORM
}

// NewPlugin is the constructor of plugin
func NewPlugin(logger logger.Logger, cfg Config, orm ORM, ethClient evmclient.Client, headBroadcaster httypes.HeadBroadcaster, contractAddress ethkey.EIP55Address) types.ReportingPlugin {
	hm := newHeadsMngr(logger, headBroadcaster)
	hm.start()

	return &plugin{
		logger:          logger,
		cfg:             cfg,
		orm:             orm,
		ethClient:       ethClient,
		headsMngr:       hm,
		contractAddress: contractAddress,
	}
}

func (p *plugin) Query(context.Context, types.ReportTimestamp) (types.Query, error) {
	currentHead := p.headsMngr.getCurrentHead()

	registry, err := p.orm.RegistryByContractAddress(p.contractAddress)
	if err != nil {
		p.logger.Error(errors.Wrap(err, "unable to load registry"))
		return nil, nil
	}

	var activeUpkeeps []keeper.UpkeepRegistration
	if p.cfg.KeeperTurnFlagEnabled() {
		var turnBinary string
		if turnBinary, err = p.turnBlockHashBinary(registry, currentHead, p.cfg.KeeperTurnLookBack()); err != nil {
			return nil, errors.Wrap(err, "unable to get turn block number hash")
		}

		activeUpkeeps, err = p.orm.NewEligibleUpkeepsForRegistry(
			p.contractAddress,
			currentHead.Number,
			p.cfg.KeeperMaximumGracePeriod(),
			turnBinary)
		if err != nil {
			return nil, errors.Wrap(err, "unable to load active registrations")
		}
	} else {
		activeUpkeeps, err = p.orm.EligibleUpkeepsForRegistry(
			p.contractAddress,
			currentHead.Number,
			p.cfg.KeeperMaximumGracePeriod(),
		)
		if err != nil {
			return nil, errors.Wrap(err, "unable to load active registrations")
		}
	}

	p.logger.Info("Query(): ", currentHead, activeUpkeeps)
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

func (p *plugin) turnBlockHashBinary(registry keeper.Registry, head *evmtypes.Head, lookback int64) (string, error) {
	turnBlock := head.Number - (head.Number % int64(registry.BlockCountPerTurn)) - lookback
	block, err := p.ethClient.BlockByNumber(context.Background(), big.NewInt(turnBlock))
	if err != nil {
		return "", err
	}
	hashAtHeight := block.Hash()
	binaryString := fmt.Sprintf("%b", hashAtHeight.Big())
	return binaryString, nil
}
