package v2

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"
	"github.com/smartcontractkit/chainlink-evm/pkg/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncer/versioning"
)

const (
	allowlistedRequestsSyncerName = "AllowlistedRequestsSyncer"

	GetActiveAllowlistedRequestsReverseMethodName = "getActiveAllowlistedRequestsReverse"
	TotalAllowlistedRequestsMethodName            = "totalAllowlistedRequests"

	defaultTickIntervalForAllowlistedRequests = 5 * time.Second
)

// AllowlistedRequestsSyncer periodically fetches allowlisted vault requests
// from the WorkflowRegistry contract and provides them to consumers.
type AllowlistedRequestsSyncer interface {
	services.Service
	GetAllowlistedRequests(ctx context.Context) []workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest
}

type allowlistedRequestsSyncer struct {
	services.StateMachine

	stopCh services.StopChan
	wg     sync.WaitGroup

	lggr                    logger.Logger
	workflowRegistryAddress string

	lastSeenAllowlistedRequestsCount *big.Int
	allowListedRequests              []workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest
	allowListedMu                    sync.RWMutex

	contractReaderFn versioning.ContractReaderFactory
	contractReader   types.ContractReader

	ticker <-chan time.Time
}

// AllowlistedRequestsSyncerOption is a functional option for configuring the syncer.
type AllowlistedRequestsSyncerOption func(*allowlistedRequestsSyncer)

// WithAllowlistedRequestsTicker overrides the default tick interval for testing.
func WithAllowlistedRequestsTicker(ticker <-chan time.Time) AllowlistedRequestsSyncerOption {
	return func(s *allowlistedRequestsSyncer) {
		s.ticker = ticker
	}
}

// NewAllowlistedRequestsSyncer creates a new syncer that periodically fetches
// allowlisted requests from the WorkflowRegistry contract.
func NewAllowlistedRequestsSyncer(
	lggr logger.Logger,
	contractReaderFn versioning.ContractReaderFactory,
	addr string,
	opts ...AllowlistedRequestsSyncerOption,
) *allowlistedRequestsSyncer {
	s := &allowlistedRequestsSyncer{
		lggr:                             logger.Named(lggr, allowlistedRequestsSyncerName),
		contractReaderFn:                 contractReaderFn,
		workflowRegistryAddress:          addr,
		lastSeenAllowlistedRequestsCount: big.NewInt(0),
		stopCh:                           make(services.StopChan),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *allowlistedRequestsSyncer) Start(_ context.Context) error {
	return s.StartOnce(s.Name(), func() error {
		ctx, cancel := s.stopCh.NewCtx()
		initDoneCh := make(chan struct{})

		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			defer close(initDoneCh)

			ticker := s.getTicker()
			for s.contractReader == nil {
				select {
				case <-ctx.Done():
					return
				case <-ticker:
					reader, err := s.newContractReader(ctx)
					if err != nil {
						s.lggr.Infow("contract reader unavailable", "error", err.Error())
						continue
					}
					s.contractReader = reader
				}
			}
		}()

		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			defer cancel()
			select {
			case <-initDoneCh:
			case <-ctx.Done():
				return
			}
			s.syncLoop(ctx)
		}()

		return nil
	})
}

func (s *allowlistedRequestsSyncer) Close() error {
	return s.StopOnce(s.Name(), func() error {
		close(s.stopCh)
		s.wg.Wait()
		return nil
	})
}

func (s *allowlistedRequestsSyncer) Ready() error {
	return nil
}

func (s *allowlistedRequestsSyncer) HealthReport() map[string]error {
	return map[string]error{s.Name(): s.Healthy()}
}

func (s *allowlistedRequestsSyncer) Name() string {
	return allowlistedRequestsSyncerName
}

func (s *allowlistedRequestsSyncer) GetAllowlistedRequests(_ context.Context) []workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest {
	s.allowListedMu.RLock()
	defer s.allowListedMu.RUnlock()
	out := make([]workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest, len(s.allowListedRequests))
	copy(out, s.allowListedRequests)
	return out
}

func (s *allowlistedRequestsSyncer) syncLoop(ctx context.Context) {
	ticker := s.getTicker()
	s.lggr.Debug("starting syncAllowlistedRequests")
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker:
			newRequests, totalCount, head, err := s.fetchAllowlistedRequests(ctx, s.contractReader)
			if err != nil {
				s.lggr.Errorw("failed to fetch allowlisted requests", "err", err)
				continue
			}
			s.allowListedMu.Lock()
			active := []workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest{}
			expiredCount := 0
			for _, req := range s.allowListedRequests {
				if int64(req.ExpiryTimestamp) > time.Now().Unix() {
					active = append(active, req)
				} else {
					expiredCount++
				}
			}
			active = append(active, newRequests...)
			s.allowListedRequests = active
			s.lastSeenAllowlistedRequestsCount = totalCount
			s.lggr.Debugw("synced allowlisted requests",
				"newRequestsNum", len(newRequests),
				"expiredRequestsNum", expiredCount,
				"activeRequestsNum", len(s.allowListedRequests),
				"lastSeenOnchainRequestsNum", s.lastSeenAllowlistedRequestsCount,
				"blockHeight", head.Height,
			)
			s.allowListedMu.Unlock()
		}
	}
}

