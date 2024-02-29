package allowlist

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"regexp"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"golang.org/x/mod/semver"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/functions/generated/functions_allow_list"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/functions/generated/functions_router"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const (
	defaultStoredAllowlistBatchSize      = 1000
	defaultOnchainAllowlistBatchSize     = 100
	defaultFetchingDelayInRangeSec       = 1
	tosContractMinBatchProcessingVersion = "v1.1.0"
)

type OnchainAllowlistConfig struct {
	// ContractAddress is required
	ContractAddress    common.Address `json:"contractAddress"`
	ContractVersion    uint32         `json:"contractVersion"`
	BlockConfirmations uint           `json:"blockConfirmations"`
	// UpdateFrequencySec can be zero to disable periodic updates
	UpdateFrequencySec        uint `json:"updateFrequencySec"`
	UpdateTimeoutSec          uint `json:"updateTimeoutSec"`
	StoredAllowlistBatchSize  uint `json:"storedAllowlistBatchSize"`
	OnchainAllowlistBatchSize uint `json:"onchainAllowlistBatchSize"`
	// FetchingDelayInRangeSec prevents RPC client being rate limited when fetching the allowlist in ranges.
	FetchingDelayInRangeSec uint `json:"fetchingDelayInRangeSec"`
}

// OnchainAllowlist maintains an allowlist of addresses fetched from the blockchain (EVM-only).
// Use UpdateFromContract() for a one-time update or set OnchainAllowlistConfig.UpdateFrequencySec
// for repeated updates.
// All methods are thread-safe.
//
//go:generate mockery --quiet --name OnchainAllowlist --output ./mocks/ --case=underscore
type OnchainAllowlist interface {
	job.ServiceCtx

	Allow(common.Address) bool
	UpdateFromContract(ctx context.Context) error
}

type onchainAllowlist struct {
	services.StateMachine

	config             OnchainAllowlistConfig
	allowlist          atomic.Pointer[map[common.Address]struct{}]
	orm                ORM
	client             evmclient.Client
	contractV1         *functions_router.FunctionsRouter
	blockConfirmations *big.Int
	lggr               logger.Logger
	closeWait          sync.WaitGroup
	stopCh             services.StopChan
}

func NewOnchainAllowlist(client evmclient.Client, config OnchainAllowlistConfig, orm ORM, lggr logger.Logger) (OnchainAllowlist, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}
	if lggr == nil {
		return nil, errors.New("logger is nil")
	}
	if config.ContractVersion != 1 {
		return nil, fmt.Errorf("unsupported contract version %d", config.ContractVersion)
	}

	if config.StoredAllowlistBatchSize == 0 {
		lggr.Info("StoredAllowlistBatchSize not specified, using default size: ", defaultStoredAllowlistBatchSize)
		config.StoredAllowlistBatchSize = defaultStoredAllowlistBatchSize
	}

	if config.OnchainAllowlistBatchSize == 0 {
		lggr.Info("OnchainAllowlistBatchSize not specified, using default size: ", defaultOnchainAllowlistBatchSize)
		config.OnchainAllowlistBatchSize = defaultOnchainAllowlistBatchSize
	}

	if config.FetchingDelayInRangeSec == 0 {
		lggr.Info("FetchingDelayInRangeSec not specified, using default delay: ", defaultFetchingDelayInRangeSec)
		config.FetchingDelayInRangeSec = defaultFetchingDelayInRangeSec
	}

	if config.UpdateFrequencySec != 0 && config.FetchingDelayInRangeSec >= config.UpdateFrequencySec {
		return nil, fmt.Errorf("to avoid updates overlapping FetchingDelayInRangeSec:%d should be less than UpdateFrequencySec:%d", config.FetchingDelayInRangeSec, config.UpdateFrequencySec)
	}

	contractV1, err := functions_router.NewFunctionsRouter(config.ContractAddress, client)
	if err != nil {
		return nil, fmt.Errorf("unexpected error during functions_router.NewFunctionsRouter: %s", err)
	}
	allowlist := &onchainAllowlist{
		config:             config,
		orm:                orm,
		client:             client,
		contractV1:         contractV1,
		blockConfirmations: big.NewInt(int64(config.BlockConfirmations)),
		lggr:               lggr.Named("OnchainAllowlist"),
		stopCh:             make(services.StopChan),
	}
	emptyMap := make(map[common.Address]struct{})
	allowlist.allowlist.Store(&emptyMap)
	return allowlist, nil
}

