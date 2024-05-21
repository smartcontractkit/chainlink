package functions

import (
	"context"
	"database/sql"
	"encoding/binary"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

type FunctionsPluginType int

const (
	FunctionsPlugin FunctionsPluginType = iota
	ThresholdPlugin
	S4Plugin
)

type configPoller struct {
	lggr               logger.Logger
	destChainLogPoller logpoller.LogPoller
	targetContract     atomic.Pointer[common.Address]
	pluginType         FunctionsPluginType
}

var _ types.ConfigPoller = &configPoller{}
var _ types.RouteUpdateSubscriber = &configPoller{}

// ConfigSet Common to all OCR2 evm based contracts: https://github.com/smartcontractkit/libocr/blob/master/contract2/dev/OCR2Abstract.sol
var ConfigSet common.Hash

var defaultABI abi.ABI

const configSetEventName = "ConfigSet"

func init() {
	var err error
	abiPointer, err := ocr2aggregator.OCR2AggregatorMetaData.GetAbi()
	if err != nil {
		panic(err)
	}
	defaultABI = *abiPointer
	ConfigSet = defaultABI.Events[configSetEventName].ID
}

func unpackLogData(d []byte) (*ocr2aggregator.OCR2AggregatorConfigSet, error) {
	unpacked := new(ocr2aggregator.OCR2AggregatorConfigSet)
	err := defaultABI.UnpackIntoInterface(unpacked, configSetEventName, d)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unpack log data")
	}
	return unpacked, nil
}

func configFromLog(logData []byte, pluginType FunctionsPluginType) (ocrtypes.ContractConfig, error) {
	unpacked, err := unpackLogData(logData)
	if err != nil {
		return ocrtypes.ContractConfig{}, err
	}

	var transmitAccounts []ocrtypes.Account
	for _, addr := range unpacked.Transmitters {
		transmitAccounts = append(transmitAccounts, ocrtypes.Account(addr.String()))
	}
	var signers []ocrtypes.OnchainPublicKey
	for _, addr := range unpacked.Signers {
		addr := addr
		signers = append(signers, addr[:])
	}

	// Replace the first two bytes of the config digest with the plugin type to avoid duplicate config digests between Functions plugins
	switch pluginType {
	case FunctionsPlugin:
		// FunctionsPluginType should already have the correct prefix, so this is a no-op
	case ThresholdPlugin:
		binary.BigEndian.PutUint16(unpacked.ConfigDigest[:2], uint16(ThresholdDigestPrefix))
	case S4Plugin:
		binary.BigEndian.PutUint16(unpacked.ConfigDigest[:2], uint16(S4DigestPrefix))
	default:
		return ocrtypes.ContractConfig{}, errors.New("unknown plugin type")
	}

	return ocrtypes.ContractConfig{
		ConfigDigest:          unpacked.ConfigDigest,
		ConfigCount:           unpacked.ConfigCount,
		Signers:               signers,
		Transmitters:          transmitAccounts,
		F:                     unpacked.F,
		OnchainConfig:         unpacked.OnchainConfig,
		OffchainConfigVersion: unpacked.OffchainConfigVersion,
		OffchainConfig:        unpacked.OffchainConfig,
	}, nil
}

func configPollerFilterName(addr common.Address) string {
	return logpoller.FilterName("FunctionsOCR2ConfigPoller", addr.String())
}

func NewFunctionsConfigPoller(pluginType FunctionsPluginType, destChainPoller logpoller.LogPoller, lggr logger.Logger) (*configPoller, error) {
	cp := &configPoller{
		lggr:               lggr,
		destChainLogPoller: destChainPoller,
		pluginType:         pluginType,
	}
	return cp, nil
}

func (cp *configPoller) Start() {}

func (cp *configPoller) Close() error {
	return nil
}

func (cp *configPoller) Notify() <-chan struct{} {
	return nil
}

func (cp *configPoller) Replay(ctx context.Context, fromBlock int64) error {
	return cp.destChainLogPoller.Replay(ctx, fromBlock)
}

func (cp *configPoller) LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest ocrtypes.ConfigDigest, err error) {
	contractAddr := cp.targetContract.Load()
	if contractAddr == nil {
		return 0, ocrtypes.ConfigDigest{}, nil
	}

	latest, err := cp.destChainLogPoller.LatestLogByEventSigWithConfs(ctx, ConfigSet, *contractAddr, 1)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ocrtypes.ConfigDigest{}, nil
		}
		return 0, ocrtypes.ConfigDigest{}, err
	}
	latestConfigSet, err := configFromLog(latest.Data, cp.pluginType)
	if err != nil {
		return 0, ocrtypes.ConfigDigest{}, err
	}
	return uint64(latest.BlockNumber), latestConfigSet.ConfigDigest, nil
}

func (cp *configPoller) LatestConfig(ctx context.Context, changedInBlock uint64) (ocrtypes.ContractConfig, error) {
	// NOTE: if targetContract changes between invocations of LatestConfigDetails() and LatestConfig()
	// (unlikely), we'll return an error here and libocr will re-try.
	contractAddr := cp.targetContract.Load()
	if contractAddr == nil {
		return ocrtypes.ContractConfig{}, errors.New("no target contract address set yet")
	}

	lgs, err := cp.destChainLogPoller.Logs(ctx, int64(changedInBlock), int64(changedInBlock), ConfigSet, *contractAddr)
	if err != nil {
		return ocrtypes.ContractConfig{}, err
	}
	if len(lgs) == 0 {
		return ocrtypes.ContractConfig{}, errors.New("no logs found")
	}
	latestConfigSet, err := configFromLog(lgs[len(lgs)-1].Data, cp.pluginType)
	if err != nil {
		return ocrtypes.ContractConfig{}, err
	}
	cp.lggr.Infow("LatestConfig", "latestConfig", latestConfigSet)
	return latestConfigSet, nil
}

func (cp *configPoller) LatestBlockHeight(ctx context.Context) (blockHeight uint64, err error) {
	latest, err := cp.destChainLogPoller.LatestBlock(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return uint64(latest.BlockNumber), nil
}

// called from LogPollerWrapper in a separate goroutine
func (cp *configPoller) UpdateRoutes(ctx context.Context, activeCoordinator common.Address, proposedCoordinator common.Address) error {
	cp.targetContract.Store(&activeCoordinator)
	// Register filters for both active and proposed
	err := cp.destChainLogPoller.RegisterFilter(ctx, logpoller.Filter{Name: configPollerFilterName(activeCoordinator), EventSigs: []common.Hash{ConfigSet}, Addresses: []common.Address{activeCoordinator}})
	if err != nil {
		return err
	}
	err = cp.destChainLogPoller.RegisterFilter(ctx, logpoller.Filter{Name: configPollerFilterName(proposedCoordinator), EventSigs: []common.Hash{ConfigSet}, Addresses: []common.Address{activeCoordinator}})
	if err != nil {
		return err
	}
	// TODO: unregister old filter (needs refactor to get pg.Queryer)
	return nil
}
