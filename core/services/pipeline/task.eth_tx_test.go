package pipeline_test

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	bptxmmocks "github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager/mocks"
	keystoremocks "github.com/smartcontractkit/chainlink/core/services/keystore/mocks"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"

	"gopkg.in/guregu/null.v4"
)

func TestETHTxTask(t *testing.T) {
	tests := []struct {
		name                  string
		from                  string
		to                    string
		data                  string
		gasLimit              string
		txMeta                string
		vars                  pipeline.Vars
		inputs                []pipeline.Result
		setupClientMocks      func(config *configtest.TestGeneralConfig, keyStore *keystoremocks.Eth, txManager *bptxmmocks.TxManager)
		expected              interface{}
		expectedErrorCause    error
		expectedErrorContains string
	}{
		{
			"happy (no vars)",
			`[ "0x882969652440ccf14a5dbb9bd53eb21cb1e11e5c" ]`,
			"0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF",
			"foobar",
			"12345",
			`{ "jobID": 321, "requestID": "0x5198616554d738d9485d1a7cf53b2f33e09c3bbc8fe9ac0020bd672cd2bc15d2", "requestTxHash": "0xc524fafafcaec40652b1f84fca09c231185437d008d195fccf2f51e64b7062f8" }`,
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
			nil, nil, "",
		},
		{
			"happy (with vars)",
			`[ $(fromAddr) ]`,
			"$(toAddr)",
			"$(data)",
			"$(gasLimit)",
			`{ "jobID": $(jobID), "requestID": $(requestID), "requestTxHash": $(requestTxHash) }`,
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
			nil, nil, "",
		},
		{
			"happy (with vars 2)",
			`$(fromAddrs)`,
			"$(toAddr)",
			"$(data)",
			"$(gasLimit)",
			`$(requestData)`,
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
			nil, nil, "",
		},
		{
			"happy (no `from`, keystore has key)",
			``,
			"$(toAddr)",
			"$(data)",
			"$(gasLimit)",
			`$(requestData)`,
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
			nil, nil, "",
		},
		{
			"happy (missing keys in txMeta)",
			`[ "0x882969652440ccf14a5dbb9bd53eb21cb1e11e5c" ]`,
			"0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF",
			"foobar",
			"12345",
			`{}`,
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
			nil, nil, "",
		},
		{
			"happy (missing gasLimit takes config default)",
			`[ "0x882969652440ccf14a5dbb9bd53eb21cb1e11e5c" ]`,
			"0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF",
			"foobar",
			"",
			`{ "jobID": 321, "requestID": "0x5198616554d738d9485d1a7cf53b2f33e09c3bbc8fe9ac0020bd672cd2bc15d2", "requestTxHash": "0xc524fafafcaec40652b1f84fca09c231185437d008d195fccf2f51e64b7062f8" }`,
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
			nil, nil, "",
		},
		{
			"error from keystore",
			``,
			"$(toAddr)",
			"$(data)",
			"$(gasLimit)",
			`$(requestData)`,
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
			nil, pipeline.ErrTaskRunFailed, "while querying keystore",
		},
		{
			"error from tx manager",
			`[ "0x882969652440ccf14a5dbb9bd53eb21cb1e11e5c" ]`,
			"0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF",
			"foobar",
			"12345",
			`{ "jobID": 321, "requestID": "0x5198616554d738d9485d1a7cf53b2f33e09c3bbc8fe9ac0020bd672cd2bc15d2", "requestTxHash": "0xc524fafafcaec40652b1f84fca09c231185437d008d195fccf2f51e64b7062f8" }`,
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
			nil, pipeline.ErrTaskRunFailed, "while creating transaction",
		},
		{
			"extra keys in txMeta",
			`[ "0x882969652440ccf14a5dbb9bd53eb21cb1e11e5c" ]`,
			"0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF",
			"foobar",
			"12345",
			`{ "jobID": 321, "requestID": "0x5198616554d738d9485d1a7cf53b2f33e09c3bbc8fe9ac0020bd672cd2bc15d2", "requestTxHash": "0xc524fafafcaec40652b1f84fca09c231185437d008d195fccf2f51e64b7062f8", "foo": "bar" }`,
			pipeline.NewVarsFrom(nil),
			nil,
			func(config *configtest.TestGeneralConfig, keyStore *keystoremocks.Eth, txManager *bptxmmocks.TxManager) {
				config.Overrides.GlobalEvmGasLimitDefault = null.IntFrom(999)
			},
			nil, pipeline.ErrBadInput, "txMeta",
		},
		{
			"bad values in txMeta",
			`[ "0x882969652440ccf14a5dbb9bd53eb21cb1e11e5c" ]`,
			"0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF",
			"foobar",
			"12345",
			`{ "jobID": "asdf", "requestID": 123, "requestTxHash": true }`,
			pipeline.NewVarsFrom(nil),
			nil,
			func(config *configtest.TestGeneralConfig, keyStore *keystoremocks.Eth, txManager *bptxmmocks.TxManager) {
				config.Overrides.GlobalEvmGasLimitDefault = null.IntFrom(999)
			},
			nil, pipeline.ErrBadInput, "txMeta",
		},
		{
			"missing `to`",
			`[ "0x882969652440ccf14a5dbb9bd53eb21cb1e11e5c" ]`,
			"",
			"foobar",
			"12345",
			`{ "jobID": 321, "requestID": "0x5198616554d738d9485d1a7cf53b2f33e09c3bbc8fe9ac0020bd672cd2bc15d2", "requestTxHash": "0xc524fafafcaec40652b1f84fca09c231185437d008d195fccf2f51e64b7062f8" }`,
			pipeline.NewVarsFrom(nil),
			nil,
			func(config *configtest.TestGeneralConfig, keyStore *keystoremocks.Eth, txManager *bptxmmocks.TxManager) {
				config.Overrides.GlobalEvmGasLimitDefault = null.IntFrom(999)
			},
			nil, pipeline.ErrParameterEmpty, "to",
		},
		{
			"errored input",
			`[ "0x882969652440ccf14a5dbb9bd53eb21cb1e11e5c" ]`,
			"0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF",
			"foobar",
			"12345",
			`{ "jobID": 321, "requestID": "0x5198616554d738d9485d1a7cf53b2f33e09c3bbc8fe9ac0020bd672cd2bc15d2", "requestTxHash": "0xc524fafafcaec40652b1f84fca09c231185437d008d195fccf2f51e64b7062f8" }`,
			pipeline.NewVarsFrom(nil),
			[]pipeline.Result{{Error: errors.New("uh oh")}},
			func(config *configtest.TestGeneralConfig, keyStore *keystoremocks.Eth, txManager *bptxmmocks.TxManager) {
			},
			nil, pipeline.ErrTooManyErrors, "task inputs",
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
				MinConfirmations: "0",
			}

			keyStore := new(keystoremocks.Eth)
			txManager := new(bptxmmocks.TxManager)
			db := pgtest.NewGormDB(t)
			cfg := configtest.NewTestGeneralConfig(t)

			cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: db, GeneralConfig: cfg, TxManager: txManager, KeyStore: keyStore})

			test.setupClientMocks(cfg, keyStore, txManager)
			task.HelperSetDependencies(db, cc, keyStore)

			result := task.Run(context.Background(), test.vars, test.inputs)

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