func (a *onchainAllowlist) Start(ctx context.Context) error {
	return a.StartOnce("OnchainAllowlist", func() error {
		a.lggr.Info("starting onchain allowlist")
		if a.config.UpdateFrequencySec == 0 || a.config.UpdateTimeoutSec == 0 {
			a.lggr.Info("OnchainAllowlist periodic updates are disabled")
			return nil
		}

		a.loadStoredAllowedSenderList()

		updateOnce := func() {
			timeoutCtx, cancel := utils.ContextFromChanWithTimeout(a.stopCh, time.Duration(a.config.UpdateTimeoutSec)*time.Second)
			if err := a.UpdateFromContract(timeoutCtx); err != nil {
				a.lggr.Errorw("error calling UpdateFromContract", "err", err)
			}
			cancel()
		}

		a.closeWait.Add(1)
		go func() {
			defer a.closeWait.Done()
			// update immediately after start to populate the allowlist without waiting UpdateFrequencySec seconds
			updateOnce()
			ticker := time.NewTicker(time.Duration(a.config.UpdateFrequencySec) * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-a.stopCh:
					return
				case <-ticker.C:
					updateOnce()
				}
			}
		}()
		return nil
	})
}

func (a *onchainAllowlist) Close() error {
	return a.StopOnce("OnchainAllowlist", func() (err error) {
		a.lggr.Info("closing onchain allowlist")
		close(a.stopCh)
		a.closeWait.Wait()
		return nil
	})
}

func (a *onchainAllowlist) Allow(address common.Address) bool {
	allowlist := *a.allowlist.Load()
	_, ok := allowlist[address]
	return ok
}

func (a *onchainAllowlist) UpdateFromContract(ctx context.Context) error {
	latestBlockHeight, err := a.client.LatestBlockHeight(ctx)
	if err != nil {
		return errors.Wrap(err, "error calling LatestBlockHeight")
	}
	if latestBlockHeight == nil {
		return errors.New("LatestBlockHeight returned nil")
	}
	blockNum := big.NewInt(0).Sub(latestBlockHeight, a.blockConfirmations)
	return a.updateFromContractV1(ctx, blockNum)
}

func (a *onchainAllowlist) updateFromContractV1(ctx context.Context, blockNum *big.Int) error {
	tosID, err := a.contractV1.GetAllowListId(&bind.CallOpts{
		Pending: false,
		Context: ctx,
	})
	if err != nil {
		return errors.Wrap(err, "unexpected error during functions_router.GetAllowListId")
	}
	a.lggr.Debugw("successfully fetched allowlist route ID", "id", hex.EncodeToString(tosID[:]))
	if tosID == [32]byte{} {
		return errors.New("allowlist route ID has not been set")
	}
	tosAddress, err := a.contractV1.GetContractById(&bind.CallOpts{
		Pending: false,
		Context: ctx,
	}, tosID)
	if err != nil {
		return errors.Wrap(err, "unexpected error during functions_router.GetContractById")
	}
	a.lggr.Debugw("successfully fetched allowlist contract address", "address", tosAddress)
	tosContract, err := functions_allow_list.NewTermsOfServiceAllowList(tosAddress, a.client)
	if err != nil {
		return errors.Wrap(err, "unexpected error during functions_allow_list.NewTermsOfServiceAllowList")
	}

	var allowedSenderList []common.Address
	typeAndVersion, err := tosContract.TypeAndVersion(&bind.CallOpts{
		Pending:     false,
		BlockNumber: blockNum,
		Context:     ctx,
	})
	if err != nil {
		return errors.Wrap(err, "failed to fetch the tos contract type and version")
	}

	currentVersion, err := ExtractContractVersion(typeAndVersion)
	if err != nil {
		return fmt.Errorf("failed to extract version: %w", err)
	}

	if semver.Compare(tosContractMinBatchProcessingVersion, currentVersion) <= 0 {
		err = a.syncBlockedSenders(ctx, tosContract, blockNum)
		if err != nil {
			return errors.Wrap(err, "failed to sync the stored allowed and blocked senders")
		}

		allowedSenderList, err = a.getAllowedSendersBatched(ctx, tosContract, blockNum)
		if err != nil {
			return errors.Wrap(err, "failed to get allowed senders in rage")
		}
	} else {
		allowedSenderList, err = tosContract.GetAllAllowedSenders(&bind.CallOpts{
			Pending:     false,
			BlockNumber: blockNum,
			Context:     ctx,
		})
		if err != nil {
			return errors.Wrap(err, "error calling GetAllAllowedSenders")
		}

		err = a.orm.PurgeAllowedSenders()
		if err != nil {
			a.lggr.Errorf("failed to purge allowedSenderList: %w", err)
		}

		err = a.orm.CreateAllowedSenders(allowedSenderList)
		if err != nil {
			a.lggr.Errorf("failed to update stored allowedSenderList: %w", err)
		}
	}

	a.update(allowedSenderList)
	return nil
}

