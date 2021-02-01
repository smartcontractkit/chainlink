package log_test

import (
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/log"
	"github.com/smartcontractkit/chainlink/core/services/log/mocks"
)

func TestDecodingLogListener(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	contract, err := eth.GetV6ContractCodec("FluxAggregator")
	require.NoError(t, err)

	logTypes := map[common.Hash]interface{}{
		eth.MustGetV6ContractEventID("FluxAggregator", "NewRound"): &LogNewRound{},
	}

	var decodedLog interface{}

	listener := simpleLogListener{
		func(lb log.Broadcast, innerErr error) {
			err = innerErr
			decodedLog = lb.DecodedLog()
		},
		createJob(t, store).ID,
	}

	decodingListener := log.NewDecodingListener(contract, logTypes, &listener)
	rawLog := cltest.LogFromFixture(t, "../testdata/new_round_log.json")
	logBroadcast := new(mocks.Broadcast)

	logBroadcast.On("RawLog").Return(rawLog)
	logBroadcast.On("SetDecodedLog", mock.Anything).Run(func(args mock.Arguments) {
		logBroadcast.On("DecodedLog").Return(args.Get(0))
	})

	decodingListener.HandleLog(logBroadcast, nil)
	require.NoError(t, err)
	newRoundLog := decodedLog.(*LogNewRound)

	require.Equal(t, newRoundLog.Log, rawLog)
	require.True(t, newRoundLog.RoundId.Cmp(big.NewInt(1)) == 0)
	require.Equal(t, newRoundLog.StartedBy, common.HexToAddress("f17f52151ebef6c7334fad080c5704d77216b732"))
	require.True(t, newRoundLog.StartedAt.Cmp(big.NewInt(15)) == 0)

	expectedErr := errors.New("oh no!")
	nilLb := new(mocks.Broadcast)

	logBroadcast.On("Log").Return(nil).Once()
	decodingListener.HandleLog(nilLb, expectedErr)
	require.Equal(t, err, expectedErr)
}
