package evm

import (
	"context"
	goerrors "errors"
	"fmt"
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	coreTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"go.uber.org/multierr"

	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated"
	iregistry21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/models"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/core"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/encoding"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/logprovider"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const (
	// defaultAllowListExpiration decides how long an upkeep's allow list info will be valid for.
	defaultAllowListExpiration = 20 * time.Minute
	// allowListCleanupInterval decides when the expired items in allowList cache will be deleted.
	allowListCleanupInterval   = 5 * time.Minute
	logTriggerRefreshBatchSize = 32
)

var (
	RegistryServiceName = "AutomationRegistry"

	ErrLogReadFailure              = fmt.Errorf("failure reading logs")
	ErrHeadNotAvailable            = fmt.Errorf("head not available")
	ErrInitializationFailure       = fmt.Errorf("failed to initialize registry")
	ErrContextCancelled            = fmt.Errorf("context was cancelled")
	ErrABINotParsable              = fmt.Errorf("error parsing abi")
	ActiveUpkeepIDBatchSize  int64 = 1000
	// This is the interval at which active upkeep list is fully refreshed from chain
	refreshInterval = 15 * time.Minute
	// This is the lookback for polling upkeep state event logs from latest block
	logEventLookback int64 = 250
)

//go:generate mockery --quiet --name Registry --output ./mocks/ --case=underscore
type Registry interface {
	GetUpkeep(opts *bind.CallOpts, id *big.Int) (encoding.UpkeepInfo, error)
	GetState(opts *bind.CallOpts) (iregistry21.GetState, error)
	GetActiveUpkeepIDs(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error)
	GetUpkeepPrivilegeConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error)
	GetUpkeepTriggerConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error)
	CheckCallback(opts *bind.CallOpts, id *big.Int, values [][]byte, extraData []byte) (iregistry21.CheckCallback, error)
	ParseLog(log coreTypes.Log) (generated.AbigenLog, error)
}

//go:generate mockery --quiet --name HttpClient --output ./mocks/ --case=underscore
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewEvmRegistry(
	lggr logger.Logger,
	addr common.Address,
	client evm.Chain,
	registry *iregistry21.IKeeperRegistryMaster,
	mc *models.MercuryCredentials,
	al ActiveUpkeepList,
	logEventProvider logprovider.LogEventProvider,
	packer encoding.Packer,
	blockSub *BlockSubscriber,
	finalityDepth uint32,
) *EvmRegistry {
	return &EvmRegistry{
		ctx:          context.Background(),
		threadCtrl:   utils.NewThreadControl(),
		lggr:         lggr.Named(RegistryServiceName),
		poller:       client.LogPoller(),
		addr:         addr,
		client:       client.Client(),
		logProcessed: make(map[string]bool),
		registry:     registry,
		abi:          core.RegistryABI,
		active:       al,
		packer:       packer,
		headFunc:     func(ocr2keepers.BlockKey) {},
		chLog:        make(chan logpoller.Log, 1000),
		mercury: &MercuryConfig{
			cred:           mc,
			abi:            core.StreamsCompatibleABI,
			allowListCache: cache.New(defaultAllowListExpiration, allowListCleanupInterval),
		},
		hc:               http.DefaultClient,
		logEventProvider: logEventProvider,
		bs:               blockSub,
		finalityDepth:    finalityDepth,
	}
}

var upkeepStateEvents = []common.Hash{
	iregistry21.IKeeperRegistryMasterUpkeepRegistered{}.Topic(),       // adds new upkeep id to registry
	iregistry21.IKeeperRegistryMasterUpkeepReceived{}.Topic(),         // adds new upkeep id to registry via migration
	iregistry21.IKeeperRegistryMasterUpkeepUnpaused{}.Topic(),         // unpauses an upkeep
	iregistry21.IKeeperRegistryMasterUpkeepPaused{}.Topic(),           // pauses an upkeep
	iregistry21.IKeeperRegistryMasterUpkeepMigrated{}.Topic(),         // migrated an upkeep, equivalent to cancel from this registry's perspective
	iregistry21.IKeeperRegistryMasterUpkeepCanceled{}.Topic(),         // cancels an upkeep
	iregistry21.IKeeperRegistryMasterUpkeepTriggerConfigSet{}.Topic(), // trigger config was changed
}

type MercuryConfig struct {
	cred *models.MercuryCredentials
	abi  abi.ABI
	// allowListCache stores the upkeeps privileges. In 2.1, this only includes a JSON bytes for allowed to use mercury
	allowListCache *cache.Cache
}

