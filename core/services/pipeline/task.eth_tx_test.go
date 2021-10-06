package pipeline_test

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	clnull "github.com/smartcontractkit/chainlink/core/null"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	bptxmmocks "github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager/mocks"
	keystoremocks "github.com/smartcontractkit/chainlink/core/services/keystore/mocks"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

func TestETHTxTask(t *testing.T) {
	tests := []struct {
		name                  string
		from                  string
		to                    string
		data                  string
		gasLimit              string
		txMeta                string
		minConfirmations      string
		vars                  pipeline.Vars
		inputs                []pipeline.Result
		setupClientMocks      func(config *configtest.TestGeneralConfig, keyStore *keystoremocks.Eth, txManager *bptxmmocks.TxManager)
		expected              interface{}
		expectedErrorCause    error
		expectedErrorContains string
		expectedRunInfo       pipeline.RunInfo
	}{
		{
			"happy (no vars)",
			`[ "0x882969652440ccf14a5dbb9bd53eb21cb1e11e5c" ]`,
			"0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF",
			"foobar",
			"12345",
			`{ "jobID": 321, "requestID": "0x5198616554d738d9485d1a7cf53b2f33e09c3bbc8fe9ac0020bd672cd2bc15d2", "requestTxHash": "0xc524fafafcaec40652b1f84fca09c231185437d008d195fccf2f51e64b7062f8" }`,
			`0`,
			pipeline.NewVarsFrom(nil),
			nil,
			func(config *configtest.TestGeneralConfig, keyStore *keystoremocks.Eth, txManager *bptxmmocks.TxManager) {
				config.Overrides.GlobalEvmGasLimitDefault = null.IntFrom(999)
				from := common.HexToAddress("0x882969652440ccf14a5dbb9bd53eb21cb1e11e5c")
				to := common.HexToAddress("0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF")
				data := []byte("foobar")
				gasLimit := uint64(12345)
				txMeta := &bulletprooftxmanager.EthTxMeta{JobID: 321, RequestID: common.HexToHash("0x5198616554d738d9485d1a7cf53b2f33e09c3bbc8fe9ac0020bd672cd2bc15d2"), RequestTxHash: common.HexToHash("0xc524fafafcaec40652b1f84fca09c231185437d008d195fccf2f51e64b7062f8")}
				keyStore.On("GetRoundRobinAddress", from).Return(from, nil)
				txManager.On("CreateEthTransaction", mock.Anything, bulletprooftxmanager.NewTx{
					FromAddress:    from,
					ToAddress:      to,
					EncodedPayload: data,
					GasLimit:       gasLimit,
					Meta:           txMeta,
					Strategy:       bulletprooftxmanager.SendEveryStrategy{},
				}).Return(bulletprooftxmanager.EthTx{}, nil)
			},
			nil, nil, "", pipeline.RunInfo{},
		},
		{
			"happy (with vars)",
			`[ $(fromAddr) ]`,
			"$(toAddr)",
			"$(data)",
			"$(gasLimit)",
			`{ "jobID": $(jobID), "requestID": $(requestID), "requestTxHash": $(requestTxHash) }`,
			`0`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"fromAddr":      common.HexToAddress("0x882969652440ccf14a5dbb9bd53eb21cb1e11e5c"),
				"toAddr":        common.HexToAddress("0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF"),
				"data":          []byte("foobar"),
				"gasLimit":      uint64(12345),
				"jobID":         int32(321),
				"requestID":     common.HexToHash("0x5198616554d738d9485d1a7cf53b2f33e09c3bbc8fe9ac0020bd672cd2bc15d2"),
				"requestTxHash": common.HexToHash("0xc524fafafcaec40652b1f84fca09c231185437d008d195fccf2f51e64b7062f8"),
			}),
			nil,
			func(config *configtest.TestGeneralConfig, keyStore *keystoremocks.Eth, txManager *bptxmmocks.TxManager) {
				config.Overrides.GlobalEvmGasLimitDefault = null.IntFrom(999)
				from := common.HexToAddress("0x882969652440ccf14a5dbb9bd53eb21cb1e11e5c")
				to := common.HexToAddress("0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF")
				data := []byte("foobar")
				gasLimit := uint64(12345)
				txMeta := &bulletprooftxmanager.EthTxMeta{JobID: 321, RequestID: common.HexToHash("0x5198616554d738d9485d1a7cf53b2f33e09c3bbc8fe9ac0020bd672cd2bc15d2"), RequestTxHash: common.HexToHash("0xc524fafafcaec40652b1f84fca09c231185437d008d195fccf2f51e64b7062f8")}
				keyStore.On("GetRoundRobinAddress", from).Return(from, nil)
				txManager.On("CreateEthTransaction", mock.Anything, bulletprooftxmanager.NewTx{
					FromAddress:    from,
					ToAddress:      to,
					EncodedPayload: data,
					GasLimit:       gasLimit,
					Meta:           txMeta,
					Strategy:       bulletprooftxmanager.SendEveryStrategy{},
				}).Return(bulletprooftxmanager.EthTx{}, nil)
			},
			nil, nil, "", pipeline.RunInfo{},
		},
		{
			"happy (with vars 2)",
			`$(fromAddrs)`,
			"$(toAddr)",
			"$(data)",
			"$(gasLimit)",
			`$(requestData)`,
			`0`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"fromAddrs": []common.Address{common.HexToAddress("0x882969652440ccf14a5dbb9bd53eb21cb1e11e5c")},
				"toAddr":    common.HexToAddress("0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF"),
				"data":      []byte("foobar"),
				"gasLimit":  uint64(12345),
				"requestData": map[string]interface{}{
					"jobID":         int32(321),
					"requestID":     common.HexToHash("0x5198616554d738d9485d1a7cf53b2f33e09c3bbc8fe9ac0020bd672cd2bc15d2"),
					"requestTxHash": common.HexToHash("0xc524fafafcaec40652b1f84fca09c231185437d008d195fccf2f51e64b7062f8"),
				},
			}),
			nil,
			func(config *configtest.TestGeneralConfig, keyStore *keystoremocks.Eth, txManager *bptxmmocks.TxManager) {
				config.Overrides.GlobalEvmGasLimitDefault = null.IntFrom(999)
				from := common.HexToAddress("0x882969652440ccf14a5dbb9bd53eb21cb1e11e5c")
				to := common.HexToAddress("0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF")
				data := []byte("foobar")
				gasLimit := uint64(12345)
				txMeta := &bulletprooftxmanager.EthTxMeta{JobID: 321, RequestID: common.HexToHash("0x5198616554d738d9485d1a7cf53b2f33e09c3bbc8fe9ac0020bd672cd2bc15d2"), RequestTxHash: common.HexToHash("0xc524fafafcaec40652b1f84fca09c231185437d008d195fccf2f51e64b7062f8")}
				keyStore.On("GetRoundRobinAddress", from).Return(from, nil)
				txManager.On("CreateEthTransaction", mock.Anything, bulletprooftxmanager.NewTx{
					FromAddress:    from,
					ToAddress:      to,
					EncodedPayload: data,
					GasLimit:       gasLimit,
					Meta:           txMeta,
					Strategy:       bulletprooftxmanager.SendEveryStrategy{},
				}).Return(bulletprooftxmanager.EthTx{}, nil)
			},
			nil, nil, "", pipeline.RunInfo{},
		},
		{
			"happy (no `from`, keystore has key)",
			``,
			"$(toAddr)",
			"$(data)",
			"$(gasLimit)",
			`$(requestData)`,
			`0`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"fromAddrs": []common.Address{common.HexToAddress("0x882969652440ccf14a5dbb9bd53eb21cb1e11e5c")},
				"toAddr":    common.HexToAddress("0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF"),
				"data":      []byte("foobar"),
				"gasLimit":  uint64(12345),
				"requestData": map[string]interface{}{
					"jobID":         int32(321),
					"requestID":     common.HexToHash("0x5198616554d738d9485d1a7cf53b2f33e09c3bbc8fe9ac0020bd672cd2bc15d2"),
					"requestTxHash": common.HexToHash("0xc524fafafcaec40652b1f84fca09c231185437d008d195fccf2f51e64b7062f8"),
				},
			}),
			nil,
			func(config *configtest.TestGeneralConfig, keyStore *keystoremocks.Eth, txManager *bptxmmocks.TxManager) {
				config.Overrides.GlobalEvmGasLimitDefault = null.IntFrom(999)
				from := common.HexToAddress("0x882969652440ccf14a5dbb9bd53eb21cb1e11e5c")
				to := common.HexToAddress("0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF")
				data := []byte("foobar")
				gasLimit := uint64(12345)
				txMeta := &bulletprooftxmanager.EthTxMeta{JobID: 321, RequestID: common.HexToHash("0x5198616554d738d9485d1a7cf53b2f33e09c3bbc8fe9ac0020bd672cd2bc15d2"), RequestTxHash: common.HexToHash("0xc524fafafcaec40652b1f84fca09c231185437d008d195fccf2f51e64b7062f8")}
				keyStore.On("GetRoundRobinAddress").Return(from, nil)
				txManager.On("CreateEthTransaction", mock.Anything, bulletprooftxmanager.NewTx{
					FromAddress:    from,
					ToAddress:      to,
					EncodedPayload: data,
					GasLimit:       gasLimit,
					Meta:           txMeta,
					Strategy:       bulletprooftxmanager.SendEveryStrategy{},
				}).Return(bulletprooftxmanager.EthTx{}, nil)
			},
			nil, nil, "", pipeline.RunInfo{},
		},
		{
			"happy (missing keys in txMeta)",
			`[ "0x882969652440ccf14a5dbb9bd53eb21cb1e11e5c" ]`,
			"0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF",
			"foobar",
			"12345",
			`{}`,
			`0`,
			pipeline.NewVarsFrom(nil),
			nil,
			func(config *configtest.TestGeneralConfig, keyStore *keystoremocks.Eth, txManager *bptxmmocks.TxManager) {
				config.Overrides.GlobalEvmGasLimitDefault = null.IntFrom(999)
				from := common.HexToAddress("0x882969652440ccf14a5dbb9bd53eb21cb1e11e5c")
				to := common.HexToAddress("0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF")
				data := []byte("foobar")
				gasLimit := uint64(12345)
				txMeta := &bulletprooftxmanager.EthTxMeta{}
				keyStore.On("GetRoundRobinAddress", from).Return(from, nil)
				txManager.On("CreateEthTransaction", mock.Anything, bulletprooftxmanager.NewTx{
					FromAddress:    from,
					ToAddress:      to,
					EncodedPayload: data,
					GasLimit:       gasLimit,
					Meta:           txMeta,
					Strategy:       bulletprooftxmanager.SendEveryStrategy{},
				}).Return(bulletprooftxmanager.EthTx{}, nil)
			},
			nil, nil, "", pipeline.RunInfo{},
		},
		{
			"happy (missing gasLimit takes config default)",
			`[ "0x882969652440ccf14a5dbb9bd53eb21cb1e11e5c" ]`,
			"0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF",
			"foobar",
			"",
			`{ "jobID": 321, "requestID": "0x5198616554d738d9485d1a7cf53b2f33e09c3bbc8fe9ac0020bd672cd2bc15d2", "requestTxHash": "0xc524fafafcaec40652b1f84fca09c231185437d008d195fccf2f51e64b7062f8" }`,
			`0`,
			pipeline.NewVarsFrom(nil),
			nil,
			func(config *configtest.TestGeneralConfig, keyStore *keystoremocks.Eth, txManager *bptxmmocks.TxManager) {
				config.Overrides.GlobalEvmGasLimitDefault = null.IntFrom(999)
				from := common.HexToAddress("0x882969652440ccf14a5dbb9bd53eb21cb1e11e5c")
				to := common.HexToAddress("0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF")
				data := []byte("foobar")
				gasLimit := uint64(999)
				txMeta := &bulletprooftxmanager.EthTxMeta{JobID: 321, RequestID: common.HexToHash("0x5198616554d738d9485d1a7cf53b2f33e09c3bbc8fe9ac0020bd672cd2bc15d2"), RequestTxHash: common.HexToHash("0xc524fafafcaec40652b1f84fca09c231185437d008d195fccf2f51e64b7062f8")}
				keyStore.On("GetRoundRobinAddress", from).Return(from, nil)
				txManager.On("CreateEthTransaction", mock.Anything, bulletprooftxmanager.NewTx{
					FromAddress:    from,
					ToAddress:      to,
					EncodedPayload: data,
					GasLimit:       gasLimit,
					Meta:           txMeta,
					Strategy:       bulletprooftxmanager.SendEveryStrategy{},
				}).Return(bulletprooftxmanager.EthTx{}, nil)
			},
			nil, nil, "", pipeline.RunInfo{},
		},
		{
			"error from keystore",
			``,
			"$(toAddr)",
			"$(data)",
			"$(gasLimit)",
			`$(requestData)`,
			`0`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"fromAddrs": []common.Address{common.HexToAddress("0x882969652440ccf14a5dbb9bd53eb21cb1e11e5c")},
				"toAddr":    common.HexToAddress("0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF"),
				"data":      []byte("foobar"),
				"gasLimit":  uint64(12345),
				"requestData": map[string]interface{}{
					"jobID":         int32(321),
					"requestID":     common.HexToHash("0x5198616554d738d9485d1a7cf53b2f33e09c3bbc8fe9ac0020bd672cd2bc15d2"),
					"requestTxHash": common.HexToHash("0xc524fafafcaec40652b1f84fca09c231185437d008d195fccf2f51e64b7062f8"),
				},
			}),
			nil,
			func(config *configtest.TestGeneralConfig, keyStore *keystoremocks.Eth, txManager *bptxmmocks.TxManager) {
				config.Overrides.GlobalEvmGasLimitDefault = null.IntFrom(999)
				keyStore.On("GetRoundRobinAddress").Return(nil, errors.New("uh oh"))
			},
			nil, pipeline.ErrTaskRunFailed, "while querying keystore", pipeline.RunInfo{IsRetryable: true},
		},
		{
			"error from tx manager",
			`[ "0x882969652440ccf14a5dbb9bd53eb21cb1e11e5c" ]`,
			"0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF",
			"foobar",
			"12345",
			`{ "jobID": 321, "requestID": "0x5198616554d738d9485d1a7cf53b2f33e09c3bbc8fe9ac0020bd672cd2bc15d2", "requestTxHash": "0xc524fafafcaec40652b1f84fca09c231185437d008d195fccf2f51e64b7062f8" }`,
			`0`,
			pipeline.NewVarsFrom(nil),
			nil,
			func(config *configtest.TestGeneralConfig, keyStore *keystoremocks.Eth, txManager *bptxmmocks.TxManager) {
				config.Overrides.GlobalEvmGasLimitDefault = null.IntFrom(999)
				from := common.HexToAddress("0x882969652440ccf14a5dbb9bd53eb21cb1e11e5c")
				to := common.HexToAddress("0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF")
				data := []byte("foobar")
				gasLimit := uint64(12345)
				txMeta := &bulletprooftxmanager.EthTxMeta{JobID: 321, RequestID: common.HexToHash("0x5198616554d738d9485d1a7cf53b2f33e09c3bbc8fe9ac0020bd672cd2bc15d2"), RequestTxHash: common.HexToHash("0xc524fafafcaec40652b1f84fca09c231185437d008d195fccf2f51e64b7062f8")}
				keyStore.On("GetRoundRobinAddress", from).Return(from, nil)
				txManager.On("CreateEthTransaction", mock.Anything, bulletprooftxmanager.NewTx{
					FromAddress:    from,
					ToAddress:      to,
					EncodedPayload: data,
					GasLimit:       gasLimit,
					Meta:           txMeta,
					Strategy:       bulletprooftxmanager.SendEveryStrategy{},
				}).Return(bulletprooftxmanager.EthTx{}, errors.New("uh oh"))
			},
			nil, pipeline.ErrTaskRunFailed, "while creating transaction", pipeline.RunInfo{IsRetryable: true},
		},
		{
			"extra keys in txMeta",
			`[ "0x882969652440ccf14a5dbb9bd53eb21cb1e11e5c" ]`,
			"0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF",
			"foobar",
			"12345",
			`{ "jobID": 321, "requestID": "0x5198616554d738d9485d1a7cf53b2f33e09c3bbc8fe9ac0020bd672cd2bc15d2", "requestTxHash": "0xc524fafafcaec40652b1f84fca09c231185437d008d195fccf2f51e64b7062f8", "foo": "bar" }`,
			`0`,
			pipeline.NewVarsFrom(nil),
			nil,
			func(config *configtest.TestGeneralConfig, keyStore *keystoremocks.Eth, txManager *bptxmmocks.TxManager) {
				config.Overrides.GlobalEvmGasLimitDefault = null.IntFrom(999)
			},
			nil, pipeline.ErrBadInput, "txMeta", pipeline.RunInfo{},
		},
		{
			"bad values in txMeta",
			`[ "0x882969652440ccf14a5dbb9bd53eb21cb1e11e5c" ]`,
			"0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF",
			"foobar",
			"12345",
			`{ "jobID": "asdf", "requestID": 123, "requestTxHash": true }`,
			`0`,
			pipeline.NewVarsFrom(nil),
			nil,
			func(config *configtest.TestGeneralConfig, keyStore *keystoremocks.Eth, txManager *bptxmmocks.TxManager) {
				config.Overrides.GlobalEvmGasLimitDefault = null.IntFrom(999)
			},
			nil, pipeline.ErrBadInput, "txMeta", pipeline.RunInfo{},
		},
		{
			"missing `to`",
			`[ "0x882969652440ccf14a5dbb9bd53eb21cb1e11e5c" ]`,
			"",
			"foobar",
			"12345",
			`{ "jobID": 321, "requestID": "0x5198616554d738d9485d1a7cf53b2f33e09c3bbc8fe9ac0020bd672cd2bc15d2", "requestTxHash": "0xc524fafafcaec40652b1f84fca09c231185437d008d195fccf2f51e64b7062f8" }`,
			`0`,
			pipeline.NewVarsFrom(nil),
			nil,
			func(config *configtest.TestGeneralConfig, keyStore *keystoremocks.Eth, txManager *bptxmmocks.TxManager) {
				config.Overrides.GlobalEvmGasLimitDefault = null.IntFrom(999)
			},
			nil, pipeline.ErrParameterEmpty, "to", pipeline.RunInfo{},
		},
		{
			"errored input",
			`[ "0x882969652440ccf14a5dbb9bd53eb21cb1e11e5c" ]`,
			"0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF",
			"foobar",
			"12345",
			`{ "jobID": 321, "requestID": "0x5198616554d738d9485d1a7cf53b2f33e09c3bbc8fe9ac0020bd672cd2bc15d2", "requestTxHash": "0xc524fafafcaec40652b1f84fca09c231185437d008d195fccf2f51e64b7062f8" }`,
			`0`,
			pipeline.NewVarsFrom(nil),
			[]pipeline.Result{{Error: errors.New("uh oh")}},
			func(config *configtest.TestGeneralConfig, keyStore *keystoremocks.Eth, txManager *bptxmmocks.TxManager) {
			},
			nil, pipeline.ErrTooManyErrors, "task inputs", pipeline.RunInfo{},
		},
		{
			"async mode (with > 0 minConfirmations)",
			`[ "0x882969652440ccf14a5dbb9bd53eb21cb1e11e5c" ]`,
			"0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF",
			"foobar",
			"12345",
			`{ "jobID": 321, "requestID": "0x5198616554d738d9485d1a7cf53b2f33e09c3bbc8fe9ac0020bd672cd2bc15d2", "requestTxHash": "0xc524fafafcaec40652b1f84fca09c231185437d008d195fccf2f51e64b7062f8" }`,
			`3`,
			pipeline.NewVarsFrom(nil),
			nil,
			func(config *configtest.TestGeneralConfig, keyStore *keystoremocks.Eth, txManager *bptxmmocks.TxManager) {
				config.Overrides.GlobalEvmGasLimitDefault = null.IntFrom(999)
				from := common.HexToAddress("0x882969652440ccf14a5dbb9bd53eb21cb1e11e5c")
				keyStore.On("GetRoundRobinAddress", from).Return(from, nil)
				txManager.On("CreateEthTransaction", mock.Anything, mock.MatchedBy(func(tx bulletprooftxmanager.NewTx) bool {
					return tx.MinConfirmations == clnull.Uint32From(3) && tx.PipelineTaskRunID != nil
				})).Return(bulletprooftxmanager.EthTx{}, nil)
			},
			nil, nil, "", pipeline.RunInfo{IsPending: true},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			task := pipeline.ETHTxTask{
				BaseTask:         pipeline.NewBaseTask(0, "ethtx", nil, nil, 0),
				From:             test.from,
				To:               test.to,
				Data:             test.data,
				GasLimit:         test.gasLimit,
				TxMeta:           test.txMeta,
				MinConfirmations: test.minConfirmations,
			}

			keyStore := new(keystoremocks.Eth)
			keyStore.Test(t)
			txManager := new(bptxmmocks.TxManager)
			txManager.Test(t)
			db := pgtest.NewGormDB(t)
			cfg := configtest.NewTestGeneralConfig(t)

			cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: db, GeneralConfig: cfg, TxManager: txManager, KeyStore: keyStore})

			test.setupClientMocks(cfg, keyStore, txManager)
			task.HelperSetDependencies(db, cc, keyStore)

			result, runInfo := task.Run(context.Background(), test.vars, test.inputs)
			assert.Equal(t, test.expectedRunInfo, runInfo)

			if test.expectedErrorCause != nil {
				require.Equal(t, test.expectedErrorCause, errors.Cause(result.Error))
				require.Nil(t, result.Value)
				if test.expectedErrorContains != "" {
					require.Contains(t, result.Error.Error(), test.expectedErrorContains)
				}
			} else {
				require.NoError(t, result.Error)
				require.Equal(t, test.expected, result.Value)
			}

			keyStore.AssertExpectations(t)
			txManager.AssertExpectations(t)
		})
	}
}
