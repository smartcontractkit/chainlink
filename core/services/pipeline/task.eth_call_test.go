package pipeline_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	evmmocks "github.com/smartcontractkit/chainlink/core/chains/evm/mocks"
	txmmocks "github.com/smartcontractkit/chainlink/core/chains/evm/txmgr/mocks"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	keystoremocks "github.com/smartcontractkit/chainlink/core/services/keystore/mocks"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	pipelinemocks "github.com/smartcontractkit/chainlink/core/services/pipeline/mocks"
)

func TestETHCallTask(t *testing.T) {
	tests := []struct {
		name                  string
		contract              string
		data                  string
		evmChainID            string
		vars                  pipeline.Vars
		inputs                []pipeline.Result
		setupClientMocks      func(ethClient *evmmocks.Client, config *pipelinemocks.Config)
		expected              interface{}
		expectedErrorCause    error
		expectedErrorContains string
	}{
		{
			"happy",
			"0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF",
			"$(foo)",
			"",
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": []byte("foo bar"),
			}),
			nil,
			func(ethClient *evmmocks.Client, config *pipelinemocks.Config) {
				contractAddr := common.HexToAddress("0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF")
				ethClient.
					On("CallContract", mock.Anything, ethereum.CallMsg{To: &contractAddr, Data: []byte("foo bar")}, (*big.Int)(nil)).
					Return([]byte("baz quux"), nil)
			},
			[]byte("baz quux"), nil, "",
		},
		{
			"bad contract address",
			"0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbee",
			"$(foo)",
			"",
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": []byte("foo bar"),
			}),
			nil,
			func(ethClient *evmmocks.Client, config *pipelinemocks.Config) {},
			nil, pipeline.ErrBadInput, "contract",
		},
		{
			"missing data var",
			"0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF",
			"$(foo)",
			"",
			pipeline.NewVarsFrom(map[string]interface{}{
				"zork": []byte("foo bar"),
			}),
			nil,
			func(ethClient *evmmocks.Client, config *pipelinemocks.Config) {},
			nil, pipeline.ErrKeypathNotFound, "data",
		},
		{
			"no data",
			"0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF",
			"$(foo)",
			"",
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": []byte(nil),
			}),
			nil,
			func(ethClient *evmmocks.Client, config *pipelinemocks.Config) {},
			nil, pipeline.ErrBadInput, "data",
		},
		{
			"errored input",
			"0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF",
			"$(foo)",
			"",
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": []byte("foo bar"),
			}),
			[]pipeline.Result{{Error: errors.New("uh oh")}},
			func(ethClient *evmmocks.Client, config *pipelinemocks.Config) {},
			nil, pipeline.ErrTooManyErrors, "task inputs",
		},
		{
			"missing chainID",
			"0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF",
			"$(foo)",
			"$(evmChainID)",
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo":        []byte("foo bar"),
				"evmChainID": "123",
			}),
			nil,
			func(ethClient *evmmocks.Client, config *pipelinemocks.Config) {
				contractAddr := common.HexToAddress("0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF")
				ethClient.
					On("CallContract", mock.Anything, ethereum.CallMsg{To: &contractAddr, Data: []byte("foo bar")}, (*big.Int)(nil)).
					Return([]byte("baz quux"), nil)
			},
			nil, nil, "chain not found",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			task := pipeline.ETHCallTask{
				BaseTask:   pipeline.NewBaseTask(0, "ethcall", nil, nil, 0),
				Contract:   test.contract,
				Data:       test.data,
				EVMChainID: test.evmChainID,
			}

			ethClient := new(evmmocks.Client)
			config := new(pipelinemocks.Config)
			test.setupClientMocks(ethClient, config)

			cfg := configtest.NewTestGeneralConfig(t)

			keyStore := new(keystoremocks.Eth)
			keyStore.Test(t)
			txManager := new(txmmocks.TxManager)
			txManager.Test(t)
			db := pgtest.NewSqlxDB(t)

			var cc evm.ChainSet
			if test.expectedErrorCause != nil || test.expectedErrorContains != "" {
				cc = evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: db, GeneralConfig: cfg, TxManager: txManager, KeyStore: keyStore})
			} else {
				cc = cltest.NewChainSetMockWithOneChain(t, ethClient, evmtest.NewChainScopedConfig(t, cfg))
			}

			task.HelperSetDependencies(cc, cfg)

			result, runInfo := task.Run(context.Background(), logger.TestLogger(t), test.vars, test.inputs)
			assert.False(t, runInfo.IsPending)
			assert.False(t, runInfo.IsRetryable)

			if test.expectedErrorCause != nil || test.expectedErrorContains != "" {
				require.Nil(t, result.Value)
				if test.expectedErrorCause != nil {
					require.Equal(t, test.expectedErrorCause, errors.Cause(result.Error))
				}
				if test.expectedErrorContains != "" {
					require.Contains(t, result.Error.Error(), test.expectedErrorContains)
				}
			} else {
				require.NoError(t, result.Error)
				require.Equal(t, test.expected, result.Value)
			}
		})
	}
}