type EvmRegistry struct {
	utils.StartStopOnce
	threadCtrl       utils.ThreadControl
	lggr             logger.Logger
	poller           logpoller.LogPoller
	addr             common.Address
	client           client.Client
	chainID          uint64
	registry         Registry
	abi              abi.ABI
	packer           encoding.Packer
	chLog            chan logpoller.Log
	mu               sync.RWMutex
	logProcessed     map[string]bool
	active           ActiveUpkeepList
	lastPollBlock    int64
	ctx              context.Context
	headFunc         func(ocr2keepers.BlockKey)
	mercury          *MercuryConfig
	hc               HttpClient
	bs               *BlockSubscriber
	logEventProvider logprovider.LogEventProvider
	finalityDepth    uint32
}

func (r *EvmRegistry) Name() string {
	return r.lggr.Name()
}

func (r *EvmRegistry) Start(ctx context.Context) error {
	return r.StartOnce(RegistryServiceName, func() error {
		if err := r.registerEvents(r.chainID, r.addr); err != nil {
			return fmt.Errorf("logPoller error while registering automation events: %w", err)
		}

		r.threadCtrl.Go(func(ctx context.Context) {
			lggr := r.lggr.With("where", "upkeeps_referesh")
			err := r.refreshActiveUpkeeps()
			if err != nil {
				lggr.Errorf("failed to initialize upkeeps", err)
			}

			ticker := time.NewTicker(refreshInterval)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					err = r.refreshActiveUpkeeps()
					if err != nil {
						lggr.Errorf("failed to refresh upkeeps", err)
					}
				case <-ctx.Done():
					return
				}
			}
		})

		r.threadCtrl.Go(func(ctx context.Context) {
			lggr := r.lggr.With("where", "logs_polling")
			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					err := r.pollUpkeepStateLogs()
					if err != nil {
						lggr.Errorf("failed to poll logs for upkeeps", err)
					}
				case <-ctx.Done():
					return
				}
			}
		})

		r.threadCtrl.Go(func(ctx context.Context) {
			lggr := r.lggr.With("where", "logs_processing")
			ch := r.chLog

			for {
				select {
				case l := <-ch:
					err := r.processUpkeepStateLog(l)
					if err != nil {
						lggr.Errorf("failed to process log for upkeep", err)
					}
				case <-ctx.Done():
					return
				}
			}
		})

		return nil
	})
}

func (r *EvmRegistry) Close() error {
	return r.StopOnce(RegistryServiceName, func() error {
		r.threadCtrl.Close()
		return nil
	})
}

func (r *EvmRegistry) HealthReport() map[string]error {
	return map[string]error{RegistryServiceName: r.Healthy()}
}

func (r *EvmRegistry) refreshActiveUpkeeps() error {
	// Allow for max timeout of refreshInterval
	ctx, cancel := context.WithTimeout(r.ctx, refreshInterval)
	defer cancel()

	r.lggr.Debugf("Refreshing active upkeeps list")
	// get active upkeep ids from contract
	ids, err := r.getLatestIDsFromContract(ctx)
	if err != nil {
		return fmt.Errorf("failed to get active upkeep ids from contract during refresh: %s", err)
	}
	r.active.Reset(ids...)

	var logTriggerIDs []*big.Int
	for _, id := range ids {
		uid := &ocr2keepers.UpkeepIdentifier{}
		if ok := uid.FromBigInt(id); !ok {
			r.lggr.Warnf("failed to parse upkeep id %s", id.String())
			continue
		}
		switch core.GetUpkeepType(*uid) {
		case ocr2keepers.LogTrigger:
			logTriggerIDs = append(logTriggerIDs, id)
		default:
		}
	}

	_, err = r.logEventProvider.RefreshActiveUpkeeps(logTriggerIDs...)
	if err != nil {
		return fmt.Errorf("failed to refresh active upkeep ids in log event provider: %w", err)
	}

	// Try to refersh log trigger config for all log upkeeps
	return r.refreshLogTriggerUpkeeps(logTriggerIDs)
}

