package logevent_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	commonmocks "github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/triggers/logevent"
	coretestutils "github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

// Test for Log Event Trigger Capability happy path for EVM
func TestLogEventTriggerEVMHappyPath(t *testing.T) {
	th := testutils.NewContractReaderTH(t)

	logEventConfig := logevent.Config{
		ChainID:        th.BackendTH.ChainID.Uint64(),
		Network:        "evm",
		LookbackBlocks: 1000,
		PollPeriod:     1000,
	}

	// Create a new contract reader to return from mock relayer
	ctx := coretestutils.Context(t)

	// Fetch latest head from simulated backend to return from mock relayer
	height, err := th.BackendTH.EVMClient.LatestBlockHeight(ctx)
	require.NoError(t, err)
	block, err := th.BackendTH.EVMClient.BlockByNumber(ctx, height)
	require.NoError(t, err)

	// Mock relayer to return a New ContractReader instead of gRPC client of a ContractReader
	relayer := commonmocks.NewRelayer(t)
	relayer.On("NewContractReader", mock.Anything, th.LogEmitterContractReaderCfg).Return(th.LogEmitterContractReader, nil).Once()
	relayer.On("LatestHead", mock.Anything).Return(commontypes.Head{
		Height:    height.String(),
		Hash:      block.Hash().Bytes(),
		Timestamp: block.Time(),
	}, nil).Once()

	// Create Log Event Trigger Service and register trigger
	logEventTriggerService := logevent.NewLogEventTriggerService(logevent.Params{
		Logger:         th.BackendTH.Lggr,
		Relayer:        relayer,
		LogEventConfig: logEventConfig,
	})
	log1Ch, err := logEventTriggerService.RegisterTrigger(ctx, th.LogEmitterRegRequest)
	require.NoError(t, err)

	expectedLogVal := int64(10)

	// Send a blockchain transaction that emits logs
	go func() {
		_, err =
			th.LogEmitterContract.EmitLog1(th.BackendTH.ContractsOwner, []*big.Int{big.NewInt(expectedLogVal)})
		require.NoError(t, err)
		th.BackendTH.Backend.Commit()
		th.BackendTH.Backend.Commit()
	}()

	// Wait for logs with a timeout
	_, output, err := testutils.WaitForLog(th.BackendTH.Lggr, log1Ch, 5*time.Second)
	require.NoError(t, err)
	err = testutils.PrintMap(output, "EmitLog", th.BackendTH.Lggr)
	require.NoError(t, err)
	// Verify if valid cursor is returned
	cursor, err := testutils.GetStrVal(output, "Cursor")
	require.NoError(t, err)
	require.True(t, len(cursor) > 60)
	// Verify if Arg0 is correct
	actualLogVal, err := testutils.GetBigIntValL2(output, "Data", "Arg0")
	require.NoError(t, err)
	require.Equal(t, expectedLogVal, actualLogVal.Int64())

	logEventTriggerService.Close()
}
