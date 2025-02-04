package llo

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/configurator"
)

type InstanceType string

const (
	InstanceTypeBlue  InstanceType = InstanceType("Blue")
	InstanceTypeGreen InstanceType = InstanceType("Green")
)

var (
	NoLimitSortAsc = query.NewLimitAndSort(query.Limit{}, query.NewSortBySequence(query.Asc))
)

type ConfigPollerService interface {
	services.Service
	ocrtypes.ContractConfigTracker
}

type LogPoller interface {
	LatestBlock(ctx context.Context) (logpoller.LogPollerBlock, error)
	FilteredLogs(ctx context.Context, filter []query.Expression, limitAndSort query.LimitAndSort, queryName string) ([]logpoller.Log, error)
}

// ConfigCache is most likely the global RetirementReportCache. Every config
// ever seen by this tracker will be stored at least once in the cache.
type ConfigCache interface {
	StoreConfig(ctx context.Context, cd ocrtypes.ConfigDigest, signers [][]byte, f uint8) error
}

type configPoller struct {
	services.Service
	eng *services.Engine

	lp          LogPoller
	cc          ConfigCache
	addr        common.Address
	donID       uint32
	donIDTopic  [32]byte
	filterExprs []query.Expression

	fromBlock int64
	mu        sync.RWMutex

	instanceType InstanceType
}

func DonIDToBytes32(donID uint32) [32]byte {
	var b [32]byte
	copy(b[:], common.LeftPadBytes(big.NewInt(int64(donID)).Bytes(), 32))
	return b
}

// NewConfigPoller creates a new LLOConfigPoller
func NewConfigPoller(lggr logger.Logger, lp LogPoller, cc ConfigCache, addr common.Address, donID uint32, instanceType InstanceType, fromBlock uint64) ConfigPollerService {
	return newConfigPoller(lggr, lp, cc, addr, donID, instanceType, fromBlock)
}

func newConfigPoller(lggr logger.Logger, lp LogPoller, cc ConfigCache, addr common.Address, donID uint32, instanceType InstanceType, fromBlock uint64) *configPoller {
	donIDTopic := DonIDToBytes32(donID)
	exprs := []query.Expression{
		logpoller.NewAddressFilter(addr),
		query.Or(
			logpoller.NewEventSigFilter(ProductionConfigSet),
			logpoller.NewEventSigFilter(StagingConfigSet),
		),
		logpoller.NewEventByTopicFilter(1, []logpoller.HashedValueComparator{
			{Values: []common.Hash{donIDTopic}, Operator: primitives.Eq},
		}),
		// NOTE: Optimize for fast config switches. On Arbitrum, finalization
		// can take tens of minutes
		// (https://grafana.ops.prod.cldev.sh/d/e0453cc9-4b4a-41e1-9f01-7c21de805b39/blockchain-finality-and-gas?orgId=1&var-env=All&var-network_name=ethereum-testnet-sepolia-arbitrum-1&var-network_name=ethereum-mainnet-arbitrum-1&from=1732460992641&to=1732547392641)
		query.Confidence(primitives.Unconfirmed),
	}
	cp := &configPoller{
		lp:           lp,
		cc:           cc,
		addr:         addr,
		donID:        donID,
		donIDTopic:   DonIDToBytes32(donID),
		filterExprs:  exprs,
		instanceType: instanceType,
		fromBlock:    int64(fromBlock),
	}
	cp.Service, cp.eng = services.Config{
		Name: "LLOConfigPoller",
	}.NewServiceEngine(logger.Sugared(lggr).Named(string(instanceType)).With("instanceType", instanceType))

	return cp
}

func (cp *configPoller) Notify() <-chan struct{} {
	return nil // rely on libocr's builtin config polling
}

// LatestConfigDetails returns the latest config details from the logs
func (cp *configPoller) LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest ocrtypes.ConfigDigest, err error) {
	latestConfig, log, err := cp.latestConfig(ctx, cp.readFromBlock(), math.MaxInt64)
	if err != nil {
		return 0, ocrtypes.ConfigDigest{}, fmt.Errorf("failed to get latest config: %w", err)
	}
	// Slight optimization, since we only care about the latest log, we can
	// avoid re-scanning from the original fromBlock every time by setting
	// fromBlock to the latest seen log here.
	//
	// This should always be safe even if LatestConfigDetails is called
	// concurrently.
	cp.setFromBlock(log.BlockNumber)
	return uint64(log.BlockNumber), latestConfig.ConfigDigest, nil // #nosec G115 // log.BlockNumber will never be negative
}

