package client

import (
	"errors"
	"math"
	"testing"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	commonclient "github.com/smartcontractkit/chainlink/v2/common/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestRPCClient_MakeLogsValid(t *testing.T) {
	testCases := []struct {
		Name             string
		TxIndex          uint
		LogIndex         uint
		ExpectedLogIndex uint
		ExpectedError    error
	}{
		{
			Name:             "TxIndex = 0 LogIndex = 0",
			TxIndex:          0,
			LogIndex:         0,
			ExpectedLogIndex: 0,
			ExpectedError:    nil,
		},
		{
			Name:             "TxIndex = 0 LogIndex = 1",
			TxIndex:          0,
			LogIndex:         1,
			ExpectedLogIndex: 1,
			ExpectedError:    nil,
		},
		{
			Name:             "TxIndex = 0 LogIndex = MaxUint32",
			TxIndex:          0,
			LogIndex:         math.MaxUint32,
			ExpectedLogIndex: math.MaxUint32,
			ExpectedError:    nil,
		},
		{
			Name:             "LogIndex = MaxUint32 + 1 => returns an error",
			TxIndex:          0,
			LogIndex:         math.MaxUint32 + 1,
			ExpectedLogIndex: 0,
			ExpectedError:    errors.New("log's index 4294967296 of tx 0x0000000000000000000000000000000000000000000000000000000000000000 exceeds max supported value of 4294967295"),
		},
		{
			Name:             "TxIndex = 1 LogIndex = 0",
			TxIndex:          1,
			LogIndex:         0,
			ExpectedLogIndex: math.MaxUint32 + 1,
			ExpectedError:    nil,
		},
		{
			Name:             "TxIndex = MaxUint32 LogIndex = MaxUint32",
			TxIndex:          math.MaxUint32,
			LogIndex:         math.MaxUint32,
			ExpectedLogIndex: math.MaxUint64,
			ExpectedError:    nil,
		},
		{
			Name:             "TxIndex = MaxUint32 + 1 => returns an error",
			TxIndex:          math.MaxUint32 + 1,
			LogIndex:         0,
			ExpectedLogIndex: 0,
			ExpectedError:    errors.New("TxIndex of tx 0x0000000000000000000000000000000000000000000000000000000000000000 exceeds max supported value of 4294967295"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			rpc := newRPCClient(logger.TestLogger(t), nil, nil, "eth-primary-rpc-0", 0, nil, commonclient.Primary, 0, 0, commonclient.QueryTimeout, commonclient.QueryTimeout, "")
			log, err := rpc.makeLogValid(ethtypes.Log{TxIndex: tc.TxIndex, Index: tc.LogIndex})
			// non sei should return as is
			require.NoError(t, err)
			require.Equal(t, tc.TxIndex, log.TxIndex)
			require.Equal(t, tc.LogIndex, log.Index)
			seiRPC := newRPCClient(logger.TestLogger(t), nil, nil, "eth-primary-rpc-0", 0, nil, commonclient.Primary, 0, 0, commonclient.QueryTimeout, commonclient.QueryTimeout, chaintype.ChainSei)
			log, err = seiRPC.makeLogValid(ethtypes.Log{TxIndex: tc.TxIndex, Index: tc.LogIndex})
			if tc.ExpectedError != nil {
				require.EqualError(t, err, tc.ExpectedError.Error())
				return
			}

			require.Equal(t, tc.ExpectedLogIndex, log.Index)
			require.Equal(t, tc.TxIndex, log.TxIndex)
		})
	}
}
