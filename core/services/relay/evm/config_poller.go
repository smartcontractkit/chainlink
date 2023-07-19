package evm

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
	evmRelayTypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var (
	// ConfigSet Common to all OCR2 evm based contracts: https://github.com/smartcontractkit/libocr/blob/master/contract2/dev/OCR2Abstract.sol
	ConfigSet common.Hash

	defaultABI abi.ABI
)

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

func configFromLog(logData []byte) (ocrtypes.ContractConfig, error) {
	unpacked, err := unpackLogData(logData)
	if err != nil {
		return ocrtypes.ContractConfig{}, err
	}

	var transmitAccounts []ocrtypes.Account
	for _, addr := range unpacked.Transmitters {
		transmitAccounts = append(transmitAccounts, ocrtypes.Account(addr.Hex()))
	}
	var signers []ocrtypes.OnchainPublicKey
	for _, addr := range unpacked.Signers {
		addr := addr
		signers = append(signers, addr[:])
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

type configPoller struct {
	utils.StartStopOnce

	lggr               logger.Logger
	filterName         string
	destChainLogPoller logpoller.LogPoller
	addr               common.Address

	client        types.ContractCaller
	contract      *ocr2aggregator.OCR2Aggregator
	persistConfig atomic.Bool
	wg            sync.WaitGroup
	chDone        utils.StopChan

	failedRPCContractCalls prometheus.Counter
}

func configPollerFilterName(addr common.Address) string {
	return logpoller.FilterName("OCR2ConfigPoller", addr.String())
}

func NewConfigPoller(lggr logger.Logger, client client.Client, destChainPoller logpoller.LogPoller, addr common.Address) (evmRelayTypes.ConfigPoller, error) {
	return newConfigPoller(lggr, client, destChainPoller, addr)
}

func newConfigPoller(lggr logger.Logger, client client.Client, destChainPoller logpoller.LogPoller, addr common.Address) (*configPoller, error) {
	err := destChainPoller.RegisterFilter(logpoller.Filter{Name: configPollerFilterName(addr), EventSigs: []common.Hash{ConfigSet}, Addresses: []common.Address{addr}})
	if err != nil {
		return nil, err
	}

	contract, err := ocr2aggregator.NewOCR2Aggregator(addr, client)
	if err != nil {
		return nil, err
	}

	cp := &configPoller{
		lggr:                   lggr,
		filterName:             configPollerFilterName(addr),
		destChainLogPoller:     destChainPoller,
		addr:                   addr,
		client:                 client,
		contract:               contract,
		chDone:                 make(chan struct{}),
		failedRPCContractCalls: types.FailedRPCContractCalls.WithLabelValues(client.ConfiguredChainID().String(), addr.Hex(), ""),
	}

	return cp, nil
}

func (cp *configPoller) Start() {
	err := cp.StartOnce("OCR2ConfigPoller", func() error {
		cp.wg.Add(1)
		go cp.enablePersistConfig()
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func (cp *configPoller) Close() error {
	return cp.StopOnce("OCR2ConfigPoller", func() error {
		close(cp.chDone)
		cp.wg.Wait()
		return nil
	})
}

// Notify noop method
func (cp *configPoller) Notify() <-chan struct{} {
	return nil
}

// Replay abstracts the logpoller.LogPoller Replay() implementation
func (cp *configPoller) Replay(ctx context.Context, fromBlock int64) error {
	return cp.destChainLogPoller.Replay(ctx, fromBlock)
}

// LatestConfigDetails returns the latest config details from the logs
func (cp *configPoller) LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest ocrtypes.ConfigDigest, err error) {
	latest, err := cp.destChainLogPoller.LatestLogByEventSigWithConfs(ConfigSet, cp.addr, 1, pg.WithParentCtx(ctx))
	if err != nil {
		// If contract is not configured, or logs have been pruned, we will not have the log.
		if errors.Is(err, sql.ErrNoRows) {
			if cp.persistConfig.Load() {
				// Fallback to RPC call in case logs have been pruned
				return cp.callLatestConfigDetails(ctx)
			}
			return 0, ocrtypes.ConfigDigest{}, nil
		}
		return 0, ocrtypes.ConfigDigest{}, err
	}
	latestConfigSet, err := configFromLog(latest.Data)
	if err != nil {
		return 0, ocrtypes.ConfigDigest{}, err
	}
	return uint64(latest.BlockNumber), latestConfigSet.ConfigDigest, nil
}

// LatestConfig returns the latest config from the logs on a certain block
func (cp *configPoller) LatestConfig(ctx context.Context, changedInBlock uint64) (ocrtypes.ContractConfig, error) {
	lgs, err := cp.destChainLogPoller.Logs(int64(changedInBlock), int64(changedInBlock), ConfigSet, cp.addr, pg.WithParentCtx(ctx))
	if err != nil {
		return ocrtypes.ContractConfig{}, err
	}
	if len(lgs) == 0 {
		if cp.persistConfig.Load() {
			minedInBlock, cfg, err := cp.callLatestConfig(ctx)
			if err != nil {
				return cfg, err
			}
			if cfg.ConfigDigest != (ocrtypes.ConfigDigest{}) && changedInBlock != minedInBlock {
				return cfg, fmt.Errorf("block number mismatch: expected to find config changed in block %d but the config was changed in block %d", changedInBlock, minedInBlock)
			}
			return cfg, err
		}
		return ocrtypes.ContractConfig{}, fmt.Errorf("missing config on contract %s (chain %s) at block %d", cp.addr.Hex(), cp.client.ConfiguredChainID().String(), changedInBlock)
	}
	latestConfigSet, err := configFromLog(lgs[len(lgs)-1].Data)
	if err != nil {
		return ocrtypes.ContractConfig{}, err
	}
	cp.lggr.Infow("LatestConfig", "latestConfig", latestConfigSet)
	return latestConfigSet, nil
}

// LatestBlockHeight returns the latest block height from the logs
func (cp *configPoller) LatestBlockHeight(ctx context.Context) (blockHeight uint64, err error) {
	latest, err := cp.destChainLogPoller.LatestBlock(pg.WithParentCtx(ctx))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return uint64(latest), nil
}

// enablePersistConfig runs in parallel so that we can attempt to use logs for config even if RPC calls are failing
func (cp *configPoller) enablePersistConfig() {
	defer cp.wg.Done()
	ctx, cancel := cp.chDone.Ctx(context.Background())
	defer cancel()
	b := types.NewRPCCallBackoff()
	for {
		enabled, err := cp.callIsConfigPersisted(ctx)
		if err == nil {
			cp.persistConfig.Store(enabled)
			return
		} else {
			cp.lggr.Warnw("Failed to determine whether config persistence is enabled, retrying", "err", err)
		}
		select {
		case <-time.After(b.Duration()):
			// keep trying for as long as it takes, with exponential backoff
		case <-cp.chDone:
			return
		}
	}
}

func (cp *configPoller) callIsConfigPersisted(ctx context.Context) (persistConfig bool, err error) {
	persistConfig, err = cp.contract.PersistConfig(&bind.CallOpts{Context: ctx})
	if err != nil {
		if methodNotImplemented(err) {
			return false, nil
		}
		cp.failedRPCContractCalls.Inc()
		return
	}
	return persistConfig, nil
}

func methodNotImplemented(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "execution reverted")
}

func (cp *configPoller) callLatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest ocrtypes.ConfigDigest, err error) {
	details, err := cp.contract.LatestConfigDetails(&bind.CallOpts{
		Context: ctx,
	})
	if err != nil {
		cp.failedRPCContractCalls.Inc()
	}
	return uint64(details.BlockNumber), details.ConfigDigest, err
}

// Some chains "manage" state bloat by deleting older logs. This RPC call
// allows us work around such restrictions.
func (cp *configPoller) callLatestConfig(ctx context.Context) (changedInBlock uint64, cfg ocrtypes.ContractConfig, err error) {
	ocr2AbstractConfig, err := cp.contract.LatestConfig(&bind.CallOpts{
		Context: ctx,
	})
	if err != nil {
		cp.failedRPCContractCalls.Inc()
		return
	}
	signers := make([]ocrtypes.OnchainPublicKey, len(ocr2AbstractConfig.Signers))
	for i := range signers {
		signers[i] = ocr2AbstractConfig.Signers[i].Bytes()
	}
	transmitters := make([]ocrtypes.Account, len(ocr2AbstractConfig.Transmitters))
	for i := range transmitters {
		transmitters[i] = ocrtypes.Account(ocr2AbstractConfig.Transmitters[i].Hex())
	}
	return uint64(ocr2AbstractConfig.CurrentConfigBlockNumber), ocrtypes.ContractConfig{
		ConfigDigest:          ocr2AbstractConfig.ConfigDigest,
		ConfigCount:           ocr2AbstractConfig.ConfigCount,
		Signers:               signers,
		Transmitters:          transmitters,
		F:                     ocr2AbstractConfig.F,
		OnchainConfig:         ocr2AbstractConfig.OnchainConfig,
		OffchainConfigVersion: ocr2AbstractConfig.OffchainConfigVersion,
		OffchainConfig:        ocr2AbstractConfig.OffchainConfig,
	}, err
}
