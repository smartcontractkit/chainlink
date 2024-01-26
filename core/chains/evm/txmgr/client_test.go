package txmgr_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
)

func TestClient_BatchGetReceiptsWithFinalizedHeight(t *testing.T) {
	client := evmtest.NewEthClientMockWithDefaultChain(t)
	txmClient := txmgr.NewEvmTxmClient(client)
	testCases := []struct {
		Name string
		// inputs
		UseFinalityTag bool
		FinalityDepth  uint32

		// RPC response
		RPCErr     error
		BlockError error
		RPCHead    evmtypes.Head

		// Call Results
		ExpectedBlock *big.Int
		ExpectedErr   error
	}{
		{
			Name:           "returns error if call fails",
			UseFinalityTag: true,
			FinalityDepth:  10,

			RPCErr: errors.New("failed to call RPC"),

			ExpectedErr: errors.New("failed to call RPC"),
		},
		{
			Name: "returns error if fail to fetch block",

			BlockError: errors.New("failed to get bock"),

			ExpectedErr: errors.New("failed to get block"),
		},
		{
			Name:           "Returns block as is, if we are using finality tag",
			UseFinalityTag: true,
			FinalityDepth:  10,

			ExpectedBlock: big.NewInt(100),

			RPCHead: evmtypes.Head{Number: 100},
		},
		{
			Name:           "Subtracts finality depth if finality tag is disabled",
			UseFinalityTag: false,
			FinalityDepth:  10,

			ExpectedBlock: big.NewInt(90),

			RPCHead: evmtypes.Head{Number: 100},
		},
	}

	for _, testCase := range testCases {
		client.On("BatchCallContext", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			reqs := args.Get(1).([]rpc.BatchElem)
			blockReq := reqs[len(reqs)-1]
			blockReq.Error = testCase.BlockError
			reqHead := blockReq.Result.(*evmtypes.Head)
			*reqHead = testCase.RPCHead
		}).Return(testCase.ExpectedErr).Once()
		block, _, _, err := txmClient.BatchGetReceiptsWithFinalizedHeight(testutils.Context(t), nil, testCase.UseFinalityTag, testCase.FinalityDepth)
		if testCase.ExpectedErr != nil {
			assert.Error(t, testCase.ExpectedErr, err)
		} else {
			assert.NoError(t, err)
		}

		assert.Equal(t, testCase.ExpectedBlock, block)
	}
}