func (a *onchainAllowlist) getAllowedSendersBatched(ctx context.Context, tosContract *functions_allow_list.TermsOfServiceAllowList, blockNum *big.Int) ([]common.Address, error) {
	allowedSenderList := make([]common.Address, 0)
	count, err := tosContract.GetAllowedSendersCount(&bind.CallOpts{
		Pending:     false,
		BlockNumber: blockNum,
		Context:     ctx,
	})
	if err != nil {
		return nil, errors.Wrap(err, "unexpected error during functions_allow_list.GetAllowedSendersCount")
	}

	throttleTicker := time.NewTicker(time.Duration(a.config.FetchingDelayInRangeSec) * time.Second)
	for idxStart := uint64(0); idxStart < count; idxStart += uint64(a.config.OnchainAllowlistBatchSize) {
		<-throttleTicker.C

		idxEnd := idxStart + uint64(a.config.OnchainAllowlistBatchSize)
		if idxEnd >= count {
			idxEnd = count - 1
		}

		allowedSendersBatch, err := tosContract.GetAllowedSendersInRange(&bind.CallOpts{
			Pending:     false,
			BlockNumber: blockNum,
			Context:     ctx,
		}, idxStart, idxEnd)
		if err != nil {
			return nil, errors.Wrap(err, "error calling GetAllowedSendersInRange")
		}

		allowedSenderList = append(allowedSenderList, allowedSendersBatch...)
		err = a.orm.CreateAllowedSenders(allowedSendersBatch)
		if err != nil {
			a.lggr.Errorf("failed to update stored allowedSenderList: %w", err)
		}
	}
	throttleTicker.Stop()

	return allowedSenderList, nil
}

// syncBlockedSenders fetches the list of blocked addresses from the contract in batches
// and removes the addresses from the functions_allowlist table if present
func (a *onchainAllowlist) syncBlockedSenders(ctx context.Context, tosContract *functions_allow_list.TermsOfServiceAllowList, blockNum *big.Int) error {
	count, err := tosContract.GetBlockedSendersCount(&bind.CallOpts{
		Pending:     false,
		BlockNumber: blockNum,
		Context:     ctx,
	})
	if err != nil {
		return errors.Wrap(err, "unexpected error during functions_allow_list.GetBlockedSendersCount")
	}

	throttleTicker := time.NewTicker(time.Duration(a.config.FetchingDelayInRangeSec) * time.Second)
	for idxStart := uint64(0); idxStart < count; idxStart += uint64(a.config.OnchainAllowlistBatchSize) {
		<-throttleTicker.C

		idxEnd := idxStart + uint64(a.config.OnchainAllowlistBatchSize)
		if idxEnd >= count {
			idxEnd = count - 1
		}

		blockedSendersBatch, err := tosContract.GetBlockedSendersInRange(&bind.CallOpts{
			Pending:     false,
			BlockNumber: blockNum,
			Context:     ctx,
		}, idxStart, idxEnd)
		if err != nil {
			return errors.Wrap(err, "error calling GetAllowedSendersInRange")
		}

		err = a.orm.DeleteAllowedSenders(blockedSendersBatch)
		if err != nil {
			a.lggr.Errorf("failed to delete blocked address from allowed list in storage: %w", err)
		}
	}
	throttleTicker.Stop()

	return nil
}

func (a *onchainAllowlist) update(addrList []common.Address) {
	newAllowlist := make(map[common.Address]struct{})
	for _, addr := range addrList {
		newAllowlist[addr] = struct{}{}
	}
	a.allowlist.Store(&newAllowlist)
	a.lggr.Infow("allowlist updated successfully", "len", len(addrList))
}

func (a *onchainAllowlist) loadStoredAllowedSenderList() {
	allowedList := make([]common.Address, 0)
	offset := uint(0)
	for {
		asBatch, err := a.orm.GetAllowedSenders(offset, a.config.StoredAllowlistBatchSize)
		if err != nil {
			a.lggr.Errorf("failed to get stored allowed senders: %w", err)
			break
		}

		allowedList = append(allowedList, asBatch...)

		if len(asBatch) < int(a.config.StoredAllowlistBatchSize) {
			break
		}
		offset += a.config.StoredAllowlistBatchSize
	}

	a.update(allowedList)
}

func ExtractContractVersion(str string) (string, error) {
	pattern := `v(\d+).(\d+).(\d+)`
	re := regexp.MustCompile(pattern)

	match := re.FindStringSubmatch(str)
	if len(match) != 4 {
		return "", fmt.Errorf("version not found in string: %s", str)
	}
	return fmt.Sprintf("v%s.%s.%s", match[1], match[2], match[3]), nil
}