func (cp *configPoller) readFromBlock() int64 {
	cp.mu.RLock()
	defer cp.mu.RUnlock()
	return cp.fromBlock
}

func (cp *configPoller) setFromBlock(fromBlock int64) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	if fromBlock > cp.fromBlock {
		cp.fromBlock = fromBlock
	}
}

func (cp *configPoller) latestConfig(ctx context.Context, fromBlock, toBlock int64) (latestConfig FullConfigFromLog, latestLog logpoller.Log, err error) {
	// Get all configset logs and run through them forwards
	// NOTE: It's useful to get _all_ logs rather than just the latest since
	// they are stored in the ConfigCache
	exprs := make([]query.Expression, 0, len(cp.filterExprs)+2)
	exprs = append(exprs, cp.filterExprs...)
	exprs = append(exprs,
		query.Block(strconv.FormatInt(fromBlock, 10), primitives.Gte),
		query.Block(strconv.FormatInt(toBlock, 10), primitives.Lte),
	)
	logs, err := cp.lp.FilteredLogs(ctx, exprs, NoLimitSortAsc, "LLOConfigPoller - latestConfig")
	if err != nil {
		return latestConfig, latestLog, fmt.Errorf("failed to get logs: %w", err)
	}
	for _, log := range logs {
		if !bytes.Equal(log.Topics[1], cp.donIDTopic[:]) {
			// skip logs for other donIDs, shouldn't happen given the
			// FilterLogs call, but belts and braces
			continue
		}
		switch log.EventSig {
		case ProductionConfigSet:
			event, err := DecodeProductionConfigSetLog(log.Data)
			if err != nil {
				return latestConfig, log, fmt.Errorf("failed to unpack ProductionConfigSet log data: %w", err)
			}

			if err = cp.cc.StoreConfig(ctx, event.ConfigDigest, event.Signers, event.F); err != nil {
				cp.eng.Errorf("failed to store production config: %v", err)
			}

			isProduction := (cp.instanceType != InstanceTypeBlue) == event.IsGreenProduction
			if isProduction {
				latestLog = log
				latestConfig, err = FullConfigFromProductionConfigSet(event)
				if err != nil {
					return latestConfig, latestLog, fmt.Errorf("FullConfigFromProductionConfigSet failed: %w", err)
				}
			}
		case StagingConfigSet:
			event, err := DecodeStagingConfigSetLog(log.Data)
			if err != nil {
				return latestConfig, latestLog, fmt.Errorf("failed to unpack ProductionConfigSet log data: %w", err)
			}

			if err = cp.cc.StoreConfig(ctx, event.ConfigDigest, event.Signers, event.F); err != nil {
				cp.eng.Errorf("failed to store staging config: %v", err)
			}

			isProduction := (cp.instanceType != InstanceTypeBlue) == event.IsGreenProduction
			if !isProduction {
				latestLog = log
				latestConfig, err = FullConfigFromStagingConfigSet(event)
				if err != nil {
					return latestConfig, latestLog, fmt.Errorf("FullConfigFromStagingConfigSet failed: %w", err)
				}
			}
		default:
			// ignore unknown log types
			continue
		}
	}

	return
}

// LatestConfig returns the latest config from the logs starting from a certain block
func (cp *configPoller) LatestConfig(ctx context.Context, changedInBlock uint64) (ocrtypes.ContractConfig, error) {
	cfg, latestLog, err := cp.latestConfig(ctx, int64(changedInBlock), math.MaxInt64) // #nosec G115
	if err != nil {
		return ocrtypes.ContractConfig{}, fmt.Errorf("failed to get latest config: %w", err)
	}
	cp.eng.Infow("LatestConfig fetched", "config", cfg.ContractConfig, "txHash", latestLog.TxHash, "blockNumber", latestLog.BlockNumber, "blockHash", latestLog.BlockHash, "logIndex", latestLog.LogIndex, "instanceType", cp.instanceType, "donID", cp.donID, "changedInBlock", changedInBlock)
	return cfg.ContractConfig, nil
}

