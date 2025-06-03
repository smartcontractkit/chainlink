package evm

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-evm/pkg/client"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	evmRelayTypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocrconfigurationstoreevmsimple"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

// configPollerEVMSimple polls config from OCRConfigurationStoreEVMSimple contract
// using logpoller.LogPoller to watch NewConfiguration events and calling ReadConfig.
type configPollerEVMSimple struct {
	services.StateMachine

	lggr                logger.Logger
	filterName          string
	logPoller           logpoller.LogPoller
	configStoreContract *ocrconfigurationstoreevmsimple.OCRConfigurationStoreEVMSimple
	address             common.Address
	logDecoder          LogDecoder
	eventName           string
	abi                 *abi.ABI
}

func configPollerEVMSimpleFilterName(addr common.Address) string {
	return logpoller.FilterName("OCRConfigPollerEVMSimple", addr.String())
}

type ConfigPollerEVMSimpleConfig struct {
	LogPoller  logpoller.LogPoller
	Address    common.Address
	LogDecoder LogDecoder
	Client     client.Client
}

func NewConfigPollerEVMSimple(ctx context.Context, lggr logger.Logger, cfg ConfigPollerEVMSimpleConfig) (evmRelayTypes.ConfigPoller, error) {
	/**

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
		ld:                     ld,
	}

	if configStoreAddr != nil {
		cp.configStoreContractAddr = configStoreAddr
		cp.configStoreContract, err = ocrconfigurationstoreevmsimple.NewOCRConfigurationStoreEVMSimple(*configStoreAddr, client)
		if err != nil {
			return nil, err
		}
	}

	return cp, nil
	*/

	// Register filter for NewConfiguration events
	err := cfg.LogPoller.RegisterFilter(ctx, logpoller.Filter{
		Name:      configPollerEVMSimpleFilterName(cfg.Address),
		EventSigs: []common.Hash{cfg.LogDecoder.EventSig()},
		Addresses: []common.Address{cfg.Address},
	})
	if err != nil {
		return nil, err
	}

	lggr.Infof("TRACE Registered filter for event sig %s for contract %v", cfg.LogDecoder.EventSig(), cfg.Address.Hex())

	const eventName = "NewConfiguration"
	abi, err := ocrconfigurationstoreevmsimple.OCRConfigurationStoreEVMSimpleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	// Create caller for ReadConfig calls
	// caller, err := ocrconfigurationstoreevmsimple.NewOCRConfigurationStoreEVMSimpleCaller(cfg.Address, cfg.Client)
	// if err != nil {
	// 	return nil, err
	// }

	configStoreContract, err := ocrconfigurationstoreevmsimple.NewOCRConfigurationStoreEVMSimple(cfg.Address, cfg.Client)
	if err != nil {
		return nil, err
	}

	return &configPollerEVMSimple{
		lggr:                lggr,
		filterName:          configPollerEVMSimpleFilterName(cfg.Address),
		logPoller:           cfg.LogPoller,
		configStoreContract: configStoreContract,
		address:             cfg.Address,
		logDecoder:          cfg.LogDecoder,
		eventName:           eventName,
		abi:                 abi,
	}, nil
}

func (cp *configPollerEVMSimple) Start() {
	cp.lggr.Infof("Starting config poller for contract %s", cp.address.Hex())
}

func (cp *configPollerEVMSimple) Close() error {
	return nil
}

// Notify noop method - logpoller handles notifications
func (cp *configPollerEVMSimple) Notify() <-chan struct{} {
	return nil
}

func (cp *configPollerEVMSimple) LatestConfig(ctx context.Context, changedInBlock uint64) (ocrtypes.ContractConfig, error) {
	cp.lggr.Infof("TRACE configPollerEVMSimple LatestConfig called for block %d", changedInBlock)

	lgs, err := cp.logPoller.Logs(ctx, int64(changedInBlock), int64(changedInBlock), cp.logDecoder.EventSig(), cp.address)
	if err != nil {
		return ocrtypes.ContractConfig{}, err
	}
	if len(lgs) == 0 {
		return ocrtypes.ContractConfig{}, errors.New("no logs found for config")
	}
	latestConfigSet, err := cp.logDecoder.Decode(lgs[len(lgs)-1].Data)
	if err != nil {
		return ocrtypes.ContractConfig{}, err
	}
	cp.lggr.Infof("TRACE configPollerEVMSimple latestConfigSet %s", latestConfigSet)
	cp.lggr.Infow("LatestConfig", "latestConfig", latestConfigSet)
	return latestConfigSet, nil
}

