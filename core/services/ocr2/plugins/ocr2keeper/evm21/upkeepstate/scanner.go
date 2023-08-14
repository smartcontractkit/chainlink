package upkeepstate

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	iregistry21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

type PerformedLogsScanner interface {
	WorkIDsInRange(ctx context.Context, start, end int64) ([]string, error)
}

type performedEventsScanner struct {
	lggr            logger.Logger
	poller          logpoller.LogPoller
	registryAddress common.Address
}

func NewPerformedEventsScanner(
	lggr logger.Logger,
	poller logpoller.LogPoller,
	registryAddress common.Address,
) *performedEventsScanner {
	return &performedEventsScanner{
		lggr:            lggr,
		poller:          poller,
		registryAddress: registryAddress,
	}
}

func (s *performedEventsScanner) Start(_ context.Context) error {
	return s.poller.RegisterFilter(logpoller.Filter{
		Name: dedupFilterName(s.registryAddress),
		EventSigs: []common.Hash{
			// listening to dedup key added event
			iregistry21.IKeeperRegistryMasterDedupKeyAdded{}.Topic(),
		},
		Addresses: []common.Address{s.registryAddress},
	})
}

// implements io.Closer, does nothing upon close
func (s *performedEventsScanner) Close() error {
	return nil
}

func (s *performedEventsScanner) WorkIDsInRange(ctx context.Context, start, end int64) ([]string, error) {
	logs, err := s.poller.LogsWithSigs(
		start,
		end,
		[]common.Hash{
			iregistry21.IKeeperRegistryMasterDedupKeyAdded{}.Topic(),
		},
		s.registryAddress,
		pg.WithParentCtx(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("error fetching logs: %w", err)
	}

	return s.logsToWorkIDs(logs), nil
}

func (s *performedEventsScanner) logsToWorkIDs(logs []logpoller.Log) []string {
	workIDs := make([]string, 0)
	for _, log := range logs {
		topics := log.GetTopics()
		if len(topics) < 2 {
			s.lggr.Debugw("unexpected log topics", "topics", topics)
			continue
		}
		workIDs = append(workIDs, hexutil.Encode(topics[1].Bytes()))
	}
	return workIDs
}

func dedupFilterName(addr common.Address) string {
	return logpoller.FilterName("KeepersRegistry UpkeepStates Deduped", addr)
}
