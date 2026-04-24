package ccipdata

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
)

const (
	V1_0_0 = "1.0.0"
	V1_1_0 = "1.1.0"
	V1_2_0 = "1.2.0"
	V1_4_0 = "1.4.0"
	V1_5_0 = "1.5.0"
	//nolint:revive // was like this before
	V1_6_0 = "1.6.0"
)

const (
	// CommitExecLogsRetention defines the duration for which logs critical for Commit/Exec plugins processing are retained.
	// Although Exec relies on permissionlessExecThreshold which is lower than 24hours for picking eligible CommitRoots,
	// Commit still can reach to older logs because it filters them by sequence numbers. For instance, in case of RMN curse on chain,
	// we might have logs waiting in OnRamp to be committed first. When outage takes days we still would
	// be able to bring back processing without replaying any logs from chain. You can read that param as
	// "how long CCIP can be down and still be able to process all the messages after getting back to life".
	// Breaching this threshold would require replaying chain using LogPoller from the beginning of the outage.
	CommitExecLogsRetention = 30 * 24 * time.Hour // 30 days
	// CacheEvictionLogsRetention defines the duration for which logs used for caching on-chain data are kept.
	// Restarting node clears the cache entirely and rebuilds it from scratch by fetching data from chain,
	// so we don't need to keep these logs for very long. All events relying on cache.NewLogpollerEventsBased should use this retention.
	CacheEvictionLogsRetention = 7 * 24 * time.Hour // 7 days
	// PriceUpdatesLogsRetention defines the duration for which logs with price updates are kept.
	// These logs are emitted whenever the token price or gas price is updated and Commit scans very small time windows (e.g. 2 hours)
	PriceUpdatesLogsRetention = 1 * 24 * time.Hour // 1 day
)

type Event[T any] struct {
	Data T
	cciptypes.TxMeta
}

func LogsConfirmations(finalized bool) evmtypes.Confirmations {
	if finalized {
		return evmtypes.Finalized
	}
	return evmtypes.Unconfirmed
}

func ParseLogs[T any](logs []logpoller.Log, lggr logger.Logger, parseFunc func(log types.Log) (*T, error)) ([]Event[T], error) {
	reqs := make([]Event[T], 0, len(logs))

	for _, log := range logs {
		data, err := parseFunc(log.ToGethLog())
		if err != nil {
			lggr.Errorw("Unable to parse log", "err", err)
			continue
		}
		reqs = append(reqs, Event[T]{
			Data: *data,
			TxMeta: cciptypes.TxMeta{
				BlockTimestampUnixMilli: log.BlockTimestamp.UnixMilli(),
				BlockNumber:             uint64(log.BlockNumber),
				TxHash:                  log.TxHash.String(),
				LogIndex:                uint64(log.LogIndex),
			},
		})
	}

	if len(logs) != len(reqs) {
		return nil, fmt.Errorf("%d logs were not parsed", len(logs)-len(reqs))
	}
	return reqs, nil
}
