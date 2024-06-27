package upkeepstate

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	iregistry21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/logprovider"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

var (
	_ PerformedLogsScanner = &performedEventsScanner{}

	workIDsBatchSize = 25
)

type PerformedLogsScanner interface {
	ScanWorkIDs(ctx context.Context, workIDs ...string) ([]string, error)

	Start(context.Context) error
	io.Closer
}

type performedEventsScanner struct {
	lggr            logger.Logger
	poller          logpoller.LogPoller
	registryAddress common.Address

	finalityDepth uint32
}

func NewPerformedEventsScanner(
	lggr logger.Logger,
	poller logpoller.LogPoller,
	registryAddress common.Address,
	finalityDepth uint32,
) *performedEventsScanner {
	return &performedEventsScanner{
		lggr:            lggr.Named("EventsScanner"),
		poller:          poller,
		registryAddress: registryAddress,
		finalityDepth:   finalityDepth,
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
		Retention: logprovider.LogRetention,
	})
}

// Close implements io.Closer and does nothing
func (s *performedEventsScanner) Close() error {
	return nil
}

func (s *performedEventsScanner) ScanWorkIDs(ctx context.Context, workID ...string) ([]string, error) {
	var ids []common.Hash
	for _, id := range workID {
		ids = append(ids, common.HexToHash(id))
	}
	logs := make([]logpoller.Log, 0)
	for i := 0; i < len(ids); i += workIDsBatchSize {
		end := i + workIDsBatchSize
		if end > len(ids) {
			end = len(ids)
		}
		batch := ids[i:end]
		batchLogs, err := s.poller.IndexedLogs(iregistry21.IKeeperRegistryMasterDedupKeyAdded{}.Topic(), s.registryAddress, 1, batch, logpoller.Confirmations(s.finalityDepth), pg.WithParentCtx(ctx))
		if err != nil {
			return nil, fmt.Errorf("error fetching logs: %w", err)
		}
		logs = append(logs, batchLogs...)
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
		workIDs = append(workIDs, hex.EncodeToString(topics[1].Bytes()))
	}
	return workIDs
}

func dedupFilterName(addr common.Address) string {
	return logpoller.FilterName("KeepersRegistry UpkeepStates Deduped", addr)
}
