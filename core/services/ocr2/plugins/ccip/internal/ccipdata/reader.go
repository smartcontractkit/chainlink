package ccipdata

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/cciptypes"
)

const (
	V1_0_0 = "1.0.0"
	V1_1_0 = "1.1.0"
	V1_2_0 = "1.2.0"
	V1_4_0 = "1.4.0"
	V1_5_0 = "1.5.0-dev"
)

type Event[T any] struct {
	Data T
	cciptypes.TxMeta
}

func LogsConfirmations(finalized bool) logpoller.Confirmations {
	if finalized {
		return logpoller.Finalized
	}
	return logpoller.Unconfirmed
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