// refreshLogTriggerUpkeeps refreshes the active upkeep ids for log trigger upkeeps
func (r *EvmRegistry) refreshLogTriggerUpkeeps(ids []*big.Int) error {
	var err error
	for i := 0; i < len(ids); i += logTriggerRefreshBatchSize {
		end := i + logTriggerRefreshBatchSize
		if end > len(ids) {
			end = len(ids)
		}
		idBatch := ids[i:end]

		if batchErr := r.refreshLogTriggerUpkeepsBatch(idBatch); batchErr != nil {
			multierr.AppendInto(&err, batchErr)
		}

		time.Sleep(500 * time.Millisecond)
	}

	return err
}

func (r *EvmRegistry) refreshLogTriggerUpkeepsBatch(logTriggerIDs []*big.Int) error {
	var logTriggerHashes []common.Hash
	for _, id := range logTriggerIDs {
		logTriggerHashes = append(logTriggerHashes, common.BigToHash(id))
	}

	unpausedLogs, err := r.poller.IndexedLogs(iregistry21.IKeeperRegistryMasterUpkeepUnpaused{}.Topic(), r.addr, 1, logTriggerHashes, logpoller.Confirmations(r.finalityDepth), pg.WithParentCtx(r.ctx))
	if err != nil {
		return err
	}
	configSetLogs, err := r.poller.IndexedLogs(iregistry21.IKeeperRegistryMasterUpkeepTriggerConfigSet{}.Topic(), r.addr, 1, logTriggerHashes, logpoller.Confirmations(r.finalityDepth), pg.WithParentCtx(r.ctx))
	if err != nil {
		return err
	}

	logs := append(unpausedLogs, configSetLogs...)

	configSetBlockNumbers := map[string]uint64{}
	unpausedBlockNumbers := map[string]uint64{}
	perUpkeepConfig := map[string][]byte{}

	for _, log := range logs {
		rawLog := log.ToGethLog()
		abilog, err := r.registry.ParseLog(rawLog)
		if err != nil {
			return err
		}
		switch l := abilog.(type) {
		case *iregistry21.IKeeperRegistryMasterUpkeepTriggerConfigSet:
			if rawLog.BlockNumber > configSetBlockNumbers[l.Id.String()] {
				configSetBlockNumbers[l.Id.String()] = rawLog.BlockNumber
				perUpkeepConfig[l.Id.String()] = l.TriggerConfig
			}
		case *iregistry21.IKeeperRegistryMasterUpkeepUnpaused:
			if rawLog.BlockNumber > unpausedBlockNumbers[l.Id.String()] {
				unpausedBlockNumbers[l.Id.String()] = rawLog.BlockNumber
			}
		}
	}

	var merr error
	for _, id := range logTriggerIDs {
		logBlock, ok := configSetBlockNumbers[id.String()]
		if !ok {
			r.lggr.Warnf("unable to find finalized config set block number for %s, using 0 as config start block", id.String())
			// Use zero as config update block so it can be updated if an actual event is found later
			logBlock = 0
		}

		config, ok := perUpkeepConfig[id.String()]
		if !ok {
			r.lggr.Warnf("unable to find per finalized log config for %s, will fetch latest config from chain", id.String())
			// Set it to empty bytes so that latest config is fetched within r.updateTriggerConfig
			config = []byte{}
		}

		// In case an upkeep was paused then unpaused after a config set event, start the config from the unpaused block number
		if unpausedBlockNumbers[id.String()] > logBlock {
			logBlock = unpausedBlockNumbers[id.String()]
		}
		if err := r.updateTriggerConfig(id, config, logBlock); err != nil {
			merr = goerrors.Join(merr, fmt.Errorf("failed to update trigger config for upkeep id %s: %w", id.String(), err))
		}
	}

	return merr
}

func (r *EvmRegistry) pollUpkeepStateLogs() error {
	var latest int64
	var end int64
	var err error

	if end, err = r.poller.LatestBlock(pg.WithParentCtx(r.ctx)); err != nil {
		return fmt.Errorf("%w: %s", ErrHeadNotAvailable, err)
	}

	r.mu.Lock()
	latest = r.lastPollBlock
	r.lastPollBlock = end
	r.mu.Unlock()

	// if start and end are the same, no polling needs to be done
	if latest == 0 || latest == end {
		return nil
	}

	var logs []logpoller.Log
	if logs, err = r.poller.LogsWithSigs(
		end-logEventLookback,
		end,
		upkeepStateEvents,
		r.addr,
		pg.WithParentCtx(r.ctx),
	); err != nil {
		return fmt.Errorf("%w: %s", ErrLogReadFailure, err)
	}

	for _, log := range logs {
		r.chLog <- log
	}

	return nil
}

