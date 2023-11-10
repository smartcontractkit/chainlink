package log

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func TestBroadcaster_OnNewLog(t *testing.T) {
	addressToWatch := common.Address{1}
	testCases := []struct {
		setup             func() (lb *broadcaster)
		name              string
		logs              []types.Log
		shouldResubscribe []bool
		assertions        func(t *testing.T, lb *broadcaster)
	}{
		{
			func() *broadcaster {
				lggr := logger.TestLogger(t)
				regs := newRegistrations(lggr, *big.NewInt(1337))
				listener := &mockListener{}
				listener.On("JobID").Return(int32(1))
				o := &mockORM{}
				o.On("SetPendingMinBlock", mock.AnythingOfType("*int64"), mock.Anything).Return(nil)
				o.On("SetPendingMinBlock", mock.AnythingOfType("*int64"), mock.Anything).Return(nil)
				sub := &subscriber{
					listener: listener,
					opts: ListenerOpts{
						Contract:                 addressToWatch,
						MinIncomingConfirmations: 1,
						LogsWithTopics: map[common.Hash][][]Topic{
							common.Hash{1}: {},
						},
					},
				}
				regs.addSubscriber(sub)
				return &broadcaster{
					logger:        lggr,
					logPool:       newLogPool(lggr),
					registrations: regs,
					orm:           o,
				}
			},
			"new log",
			[]types.Log{
				{
					Address:     addressToWatch,
					BlockNumber: 1,
					BlockHash:   testutils.Random32Byte(),
					Data:        hexutil.MustDecode("0xdeadbeef"),
					TxHash:      testutils.Random32Byte(),
					Index:       0,
					TxIndex:     1,
					Topics:      []common.Hash{{1}},
					Removed:     false,
				},
				{
					Address:     addressToWatch,
					BlockNumber: 2,
					BlockHash:   testutils.Random32Byte(),
					Data:        hexutil.MustDecode("0xdeadbeef"),
					TxHash:      testutils.Random32Byte(),
					Index:       1,
					TxIndex:     2,
					Topics:      []common.Hash{{1}},
					Removed:     false,
				},
			},
			[]bool{false, false},
			func(t *testing.T, lb *broadcaster) {
				// log pool should have both logs, as they are different
				require.Len(t, lb.logPool.logsByBlockHash, 2)
				// log pool should have both blockhashes, as they are different
				require.Len(t, lb.logPool.hashesByBlockNumbers, 2)
			},
		},
		{
			func() *broadcaster {
				lggr := logger.TestLogger(t)
				regs := newRegistrations(lggr, *big.NewInt(1337))
				listener := &mockListener{}
				listener.On("JobID").Return(int32(1))
				o := &mockORM{}
				o.On("SetPendingMinBlock", mock.AnythingOfType("*int64"), mock.Anything).Return(nil).Twice()
				sub := &subscriber{
					listener: listener,
					opts: ListenerOpts{
						Contract:                 addressToWatch,
						MinIncomingConfirmations: 1,
						LogsWithTopics: map[common.Hash][][]Topic{
							common.Hash{1}: {},
						},
					},
				}
				regs.addSubscriber(sub)
				return &broadcaster{
					logger:        lggr,
					logPool:       newLogPool(lggr),
					registrations: regs,
					orm:           o,
				}
			},
			"removed log, valid blockhash and block num",
			[]types.Log{
				{
					Address:     addressToWatch,
					BlockNumber: 1,
					BlockHash:   common.Hash{1, 2, 3},
					Data:        hexutil.MustDecode("0xdeadbeef"),
					TxHash:      common.Hash{1},
					Index:       0,
					TxIndex:     1,
					Topics:      []common.Hash{{1}},
					Removed:     false,
				},
				{
					Address:     addressToWatch,
					BlockNumber: 1,
					BlockHash:   common.Hash{1, 2, 3},
					Data:        hexutil.MustDecode("0xdeadbeef"),
					TxHash:      common.Hash{1},
					Index:       0,
					TxIndex:     1,
					Topics:      []common.Hash{{1}},
					Removed:     true,
				},
				{
					Address:     addressToWatch,
					BlockNumber: 1,
					BlockHash:   common.Hash{1, 2, 3, 4},
					Data:        hexutil.MustDecode("0xdeadbeef"),
					TxHash:      common.Hash{1},
					Index:       0,
					TxIndex:     1,
					Topics:      []common.Hash{{1}},
					Removed:     false,
				},
			},
			[]bool{false, false, false},
			func(t *testing.T, lb *broadcaster) {
				// log pool should only have a single log, which is the one after the removed
				// one
				require.Len(t, lb.logPool.logsByBlockHash, 1)
				// log pool should only have a single blockhash, which is the one after the
				// removed one
				require.Len(t, lb.logPool.hashesByBlockNumbers, 1)
			},
		},
		{
			func() *broadcaster {
				lggr := logger.TestLogger(t)
				regs := newRegistrations(lggr, *big.NewInt(1337))
				return &broadcaster{
					logger:        lggr,
					logPool:       newLogPool(lggr),
					registrations: regs,
				}
			},
			"unregistered address",
			[]types.Log{
				{
					Address:     addressToWatch,
					BlockNumber: 1,
					BlockHash:   common.Hash{1, 2, 3},
					Data:        hexutil.MustDecode("0xdeadbeef"),
					TxHash:      common.Hash{1},
					Index:       0,
					TxIndex:     1,
					Topics:      []common.Hash{{1}},
					Removed:     false,
				},
			},
			[]bool{false, false, false},
			func(t *testing.T, lb *broadcaster) {
				// log pool should have nothing since address is not registered
				require.Len(t, lb.logPool.logsByBlockHash, 0)
				require.Len(t, lb.logPool.hashesByBlockNumbers, 0)
			},
		},
		{
			func() *broadcaster {
				lggr := logger.TestLogger(t)
				regs := newRegistrations(lggr, *big.NewInt(1337))
				listener := &mockListener{}
				listener.On("JobID").Return(int32(1))
				o := &mockORM{}
				o.On("SetPendingMinBlock", mock.AnythingOfType("*int64"), mock.Anything).Return(nil).Once()
				sub := &subscriber{
					listener: listener,
					opts: ListenerOpts{
						Contract:                 addressToWatch,
						MinIncomingConfirmations: 1,
						LogsWithTopics: map[common.Hash][][]Topic{
							common.Hash{1}: {},
						},
					},
				}
				regs.addSubscriber(sub)
				return &broadcaster{
					logger:        lggr,
					logPool:       newLogPool(lggr),
					registrations: regs,
					orm:           o,
				}
			},
			"removed log, invalid blockhash and block num",
			[]types.Log{
				// original log, looks good
				{
					Address:     addressToWatch,
					BlockNumber: 1,
					BlockHash:   common.Hash{1, 2, 3},
					Data:        hexutil.MustDecode("0xdeadbeef"),
					TxHash:      common.Hash{1},
					Index:       0,
					TxIndex:     1,
					Topics:      []common.Hash{{1}},
					Removed:     false,
				},
				// removed log, zero blockhash
				{
					Address:     addressToWatch,
					BlockNumber: 0,
					BlockHash:   common.Hash{},
					Data:        hexutil.MustDecode("0xdeadbeef"),
					TxHash:      common.Hash{1},
					Index:       0,
					TxIndex:     0,
					Topics:      []common.Hash{{1}},
					Removed:     true,
				},
			},
			[]bool{false, true},
			func(t *testing.T, lb *broadcaster) {
				// log pool should only have a single log, but should be cleared after
				// we resubscribe
				require.Len(t, lb.logPool.logsByBlockHash, 1)
				require.Len(t, lb.logPool.hashesByBlockNumbers, 1)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			lb := tc.setup()
			for i := range tc.logs {
				shouldResubscribe := lb.onNewLog(tc.logs[i])
				require.Equal(tt, tc.shouldResubscribe[i], shouldResubscribe)
			}
			tc.assertions(tt, lb)
		})
	}
}

