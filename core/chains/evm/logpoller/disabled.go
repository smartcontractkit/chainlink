package logpoller

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"
	pkgerrors "github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

var (
	ErrDisabled                 = pkgerrors.New("log poller disabled")
	LogPollerDisabled LogPoller = disabled{}
)

type disabled struct{}

func (disabled) Name() string { return "disabledLogPoller" }

func (disabled) Start(ctx context.Context) error { return ErrDisabled }

func (disabled) Close() error { return ErrDisabled }

func (disabled) Healthy() error {
	return ErrDisabled
}

func (disabled) Ready() error { return ErrDisabled }

func (disabled) HealthReport() map[string]error {
	return map[string]error{"disabledLogPoller": ErrDisabled}
}

func (disabled) Replay(ctx context.Context, fromBlock int64) error { return ErrDisabled }

func (disabled) ReplayAsync(fromBlock int64) {}

func (disabled) RegisterFilter(ctx context.Context, filter Filter) error { return ErrDisabled }

func (disabled) UnregisterFilter(ctx context.Context, name string) error { return ErrDisabled }

func (disabled) HasFilter(name string) bool { return false }

func (disabled) GetFilters() map[string]Filter { return nil }

func (disabled) LatestBlock(ctx context.Context) (LogPollerBlock, error) {
	return LogPollerBlock{}, ErrDisabled
}

func (disabled) GetBlocksRange(ctx context.Context, numbers []uint64) ([]LogPollerBlock, error) {
	return nil, ErrDisabled
}

func (disabled) Logs(ctx context.Context, start, end int64, eventSig common.Hash, address common.Address) ([]Log, error) {
	return nil, ErrDisabled
}

func (disabled) LogsWithSigs(ctx context.Context, start, end int64, eventSigs []common.Hash, address common.Address) ([]Log, error) {
	return nil, ErrDisabled
}

func (disabled) LatestLogByEventSigWithConfs(ctx context.Context, eventSig common.Hash, address common.Address, confs evmtypes.Confirmations) (*Log, error) {
	return nil, ErrDisabled
}

func (disabled) LatestLogEventSigsAddrsWithConfs(ctx context.Context, fromBlock int64, eventSigs []common.Hash, addresses []common.Address, confs evmtypes.Confirmations) ([]Log, error) {
	return nil, ErrDisabled
}

func (disabled) IndexedLogs(ctx context.Context, eventSig common.Hash, address common.Address, topicIndex int, topicValues []common.Hash, confs evmtypes.Confirmations) ([]Log, error) {
	return nil, ErrDisabled
}

func (disabled) IndexedLogsByBlockRange(ctx context.Context, start, end int64, eventSig common.Hash, address common.Address, topicIndex int, topicValues []common.Hash) ([]Log, error) {
	return nil, ErrDisabled
}

func (d disabled) IndexedLogsByTxHash(ctx context.Context, eventSig common.Hash, address common.Address, txHash common.Hash) ([]Log, error) {
	return nil, ErrDisabled
}

func (disabled) IndexedLogsTopicGreaterThan(ctx context.Context, eventSig common.Hash, address common.Address, topicIndex int, topicValueMin common.Hash, confs evmtypes.Confirmations) ([]Log, error) {
	return nil, ErrDisabled
}

func (disabled) IndexedLogsTopicRange(ctx context.Context, eventSig common.Hash, address common.Address, topicIndex int, topicValueMin common.Hash, topicValueMax common.Hash, confs evmtypes.Confirmations) ([]Log, error) {
	return nil, ErrDisabled
}

func (disabled) LogsDataWordRange(ctx context.Context, eventSig common.Hash, address common.Address, wordIndex int, wordValueMin, wordValueMax common.Hash, confs evmtypes.Confirmations) ([]Log, error) {
	return nil, ErrDisabled
}

func (disabled) LogsDataWordGreaterThan(ctx context.Context, eventSig common.Hash, address common.Address, wordIndex int, wordValueMin common.Hash, confs evmtypes.Confirmations) ([]Log, error) {
	return nil, ErrDisabled
}

func (d disabled) IndexedLogsWithSigsExcluding(ctx context.Context, address common.Address, eventSigA, eventSigB common.Hash, topicIndex int, fromBlock, toBlock int64, confs evmtypes.Confirmations) ([]Log, error) {
	return nil, ErrDisabled
}

func (d disabled) LogsCreatedAfter(ctx context.Context, eventSig common.Hash, address common.Address, time time.Time, confs evmtypes.Confirmations) ([]Log, error) {
	return nil, ErrDisabled
}

func (d disabled) IndexedLogsCreatedAfter(ctx context.Context, eventSig common.Hash, address common.Address, topicIndex int, topicValues []common.Hash, after time.Time, confs evmtypes.Confirmations) ([]Log, error) {
	return nil, ErrDisabled
}

func (d disabled) LatestBlockByEventSigsAddrsWithConfs(ctx context.Context, fromBlock int64, eventSigs []common.Hash, addresses []common.Address, confs evmtypes.Confirmations) (int64, error) {
	return 0, ErrDisabled
}

func (d disabled) LogsDataWordBetween(ctx context.Context, eventSig common.Hash, address common.Address, wordIndexMin, wordIndexMax int, wordValue common.Hash, confs evmtypes.Confirmations) ([]Log, error) {
	return nil, ErrDisabled
}

func (d disabled) FilteredLogs(_ context.Context, _ query.KeyFilter, _ query.LimitAndSort, _ string) ([]Log, error) {
	return nil, ErrDisabled
}

func (d disabled) FindLCA(ctx context.Context) (*LogPollerBlock, error) {
	return nil, ErrDisabled
}

func (d disabled) DeleteLogsAndBlocksAfter(ctx context.Context, start int64) error {
	return ErrDisabled
}