func (r *EvmRegistry) processUpkeepStateLog(l logpoller.Log) error {
	lid := fmt.Sprintf("%s%d", l.TxHash.String(), l.LogIndex)
	r.mu.Lock()
	if _, ok := r.logProcessed[lid]; ok {
		r.mu.Unlock()
		return nil
	}
	r.logProcessed[lid] = true
	r.mu.Unlock()
	txHash := l.TxHash.String()

	rawLog := l.ToGethLog()
	abilog, err := r.registry.ParseLog(rawLog)
	if err != nil {
		return err
	}

	switch l := abilog.(type) {
	case *iregistry21.IKeeperRegistryMasterUpkeepPaused:
		r.lggr.Debugf("KeeperRegistryUpkeepPaused log detected for upkeep ID %s in transaction %s", l.Id.String(), txHash)
		r.removeFromActive(l.Id)
	case *iregistry21.IKeeperRegistryMasterUpkeepCanceled:
		r.lggr.Debugf("KeeperRegistryUpkeepCanceled log detected for upkeep ID %s in transaction %s", l.Id.String(), txHash)
		r.removeFromActive(l.Id)
	case *iregistry21.IKeeperRegistryMasterUpkeepMigrated:
		r.lggr.Debugf("KeeperRegistryMasterUpkeepMigrated log detected for upkeep ID %s in transaction %s", l.Id.String(), txHash)
		r.removeFromActive(l.Id)
	case *iregistry21.IKeeperRegistryMasterUpkeepTriggerConfigSet:
		r.lggr.Debugf("KeeperRegistryUpkeepTriggerConfigSet log detected for upkeep ID %s in transaction %s", l.Id.String(), txHash)
		if err := r.updateTriggerConfig(l.Id, l.TriggerConfig, rawLog.BlockNumber); err != nil {
			r.lggr.Warnf("failed to update trigger config upon KeeperRegistryMasterUpkeepTriggerConfigSet for upkeep ID %s: %s", l.Id.String(), err)
		}
	case *iregistry21.IKeeperRegistryMasterUpkeepRegistered:
		uid := &ocr2keepers.UpkeepIdentifier{}
		uid.FromBigInt(l.Id)
		trigger := core.GetUpkeepType(*uid)
		r.lggr.Debugf("KeeperRegistryUpkeepRegistered log detected for upkeep ID %s (trigger=%d) in transaction %s", l.Id.String(), trigger, txHash)
		r.active.Add(l.Id)
		if err := r.updateTriggerConfig(l.Id, nil, rawLog.BlockNumber); err != nil {
			r.lggr.Warnf("failed to update trigger config upon KeeperRegistryMasterUpkeepRegistered for upkeep ID %s: %s", err)
		}
	case *iregistry21.IKeeperRegistryMasterUpkeepReceived:
		r.lggr.Debugf("KeeperRegistryUpkeepReceived log detected for upkeep ID %s in transaction %s", l.Id.String(), txHash)
		r.active.Add(l.Id)
		if err := r.updateTriggerConfig(l.Id, nil, rawLog.BlockNumber); err != nil {
			r.lggr.Warnf("failed to update trigger config upon KeeperRegistryMasterUpkeepReceived for upkeep ID %s: %s", err)
		}
	case *iregistry21.IKeeperRegistryMasterUpkeepUnpaused:
		r.lggr.Debugf("KeeperRegistryUpkeepUnpaused log detected for upkeep ID %s in transaction %s", l.Id.String(), txHash)
		r.active.Add(l.Id)
		if err := r.updateTriggerConfig(l.Id, nil, rawLog.BlockNumber); err != nil {
			r.lggr.Warnf("failed to update trigger config upon KeeperRegistryMasterUpkeepUnpaused for upkeep ID %s: %s", err)
		}
	default:
		r.lggr.Debugf("Unknown log detected for log %+v in transaction %s", l, txHash)
	}

	return nil
}

func RegistryUpkeepFilterName(addr common.Address) string {
	return logpoller.FilterName("KeeperRegistry Events", addr.String())
}

// registerEvents registers upkeep state events from keeper registry on log poller
func (r *EvmRegistry) registerEvents(_ uint64, addr common.Address) error {
	// Add log filters for the log poller so that it can poll and find the logs that we need
	return r.poller.RegisterFilter(logpoller.Filter{
		Name:      RegistryUpkeepFilterName(addr),
		EventSigs: upkeepStateEvents,
		Addresses: []common.Address{addr},
	})
}

