package evm

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	coreTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/ocr2keepers/pkg/v3/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	types3 "github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_utils_2_1"
	iregistry21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/core"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/encoding"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/logprovider"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func TestPollLogs(t *testing.T) {
	tests := []struct {
		Name             string
		LastPoll         int64
		Address          common.Address
		ExpectedLastPoll int64
		ExpectedErr      error
		LatestBlock      *struct {
			OutputBlock int64
			OutputErr   error
		}
		LogsWithSigs *struct {
			InputStart int64
			InputEnd   int64
			OutputLogs []logpoller.Log
			OutputErr  error
		}
	}{
		{
			Name:        "LatestBlockError",
			ExpectedErr: ErrHeadNotAvailable,
			LatestBlock: &struct {
				OutputBlock int64
				OutputErr   error
			}{
				OutputBlock: 0,
				OutputErr:   fmt.Errorf("test error output"),
			},
		},
		{
			Name:             "LastHeadPollIsLatestHead",
			LastPoll:         500,
			ExpectedLastPoll: 500,
			ExpectedErr:      nil,
			LatestBlock: &struct {
				OutputBlock int64
				OutputErr   error
			}{
				OutputBlock: 500,
				OutputErr:   nil,
			},
		},
		{
			Name:             "LastHeadPollNotInitialized",
			LastPoll:         0,
			ExpectedLastPoll: 500,
			ExpectedErr:      nil,
			LatestBlock: &struct {
				OutputBlock int64
				OutputErr   error
			}{
				OutputBlock: 500,
				OutputErr:   nil,
			},
		},
		{
			Name:             "LogPollError",
			LastPoll:         480,
			Address:          common.BigToAddress(big.NewInt(1)),
			ExpectedLastPoll: 500,
			ExpectedErr:      ErrLogReadFailure,
			LatestBlock: &struct {
				OutputBlock int64
				OutputErr   error
			}{
				OutputBlock: 500,
				OutputErr:   nil,
			},
			LogsWithSigs: &struct {
				InputStart int64
				InputEnd   int64
				OutputLogs []logpoller.Log
				OutputErr  error
			}{
				InputStart: 250,
				InputEnd:   500,
				OutputLogs: []logpoller.Log{},
				OutputErr:  fmt.Errorf("test output error"),
			},
		},
		{
			Name:             "LogPollSuccess",
			LastPoll:         480,
			Address:          common.BigToAddress(big.NewInt(1)),
			ExpectedLastPoll: 500,
			ExpectedErr:      nil,
			LatestBlock: &struct {
				OutputBlock int64
				OutputErr   error
			}{
				OutputBlock: 500,
				OutputErr:   nil,
			},
			LogsWithSigs: &struct {
				InputStart int64
				InputEnd   int64
				OutputLogs []logpoller.Log
				OutputErr  error
			}{
				InputStart: 250,
				InputEnd:   500,
				OutputLogs: []logpoller.Log{
					{EvmChainId: utils.NewBig(big.NewInt(5)), LogIndex: 1},
					{EvmChainId: utils.NewBig(big.NewInt(6)), LogIndex: 2},
				},
				OutputErr: nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			mp := new(mocks.LogPoller)

			if test.LatestBlock != nil {
				mp.On("LatestBlock", mock.Anything).
					Return(logpoller.LogPollerBlock{BlockNumber: test.LatestBlock.OutputBlock}, test.LatestBlock.OutputErr)
			}

			if test.LogsWithSigs != nil {
				fc := test.LogsWithSigs
				mp.On("LogsWithSigs", fc.InputStart, fc.InputEnd, upkeepStateEvents, test.Address, mock.Anything).Return(fc.OutputLogs, fc.OutputErr)
			}

			rg := &EvmRegistry{
				addr:          test.Address,
				lastPollBlock: test.LastPoll,
				poller:        mp,
				chLog:         make(chan logpoller.Log, 10),
			}

			err := rg.pollUpkeepStateLogs()

			assert.Equal(t, test.ExpectedLastPoll, rg.lastPollBlock)
			if test.ExpectedErr != nil {
				assert.ErrorIs(t, err, test.ExpectedErr)
			} else {
				assert.Nil(t, err)
			}

			var outputLogCount int

		CheckLoop:
			for {
				chT := time.NewTimer(20 * time.Millisecond)
				select {
				case l := <-rg.chLog:
					chT.Stop()
					if test.LogsWithSigs == nil {
						assert.FailNow(t, "logs detected but no logs were expected")
					}
					outputLogCount++
					assert.Contains(t, test.LogsWithSigs.OutputLogs, l)
				case <-chT.C:
					break CheckLoop
				}
			}

			if test.LogsWithSigs != nil {
				assert.Equal(t, len(test.LogsWithSigs.OutputLogs), outputLogCount)
			}

			mp.AssertExpectations(t)
		})
	}
}

