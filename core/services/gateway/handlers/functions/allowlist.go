package functions

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/ocr2dr_oracle"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type OnchainAllowlistConfig struct {
	// ContractAddress is required
	ContractAddress    common.Address `json:"allowlistContractAddress"`
	BlockConfirmations uint           `json:"allowlistBlockConfirmations"`
	// UpdateFrequencySec can be zero to disable periodic updates
	UpdateFrequencySec uint `json:"allowlistUpdateFrequencySec"`
	UpdateTimeoutSec   uint `json:"allowlistUpdateTimeoutSec"`
}

// OnchainAllowlist maintains an allowlist of addresses fetched from the blockchain (EVM-only).
// Use UpdateFromContract() for a one-time update or set OnchainAllowlistConfig.UpdateFrequencySec for period updates.
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
	contract           *ocr2dr_oracle.OCR2DROracle
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
	contract, err := ocr2dr_oracle.NewOCR2DROracle(config.ContractAddress, client)
	if err != nil {
		return nil, fmt.Errorf("unexpected error during NewOCR2DROracle: %s", err)
	}
	allowlist := &onchainAllowlist{
		config:             config,
		client:             client,
		contract:           contract,
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
		a.closeWait.Add(1)
		go func() {
			defer a.closeWait.Done()
			ticker := time.NewTicker(time.Duration(a.config.UpdateFrequencySec))
			defer ticker.Stop()
			for {
				select {
				case <-a.stopCh:
					return
				case <-ticker.C:
					timeoutCtx, cancel := utils.ContextFromChanWithTimeout(a.stopCh, time.Duration(a.config.UpdateTimeoutSec))
					if err := a.UpdateFromContract(timeoutCtx); err != nil {
						a.lggr.Errorw("error calling UpdateFromContract", "err", err)
					}
					cancel()
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
	addrList, err := a.contract.GetAuthorizedSenders(&bind.CallOpts{
		Pending:     false,
		BlockNumber: blockNum,
		Context:     ctx,
	})
	if err != nil {
		return errors.Wrap(err, "error calling GetAuthorizedSenders")
	}
	newAllowlist := make(map[common.Address]struct{})
	for _, addr := range addrList {
		newAllowlist[addr] = struct{}{}
	}
	a.allowlist.Store(&newAllowlist)
	a.lggr.Infow("allowlist updated successfully", "len", len(addrList), "blockNumber", blockNum)
	return nil
}