func TestBroadcaster_EventLoop(t *testing.T) {
	addressToWatch := common.Address{1}
	testCases := []struct {
		name       string
		setup      func() (lb *broadcaster)
		assertions func(t *testing.T, lb *broadcaster)
		logs       []types.Log
		expected   struct {
			shouldResubscribe bool
			err               error
		}
	}{
		{
			"raw log must resubscribe",
			func() (lb *broadcaster) {
				lggr := logger.TestLogger(t)
				regs := newRegistrations(lggr, *big.NewInt(1337))
				listener := &mockListener{}
				listener.On("JobID").Return(int32(1))
				o := &mockORM{}
				o.On("SetPendingMinBlock", mock.AnythingOfType("*int64"), mock.Anything).Return(nil).Once()
				sub := &subscriber{
					listener: listener,
					opts: ListenerOpts{
						Contract:                 addressToWatch,
						MinIncomingConfirmations: 1,
						LogsWithTopics: map[common.Hash][][]Topic{
							common.Hash{1}: {},
						},
					},
				}
				regs.addSubscriber(sub)
				return &broadcaster{
					logger:                 lggr,
					logPool:                newLogPool(lggr),
					registrations:          regs,
					orm:                    o,
					newHeads:               utils.NewHighCapacityMailbox[*evmtypes.Head](),
					changeSubscriberStatus: utils.NewHighCapacityMailbox[changeSubscriberStatus](),
				}
			},
			func(t *testing.T, lb *broadcaster) {

			},
			[]types.Log{
				{
					Address:     addressToWatch,
					BlockNumber: 0,
					BlockHash:   common.Hash{0},
					Data:        hexutil.MustDecode("0xdeadbeef"),
					TxHash:      common.Hash{1},
					Index:       0,
					TxIndex:     1,
					Topics:      []common.Hash{{1}},
					Removed:     true,
				},
			},
			struct {
				shouldResubscribe bool
				err               error
			}{
				true,
				errors.New("issue processing new log"),
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			lb := tc.setup()
			chRawLogs := make(chan types.Log, len(tc.logs))
			chErr := make(chan error, 1)
			for i := range tc.logs {
				chRawLogs <- tc.logs[i]
			}
			shouldResubscribe, err := lb.eventLoop(chRawLogs, chErr)
			require.Equal(tt, tc.expected.shouldResubscribe, shouldResubscribe)
			require.EqualError(tt, err, tc.expected.err.Error())
			tc.assertions(tt, lb)
		})
	}
}
