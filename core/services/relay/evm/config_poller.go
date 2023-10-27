package evm

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocrconfigurationstoreevmsimple"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	evmRelayTypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var (
	failedRPCContractCalls = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "ocr2_failed_rpc_contract_calls",
		Help: "Running count of failed RPC contract calls by chain/contract",
	},
		[]string{"chainID", "contractAddress"},
	)
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
	client             client.Client

	aggregatorContractAddr common.Address
	aggregatorContract     *ocr2aggregator.OCR2Aggregator

	// Some chains "manage" state bloat by deleting older logs. The ConfigStore
	// contract allows us work around such restrictions.
	configStoreContractAddr *common.Address
	configStoreContract     *ocrconfigurationstoreevmsimple.OCRConfigurationStoreEVMSimple
}

func configPollerFilterName(addr common.Address) string {
	return logpoller.FilterName("OCR2ConfigPoller", addr.String())
}

func NewConfigPoller(lggr logger.Logger, client client.Client, destChainPoller logpoller.LogPoller, aggregatorContractAddr common.Address, configStoreAddr *common.Address) (evmRelayTypes.ConfigPoller, error) {
	return newConfigPoller(lggr, client, destChainPoller, aggregatorContractAddr, configStoreAddr)
}

func newConfigPoller(lggr logger.Logger, client client.Client, destChainPoller logpoller.LogPoller, aggregatorContractAddr common.Address, configStoreAddr *common.Address) (*configPoller, error) {
	err := destChainPoller.RegisterFilter(logpoller.Filter{Name: configPollerFilterName(aggregatorContractAddr), EventSigs: []common.Hash{ConfigSet}, Addresses: []common.Address{aggregatorContractAddr}})
	if err != nil {
		return nil, err
	}

	aggregatorContract, err := ocr2aggregator.NewOCR2Aggregator(aggregatorContractAddr, client)
	if err != nil {
		return nil, err
	}

	cp := &configPoller{
		lggr:                   lggr,
		filterName:             configPollerFilterName(aggregatorContractAddr),
		destChainLogPoller:     destChainPoller,
		aggregatorContractAddr: aggregatorContractAddr,
		client:                 client,
		aggregatorContract:     aggregatorContract,
	}

	if configStoreAddr != nil {
		cp.configStoreContractAddr = configStoreAddr
		cp.configStoreContract, err = ocrconfigurationstoreevmsimple.NewOCRConfigurationStoreEVMSimple(*configStoreAddr, client)
		if err != nil {
			return nil, err
		}
	}

	return cp, nil
}

func (cp *configPoller) Start() {}

func (cp *configPoller) Close() error {
	return nil
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
	latest, err := cp.destChainLogPoller.LatestLogByEventSigWithConfs(ConfigSet, cp.aggregatorContractAddr, 1, pg.WithParentCtx(ctx))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if cp.isConfigStoreAvailable() {
				// Fallback to RPC call in case logs have been pruned and configStoreContract is available
				return cp.callLatestConfigDetails(ctx)
			}
			// log not found means return zero config digest
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
	lgs, err := cp.destChainLogPoller.Logs(int64(changedInBlock), int64(changedInBlock), ConfigSet, cp.aggregatorContractAddr, pg.WithParentCtx(ctx))
	if err != nil {
		return ocrtypes.ContractConfig{}, err
	}
	if len(lgs) == 0 {
		if cp.isConfigStoreAvailable() {
			// Fallback to RPC call in case logs have been pruned
			return cp.callReadConfigFromStore(ctx)
		}
		return ocrtypes.ContractConfig{}, fmt.Errorf("no logs found for config on contract %s (chain %s) at block %d", cp.aggregatorContractAddr.Hex(), cp.client.ConfiguredChainID().String(), changedInBlock)
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
	return uint64(latest.BlockNumber), nil
}

func (cp *configPoller) isConfigStoreAvailable() bool {
	return cp.configStoreContract != nil
}

// RPC call for latest config details
func (cp *configPoller) callLatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest ocrtypes.ConfigDigest, err error) {
	details, err := cp.aggregatorContract.LatestConfigDetails(&bind.CallOpts{
		Context: ctx,
	})
	if err != nil {
		failedRPCContractCalls.WithLabelValues(cp.client.ConfiguredChainID().String(), cp.aggregatorContractAddr.Hex()).Inc()
	}
	return uint64(details.BlockNumber), details.ConfigDigest, err
}

// RPC call to read config from config store contract
func (cp *configPoller) callReadConfigFromStore(ctx context.Context) (cfg ocrtypes.ContractConfig, err error) {
	_, configDigest, err := cp.LatestConfigDetails(ctx)
	if err != nil {
		failedRPCContractCalls.WithLabelValues(cp.client.ConfiguredChainID().String(), cp.aggregatorContractAddr.Hex()).Inc()
		return cfg, fmt.Errorf("failed to get latest config details: %w", err)
	}
	if configDigest == (ocrtypes.ConfigDigest{}) {
		return cfg, fmt.Errorf("config details missing while trying to lookup config in store; no logs found for contract %s (chain %s)", cp.aggregatorContractAddr.Hex(), cp.client.ConfiguredChainID().String())
	}

	storedConfig, err := cp.configStoreContract.ReadConfig(&bind.CallOpts{
		Context: ctx,
	}, configDigest)
	if err != nil {
		failedRPCContractCalls.WithLabelValues(cp.client.ConfiguredChainID().String(), cp.configStoreContractAddr.Hex()).Inc()
		return cfg, fmt.Errorf("failed to read config from config store contract: %w", err)
	}

	signers := make([]ocrtypes.OnchainPublicKey, len(storedConfig.Signers))
	for i := range signers {
		signers[i] = storedConfig.Signers[i].Bytes()
	}
	transmitters := make([]ocrtypes.Account, len(storedConfig.Transmitters))
	for i := range transmitters {
		transmitters[i] = ocrtypes.Account(storedConfig.Transmitters[i].Hex())
	}

	return ocrtypes.ContractConfig{
		ConfigDigest:          configDigest,
		ConfigCount:           uint64(storedConfig.ConfigCount),
		Signers:               signers,
		Transmitters:          transmitters,
		F:                     storedConfig.F,
		OnchainConfig:         storedConfig.OnchainConfig,
		OffchainConfigVersion: storedConfig.OffchainConfigVersion,
		OffchainConfig:        storedConfig.OffchainConfig,
	}, err
}
