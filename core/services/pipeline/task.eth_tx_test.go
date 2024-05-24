package pipeline_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	clnull "github.com/smartcontractkit/chainlink-common/pkg/utils/null"
	txmgrcommon "github.com/smartcontractkit/chainlink/v2/common/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/chains"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	txmmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	keystoremocks "github.com/smartcontractkit/chainlink/v2/core/services/keystore/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	evmrelay "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
)

func TestETHTxTask(t *testing.T) {
	jid := int32(321)
	reqID := common.HexToHash("0x5198616554d738d9485d1a7cf53b2f33e09c3bbc8fe9ac0020bd672cd2bc15d2")
	reqTxHash := common.HexToHash("0xc524fafafcaec40652b1f84fca09c231185437d008d195fccf2f51e64b7062f8")
	specGasLimit := uint32(123)
	const defaultGasLimit uint64 = 999
	const drJobTypeGasLimit uint32 = 789

	from := common.HexToAddress("0x882969652440ccf14a5dbb9bd53eb21cb1e11e5c")
	to := common.HexToAddress("0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF")

	tests := []struct {
		name                  string
		from                  string
		to                    string
		data                  string
		gasLimit              string
		txMeta                string
		minConfirmations      string
		evmChainID            string
		transmitChecker       string
		specGasLimit          *uint32
		forwardingAllowed     bool
		vars                  pipeline.Vars
		inputs                []pipeline.Result
		setupClientMocks      func(keyStore *keystoremocks.Eth, txManager *txmmocks.MockEvmTxManager)
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
			"0",
			`{"CheckerType": "vrf_v2", "VRFCoordinatorAddress": "0x2E396ecbc8223Ebc16EC45136228AE5EDB649943"}`,
			nil,
			false,
			pipeline.NewVarsFrom(nil),
			nil,
			func(keyStore *keystoremocks.Eth, txManager *txmmocks.MockEvmTxManager) {
				data := []byte("foobar")
				gasLimit := uint64(12345)
				jobID := int32(321)
				addr := common.HexToAddress("0x2E396ecbc8223Ebc16EC45136228AE5EDB649943")
				txMeta := &txmgr.TxMeta{
					JobID:         &jobID,
					RequestID:     &reqID,
					RequestTxHash: &reqTxHash,
					FailOnRevert:  null.BoolFrom(false),
				}
				keyStore.On("GetRoundRobinAddress", mock.Anything, testutils.FixtureChainID, from).Return(from, nil)
				txManager.On("CreateTransaction", mock.Anything, txmgr.TxRequest{
					FromAddress:    from,
					ToAddress:      to,
					EncodedPayload: data,
					FeeLimit:       gasLimit,
					Meta:           txMeta,
					Strategy:       txmgrcommon.NewSendEveryStrategy(),
					Checker: txmgr.TransmitCheckerSpec{
						CheckerType:           txmgr.TransmitCheckerTypeVRFV2,
						VRFCoordinatorAddress: &addr,
					},
					SignalCallback: true,
				}).Return(txmgr.Tx{}, nil)
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
			"0",
			"",
			nil,
			false,
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
			func(keyStore *keystoremocks.Eth, txManager *txmmocks.MockEvmTxManager) {
				data := []byte("foobar")
				gasLimit := uint64(12345)
				txMeta := &txmgr.TxMeta{
					JobID:         &jid,
					RequestID:     &reqID,
					RequestTxHash: &reqTxHash,
					FailOnRevert:  null.BoolFrom(false),
				}
				keyStore.On("GetRoundRobinAddress", mock.Anything, testutils.FixtureChainID, from).Return(from, nil)
				txManager.On("CreateTransaction", mock.Anything, txmgr.TxRequest{
					FromAddress:    from,
					ToAddress:      to,
					EncodedPayload: data,
					FeeLimit:       gasLimit,
					Meta:           txMeta,
					Strategy:       txmgrcommon.NewSendEveryStrategy(),
					SignalCallback: true,
				}).Return(txmgr.Tx{}, nil)
			},
			nil, nil, "", pipeline.RunInfo{},
		},
		{
			"happy (with minConfirmations as variable expression)",
			`[ $(fromAddr) ]`,
			"$(toAddr)",
			"$(data)",
			"$(gasLimit)",
			`{ "jobID": $(jobID), "requestID": $(requestID), "requestTxHash": $(requestTxHash) }`,
			"$(minConfirmations)",
			"0",
			"",
			nil,
			false,
			pipeline.NewVarsFrom(map[string]interface{}{
				"fromAddr":         common.HexToAddress("0x882969652440ccf14a5dbb9bd53eb21cb1e11e5c"),
				"toAddr":           common.HexToAddress("0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF"),
				"data":             []byte("foobar"),
				"gasLimit":         uint64(12345),
				"minConfirmations": uint64(2),
				"jobID":            int32(321),
				"requestID":        common.HexToHash("0x5198616554d738d9485d1a7cf53b2f33e09c3bbc8fe9ac0020bd672cd2bc15d2"),
				"requestTxHash":    common.HexToHash("0xc524fafafcaec40652b1f84fca09c231185437d008d195fccf2f51e64b7062f8"),
			}),
			nil,
			func(keyStore *keystoremocks.Eth, txManager *txmmocks.MockEvmTxManager) {
				addr := common.HexToAddress("0x882969652440ccf14a5dbb9bd53eb21cb1e11e5c")
				keyStore.On("GetRoundRobinAddress", mock.Anything, testutils.FixtureChainID, addr).Return(addr, nil)
				txManager.On("CreateTransaction", mock.Anything, mock.MatchedBy(func(tx txmgr.TxRequest) bool {
					return tx.MinConfirmations == clnull.Uint32From(2)
				})).Return(txmgr.Tx{}, nil)
			},
			nil, nil, "", pipeline.RunInfo{IsPending: true},
		},
		{
			"happy (with vars 2)",
			`$(fromAddrs)`,
			"$(toAddr)",
			"$(data)",
			"$(gasLimit)",
			`$(requestData)`,
			`0`,
			"0",
			"",
			nil,
			false,
			pipeline.NewVarsFrom(map[string]interface{}{
				"fromAddrs": []common.Address{common.HexToAddress("0x882969652440ccf14a5dbb9bd53eb21cb1e11e5c")},
				"toAddr":    "0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF",
				"data":      []byte("foobar"),
				"gasLimit":  uint32(12345),
				"requestData": map[string]interface{}{
					"jobID":         int32(321),
					"requestID":     common.HexToHash("0x5198616554d738d9485d1a7cf53b2f33e09c3bbc8fe9ac0020bd672cd2bc15d2"),
					"requestTxHash": common.HexToHash("0xc524fafafcaec40652b1f84fca09c231185437d008d195fccf2f51e64b7062f8"),
				},
			}),
			nil,
			func(keyStore *keystoremocks.Eth, txManager *txmmocks.MockEvmTxManager) {
				data := []byte("foobar")
				gasLimit := uint64(12345)
				txMeta := &txmgr.TxMeta{
					JobID:         &jid,
					RequestID:     &reqID,
					RequestTxHash: &reqTxHash,
					FailOnRevert:  null.BoolFrom(false),
				}
				keyStore.On("GetRoundRobinAddress", mock.Anything, testutils.FixtureChainID, from).Return(from, nil)
				txManager.On("CreateTransaction", mock.Anything, txmgr.TxRequest{
					FromAddress:    from,
					ToAddress:      to,
					EncodedPayload: data,
					FeeLimit:       gasLimit,
					Meta:           txMeta,
					Strategy:       txmgrcommon.NewSendEveryStrategy(),
					SignalCallback: true,
				}).Return(txmgr.Tx{}, nil)
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
			"0",
			"",
			nil,
			false,
			pipeline.NewVarsFrom(map[string]interface{}{
				"fromAddrs": []common.Address{common.HexToAddress("0x882969652440ccf14a5dbb9bd53eb21cb1e11e5c")},
				"toAddr":    common.HexToAddress("0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF"),
				"data":      []byte("foobar"),
				"gasLimit":  uint32(12345),
				"requestData": map[string]interface{}{
					"jobID":         int32(321),
					"requestID":     common.HexToHash("0x5198616554d738d9485d1a7cf53b2f33e09c3bbc8fe9ac0020bd672cd2bc15d2"),
					"requestTxHash": common.HexToHash("0xc524fafafcaec40652b1f84fca09c231185437d008d195fccf2f51e64b7062f8"),
				},
			}),
			nil,
			func(keyStore *keystoremocks.Eth, txManager *txmmocks.MockEvmTxManager) {
				data := []byte("foobar")
				gasLimit := uint64(12345)
				txMeta := &txmgr.TxMeta{
					JobID:         &jid,
					RequestID:     &reqID,
					RequestTxHash: &reqTxHash,
					FailOnRevert:  null.BoolFrom(false),
				}
				keyStore.On("GetRoundRobinAddress", mock.Anything, testutils.FixtureChainID).Return(from, nil)
				txManager.On("CreateTransaction", mock.Anything, txmgr.TxRequest{
					FromAddress:    from,
					ToAddress:      to,
					EncodedPayload: data,
					FeeLimit:       gasLimit,
					Meta:           txMeta,
					Strategy:       txmgrcommon.NewSendEveryStrategy(),
					SignalCallback: true,
				}).Return(txmgr.Tx{}, nil)
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
			"0",
			"",
			nil,
			false,
			pipeline.NewVarsFrom(nil),
			nil,
			func(keyStore *keystoremocks.Eth, txManager *txmmocks.MockEvmTxManager) {
				data := []byte("foobar")
				gasLimit := uint64(12345)
				txMeta := &txmgr.TxMeta{FailOnRevert: null.BoolFrom(false)}
				keyStore.On("GetRoundRobinAddress", mock.Anything, testutils.FixtureChainID, from).Return(from, nil)
				txManager.On("CreateTransaction", mock.Anything, txmgr.TxRequest{
					FromAddress:    from,
					ToAddress:      to,
					EncodedPayload: data,
					FeeLimit:       gasLimit,
					Meta:           txMeta,
					Strategy:       txmgrcommon.NewSendEveryStrategy(),
					SignalCallback: true,
				}).Return(txmgr.Tx{}, nil)
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
			"0",
			"",
			nil, // spec does not override gas limit
			false,
			pipeline.NewVarsFrom(nil),
			nil,
			func(keyStore *keystoremocks.Eth, txManager *txmmocks.MockEvmTxManager) {
				data := []byte("foobar")
				txMeta := &txmgr.TxMeta{
					JobID:         &jid,
					RequestID:     &reqID,
					RequestTxHash: &reqTxHash,
					FailOnRevert:  null.BoolFrom(false),
				}
				keyStore.On("GetRoundRobinAddress", mock.Anything, testutils.FixtureChainID, from).Return(from, nil)
				txManager.On("CreateTransaction", mock.Anything, txmgr.TxRequest{
					FromAddress:    from,
					ToAddress:      to,
					EncodedPayload: data,
					FeeLimit:       uint64(drJobTypeGasLimit),
					Meta:           txMeta,
					Strategy:       txmgrcommon.NewSendEveryStrategy(),
					SignalCallback: true,
				}).Return(txmgr.Tx{}, nil)
			},
			nil, nil, "", pipeline.RunInfo{},
		},
		{
			"happy (missing gasLimit takes spec defined value)",
			`[ "0x882969652440ccf14a5dbb9bd53eb21cb1e11e5c" ]`,
			"0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF",
			"foobar",
			"",
			`{ "jobID": 321, "requestID": "0x5198616554d738d9485d1a7cf53b2f33e09c3bbc8fe9ac0020bd672cd2bc15d2", "requestTxHash": "0xc524fafafcaec40652b1f84fca09c231185437d008d195fccf2f51e64b7062f8" }`,
			`0`,
			"0",
			"",
			&specGasLimit,
			false,
			pipeline.NewVarsFrom(nil),
			nil,
			func(keyStore *keystoremocks.Eth, txManager *txmmocks.MockEvmTxManager) {
				data := []byte("foobar")
				txMeta := &txmgr.TxMeta{
					JobID:         &jid,
					RequestID:     &reqID,
					RequestTxHash: &reqTxHash,
					FailOnRevert:  null.BoolFrom(false),
				}
				keyStore.On("GetRoundRobinAddress", mock.Anything, testutils.FixtureChainID, from).Return(from, nil)
				txManager.On("CreateTransaction", mock.Anything, txmgr.TxRequest{
					FromAddress:    from,
					ToAddress:      to,
					EncodedPayload: data,
					FeeLimit:       uint64(specGasLimit),
					Meta:           txMeta,
					Strategy:       txmgrcommon.NewSendEveryStrategy(),
					SignalCallback: true,
				}).Return(txmgr.Tx{}, nil)
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
			"0",
			"",
			nil,
			false,
			pipeline.NewVarsFrom(map[string]interface{}{
				"fromAddrs": []common.Address{common.HexToAddress("0x882969652440ccf14a5dbb9bd53eb21cb1e11e5c")},
				"toAddr":    common.HexToAddress("0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF"),
				"data":      []byte("foobar"),
				"gasLimit":  uint32(12345),
				"requestData": map[string]interface{}{
					"jobID":         int32(321),
					"requestID":     common.HexToHash("0x5198616554d738d9485d1a7cf53b2f33e09c3bbc8fe9ac0020bd672cd2bc15d2"),
					"requestTxHash": common.HexToHash("0xc524fafafcaec40652b1f84fca09c231185437d008d195fccf2f51e64b7062f8"),
				},
			}),
			nil,
			func(keyStore *keystoremocks.Eth, txManager *txmmocks.MockEvmTxManager) {
				keyStore.On("GetRoundRobinAddress", mock.Anything, testutils.FixtureChainID).Return(nil, errors.New("uh oh"))
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
			"0",
			"",
			nil,
			false,
			pipeline.NewVarsFrom(nil),
			nil,
			func(keyStore *keystoremocks.Eth, txManager *txmmocks.MockEvmTxManager) {
				data := []byte("foobar")
				gasLimit := uint64(12345)
				txMeta := &txmgr.TxMeta{
					JobID:         &jid,
					RequestID:     &reqID,
					RequestTxHash: &reqTxHash,
					FailOnRevert:  null.BoolFrom(false),
				}
				keyStore.On("GetRoundRobinAddress", mock.Anything, testutils.FixtureChainID, from).Return(from, nil)
				txManager.On("CreateTransaction", mock.Anything, txmgr.TxRequest{
					FromAddress:    from,
					ToAddress:      to,
					EncodedPayload: data,
					FeeLimit:       gasLimit,
					Meta:           txMeta,
					Strategy:       txmgrcommon.NewSendEveryStrategy(),
					SignalCallback: true,
				}).Return(txmgr.Tx{}, errors.New("uh oh"))
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
			"0",
			"",
			nil,
			false,
			pipeline.NewVarsFrom(nil),
			nil,
			func(keyStore *keystoremocks.Eth, txManager *txmmocks.MockEvmTxManager) {},
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
			"0",
			"",
			nil,
			false,
			pipeline.NewVarsFrom(nil),
			nil,
			func(keyStore *keystoremocks.Eth, txManager *txmmocks.MockEvmTxManager) {},
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
			"0",
			"",
			nil,
			false,
			pipeline.NewVarsFrom(nil),
			nil,
			func(keyStore *keystoremocks.Eth, txManager *txmmocks.MockEvmTxManager) {},
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
			"0",
			"",
			nil,
			false,
			pipeline.NewVarsFrom(nil),
			[]pipeline.Result{{Error: errors.New("uh oh")}},
			func(keyStore *keystoremocks.Eth, txManager *txmmocks.MockEvmTxManager) {},
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
			"0",
			"",
			nil,
			false,
			pipeline.NewVarsFrom(nil),
			nil,
			func(keyStore *keystoremocks.Eth, txManager *txmmocks.MockEvmTxManager) {
				from := common.HexToAddress("0x882969652440ccf14a5dbb9bd53eb21cb1e11e5c")
				keyStore.On("GetRoundRobinAddress", mock.Anything, testutils.FixtureChainID, from).Return(from, nil)
				txManager.On("CreateTransaction", mock.Anything, mock.MatchedBy(func(tx txmgr.TxRequest) bool {
					return tx.MinConfirmations == clnull.Uint32From(3) && tx.PipelineTaskRunID != nil
				})).Return(txmgr.Tx{}, nil)
			},
			nil, nil, "", pipeline.RunInfo{IsPending: true},
		},
		{
			"non-existent chain-id",
			`[ $(fromAddr) ]`,
			"$(toAddr)",
			"$(data)",
			"$(gasLimit)",
			`{ "jobID": $(jobID), "requestID": $(requestID), "requestTxHash": $(requestTxHash)`,
			`0`,
			"$(evmChainID)",
			"",
			nil,
			false,
			pipeline.NewVarsFrom(map[string]interface{}{
				"fromAddr":      common.HexToAddress("0x882969652440ccf14a5dbb9bd53eb21cb1e11e5c"),
				"toAddr":        common.HexToAddress("0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF"),
				"data":          []byte("foobar"),
				"gasLimit":      uint32(12345),
				"jobID":         int32(321),
				"requestID":     common.HexToHash("0x5198616554d738d9485d1a7cf53b2f33e09c3bbc8fe9ac0020bd672cd2bc15d2"),
				"requestTxHash": common.HexToHash("0xc524fafafcaec40652b1f84fca09c231185437d008d195fccf2f51e64b7062f8"),
				"evmChainID":    "123",
			}),
			nil,
			func(keyStore *keystoremocks.Eth, txManager *txmmocks.MockEvmTxManager) {
			},
			nil, nil, chains.ErrNoSuchChainID.Error(), pipeline.RunInfo{IsRetryable: true},
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
				EVMChainID:       test.evmChainID,
				TransmitChecker:  test.transmitChecker,
			}

			keyStore := keystoremocks.NewEth(t)
			txManager := txmmocks.NewMockEvmTxManager(t)
			db := pgtest.NewSqlxDB(t)
			cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
				c.EVM[0].GasEstimator.LimitDefault = ptr(defaultGasLimit)
				c.EVM[0].GasEstimator.LimitJobType.DR = ptr(drJobTypeGasLimit)
			})
			lggr := logger.TestLogger(t)

			relayExtenders := evmtest.NewChainRelayExtenders(t, evmtest.TestChainOpts{DB: db, GeneralConfig: cfg,
				TxManager: txManager, KeyStore: keyStore})
			legacyChains := evmrelay.NewLegacyChainsFromRelayerExtenders(relayExtenders)

			test.setupClientMocks(keyStore, txManager)
			task.HelperSetDependencies(legacyChains, keyStore, test.specGasLimit, pipeline.DirectRequestJobType)

			result, runInfo := task.Run(testutils.Context(t), lggr, test.vars, test.inputs)
			assert.Equal(t, test.expectedRunInfo, runInfo)

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

func ptr[T any](t T) *T { return &t }
