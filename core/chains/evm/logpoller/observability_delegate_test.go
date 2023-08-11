package logpoller_test

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	mocklp "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

func TestObservedLoggerDelegateToOrigin(t *testing.T) {
	eventSig := common.HexToHash("0xhash")
	address := common.HexToAddress("0xaddress")
	fromBlock := int64(20)
	confs := 10
	qopts := pg.WithParentCtx(testutils.Context(t))
	mockQOpt := mock.AnythingOfType("pg.QOpt")
	after := time.Now()

	delegateMockPoller := mocklp.NewLogPoller(t)
	lp := logpoller.NewTestObservedLogPoller(delegateMockPoller)

	delegateMockPoller.On("LogsCreatedAfter", eventSig, address, after, confs, mockQOpt).Return(nil, nil).Once()
	delegateMockPoller.On("LatestLogByEventSigWithConfs", eventSig, address, confs, mockQOpt).Return(nil, nil).Once()
	delegateMockPoller.On("LatestLogEventSigsAddrsWithConfs", fromBlock, []common.Hash{eventSig}, []common.Address{address}, confs, mockQOpt).Return(nil, nil).Once()
	delegateMockPoller.On("LatestBlockByEventSigsAddrsWithConfs", fromBlock, []common.Hash{eventSig}, []common.Address{address}, confs, mockQOpt).Return(int64(20), nil).Once()

	_, err := lp.LogsCreatedAfter(eventSig, address, after, confs, qopts)
	require.NoError(t, err)

	_, err = lp.LatestLogByEventSigWithConfs(eventSig, address, confs, qopts)
	require.NoError(t, err)

	_, err = lp.LatestLogEventSigsAddrsWithConfs(fromBlock, []common.Hash{eventSig}, []common.Address{address}, confs, qopts)
	require.NoError(t, err)

	_, err = lp.LatestBlockByEventSigsAddrsWithConfs(fromBlock, []common.Hash{eventSig}, []common.Address{address}, confs, qopts)
	require.NoError(t, err)
}
