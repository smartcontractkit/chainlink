package functions

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/functions/generated/functions_allow_list"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/functions/generated/functions_router"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type OnchainAllowlistConfig struct {
	// ContractAddress is required
	ContractAddress    common.Address `json:"contractAddress"`
	ContractVersion    uint32         `json:"contractVersion"`
	BlockConfirmations uint           `json:"blockConfirmations"`
	// UpdateFrequencySec can be zero to disable periodic updates
	UpdateFrequencySec uint `json:"updateFrequencySec"`
	UpdateTimeoutSec   uint `json:"updateTimeoutSec"`
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
	utils.StartStopOnce

	config             OnchainAllowlistConfig
	allowlist          atomic.Pointer[map[common.Address]struct{}]
	client             evmclient.Client
	contractV1         *functions_router.FunctionsRouter
	blockConfirmations *big.Int
	lggr               logger.Logger
	closeWait          sync.WaitGroup
	stopCh             utils.StopChan
}

func NewOnchainAllowlist(client evmclient.Client, config OnchainAllowlistConfig, lggr logger.Logger) (OnchainAllowlist, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}
	if lggr == nil {
		return nil, errors.New("logger is nil")
	}
	if config.ContractVersion != 1 {
		return nil, fmt.Errorf("unsupported contract version %d", config.ContractVersion)
	}
	contractV1, err := functions_router.NewFunctionsRouter(config.ContractAddress, client)
	if err != nil {
		return nil, fmt.Errorf("unexpected error during functions_router.NewFunctionsRouter: %s", err)
	}
	allowlist := &onchainAllowlist{
		config:             config,
		client:             client,
		contractV1:         contractV1,
		blockConfirmations: big.NewInt(int64(config.BlockConfirmations)),
		lggr:               lggr.Named("OnchainAllowlist"),
		stopCh:             make(utils.StopChan),
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
	addrList, err := tosContract.GetAllAllowedSenders(&bind.CallOpts{
		Pending:     false,
		BlockNumber: blockNum,
		Context:     ctx,
	})
	if err != nil {
		return errors.Wrap(err, "error calling GetAllAllowedSenders")
	}
	a.update(addrList)
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