func (s *allowlistedRequestsSyncer) getTicker() <-chan time.Time {
	if s.ticker != nil {
		return s.ticker
	}
	return time.NewTicker(defaultTickIntervalForAllowlistedRequests).C
}

func (s *allowlistedRequestsSyncer) newContractReader(ctx context.Context) (types.ContractReader, error) {
	contractReaderCfg := config.ChainReaderConfig{
		Contracts: map[string]config.ChainContractReader{
			WorkflowRegistryContractName: {
				ContractABI: workflow_registry_wrapper_v2.WorkflowRegistryABI,
				Configs: map[string]*config.ChainReaderDefinition{
					GetActiveAllowlistedRequestsReverseMethodName: {
						ChainSpecificName: GetActiveAllowlistedRequestsReverseMethodName,
						ReadType:          config.Method,
					},
					TotalAllowlistedRequestsMethodName: {
						ChainSpecificName: TotalAllowlistedRequestsMethodName,
						ReadType:          config.Method,
					},
				},
			},
		},
	}

	marshalledCfg, err := json.Marshal(contractReaderCfg)
	if err != nil {
		return nil, err
	}

	reader, err := s.contractReaderFn(ctx, marshalledCfg)
	if err != nil {
		return nil, err
	}

	bc := types.BoundContract{
		Name:    WorkflowRegistryContractName,
		Address: s.workflowRegistryAddress,
	}

	if err := reader.Bind(ctx, []types.BoundContract{bc}); err != nil {
		return nil, err
	}

	if err := reader.Start(ctx); err != nil {
		return nil, err
	}

	return reader, nil
}

func (s *allowlistedRequestsSyncer) fetchAllowlistedRequests(
	ctx context.Context,
	contractReader types.ContractReader,
) ([]workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest, *big.Int, *types.Head, error) {
	if contractReader == nil {
		return nil, nil, nil, errors.New("cannot fetch allowlisted requests: nil contract reader")
	}
	contractBinding := types.BoundContract{
		Address: s.workflowRegistryAddress,
		Name:    WorkflowRegistryContractName,
	}

	var headAtLastRead *types.Head
	var totalAllowlistedRequestsResult *big.Int
	readIdentifier := contractBinding.ReadIdentifier(TotalAllowlistedRequestsMethodName)
	headAtLastRead, err := contractReader.GetLatestValueWithHeadData(
		ctx, readIdentifier, primitives.Unconfirmed, nil, &totalAllowlistedRequestsResult,
	)
	if err != nil {
		return []workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest{}, s.lastSeenAllowlistedRequestsCount, &types.Head{Height: "0"}, errors.New("failed to get latest value with head data. error: " + err.Error())
	}

	if s.lastSeenAllowlistedRequestsCount.Cmp(totalAllowlistedRequestsResult) == 0 {
		return []workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest{}, totalAllowlistedRequestsResult, headAtLastRead, nil
	}

	var newAllowlistedRequests []workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest
	readIdentifier = contractBinding.ReadIdentifier(GetActiveAllowlistedRequestsReverseMethodName)
	endIndex := new(big.Int).Sub(totalAllowlistedRequestsResult, big.NewInt(1))
	var startIndex *big.Int

	for {
		var response struct {
			AllowlistedRequests []workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest
			SearchComplete      bool
			err                 error
		}

		startIndex = new(big.Int).Sub(endIndex, big.NewInt(MaxResultsPerQuery-1))
		if startIndex.Cmp(s.lastSeenAllowlistedRequestsCount) < 0 {
			startIndex = s.lastSeenAllowlistedRequestsCount
		}

		params := GetActiveAllowlistedRequestsReverseParams{
			EndIndex:   endIndex,
			StartIndex: startIndex,
		}
		s.lggr.Debugw("getting active allowlisted requests",
			"endIndex", endIndex,
			"startIndex", startIndex,
		)
		headAtLastRead, err = contractReader.GetLatestValueWithHeadData(
			ctx, readIdentifier, primitives.Unconfirmed, params, &response,
		)
		if err != nil {
			return []workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest{}, s.lastSeenAllowlistedRequestsCount, &types.Head{Height: "0"}, errors.New("failed to get lastest value with head data. error: " + err.Error())
		}

		s.lggr.Debugw("contract call response",
			"fetchedAllowlistedRequestsNum", len(response.AllowlistedRequests),
			"searchComplete", response.SearchComplete,
			"error", response.err,
			"blockHeight", headAtLastRead.Height)

		for _, request := range response.AllowlistedRequests {
			newAllowlistedRequests = append(newAllowlistedRequests, workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest{
				RequestDigest:   request.RequestDigest,
				Owner:           request.Owner,
				ExpiryTimestamp: request.ExpiryTimestamp,
			})
		}

		if response.SearchComplete {
			break
		}

		endIndex = endIndex.Sub(endIndex, big.NewInt(MaxResultsPerQuery))
		if endIndex.Cmp(big.NewInt(0)) < 0 {
			endIndex = big.NewInt(0)
		}
	}

	return newAllowlistedRequests, totalAllowlistedRequestsResult, headAtLastRead, nil
}