func (cp *configPollerEVMSimple) LatestBlockHeight(ctx context.Context) (blockHeight uint64, err error) {
	latest, err := cp.logPoller.LatestBlock(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return uint64(latest.BlockNumber), nil
}

func (cp *configPollerEVMSimple) Replay(ctx context.Context, fromBlock int64) error {
	return cp.logPoller.Replay(ctx, fromBlock)
}

func (cp *configPollerEVMSimple) LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest ocrtypes.ConfigDigest, err error) {
	cp.lggr.Infof("TRACE configPollerEVMSimple LatestConfigDetails called")
	latest, err := cp.logPoller.LatestLogByEventSigWithConfs(ctx, cp.logDecoder.EventSig(), cp.address, 1)
	if err != nil {
		cp.lggr.Infof("TRACE configPollerEVMSimple LatestConfigDetails error: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ocrtypes.ConfigDigest{}, err
		}
		return 0, ocrtypes.ConfigDigest{}, err
	}

	cp.lggr.Infof("TRACE configPollerEVMSimple LatestConfigDetails got log: %v", latest)
	latestConfigDigest, err := cp.decode(latest.ToGethLog())
	if err != nil {
		cp.lggr.Infof("TRACE configPollerEVMSimple LatestConfigDetails log decode error: %v", err)
		return 0, ocrtypes.ConfigDigest{}, err
	}
	cp.lggr.Infof("TRACE configPollerEVMSimple LatestConfigDetails returning configDigest: %v", latestConfigDigest)
	return uint64(latest.BlockNumber), latestConfigDigest, nil
}

func (cp *configPollerEVMSimple) decode(log types.Log) (ocrtypes.ConfigDigest, error) {
	cp.lggr.Infof("TRACE Decoding log on contract %s", cp.address.Hex())

	// Unpack the non-indexed data from logEvent.Data
	unpacked := new(ocrconfigurationstoreevmsimple.OCRConfigurationStoreEVMSimpleNewConfiguration)
	if err := cp.abi.UnpackIntoInterface(unpacked, "NewConfiguration", log.Data); err != nil {
		cp.lggr.Errorf("TRACE Failed to unpack log for event %s on contract %s: %v", cp.eventName, cp.address.Hex(), err)
		return ocrtypes.ConfigDigest{}, fmt.Errorf("failed to unpack log data: %w", err)
	}

	// Pick up the indexed fields from the log topics.
	var indexed abi.Arguments
	for _, arg := range cp.abi.Events[cp.eventName].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(unpacked, indexed, log.Topics[1:]); err != nil {
		cp.lggr.Errorf("TRACE Failed to parse indexed topics for event %s on contract %s: %v", cp.eventName, cp.address.Hex(), err)
		return ocrtypes.ConfigDigest{}, fmt.Errorf("failed to parse topics: %w", err)
	}

	if unpacked.ConfigDigest == (common.Hash{}) {
		cp.lggr.Errorf("TRACE ConfigDigest is empty for event %s on contract %s, %v", cp.eventName, cp.address.Hex(), unpacked)
		return ocrtypes.ConfigDigest{}, fmt.Errorf("config digest is empty for event %s on contract %s", cp.eventName, cp.address.Hex())
	}

	cp.lggr.Infof("TRACE Successfully decoded log for event %s on contract %s", cp.eventName, cp.address.Hex())
	return unpacked.ConfigDigest, nil
}

/**
latest, err := cp.destChainLogPoller.LatestLogByEventSigWithConfs(ctx, cp.ld.EventSig(), cp.aggregatorContractAddr, 1)
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
latestConfigSet, err := cp.ld.Decode(latest.Data)
if err != nil {
	return 0, ocrtypes.ConfigDigest{}, err
}
return uint64(latest.BlockNumber), latestConfigSet.ConfigDigest, nil
*/
