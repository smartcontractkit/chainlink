package llo

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/configurator"
)

type InstanceType string

const (
	InstanceTypeBlue  InstanceType = InstanceType("Blue")
	InstanceTypeGreen InstanceType = InstanceType("Green")
)

type ConfigPollerService interface {
	services.Service
	ocrtypes.ContractConfigTracker
}

type LogPoller interface {
	IndexedLogsByBlockRange(ctx context.Context, start, end int64, eventSig common.Hash, address common.Address, topicIndex int, topicValues []common.Hash) ([]logpoller.Log, error)
	LatestBlock(ctx context.Context) (logpoller.LogPollerBlock, error)
	LogsWithSigs(ctx context.Context, start, end int64, eventSigs []common.Hash, address common.Address) ([]logpoller.Log, error)
}

// ConfigCache is most likely the global RetirementReportCache. Every config
// ever seen by this tracker will be stored at least once in the cache.
type ConfigCache interface {
	StoreConfig(ctx context.Context, cd ocrtypes.ConfigDigest, signers [][]byte, f uint8) error
}

type configPoller struct {
	services.Service
	eng *services.Engine

	lp        LogPoller
	cc        ConfigCache
	addr      common.Address
	donID     uint32
	donIDHash [32]byte

	fromBlock uint64

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
	cp := &configPoller{
		lp:           lp,
		cc:           cc,
		addr:         addr,
		donID:        donID,
		donIDHash:    DonIDToBytes32(donID),
		instanceType: instanceType,
		fromBlock:    fromBlock,
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
	latestConfig, log, err := cp.latestConfig(ctx, int64(cp.fromBlock), math.MaxInt64) // #nosec G115
	if err != nil {
		return 0, ocrtypes.ConfigDigest{}, fmt.Errorf("failed to get latest config: %w", err)
	}
	return uint64(log.BlockNumber), latestConfig.ConfigDigest, nil
}

func (cp *configPoller) latestConfig(ctx context.Context, fromBlock, toBlock int64) (latestConfig FullConfigFromLog, latestLog logpoller.Log, err error) {
	// Get all config set logs run through them forwards
	// TODO: This could probably be optimized with a 'latestBlockNumber' cache or something to avoid reading from `fromBlock` on every call
	// TODO: Actually we only care about the latest of each type here
	// MERC-3524
	logs, err := cp.lp.LogsWithSigs(ctx, fromBlock, toBlock, []common.Hash{ProductionConfigSet, StagingConfigSet}, cp.addr)
	if err != nil {
		return latestConfig, latestLog, fmt.Errorf("failed to get logs: %w", err)
	}
	for _, log := range logs {
		// TODO: This can be optimized probably by adding donIDHash to the logpoller lookup
		// MERC-3524
		if !bytes.Equal(log.Topics[1], cp.donIDHash[:]) {
			continue
		}
		switch log.EventSig {
		case ProductionConfigSet:
			event, err := DecodeProductionConfigSetLog(log.Data)
			if err != nil {
				return latestConfig, log, fmt.Errorf("failed to unpack ProductionConfigSet log data: %w", err)
			}

			if err = cp.cc.StoreConfig(ctx, event.ConfigDigest, event.Signers, event.F); err != nil {
				cp.eng.SugaredLogger.Errorf("failed to store production config: %v", err)
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
				cp.eng.SugaredLogger.Errorf("failed to store staging config: %v", err)
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
	cfg, _, err := cp.latestConfig(ctx, int64(changedInBlock), math.MaxInt64) // #nosec G115
	if err != nil {
		return ocrtypes.ContractConfig{}, fmt.Errorf("failed to get latest config: %w", err)
	}
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
	return uint64(latest.BlockNumber), nil
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
	var transmitAccounts []ocrtypes.Account
	for _, addr := range unpacked.OffchainTransmitters {
		transmitAccounts = append(transmitAccounts, ocrtypes.Account(fmt.Sprintf("%x", addr)))
	}
	var signers []ocrtypes.OnchainPublicKey
	for _, addr := range unpacked.Signers {
		addr := addr
		signers = append(signers, addr[:])
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
	var transmitAccounts []ocrtypes.Account
	for _, addr := range unpacked.OffchainTransmitters {
		transmitAccounts = append(transmitAccounts, ocrtypes.Account(fmt.Sprintf("%x", addr)))
	}
	var signers []ocrtypes.OnchainPublicKey
	for _, addr := range unpacked.Signers {
		addr := addr
		signers = append(signers, addr[:])
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
