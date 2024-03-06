package client_test

import (
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	commonclient "github.com/smartcontractkit/chainlink/v2/common/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

func newMockRpc(t *testing.T) *mocks.RPCClient {
	mockRpc := mocks.NewRPCClient(t)
	mockRpc.On("Dial", mock.Anything).Return(nil).Once()
	mockRpc.On("Close").Return(nil).Once()
	mockRpc.On("ChainID", mock.Anything).Return(testutils.FixtureChainID, nil).Once()
	mockRpc.On("Subscribe", mock.Anything, mock.Anything, mock.Anything).Return(client.NewMockSubscription(), nil).Once()
	mockRpc.On("SetAliveLoopSub", mock.Anything).Return().Once()
	return mockRpc
}

func TestChainClient_BatchCallContext(t *testing.T) {
	t.Parallel()

	t.Run("batch requests return errors", func(t *testing.T) {
		ctx := testutils.Context(t)
		rpcError := errors.New("something went wrong")
		blockNumResp := ""
		blockNum := hexutil.EncodeBig(big.NewInt(42))
		b := []rpc.BatchElem{
			{
				Method: "eth_getBlockByNumber",
				Args:   []interface{}{blockNum, true},
				Result: &types.Block{},
			},
			{
				Method: "eth_blockNumber",
				Result: &blockNumResp,
			},
		}

		mockRpc := newMockRpc(t)
		mockRpc.On("BatchCallContext", mock.Anything, b).Run(func(args mock.Arguments) {
			reqs := args.Get(1).([]rpc.BatchElem)
			for i := 0; i < len(reqs); i++ {
				elem := &reqs[i]
				elem.Error = rpcError
			}
		}).Return(nil).Once()

		client := client.NewChainClientWithMockedRpc(t, commonclient.NodeSelectionModeRoundRobin, time.Second*0, time.Second*0, testutils.FixtureChainID, mockRpc)
		err := client.Dial(ctx)
		require.NoError(t, err)

		err = client.BatchCallContext(ctx, b)
		require.NoError(t, err)
		for _, elem := range b {
			require.ErrorIs(t, rpcError, elem.Error)
		}
	})
}
