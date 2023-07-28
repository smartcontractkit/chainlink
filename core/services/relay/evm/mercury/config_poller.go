package mercury

import (
	"context"
	"database/sql"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/mercury_verifier"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// FeedScopedConfigSet ConfigSet with FeedID for use with mercury (and multi-config DON)
var FeedScopedConfigSet common.Hash

var verifierABI abi.ABI

const (
	configSetEventName = "ConfigSet"
	feedIdTopicIndex   = 1
)

func init() {
	var err error
	verifierABI, err = abi.JSON(strings.NewReader(mercury_verifier.MercuryVerifierABI))
	if err != nil {
		panic(err)
	}
	FeedScopedConfigSet = verifierABI.Events[configSetEventName].ID
}

// FullConfigFromLog defines the contract config with the feedID
type FullConfigFromLog struct {
	ocrtypes.ContractConfig
	feedID [32]byte
}

func unpackLogData(d []byte) (*mercury_verifier.MercuryVerifierConfigSet, error) {
	unpacked := new(mercury_verifier.MercuryVerifierConfigSet)

	err := verifierABI.UnpackIntoInterface(unpacked, configSetEventName, d)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unpack log data")
	}

	return unpacked, nil
}

func configFromLog(logData []byte) (FullConfigFromLog, error) {
	unpacked, err := unpackLogData(logData)
	if err != nil {
		return FullConfigFromLog{}, err
	}

	var transmitAccounts []ocrtypes.Account
	for _, addr := range unpacked.OffchainTransmitters {
		transmitAccounts = append(transmitAccounts, ocrtypes.Account(fmt.Sprintf("%x", addr[:])))
	}
	var signers []ocrtypes.OnchainPublicKey
	for _, addr := range unpacked.Signers {
		addr := addr
		signers = append(signers, addr[:])
	}

	return FullConfigFromLog{
		feedID: unpacked.FeedId,
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

type ContractCaller interface {
	CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
	ConfiguredChainID() *big.Int
}

// ConfigPoller defines the Mercury Config Poller
type configPoller struct {
	utils.StartStopOnce

	lggr               logger.Logger
	destChainLogPoller logpoller.LogPoller
	addr               common.Address
	feedId             common.Hash
	notifyCh           chan struct{}
	subscription       pg.Subscription

	client        types.ContractCaller
	contract      *mercury_verifier.MercuryVerifier
	persistConfig atomic.Bool
	wg            sync.WaitGroup
	chDone        utils.StopChan

	failedRPCContractCalls prometheus.Counter
}

func FilterName(addr common.Address, feedID common.Hash) string {
	return logpoller.FilterName("OCR3 Mercury ConfigPoller", addr.String(), feedID.Hex())
}

// NewConfigPoller creates a new Mercury ConfigPoller
func NewConfigPoller(lggr logger.Logger, client client.Client, destChainPoller logpoller.LogPoller, addr common.Address, feedId common.Hash, eventBroadcaster pg.EventBroadcaster) (*configPoller, error) {
	err := destChainPoller.RegisterFilter(logpoller.Filter{Name: FilterName(addr, feedId), EventSigs: []common.Hash{FeedScopedConfigSet}, Addresses: []common.Address{addr}})
	if err != nil {
		return nil, err
	}

	subscription, err := eventBroadcaster.Subscribe(pg.ChannelInsertOnEVMLogs, "")
	if err != nil {
		return nil, err
	}

	contract, err := mercury_verifier.NewMercuryVerifier(addr, client)
	if err != nil {
		return nil, err
	}

	cp := &configPoller{
		lggr:                   lggr.With("addr", addr.Hex(), "feedID", feedId.Hex()),
		destChainLogPoller:     destChainPoller,
		addr:                   addr,
		feedId:                 feedId,
		notifyCh:               make(chan struct{}, 1),
		subscription:           subscription,
		client:                 client,
		contract:               contract,
		chDone:                 make(chan struct{}),
		failedRPCContractCalls: types.FailedRPCContractCalls.WithLabelValues(client.ConfiguredChainID().String(), addr.Hex(), feedId.Hex()),
	}

	return cp, nil
}

// Start the subscription to Postgres' notify events.
func (cp *configPoller) Start() {
	err := cp.StartOnce("MercuryConfigPoller", func() error {
		cp.wg.Add(2)
		go cp.startLogSubscription()
		go cp.enablePersistConfig()
		return nil
	})
	if err != nil {
		panic(err)
	}
}

// Close the subscription to Postgres' notify events.
func (cp *configPoller) Close() error {
	return cp.StopOnce("MercuryConfigPoller", func() error {
		close(cp.chDone)
		cp.subscription.Close()
		cp.wg.Wait()
		return nil
	})
}

// Notify abstracts the logpoller.LogPoller Notify() implementation
func (cp *configPoller) Notify() <-chan struct{} {
	return cp.notifyCh
}

// Replay abstracts the logpoller.LogPoller Replay() implementation
func (cp *configPoller) Replay(ctx context.Context, fromBlock int64) error {
	return cp.destChainLogPoller.Replay(ctx, fromBlock)
}

// LatestConfigDetails returns the latest config details from the logs
func (cp *configPoller) LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest ocrtypes.ConfigDigest, err error) {
	cp.lggr.Tracew("LatestConfigDetails", "eventSig", FeedScopedConfigSet, "topicIndex", feedIdTopicIndex)
	logs, err := cp.destChainLogPoller.IndexedLogs(FeedScopedConfigSet, cp.addr, feedIdTopicIndex, []common.Hash{cp.feedId}, 1, pg.WithParentCtx(ctx))
	if err != nil {
		return 0, ocrtypes.ConfigDigest{}, err
	}
	if len(logs) == 0 {
		if cp.persistConfig.Load() {
			// Fallback to RPC call in case logs have been pruned
			return cp.callLatestConfigDetails(ctx)
		}
		return 0, ocrtypes.ConfigDigest{}, nil
	}
	latest := logs[len(logs)-1]
	latestConfigSet, err := configFromLog(latest.Data)
	if err != nil {
		return 0, ocrtypes.ConfigDigest{}, err
	}
	return uint64(latest.BlockNumber), latestConfigSet.ConfigDigest, nil
}

// LatestConfig returns the latest config from the logs on a certain block
func (cp *configPoller) LatestConfig(ctx context.Context, changedInBlock uint64) (ocrtypes.ContractConfig, error) {
	cp.lggr.Tracew("LatestConfig", "changedInBlock", changedInBlock, "eventSig", FeedScopedConfigSet, "topicIndex", feedIdTopicIndex)
	lgs, err := cp.destChainLogPoller.IndexedLogsByBlockRange(int64(changedInBlock), int64(changedInBlock), FeedScopedConfigSet, cp.addr, feedIdTopicIndex, []common.Hash{cp.feedId}, pg.WithParentCtx(ctx))
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
	return latestConfigSet.ContractConfig, nil
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

func (cp *configPoller) startLogSubscription() {
	defer cp.wg.Done()

	feedIdPgHex := cp.feedId.Hex()[2:] // trim the leading 0x to make it comparable to pg's hex encoding.
	for {
		event, ok := <-cp.subscription.Events()
		if !ok {
			return
		}

		// Event payload should look like: "<address>:<topicVal1>,<topicVal2>"
		addressTopicValues := strings.Split(event.Payload, ":")
		if len(addressTopicValues) < 2 {
			cp.lggr.Warnf("invalid event from %s channel: %s", pg.ChannelInsertOnEVMLogs, event.Payload)
			continue
		}

		topicValues := strings.Split(addressTopicValues[1], ",")
		if len(topicValues) <= feedIdTopicIndex {
			continue
		}
		if topicValues[feedIdTopicIndex] != feedIdPgHex {
			continue
		}

		select {
		case cp.notifyCh <- struct{}{}:
		default:
		}
	}
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
	}, cp.feedId)
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
	}, cp.feedId)
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
		transmitters[i] = ocrtypes.Account(fmt.Sprintf("%x", ocr2AbstractConfig.Transmitters[i][:]))
	}
	return uint64(ocr2AbstractConfig.PreviousConfigBlockNumber), ocrtypes.ContractConfig{
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