func TestRegistry_refreshLogTriggerUpkeeps(t *testing.T) {
	for _, tc := range []struct {
		name             string
		ids              []*big.Int
		logEventProvider logprovider.LogEventProvider
		poller           logpoller.LogPoller
		registry         Registry
		packer           encoding.Packer
		expectsErr       bool
		wantErr          error
	}{
		{
			name: "an error is returned when fetching indexed logs for IKeeperRegistryMasterUpkeepUnpaused errors",
			ids: []*big.Int{
				core.GenUpkeepID(types.LogTrigger, "abc").BigInt(),
			},
			logEventProvider: &mockLogEventProvider{
				RefreshActiveUpkeepsFn: func(ids ...*big.Int) ([]*big.Int, error) {
					// of the ids specified in the test, only one is a valid log trigger upkeep
					assert.Equal(t, 1, len(ids))
					return ids, nil
				},
			},
			poller: &mockLogPoller{
				IndexedLogsFn: func(eventSig common.Hash, address common.Address, topicIndex int, topicValues []common.Hash, confs logpoller.Confirmations, qopts ...pg.QOpt) ([]logpoller.Log, error) {
					if eventSig == (iregistry21.IKeeperRegistryMasterUpkeepUnpaused{}.Topic()) {
						return nil, errors.New("indexed logs boom")
					}
					return nil, nil
				},
			},
			expectsErr: true,
			wantErr:    errors.New("indexed logs boom"),
		},
		{
			name: "an error is returned when fetching indexed logs for IKeeperRegistryMasterUpkeepTriggerConfigSet errors",
			ids: []*big.Int{
				core.GenUpkeepID(types.LogTrigger, "abc").BigInt(),
				core.GenUpkeepID(types.ConditionTrigger, "abc").BigInt(),
				big.NewInt(-1),
			},
			logEventProvider: &mockLogEventProvider{
				RefreshActiveUpkeepsFn: func(ids ...*big.Int) ([]*big.Int, error) {
					// of the ids specified in the test, only one is a valid log trigger upkeep
					assert.Equal(t, 1, len(ids))
					return ids, nil
				},
			},
			poller: &mockLogPoller{
				IndexedLogsFn: func(eventSig common.Hash, address common.Address, topicIndex int, topicValues []common.Hash, confs logpoller.Confirmations, qopts ...pg.QOpt) ([]logpoller.Log, error) {
					if eventSig == (iregistry21.IKeeperRegistryMasterUpkeepTriggerConfigSet{}.Topic()) {
						return nil, errors.New("indexed logs boom")
					}
					return nil, nil
				},
			},
			expectsErr: true,
			wantErr:    errors.New("indexed logs boom"),
		},
		{
			name: "an error is returned when parsing the logs using the registry errors",
			ids: []*big.Int{
				core.GenUpkeepID(types.LogTrigger, "abc").BigInt(),
				core.GenUpkeepID(types.ConditionTrigger, "abc").BigInt(),
				big.NewInt(-1),
			},
			logEventProvider: &mockLogEventProvider{
				RefreshActiveUpkeepsFn: func(ids ...*big.Int) ([]*big.Int, error) {
					// of the ids specified in the test, only one is a valid log trigger upkeep
					assert.Equal(t, 1, len(ids))
					return ids, nil
				},
			},
			poller: &mockLogPoller{
				IndexedLogsFn: func(eventSig common.Hash, address common.Address, topicIndex int, topicValues []common.Hash, confs logpoller.Confirmations, qopts ...pg.QOpt) ([]logpoller.Log, error) {
					return []logpoller.Log{
						{},
					}, nil
				},
			},
			registry: &mockRegistry{
				ParseLogFn: func(log coreTypes.Log) (generated.AbigenLog, error) {
					return nil, errors.New("parse log boom")
				},
			},
			expectsErr: true,
			wantErr:    errors.New("parse log boom"),
		},
		{
			name: "an error is returned when registering the filter errors",
			ids: []*big.Int{
				core.GenUpkeepID(types.LogTrigger, "abc").BigInt(),
				core.GenUpkeepID(types.ConditionTrigger, "abc").BigInt(),
				big.NewInt(-1),
			},
			logEventProvider: &mockLogEventProvider{
				RefreshActiveUpkeepsFn: func(ids ...*big.Int) ([]*big.Int, error) {
					// of the ids specified in the test, only one is a valid log trigger upkeep
					assert.Equal(t, 1, len(ids))
					return ids, nil
				},
				RegisterFilterFn: func(opts logprovider.FilterOptions) error {
					return errors.New("register filter boom")
				},
			},
			poller: &mockLogPoller{
				IndexedLogsFn: func(eventSig common.Hash, address common.Address, topicIndex int, topicValues []common.Hash, confs logpoller.Confirmations, qopts ...pg.QOpt) ([]logpoller.Log, error) {
					return []logpoller.Log{
						{
							BlockNumber: 1,
						},
						{
							BlockNumber: 2,
						},
					}, nil
				},
			},
			registry: &mockRegistry{
				ParseLogFn: func(log coreTypes.Log) (generated.AbigenLog, error) {
					if log.BlockNumber == 1 {
						return &iregistry21.IKeeperRegistryMasterUpkeepTriggerConfigSet{
							TriggerConfig: []byte{1, 2, 3},
							Id:            core.GenUpkeepID(types.LogTrigger, "abc").BigInt(),
						}, nil
					}
					return &iregistry21.IKeeperRegistryMasterUpkeepUnpaused{
						Id: core.GenUpkeepID(types.LogTrigger, "abc").BigInt(),
					}, nil
				},
				GetUpkeepTriggerConfigFn: func(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
					return nil, nil
				},
			},
			packer: &mockPacker{
				UnpackLogTriggerConfigFn: func(raw []byte) (automation_utils_2_1.LogTriggerConfig, error) {
					return automation_utils_2_1.LogTriggerConfig{}, nil
				},
			},
			expectsErr: true,
			wantErr:    errors.New("failed to update trigger config for upkeep id 452312848583266388373324160190187140521564213162920931037143039228013182976: failed to register log filter: register filter boom"),
		},
		{
			name: "log trigger upkeeps are refreshed without error",
			ids: []*big.Int{
				core.GenUpkeepID(types.LogTrigger, "abc").BigInt(),
				core.GenUpkeepID(types.LogTrigger, "def").BigInt(),
				core.GenUpkeepID(types.ConditionTrigger, "abc").BigInt(),
				big.NewInt(-1),
			},
			logEventProvider: &mockLogEventProvider{
				RefreshActiveUpkeepsFn: func(ids ...*big.Int) ([]*big.Int, error) {
					// of the ids specified in the test, only two are a valid log trigger upkeep
					assert.Equal(t, 2, len(ids))
					return ids, nil
				},
				RegisterFilterFn: func(opts logprovider.FilterOptions) error {
					return nil
				},
			},
			poller: &mockLogPoller{
				IndexedLogsFn: func(eventSig common.Hash, address common.Address, topicIndex int, topicValues []common.Hash, confs logpoller.Confirmations, qopts ...pg.QOpt) ([]logpoller.Log, error) {
					return []logpoller.Log{
						{
							BlockNumber: 2,
						},
						{
							BlockNumber: 1,
						},
					}, nil
				},
			},
			registry: &mockRegistry{
				ParseLogFn: func(log coreTypes.Log) (generated.AbigenLog, error) {
					if log.BlockNumber == 1 {
						return &iregistry21.IKeeperRegistryMasterUpkeepTriggerConfigSet{
							Id:            core.GenUpkeepID(types.LogTrigger, "abc").BigInt(),
							TriggerConfig: []byte{1, 2, 3},
						}, nil
					}
					return &iregistry21.IKeeperRegistryMasterUpkeepUnpaused{
						Id: core.GenUpkeepID(types.LogTrigger, "def").BigInt(),
					}, nil
				},
				GetUpkeepTriggerConfigFn: func(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
					return nil, nil
				},
			},
			packer: &mockPacker{
				UnpackLogTriggerConfigFn: func(raw []byte) (automation_utils_2_1.LogTriggerConfig, error) {
					return automation_utils_2_1.LogTriggerConfig{}, nil
				},
			},
		},
		{
			name: "log trigger upkeeps are refreshed in batch without error",
			ids: func() []*big.Int {
				res := []*big.Int{}
				for i := 0; i < logTriggerRefreshBatchSize*3; i++ {
					res = append(res, core.GenUpkeepID(types.LogTrigger, fmt.Sprintf("%d", i)).BigInt())
				}
				return res
			}(),
			logEventProvider: &mockLogEventProvider{
				RefreshActiveUpkeepsFn: func(ids ...*big.Int) ([]*big.Int, error) {
					assert.Equal(t, logTriggerRefreshBatchSize, len(ids))
					return ids, nil
				},
				RegisterFilterFn: func(opts logprovider.FilterOptions) error {
					return nil
				},
			},
			poller: &mockLogPoller{
				IndexedLogsFn: func(eventSig common.Hash, address common.Address, topicIndex int, topicValues []common.Hash, confs logpoller.Confirmations, qopts ...pg.QOpt) ([]logpoller.Log, error) {
					return []logpoller.Log{
						{
							BlockNumber: 2,
						},
						{
							BlockNumber: 1,
						},
					}, nil
				},
			},
			registry: &mockRegistry{
				ParseLogFn: func(log coreTypes.Log) (generated.AbigenLog, error) {
					if log.BlockNumber == 1 {
						return &iregistry21.IKeeperRegistryMasterUpkeepTriggerConfigSet{
							Id:            core.GenUpkeepID(types.LogTrigger, "abc").BigInt(),
							TriggerConfig: []byte{1, 2, 3},
						}, nil
					}
					return &iregistry21.IKeeperRegistryMasterUpkeepUnpaused{
						Id: core.GenUpkeepID(types.LogTrigger, "def").BigInt(),
					}, nil
				},
				GetUpkeepTriggerConfigFn: func(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
					return nil, nil
				},
			},
			packer: &mockPacker{
				UnpackLogTriggerConfigFn: func(raw []byte) (automation_utils_2_1.LogTriggerConfig, error) {
					return automation_utils_2_1.LogTriggerConfig{}, nil
				},
			},
		},
		{
			name: "log trigger upkeeps are refreshed in batch, with a partial batch without error",
			ids: func() []*big.Int {
				res := []*big.Int{}
				for i := 0; i < logTriggerRefreshBatchSize+3; i++ {
					res = append(res, core.GenUpkeepID(types.LogTrigger, fmt.Sprintf("%d", i)).BigInt())
				}
				return res
			}(),
			logEventProvider: &mockLogEventProvider{
				RefreshActiveUpkeepsFn: func(ids ...*big.Int) ([]*big.Int, error) {
					if len(ids) != logTriggerRefreshBatchSize {
						assert.Equal(t, 3, len(ids))
					}
					return ids, nil
				},
				RegisterFilterFn: func(opts logprovider.FilterOptions) error {
					return nil
				},
			},
			poller: &mockLogPoller{
				IndexedLogsFn: func(eventSig common.Hash, address common.Address, topicIndex int, topicValues []common.Hash, confs logpoller.Confirmations, qopts ...pg.QOpt) ([]logpoller.Log, error) {
					return []logpoller.Log{
						{
							BlockNumber: 2,
						},
						{
							BlockNumber: 1,
						},
					}, nil
				},
			},
			registry: &mockRegistry{
				ParseLogFn: func(log coreTypes.Log) (generated.AbigenLog, error) {
					if log.BlockNumber == 1 {
						return &iregistry21.IKeeperRegistryMasterUpkeepTriggerConfigSet{
							Id:            core.GenUpkeepID(types.LogTrigger, "abc").BigInt(),
							TriggerConfig: []byte{1, 2, 3},
						}, nil
					}
					return &iregistry21.IKeeperRegistryMasterUpkeepUnpaused{
						Id: core.GenUpkeepID(types.LogTrigger, "def").BigInt(),
					}, nil
				},
				GetUpkeepTriggerConfigFn: func(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
					return nil, nil
				},
			},
			packer: &mockPacker{
				UnpackLogTriggerConfigFn: func(raw []byte) (automation_utils_2_1.LogTriggerConfig, error) {
					return automation_utils_2_1.LogTriggerConfig{}, nil
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			lggr := logger.TestLogger(t)
			var hb types3.HeadBroadcaster
			var lp logpoller.LogPoller

			bs := NewBlockSubscriber(hb, lp, 1000, lggr)

			registry := &EvmRegistry{
				addr:             common.BigToAddress(big.NewInt(1)),
				poller:           tc.poller,
				logEventProvider: tc.logEventProvider,
				chLog:            make(chan logpoller.Log, 10),
				bs:               bs,
				registry:         tc.registry,
				packer:           tc.packer,
				lggr:             lggr,
			}

			err := registry.refreshLogTriggerUpkeeps(tc.ids)
			if tc.expectsErr {
				assert.Error(t, err)
				assert.Equal(t, err.Error(), tc.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

type mockLogEventProvider struct {
	logprovider.LogEventProvider
	RefreshActiveUpkeepsFn func(ids ...*big.Int) ([]*big.Int, error)
	RegisterFilterFn       func(opts logprovider.FilterOptions) error
}

func (p *mockLogEventProvider) RefreshActiveUpkeeps(ids ...*big.Int) ([]*big.Int, error) {
	return p.RefreshActiveUpkeepsFn(ids...)
}

func (p *mockLogEventProvider) RegisterFilter(ctx context.Context, opts logprovider.FilterOptions) error {
	return p.RegisterFilterFn(opts)
}

type mockRegistry struct {
	Registry
	GetUpkeepTriggerConfigFn func(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error)
	ParseLogFn               func(log coreTypes.Log) (generated.AbigenLog, error)
}

func (r *mockRegistry) ParseLog(log coreTypes.Log) (generated.AbigenLog, error) {
	return r.ParseLogFn(log)
}

func (r *mockRegistry) GetUpkeepTriggerConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
	return r.GetUpkeepTriggerConfigFn(opts, upkeepId)
}

type mockPacker struct {
	encoding.Packer
	UnpackLogTriggerConfigFn func(raw []byte) (automation_utils_2_1.LogTriggerConfig, error)
}

func (p *mockPacker) UnpackLogTriggerConfig(raw []byte) (automation_utils_2_1.LogTriggerConfig, error) {
	return p.UnpackLogTriggerConfigFn(raw)
}