// LatestBlockHeight returns the latest block height from the logs
func (cp *configPoller) LatestBlockHeight(ctx context.Context) (blockHeight uint64, err error) {
	latest, err := cp.lp.LatestBlock(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return uint64(latest.BlockNumber), nil // #nosec G115 // latest.BlockNumber will never be negative
}

func (cp *configPoller) InstanceType() InstanceType {
	return cp.instanceType
}

// FullConfigFromLog defines the contract config with the donID
type FullConfigFromLog struct {
	ocrtypes.ContractConfig
	donID uint32
}

func FullConfigFromProductionConfigSet(unpacked configurator.ConfiguratorProductionConfigSet) (FullConfigFromLog, error) {
	transmitAccounts := make([]ocrtypes.Account, len(unpacked.OffchainTransmitters))
	for i, addr := range unpacked.OffchainTransmitters {
		transmitAccounts[i] = ocrtypes.Account(hex.EncodeToString(addr[:]))
	}
	signers := make([]ocrtypes.OnchainPublicKey, len(unpacked.Signers))
	for i, addr := range unpacked.Signers {
		signers[i] = addr
	}

	donIDBig := common.Hash(unpacked.ConfigId).Big()
	if donIDBig.Cmp(big.NewInt(math.MaxUint32)) > 0 {
		return FullConfigFromLog{}, errors.Errorf("donID %s is too large", donIDBig)
	}
	donID := uint32(donIDBig.Uint64()) // #nosec G115

	return FullConfigFromLog{
		donID: donID,
		ContractConfig: ocrtypes.ContractConfig{
			ConfigDigest:          unpacked.ConfigDigest,
			ConfigCount:           unpacked.ConfigCount,
			Signers:               signers,
			Transmitters:          transmitAccounts,
			F:                     unpacked.F,
			OnchainConfig:         unpacked.OnchainConfig,
			OffchainConfigVersion: unpacked.OffchainConfigVersion,
			OffchainConfig:        unpacked.OffchainConfig,
		},
	}, nil
}

func FullConfigFromStagingConfigSet(unpacked configurator.ConfiguratorStagingConfigSet) (FullConfigFromLog, error) {
	transmitAccounts := make([]ocrtypes.Account, len(unpacked.OffchainTransmitters))
	for i, addr := range unpacked.OffchainTransmitters {
		transmitAccounts[i] = ocrtypes.Account(hex.EncodeToString(addr[:]))
	}
	signers := make([]ocrtypes.OnchainPublicKey, len(unpacked.Signers))
	for i, addr := range unpacked.Signers {
		signers[i] = addr
	}

	donIDBig := common.Hash(unpacked.ConfigId).Big()
	if donIDBig.Cmp(big.NewInt(math.MaxUint32)) > 0 {
		return FullConfigFromLog{}, errors.Errorf("donID %s is too large", donIDBig)
	}
	donID := uint32(donIDBig.Uint64()) // #nosec G115

	return FullConfigFromLog{
		donID: donID,
		ContractConfig: ocrtypes.ContractConfig{
			ConfigDigest:          unpacked.ConfigDigest,
			ConfigCount:           unpacked.ConfigCount,
			Signers:               signers,
			Transmitters:          transmitAccounts,
			F:                     unpacked.F,
			OnchainConfig:         unpacked.OnchainConfig,
			OffchainConfigVersion: unpacked.OffchainConfigVersion,
			OffchainConfig:        unpacked.OffchainConfig,
		},
	}, nil
}

func DecodeProductionConfigSetLog(d []byte) (configurator.ConfiguratorProductionConfigSet, error) {
	unpacked := new(configurator.ConfiguratorProductionConfigSet)

	err := configuratorABI.UnpackIntoInterface(unpacked, "ProductionConfigSet", d)
	if err != nil {
		return configurator.ConfiguratorProductionConfigSet{}, errors.Wrap(err, "failed to unpack log data")
	}
	return *unpacked, nil
}

func DecodeStagingConfigSetLog(d []byte) (configurator.ConfiguratorStagingConfigSet, error) {
	unpacked := new(configurator.ConfiguratorStagingConfigSet)

	err := configuratorABI.UnpackIntoInterface(unpacked, "StagingConfigSet", d)
	if err != nil {
		return configurator.ConfiguratorStagingConfigSet{}, errors.Wrap(err, "failed to unpack log data")
	}
	return *unpacked, nil
}
