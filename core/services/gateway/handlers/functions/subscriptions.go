package functions

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/functions/generated/functions_router"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type OnchainSubscriptionsConfig struct {
	RouterAddress      common.Address `json:"routerAddress"`
	BlockConfirmations uint           `json:"blockConfirmations"`
	QueryFrequencySec  uint           `json:"queryFrequencySec"`
	QueryTimeoutSec    uint           `json:"queryTimeoutSec"`
	QueryRangeSize     uint           `json:"queryRangeSize"`
}

// OnchainSubscriptions maintains a mirror of all subscriptions fetched from the blockchain (EVM-only).
// All methods are thread-safe.
//
//go:generate mockery --quiet --name OnchainSubscriptions --output ./mocks/ --case=underscore
type OnchainSubscriptions interface {
	job.ServiceCtx

	// GetSubscription returns a subscription for the given user address, or null if not found
	GetSubscription(common.Address) *functions_router.IFunctionsSubscriptionsSubscription
}

type onchainSubscriptions struct {
	utils.StartStopOnce

	config             OnchainSubscriptionsConfig
	subscriptions      map[common.Address]*functions_router.IFunctionsSubscriptionsSubscription
	client             evmclient.Client
	router             *functions_router.FunctionsRouter
	blockConfirmations *big.Int
	lggr               logger.Logger
	closeWait          sync.WaitGroup
	rwMutex            sync.RWMutex
	stopCh             utils.StopChan
}

func NewOnchainSubscriptions(client evmclient.Client, config OnchainSubscriptionsConfig, lggr logger.Logger) (OnchainSubscriptions, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}
	if lggr == nil {
		return nil, errors.New("logger is nil")
	}
	router, err := functions_router.NewFunctionsRouter(config.RouterAddress, client)
	if err != nil {
		return nil, fmt.Errorf("unexpected error during functions_router.NewFunctionsRouter: %s", err)
	}
	return &onchainSubscriptions{
		config:             config,
		subscriptions:      make(map[common.Address]*functions_router.IFunctionsSubscriptionsSubscription),
		client:             client,
		router:             router,
		blockConfirmations: big.NewInt(int64(config.BlockConfirmations)),
		lggr:               lggr.Named("OnchainSubscriptions"),
		stopCh:             make(utils.StopChan),
	}, nil
}

func (s *onchainSubscriptions) Start(ctx context.Context) error {
	return s.StartOnce("OnchainSubscriptions", func() error {
		s.lggr.Info("starting onchain subscriptions")
		if s.config.QueryFrequencySec == 0 {
			return errors.New("OnchainSubscriptionsConfig.UpdateFrequencySec must be greater than 0")
		}
		if s.config.QueryTimeoutSec == 0 {
			return errors.New("OnchainSubscriptionsConfig.UpdateTimeoutSec must be greater than 0")
		}
		if s.config.QueryRangeSize == 0 {
			return errors.New("OnchainSubscriptionsConfig.QueryRangeSize must be greater than 0")
		}

		s.closeWait.Add(1)
		go s.queryLoop()

		return nil
	})
}

func (s *onchainSubscriptions) Close() error {
	return s.StopOnce("OnchainSubscriptions", func() (err error) {
		s.lggr.Info("closing onchain subscriptions")
		close(s.stopCh)
		s.closeWait.Wait()
		return nil
	})
}

func (s *onchainSubscriptions) GetSubscription(address common.Address) *functions_router.IFunctionsSubscriptionsSubscription {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()
	subscription, ok := s.subscriptions[address]
	if !ok {
		return nil
	}
	return subscription
}

func (s *onchainSubscriptions) queryLoop() {
	defer s.closeWait.Done()

	ticker := time.NewTicker(time.Duration(s.config.QueryFrequencySec) * time.Second)
	defer ticker.Stop()

	var start uint64 = 1

	queryFunc := func() {
		ctx, cancel := utils.ContextFromChanWithTimeout(s.stopCh, time.Duration(s.config.QueryTimeoutSec)*time.Second)
		defer cancel()

		latestBlockHeight, err := s.client.LatestBlockHeight(ctx)
		if err != nil || latestBlockHeight == nil {
			s.lggr.Errorw("Error calling LatestBlockHeight", "err", err, "latestBlockHeight", latestBlockHeight.Int64())
			return
		}

		blockNumber := big.NewInt(0).Sub(latestBlockHeight, s.blockConfirmations)

		count, err := s.getSubscriptionsCount(ctx, blockNumber)
		if err != nil {
			s.lggr.Errorw("Error getting subscriptions count", "err", err)
			return
		}
		if count == 0 {
			s.lggr.Info("Router has no subscriptions yet")
			return
		}

		if start > count {
			start = 1
		}

		end := start + uint64(s.config.QueryRangeSize)
		if end > count {
			end = count
		}
		if err := s.querySubscriptionsRange(ctx, blockNumber, start, end); err != nil {
			s.lggr.Errorw("Error querying subscriptions", "err", err, "start", start, "end", end)
			return
		}

		start = end + 1
	}

	for {
		select {
		case <-s.stopCh:
			return
		case <-ticker.C:
			queryFunc()
		}
	}
}

func (s *onchainSubscriptions) querySubscriptionsRange(ctx context.Context, blockNumber *big.Int, start, end uint64) error {
	subscriptions, err := s.router.GetSubscriptionsInRange(&bind.CallOpts{
		Pending:     false,
		BlockNumber: blockNumber,
		Context:     ctx,
	}, start, end)
	if err != nil {
		return errors.Wrap(err, "unexpected error during functions_router.GetSubscriptionsInRange")
	}

	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()
	for _, subscription := range subscriptions {
		if subscription.Owner == utils.ZeroAddress {
			continue
		}
		subscription := subscription
		s.subscriptions[subscription.Owner] = &subscription
	}

	return nil
}

func (s *onchainSubscriptions) getSubscriptionsCount(ctx context.Context, blockNumber *big.Int) (uint64, error) {
	return s.router.GetSubscriptionCount(&bind.CallOpts{
		Pending:     false,
		BlockNumber: blockNumber,
		Context:     ctx,
	})
}
