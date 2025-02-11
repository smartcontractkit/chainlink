package pipeline_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-integrations/evm/client/clienttest"
	"github.com/smartcontractkit/chainlink/v2/core/chains"
	txmmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	keystoremocks "github.com/smartcontractkit/chainlink/v2/core/services/keystore/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	pipelinemocks "github.com/smartcontractkit/chainlink/v2/core/services/pipeline/mocks"
)

func TestETHCallTask(t *testing.T) {
	t.Parallel()
	testutils.SkipShortDB(t)

	var specGasLimit uint32 = 123
	const gasLimit uint64 = 500_000
	const drJobTypeGasLimit uint32 = 789

	tests := []struct {
		name                  string
		contract              string
		from                  string
		data                  string
		evmChainID            string
		gas                   string
		block                 string
		specGasLimit          *uint32
		vars                  pipeline.Vars
		inputs                []pipeline.Result
		setupClientMocks      func(ethClient *clienttest.Client, config *pipelinemocks.Config)
		expected              interface{}
		expectedErrorCause    error
		expectedErrorContains string
	}{
		{
			"happy with empty from",
			"0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF",
			"",
			"$(foo)",
			"0",
			"",
			"",
			nil,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": []byte("foo bar"),
			}),
			nil,
			func(ethClient *clienttest.Client, config *pipelinemocks.Config) {
				contractAddr := common.HexToAddress("0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF")
				ethClient.
					On("CallContract", mock.Anything, ethereum.CallMsg{To: &contractAddr, Gas: uint64(drJobTypeGasLimit), Data: []byte("foo bar")}, (*big.Int)(nil)).
					Return([]byte("baz quux"), nil)
			},
			[]byte("baz quux"), nil, "",
		},
		{
			"happy with gas limit per task",
			"0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF",
			"",
			"$(foo)",
			"0",
			"$(gasLimit)",
			"",
			nil,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo":      []byte("foo bar"),
				"gasLimit": 100_000,
			}),
			nil,
			func(ethClient *clienttest.Client, config *pipelinemocks.Config) {
				contractAddr := common.HexToAddress("0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF")
				ethClient.
					On("CallContract", mock.Anything, ethereum.CallMsg{To: &contractAddr, Gas: 100_000, Data: []byte("foo bar")}, (*big.Int)(nil)).
					Return([]byte("baz quux"), nil)
			},
			[]byte("baz quux"), nil, "",
		},
		{
			"happy with gas limit per spec",
			"0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF",
			"",
			"$(foo)",
			"0",
			"",
			"",
			&specGasLimit,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": []byte("foo bar"),
			}),
			nil,
			func(ethClient *clienttest.Client, config *pipelinemocks.Config) {
				contractAddr := common.HexToAddress("0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF")
				ethClient.
					On("CallContract", mock.Anything, ethereum.CallMsg{To: &contractAddr, Gas: uint64(specGasLimit), Data: []byte("foo bar")}, (*big.Int)(nil)).
					Return([]byte("baz quux"), nil)
			},
			[]byte("baz quux"), nil, "",
		},
		{
			"happy with from addr",
			"0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF",
			"0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF",
			"$(foo)",
			"0",
			"",
			"",
			nil,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": []byte("foo bar"),
			}),
			nil,
			func(ethClient *clienttest.Client, config *pipelinemocks.Config) {
				contractAddr := common.HexToAddress("0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF")
				fromAddr := common.HexToAddress("0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF")
				ethClient.
					On("CallContract", mock.Anything, ethereum.CallMsg{To: &contractAddr, Gas: uint64(drJobTypeGasLimit), From: fromAddr, Data: []byte("foo bar")}, (*big.Int)(nil)).
					Return([]byte("baz quux"), nil)
			},
			[]byte("baz quux"), nil, "",
		},
		{
			"bad from address",
			"0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF",
			"0xThisAintGonnaWork",
			"$(foo)",
			"0",
			"",
			"",
			nil,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": []byte("foo bar"),
			}),
			nil,
			func(ethClient *clienttest.Client, config *pipelinemocks.Config) {},
			nil, pipeline.ErrBadInput, "from",
		},
		{
			"bad contract address",
			"0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbee",
			"",
			"$(foo)",
			"0",
			"",
			"",
			nil,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": []byte("foo bar"),
			}),
			nil,
			func(ethClient *clienttest.Client, config *pipelinemocks.Config) {},
			nil, pipeline.ErrBadInput, "contract",
		},
		{
			"missing data var",
			"0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF",
			"",
			"$(foo)",
			"0",
			"",
			"",
			nil,
			pipeline.NewVarsFrom(map[string]interface{}{
				"zork": []byte("foo bar"),
			}),
			nil,
			func(ethClient *clienttest.Client, config *pipelinemocks.Config) {},
			nil, pipeline.ErrKeypathNotFound, "data",
		},
		{
			"no data",
			"0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF",
			"",
			"$(foo)",
			"0",
			"",
			"",
			nil,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": []byte(nil),
			}),
			nil,
			func(ethClient *clienttest.Client, config *pipelinemocks.Config) {},
			nil, pipeline.ErrBadInput, "data",
		},
		{
			"errored input",
			"0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF",
			"",
			"$(foo)",
			"0",
			"",
			"",
			nil,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": []byte("foo bar"),
			}),
			[]pipeline.Result{{Error: errors.New("uh oh")}},
			func(ethClient *clienttest.Client, config *pipelinemocks.Config) {},
			nil, pipeline.ErrTooManyErrors, "task inputs",
		},
		{
			"missing chainID",
			"0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF",
			"",
			"$(foo)",
			"$(evmChainID)",
			"",
			"",
			nil,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo":        []byte("foo bar"),
				"evmChainID": "123",
			}),
			nil,
			func(ethClient *clienttest.Client, config *pipelinemocks.Config) {
				contractAddr := common.HexToAddress("0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF")
				ethClient.
					On("CallContract", mock.Anything, ethereum.CallMsg{To: &contractAddr, Data: []byte("foo bar")}, (*big.Int)(nil)).
					Return([]byte("baz quux"), nil).Maybe()
			},
			nil, nil, chains.ErrNoSuchChainID.Error(),
		},
		{
			"simulate using latest block",
			"0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF",
			"",
			"$(foo)",
			"0",
			"",
			"latest",
			nil,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": []byte("foo bar"),
			}),
			nil,
			func(ethClient *clienttest.Client, config *pipelinemocks.Config) {
				contractAddr := common.HexToAddress("0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF")
				ethClient.
					On("CallContract", mock.Anything, ethereum.CallMsg{To: &contractAddr, Gas: uint64(drJobTypeGasLimit), Data: []byte("foo bar")}, (*big.Int)(nil)).
					Return([]byte("baz quux"), nil)
			},
			[]byte("baz quux"), nil, "",
		},
		{
			"simulate using pending block",
			"0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF",
			"",
			"$(foo)",
			"0",
			"",
			"pending",
			nil,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": []byte("foo bar"),
			}),
			nil,
			func(ethClient *clienttest.Client, config *pipelinemocks.Config) {
				contractAddr := common.HexToAddress("0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF")
				ethClient.
					On("PendingCallContract", mock.Anything, ethereum.CallMsg{To: &contractAddr, Gas: uint64(drJobTypeGasLimit), Data: []byte("foo bar")}).
					Return([]byte("baz quux"), nil)
			},
			[]byte("baz quux"), nil, "",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			task := pipeline.ETHCallTask{
				BaseTask:   pipeline.NewBaseTask(0, "ethcall", nil, nil, 0),
				Contract:   test.contract,
				From:       test.from,
				Data:       test.data,
				EVMChainID: test.evmChainID,
				Gas:        test.gas,
				Block:      test.block,
			}

			ethClient := clienttest.NewClient(t)
			config := pipelinemocks.NewConfig(t)
			test.setupClientMocks(ethClient, config)

			cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
				c.EVM[0].GasEstimator.LimitDefault = ptr(gasLimit)
				c.EVM[0].GasEstimator.LimitJobType.DR = ptr(drJobTypeGasLimit)
			})
			lggr := logger.TestLogger(t)

			keyStore := keystoremocks.NewEth(t)
			txManager := txmmocks.NewMockEvmTxManager(t)
			db := pgtest.NewSqlxDB(t)

			var legacyChains legacyevm.LegacyChainContainer
			if test.expectedErrorCause != nil || test.expectedErrorContains != "" {
				legacyChains = evmtest.NewLegacyChains(t, evmtest.TestChainOpts{
					DB:             db,
					ChainConfigs:   cfg.EVMConfigs(),
					DatabaseConfig: cfg.Database(),
					FeatureConfig:  cfg.Feature(),
					ListenerConfig: cfg.Database().Listener(),
					TxManager:      txManager,
					KeyStore:       keyStore,
				})
			} else {
				legacyChains = cltest.NewLegacyChainsWithMockChain(t, ethClient, cfg)
			}

			task.HelperSetDependencies(legacyChains, cfg.JobPipeline(), test.specGasLimit, pipeline.DirectRequestJobType)

			result, runInfo := task.Run(testutils.Context(t), lggr, test.vars, test.inputs)
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
