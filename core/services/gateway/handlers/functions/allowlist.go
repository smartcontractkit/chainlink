package functions

import (
	"context"
	"fmt"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/ocr2dr_oracle"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// OnchainAllowlist maintains an allowlist of addresses fetched from the blockchain (EVM-only).
// Use UpdateFromContract() for a one-time update or UpdatePeriodically() for periodic updates.
// All methods are thread-safe.
//
//go:generate mockery --quiet --name OnchainAllowlist --output ./mocks/ --case=underscore
type OnchainAllowlist interface {
	Allow(common.Address) bool
	UpdateFromContract(ctx context.Context) error
	UpdatePeriodically(ctx context.Context, updateFrequency time.Duration, updateTimeout time.Duration)
}

type onchainAllowlist struct {
	allowlist          atomic.Pointer[map[common.Address]struct{}]
	client             evmclient.Client
	contract           *ocr2dr_oracle.OCR2DROracle
	blockConfirmations *big.Int
	lggr               logger.Logger
}

func NewOnchainAllowlist(client evmclient.Client, contractAddress common.Address, blockConfirmations int64, lggr logger.Logger) (OnchainAllowlist, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}
	if lggr == nil {
		return nil, errors.New("logger is nil")
	}
	contract, err := ocr2dr_oracle.NewOCR2DROracle(contractAddress, client)
	if err != nil {
		return nil, fmt.Errorf("unexpected error during NewOCR2DROracle: %s", err)
	}
	allowlist := &onchainAllowlist{
		client:             client,
		contract:           contract,
		blockConfirmations: big.NewInt(blockConfirmations),
		lggr:               lggr.Named("OnchainAllowlist"),
	}
	emptyMap := make(map[common.Address]struct{})
	allowlist.allowlist.Store(&emptyMap)
	return allowlist, nil
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

func (a *onchainAllowlist) UpdatePeriodically(ctx context.Context, updateFrequency time.Duration, updateTimeout time.Duration) {
	ticker := time.NewTicker(updateFrequency)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			timeoutCtx, cancel := context.WithTimeout(ctx, updateTimeout)
			err := a.UpdateFromContract(timeoutCtx)
			if err != nil {
				a.lggr.Errorw("error calling UpdateFromContract", "err", err)
			}
			cancel()
		}
	}
}
