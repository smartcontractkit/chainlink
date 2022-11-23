package ocrcommon_test

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	txmmocks "github.com/smartcontractkit/chainlink/core/chains/evm/txmgr/mocks"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	pipelinemocks "github.com/smartcontractkit/chainlink/core/services/pipeline/mocks"
)

func Test_PipelineTransmitter_CreateEthTransaction(t *testing.T) {
	t.Parallel()

	lggr := logger.TestLogger(t)
	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore, 0)

	chainID := "12345"
	gasLimit := uint32(1000)
	effectiveTransmitterAddress := fromAddress
	toAddress := testutils.NewAddress()
	payload := []byte{1, 2, 3}
	strategy := txmmocks.NewTxStrategy(t)
	checker := txmgr.TransmitCheckerSpec{CheckerType: txmgr.TransmitCheckerTypeSimulate}
	runner := new(pipelinemocks.Runner)

	transmitter := ocrcommon.NewPipelineTransmitter(
		lggr,
		fromAddress,
		gasLimit,
		effectiveTransmitterAddress,
		strategy,
		checker,
		runner,
		job.Job{
			PipelineSpec: &pipeline.Spec{},
		},
		chainID,
	)

	runner.On("Run", mock.Anything, mock.AnythingOfType("*pipeline.Run"), mock.Anything, mock.Anything, mock.Anything).
		Return(false, nil).
		Run(func(args mock.Arguments) {
			run := args.Get(1).(*pipeline.Run)
			require.Equal(t, map[string]interface{}{
				"jobSpec": map[string]interface{}{
					"contractAddress":   toAddress.String(),
					"fromAddress":       fromAddress.String(),
					"gasLimit":          gasLimit,
					"evmChainID":        chainID,
					"forwardingAllowed": false,
					"data":              payload,
					"transmitChecker":   checker,
				},
			}, run.Inputs.Val)

			save := args.Get(3).(bool)
			require.True(t, save)

			run.State = pipeline.RunStatusCompleted
		}).Once()

	require.NoError(t, transmitter.CreateEthTransaction(testutils.Context(t), toAddress, payload))
}
