package subscriptions

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/functions/generated/functions_router"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const defaultStoreBatchSize = 100

type OnchainSubscriptionsConfig struct {
	ContractAddress    common.Address `json:"contractAddress"`
	BlockConfirmations uint           `json:"blockConfirmations"`
	UpdateFrequencySec uint           `json:"updateFrequencySec"`
	UpdateTimeoutSec   uint           `json:"updateTimeoutSec"`
	UpdateRangeSize    uint           `json:"updateRangeSize"`
	StoreBatchSize     uint           `json:"storeBatchSize"`
}

// OnchainSubscriptions maintains a mirror of all subscriptions fetched from the blockchain (EVM-only).
// All methods are thread-safe.
//
//go:generate mockery --quiet --name OnchainSubscriptions --output ./mocks/ --case=underscore
type OnchainSubscriptions interface {
	job.ServiceCtx

	// GetMaxUserBalance returns a maximum subscription balance (juels), or error if user has no subscriptions.
	GetMaxUserBalance(common.Address) (*big.Int, error)
}

type onchainSubscriptions struct {
	services.StateMachine

	config             OnchainSubscriptionsConfig
	subscriptions      UserSubscriptions
	orm                ORM
	client             evmclient.Client
	router             *functions_router.FunctionsRouter
	blockConfirmations *big.Int
	lggr               logger.Logger
	closeWait          sync.WaitGroup
	rwMutex            sync.RWMutex
	stopCh             services.StopChan
}

func NewOnchainSubscriptions(client evmclient.Client, config OnchainSubscriptionsConfig, orm ORM, lggr logger.Logger) (OnchainSubscriptions, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}
	if lggr == nil {
		return nil, errors.New("logger is nil")
	}
	router, err := functions_router.NewFunctionsRouter(config.ContractAddress, client)
	if err != nil {
		return nil, fmt.Errorf("unexpected error during functions_router.NewFunctionsRouter: %s", err)
	}

	// if StoreBatchSize is not specified use the default value
	if config.StoreBatchSize == 0 {
		lggr.Info("StoreBatchSize not specified, using default size: ", defaultStoreBatchSize)
		config.StoreBatchSize = defaultStoreBatchSize
	}

	return &onchainSubscriptions{
		config:             config,
		subscriptions:      NewUserSubscriptions(),
		orm:                orm,
		client:             client,
		router:             router,
		blockConfirmations: big.NewInt(int64(config.BlockConfirmations)),
		lggr:               lggr.Named("OnchainSubscriptions"),
		stopCh:             make(services.StopChan),
	}, nil
}

func (s *onchainSubscriptions) Start(ctx context.Context) error {
	return s.StartOnce("OnchainSubscriptions", func() error {
		s.lggr.Info("starting onchain subscriptions")
		if s.config.UpdateFrequencySec == 0 {
			return errors.New("OnchainSubscriptionsConfig.UpdateFrequencySec must be greater than 0")
		}
		if s.config.UpdateTimeoutSec == 0 {
			return errors.New("OnchainSubscriptionsConfig.UpdateTimeoutSec must be greater than 0")
		}
		if s.config.UpdateRangeSize == 0 {
			return errors.New("OnchainSubscriptionsConfig.UpdateRangeSize must be greater than 0")
		}

		s.loadStoredSubscriptions(ctx)

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

func (s *onchainSubscriptions) GetMaxUserBalance(user common.Address) (*big.Int, error) {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()
	return s.subscriptions.GetMaxUserBalance(user)
}

func (s *onchainSubscriptions) queryLoop() {
	defer s.closeWait.Done()

	ticker := time.NewTicker(time.Duration(s.config.UpdateFrequencySec) * time.Second)
	defer ticker.Stop()

	start := uint64(1)
	lastKnownCount := uint64(0)

	queryFunc := func() {
		ctx, cancel := utils.ContextFromChanWithTimeout(s.stopCh, time.Duration(s.config.UpdateTimeoutSec)*time.Second)
		defer cancel()

		latestBlockHeight, err := s.client.LatestBlockHeight(ctx)
		if err != nil || latestBlockHeight == nil {
			s.lggr.Errorw("Error calling LatestBlockHeight", "err", err, "latestBlockHeight", latestBlockHeight)
			return
		}

		blockNumber := big.NewInt(0).Sub(latestBlockHeight, s.blockConfirmations)

		if lastKnownCount == 0 || start > lastKnownCount {
			count, err := s.getSubscriptionsCount(ctx, blockNumber)
			if err != nil {
				s.lggr.Errorw("Error getting new subscriptions count", "err", err)
			} else {
				s.lggr.Infow("Updated subscriptions count", "count", count, "blockNumber", blockNumber.Int64())
				lastKnownCount = count
			}
		}

		if lastKnownCount == 0 {
			s.lggr.Info("Router has no subscriptions yet")
			return
		}

		if start > lastKnownCount {
			start = 1
		}

		end := start + uint64(s.config.UpdateRangeSize) - 1
		if end > lastKnownCount {
			end = lastKnownCount
		}
		if err := s.querySubscriptionsRange(ctx, blockNumber, start, end); err != nil {
			s.lggr.Errorw("Error querying subscriptions", "err", err, "start", start, "end", end)
			return
		}

		start = end + 1
	}

	queryFunc()

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
	s.lggr.Debugw("Querying subscriptions", "blockNumber", blockNumber, "start", start, "end", end)

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
	for i, subscription := range subscriptions {
		subscriptionId := start + uint64(i)
		subscription := subscription
		updated := s.subscriptions.UpdateSubscription(subscriptionId, &subscription)
		if updated {
			if err = s.orm.UpsertSubscription(ctx, StoredSubscription{
				SubscriptionID:                      subscriptionId,
				IFunctionsSubscriptionsSubscription: subscription,
			}); err != nil {
				s.lggr.Errorf("unexpected error updating subscription in the db: %w", err)
			}
		}
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

func (s *onchainSubscriptions) loadStoredSubscriptions(ctx context.Context) {
	offset := uint(0)
	for {
		csBatch, err := s.orm.GetSubscriptions(ctx, offset, s.config.StoreBatchSize)
		if err != nil {
			break
		}

		for _, cs := range csBatch {
			_ = s.subscriptions.UpdateSubscription(cs.SubscriptionID, &functions_router.IFunctionsSubscriptionsSubscription{
				Balance:        cs.Balance,
				Owner:          cs.Owner,
				BlockedBalance: cs.BlockedBalance,
				ProposedOwner:  cs.ProposedOwner,
				Consumers:      cs.Consumers,
				Flags:          cs.Flags,
			})
		}
		s.lggr.Debugw("Loading stored subscriptions", "offset", offset, "batch_length", len(csBatch))

		if len(csBatch) != int(s.config.StoreBatchSize) {
			break
		}
		offset += s.config.StoreBatchSize
	}
}