// removeFromActive removes an upkeepID from active list and unregisters the log filter for log upkeeps
func (r *EvmRegistry) removeFromActive(id *big.Int) {
	r.active.Remove(id)

	uid := &ocr2keepers.UpkeepIdentifier{}
	uid.FromBigInt(id)
	trigger := core.GetUpkeepType(*uid)
	switch trigger {
	case ocr2keepers.LogTrigger:
		if err := r.logEventProvider.UnregisterFilter(id); err != nil {
			r.lggr.Warnw("failed to unregister log filter", "upkeepID", id.String())
		}
		r.lggr.Debugw("unregistered log filter", "upkeepID", id.String())
	default:
	}
}

func (r *EvmRegistry) buildCallOpts(ctx context.Context, block *big.Int) *bind.CallOpts {
	opts := bind.CallOpts{
		Context: ctx,
	}

	if block == nil || block.Int64() == 0 {
		l := r.bs.latestBlock.Load()
		if l != nil && l.Number != 0 {
			opts.BlockNumber = big.NewInt(int64(l.Number))
		}
	} else {
		opts.BlockNumber = block
	}

	return &opts
}

func (r *EvmRegistry) getLatestIDsFromContract(ctx context.Context) ([]*big.Int, error) {
	opts := r.buildCallOpts(ctx, nil)

	state, err := r.registry.GetState(opts)
	if err != nil {
		n := "latest"
		if opts.BlockNumber != nil {
			n = fmt.Sprintf("%d", opts.BlockNumber.Int64())
		}

		return nil, fmt.Errorf("%w: failed to get contract state at block number '%s'", err, n)
	}

	ids := make([]*big.Int, 0, int(state.State.NumUpkeeps.Int64()))
	for int64(len(ids)) < state.State.NumUpkeeps.Int64() {
		startIndex := int64(len(ids))
		maxCount := state.State.NumUpkeeps.Int64() - startIndex

		if maxCount == 0 {
			break
		}

		if maxCount > ActiveUpkeepIDBatchSize {
			maxCount = ActiveUpkeepIDBatchSize
		}

		batchIDs, err := r.registry.GetActiveUpkeepIDs(opts, big.NewInt(startIndex), big.NewInt(maxCount))
		if err != nil {
			return nil, fmt.Errorf("%w: failed to get active upkeep IDs from index %d to %d (both inclusive)", err, startIndex, startIndex+maxCount-1)
		}

		ids = append(ids, batchIDs...)
	}

	return ids, nil
}

// updateTriggerConfig updates the trigger config for an upkeep. it will re-register a filter for this upkeep.
func (r *EvmRegistry) updateTriggerConfig(id *big.Int, cfg []byte, logBlock uint64) error {
	uid := &ocr2keepers.UpkeepIdentifier{}
	uid.FromBigInt(id)
	switch core.GetUpkeepType(*uid) {
	case ocr2keepers.LogTrigger:
		if len(cfg) == 0 {
			fetched, err := r.fetchTriggerConfig(id)
			if err != nil {
				return errors.Wrap(err, "failed to fetch log upkeep config")
			}
			cfg = fetched
		}
		parsed, err := r.packer.UnpackLogTriggerConfig(cfg)
		if err != nil {
			// Upkeep has been setup with improper config. Log a warning and ignore the upkeep.
			r.lggr.Warnw("failed to unpack log upkeep config", "upkeepID", id.String(), "err", err)
			return nil
		}
		if err := r.logEventProvider.RegisterFilter(r.ctx, logprovider.FilterOptions{
			TriggerConfig: logprovider.LogTriggerConfig(parsed),
			UpkeepID:      id,
			UpdateBlock:   logBlock,
		}); err != nil {
			return errors.Wrap(err, "failed to register log filter")
		}
		r.lggr.Debugw("registered log filter", "upkeepID", id.String(), "cfg", parsed)
	default:
	}
	return nil
}

// fetchTriggerConfig fetches trigger config in raw bytes for an upkeep.
func (r *EvmRegistry) fetchTriggerConfig(id *big.Int) ([]byte, error) {
	opts := r.buildCallOpts(r.ctx, nil)
	cfg, err := r.registry.GetUpkeepTriggerConfig(opts, id)
	if err != nil {
		r.lggr.Warnw("failed to get trigger config", "err", err)
		return nil, err
	}
	return cfg, nil
}
