package logevent_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	commonmocks "github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/triggers/logevent"
	coretestutils "github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/capabilities/testutils"
)

// Test for Log Event Trigger Capability happy path for EVM
func TestLogEventTriggerEVMHappyPath(t *testing.T) {
	t.Parallel()
	th := testutils.NewContractReaderTH(t)

	logEventConfig := logevent.Config{
		ChainID:        th.BackendTH.ChainID.String(),
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
	logEventTriggerService, err := logevent.NewTriggerService(ctx,
		th.BackendTH.Lggr,
		relayer,
		logEventConfig)
	require.NoError(t, err)

	// Start the service
	servicetest.Run(t, logEventTriggerService)

	log1Ch, err := logEventTriggerService.RegisterTrigger(ctx, th.LogEmitterRegRequest)
	require.NoError(t, err)

	emitLogTxnAndWaitForLog(t, th, log1Ch, []*big.Int{big.NewInt(10)})
}

// Test if Log Event Trigger Capability is able to receive only new logs
// by using cursor and does not receive duplicate logs
func TestLogEventTriggerCursorNewLogs(t *testing.T) {
	t.Parallel()
	th := testutils.NewContractReaderTH(t)

	logEventConfig := logevent.Config{
		ChainID:        th.BackendTH.ChainID.String(),
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
	logEventTriggerService, err := logevent.NewTriggerService(ctx,
		th.BackendTH.Lggr,
		relayer,
		logEventConfig)
	require.NoError(t, err)

	// Start the service
	servicetest.Run(t, logEventTriggerService)

	log1Ch, err := logEventTriggerService.RegisterTrigger(ctx, th.LogEmitterRegRequest)
	require.NoError(t, err)

	emitLogTxnAndWaitForLog(t, th, log1Ch, []*big.Int{big.NewInt(10)})

	// This confirms that the cursor being tracked by log event trigger capability
	// works correctly and does not send old logs again as duplicates to the
	// callback channel log1Ch
	emitLogTxnAndWaitForLog(t, th, log1Ch, []*big.Int{big.NewInt(11), big.NewInt(12)})
}

// Send a transaction to EmitLog contract to emit Log1 events with given
// input parameters and wait for those logs to be received from relayer
// and ContractReader's QueryKey APIs used by Log Event Trigger
func emitLogTxnAndWaitForLog(t *testing.T,
	th *testutils.ContractReaderTH,
	log1Ch <-chan capabilities.TriggerResponse,
	expectedLogVals []*big.Int) {
	done := make(chan struct{})
	go func() {
		defer close(done)
		_, err := th.LogEmitterContract.EmitLog1(th.BackendTH.ContractsOwner, expectedLogVals)
		assert.NoError(t, err)
		th.BackendTH.Backend.Commit()
		th.BackendTH.Backend.Commit()
		th.BackendTH.Backend.Commit()
	}()

	for _, expectedLogVal := range expectedLogVals {
		// Wait for logs with a timeout
		_, output, err := testutils.WaitForLog(th.BackendTH.Lggr, log1Ch, tests.WaitTimeout(t))
		require.NoError(t, err)
		th.BackendTH.Lggr.Infow("EmitLog", "output", output)

		// Verify if valid cursor is returned
		cursor, err := testutils.GetStrVal(output, "Cursor")
		require.NoError(t, err)
		require.True(t, len(cursor) > 60)

		// Verify if Arg0 is correct
		actualLogVal, err := testutils.GetBigIntValL2(output, "Data", "Arg0")
		require.NoError(t, err)

		require.Equal(t, expectedLogVal.Int64(), actualLogVal.Int64())
	}

	<-done
}
