package logpoller

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

var (
	ErrDisabled                 = errors.New("log poller disabled")
	LogPollerDisabled LogPoller = disabled{}
)

type disabled struct{}

func (disabled) Name() string { return "disabledLogPoller" }

func (disabled) Start(ctx context.Context) error { return ErrDisabled }

func (disabled) Close() error { return ErrDisabled }

func (disabled) Ready() error { return ErrDisabled }

func (disabled) HealthReport() map[string]error {
	return map[string]error{"disabledLogPoller": ErrDisabled}
}

func (disabled) Replay(ctx context.Context, fromBlock int64) error { return ErrDisabled }

func (disabled) ReplayAsync(fromBlock int64) {}

func (disabled) RegisterFilter(filter Filter) error { return ErrDisabled }

func (disabled) UnregisterFilter(name string, q pg.Queryer) error { return ErrDisabled }

func (disabled) LatestBlock(qopts ...pg.QOpt) (int64, error) { return -1, ErrDisabled }

func (disabled) GetBlocksRange(ctx context.Context, numbers []uint64, qopts ...pg.QOpt) ([]LogPollerBlock, error) {
	return nil, ErrDisabled
}

func (disabled) Logs(start, end int64, eventSig common.Hash, address common.Address, qopts ...pg.QOpt) ([]Log, error) {
	return nil, ErrDisabled
}

func (disabled) LogsWithSigs(start, end int64, eventSigs []common.Hash, address common.Address, qopts ...pg.QOpt) ([]Log, error) {
	return nil, ErrDisabled
}

func (disabled) LatestLogByEventSigWithConfs(eventSig common.Hash, address common.Address, confs int, qopts ...pg.QOpt) (*Log, error) {
	return nil, ErrDisabled
}

func (disabled) LatestLogEventSigsAddrsWithConfs(fromBlock int64, eventSigs []common.Hash, addresses []common.Address, confs int, qopts ...pg.QOpt) ([]Log, error) {
	return nil, ErrDisabled
}

func (disabled) IndexedLogs(eventSig common.Hash, address common.Address, topicIndex int, topicValues []common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error) {
	return nil, ErrDisabled
}

func (disabled) IndexedLogsByBlockRange(start, end int64, eventSig common.Hash, address common.Address, topicIndex int, topicValues []common.Hash, qopts ...pg.QOpt) ([]Log, error) {
	return nil, ErrDisabled
}

func (disabled) IndexedLogsTopicGreaterThan(eventSig common.Hash, address common.Address, topicIndex int, topicValueMin common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error) {
	return nil, ErrDisabled
}

func (disabled) IndexedLogsTopicRange(eventSig common.Hash, address common.Address, topicIndex int, topicValueMin common.Hash, topicValueMax common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error) {
	return nil, ErrDisabled
}

func (disabled) LogsDataWordRange(eventSig common.Hash, address common.Address, wordIndex int, wordValueMin, wordValueMax common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error) {
	return nil, ErrDisabled
}

func (disabled) LogsDataWordGreaterThan(eventSig common.Hash, address common.Address, wordIndex int, wordValueMin common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error) {
	return nil, ErrDisabled
}

func (d disabled) IndexedLogsWithSigsExcluding(address common.Address, eventSigA, eventSigB common.Hash, topicIndex int, fromBlock, toBlock int64, confs int, qopts ...pg.QOpt) ([]Log, error) {
	return nil, ErrDisabled
}
