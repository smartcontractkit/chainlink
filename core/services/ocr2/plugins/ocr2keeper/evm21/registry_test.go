package evm

import (
	"context"
	"fmt"
	"math/big"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	evmClientMocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func TestGetActiveUpkeepIDs(t *testing.T) {
	tests := []struct {
		Name         string
		LatestHead   int64
		ActiveIDs    []string
		ExpectedErr  error
		ExpectedKeys []ocr2keepers.UpkeepIdentifier
	}{
		{Name: "NoActiveIDs", LatestHead: 1, ActiveIDs: []string{}, ExpectedKeys: []ocr2keepers.UpkeepIdentifier{}},
		{Name: "AvailableActiveIDs", LatestHead: 1, ActiveIDs: []string{
			"32329108151019397958065800113404894502874153543356521479058624064899121404671",
			"5820911532554020907796191562093071158274499580927271776163559390280294438608",
		}, ExpectedKeys: []ocr2keepers.UpkeepIdentifier{
			upkeepIDFromInt("32329108151019397958065800113404894502874153543356521479058624064899121404671"),
			upkeepIDFromInt("5820911532554020907796191562093071158274499580927271776163559390280294438608"),
		}},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			actives := make(map[string]activeUpkeep)
			for _, id := range test.ActiveIDs {
				idNum := big.NewInt(0)
				idNum.SetString(id, 10)
				actives[id] = activeUpkeep{ID: idNum}
			}

			rg := &EvmRegistry{
				active: actives,
			}

			keys, err := rg.GetActiveUpkeepIDs(context.Background())

			if test.ExpectedErr != nil {
				assert.ErrorIs(t, err, test.ExpectedErr)
			} else {
				assert.Nil(t, err)
			}

			if len(test.ExpectedKeys) > 0 {
				for _, key := range keys {
					assert.Contains(t, test.ExpectedKeys, key)
				}
			} else {
				assert.Equal(t, test.ExpectedKeys, keys)
			}
		})
	}
}

func TestGetActiveUpkeepIDsByType(t *testing.T) {
	tests := []struct {
		Name         string
		LatestHead   int64
		ActiveIDs    []string
		ExpectedErr  error
		ExpectedKeys []ocr2keepers.UpkeepIdentifier
		Triggers     []uint8
	}{
		{Name: "no active ids", LatestHead: 1, ActiveIDs: []string{}, ExpectedKeys: []ocr2keepers.UpkeepIdentifier{}},
		{
			Name:       "get log upkeeps",
			LatestHead: 1,
			ActiveIDs:  []string{"8", "32329108151019397958065800113404894502874153543356521479058624064899121404671"},
			ExpectedKeys: []ocr2keepers.UpkeepIdentifier{
				upkeepIDFromInt("32329108151019397958065800113404894502874153543356521479058624064899121404671"),
			},
			Triggers: []uint8{uint8(ocr2keepers.LogTrigger)},
		},
		{
			Name:       "get conditional upkeeps",
			LatestHead: 1,
			ActiveIDs:  []string{"8", "32329108151019397958065800113404894502874153543356521479058624064899121404671"},
			ExpectedKeys: []ocr2keepers.UpkeepIdentifier{
				upkeepIDFromInt("8"),
			},
			Triggers: []uint8{uint8(ocr2keepers.ConditionTrigger)},
		},
		{
			Name:       "get multiple types of upkeeps",
			LatestHead: 1,
			ActiveIDs:  []string{"8", "32329108151019397958065800113404894502874153543356521479058624064899121404671"},
			ExpectedKeys: []ocr2keepers.UpkeepIdentifier{
				upkeepIDFromInt("8"),
				upkeepIDFromInt("32329108151019397958065800113404894502874153543356521479058624064899121404671"),
			},
			Triggers: []uint8{uint8(ocr2keepers.LogTrigger), uint8(ocr2keepers.ConditionTrigger)},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			actives := make(map[string]activeUpkeep)
			for _, id := range test.ActiveIDs {
				idNum := big.NewInt(0)
				idNum.SetString(id, 10)
				actives[id] = activeUpkeep{ID: idNum}
			}

			rg := &EvmRegistry{
				active: actives,
			}

			keys, err := rg.GetActiveUpkeepIDsByType(context.Background(), test.Triggers...)

			if test.ExpectedErr != nil {
				assert.ErrorIs(t, err, test.ExpectedErr)
			} else {
				assert.Nil(t, err)
			}

			if len(test.ExpectedKeys) > 0 {
				for _, key := range keys {
					assert.Contains(t, test.ExpectedKeys, key)
				}
			} else {
				assert.Equal(t, test.ExpectedKeys, keys)
			}
		})
	}
}

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
					Return(test.LatestBlock.OutputBlock, test.LatestBlock.OutputErr)
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

			err := rg.pollLogs()

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

func TestRegistry_GetBlockAndUpkeepId(t *testing.T) {
	r := &EvmRegistry{}
	tests := []struct {
		name       string
		input      ocr2keepers.UpkeepPayload
		wantBlock  *big.Int
		wantUpkeep *big.Int
	}{
		{
			"happy flow",
			ocr2keepers.UpkeepPayload{
				UpkeepID: upkeepIDFromInt("10"),
				Trigger: ocr2keepers.Trigger{
					BlockNumber: 1,
					BlockHash:   common.HexToHash("0x1"),
				},
			},
			big.NewInt(1),
			big.NewInt(10),
		},
		{
			"empty trigger",
			ocr2keepers.UpkeepPayload{
				UpkeepID: upkeepIDFromInt("10"),
			},
			big.NewInt(0),
			big.NewInt(10),
		},
		{
			"empty payload",
			ocr2keepers.UpkeepPayload{},
			big.NewInt(0),
			big.NewInt(0),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			block, _, upkeep := r.getBlockAndUpkeepId(tc.input.UpkeepID, tc.input.Trigger)
			assert.Equal(t, tc.wantBlock, block)
			assert.Equal(t, tc.wantUpkeep.String(), upkeep.String())
		})
	}
}

func TestRegistry_VerifyCheckBlock(t *testing.T) {
	lggr := logger.TestLogger(t)
	upkeepId := ocr2keepers.UpkeepIdentifier{}
	upkeepId.FromBigInt(big.NewInt(12345))
	tests := []struct {
		name        string
		checkBlock  *big.Int
		latestBlock *big.Int
		upkeepId    *big.Int
		checkHash   common.Hash
		payload     ocr2keepers.UpkeepPayload
		blocks      map[int64]string
		state       uint8
		retryable   bool
		makeEthCall bool
	}{
		{
			name:        "check block number too told",
			checkBlock:  big.NewInt(500),
			latestBlock: big.NewInt(800),
			upkeepId:    big.NewInt(12345),
			checkHash:   common.HexToHash("0x5bff03de234fe771ac0d685f9ee0fb0b757ea02ec9e6f10e8e2ee806db1b6b83"),
			payload: ocr2keepers.UpkeepPayload{
				UpkeepID: upkeepId,
				Trigger:  ocr2keepers.NewTrigger(500, common.HexToHash("0x5bff03de234fe771ac0d685f9ee0fb0b757ea02ec9e6f10e8e2ee806db1b6b83")),
				WorkID:   "work",
			},
			state:     CheckBlockTooOld,
			retryable: false,
		},
		{
			name:        "check block number invalid",
			checkBlock:  big.NewInt(500),
			latestBlock: big.NewInt(560),
			upkeepId:    big.NewInt(12345),
			checkHash:   common.HexToHash("0x5bff03de234fe771ac0d685f9ee0fb0b757ea02ec9e6f10e8e2ee806db1b6b83"),
			payload: ocr2keepers.UpkeepPayload{
				UpkeepID: upkeepId,
				Trigger:  ocr2keepers.NewTrigger(500, common.HexToHash("0x5bff03de234fe771ac0d685f9ee0fb0b757ea02ec9e6f10e8e2ee806db1b6b83")),
				WorkID:   "work",
			},
			state:       CheckBlockInvalid,
			retryable:   true,
			makeEthCall: true,
		},
		{
			name:        "check block hash does not match",
			checkBlock:  big.NewInt(500),
			latestBlock: big.NewInt(560),
			upkeepId:    big.NewInt(12345),
			checkHash:   common.HexToHash("0x5bff03de234fe771ac0d685f9ee0fb0b757ea02ec9e6f10e8e2ee806db1b6b83"),
			payload: ocr2keepers.UpkeepPayload{
				UpkeepID: upkeepId,
				Trigger:  ocr2keepers.NewTrigger(500, common.HexToHash("0x5bff03de234fe771ac0d685f9ee0fb0b757ea02ec9e6f10e8e2ee806db1b6b83")),
				WorkID:   "work",
			},
			blocks: map[int64]string{
				500: "0xa518faeadcc423338c62572da84dda35fe44b34f521ce88f6081b703b250cca4",
			},
			state:     CheckBlockInvalid,
			retryable: false,
		},
		{
			name:        "check block is valid",
			checkBlock:  big.NewInt(500),
			latestBlock: big.NewInt(560),
			upkeepId:    big.NewInt(12345),
			checkHash:   common.HexToHash("0x5bff03de234fe771ac0d685f9ee0fb0b757ea02ec9e6f10e8e2ee806db1b6b83"),
			payload: ocr2keepers.UpkeepPayload{
				UpkeepID: upkeepId,
				Trigger:  ocr2keepers.NewTrigger(500, common.HexToHash("0x5bff03de234fe771ac0d685f9ee0fb0b757ea02ec9e6f10e8e2ee806db1b6b83")),
				WorkID:   "work",
			},
			blocks: map[int64]string{
				500: "0x5bff03de234fe771ac0d685f9ee0fb0b757ea02ec9e6f10e8e2ee806db1b6b83",
			},
			state:     NoPipelineError,
			retryable: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lb := atomic.Int64{}
			lb.Store(tc.latestBlock.Int64())
			bs := &BlockSubscriber{
				latestBlock: lb,
				blocks:      tc.blocks,
			}
			e := &EvmRegistry{
				lggr: lggr,
				bs:   bs,
			}
			if tc.makeEthCall {
				client := new(evmClientMocks.Client)
				client.On("BlockByNumber", mock.Anything, tc.checkBlock).Return(nil, fmt.Errorf("error"))
				e.client = client
			}

			state, retryable := e.verifyCheckBlock(context.Background(), tc.checkBlock, tc.upkeepId, tc.checkHash)
			assert.Equal(t, tc.state, state)
			assert.Equal(t, tc.retryable, retryable)
		})
	}
}

func TestRegistry_VerifyLogExists(t *testing.T) {
	lggr := logger.TestLogger(t)
	upkeepId := ocr2keepers.UpkeepIdentifier{}
	upkeepId.FromBigInt(big.NewInt(12345))

	extension := &ocr2keepers.LogTriggerExtension{
		TxHash:      common.HexToHash("0xc8def8abdcf3a4eaaf6cc13bff3e4e2a7168d86ea41dbbf97451235aa76c3651"),
		Index:       0,
		BlockHash:   common.HexToHash("0x3df0e926f3e21ec1195ffe007a2899214905eb02e768aa89ce0b94accd7f3d71"),
		BlockNumber: 500,
	}
	extension1 := &ocr2keepers.LogTriggerExtension{
		TxHash:      common.HexToHash("0xc8def8abdcf3a4eaaf6cc13bff3e4e2a7168d86ea41dbbf97451235aa76c3651"),
		Index:       0,
		BlockHash:   common.HexToHash("0x3df0e926f3e21ec1195ffe007a2899214905eb02e768aa89ce0b94accd7f3d71"),
		BlockNumber: 0,
	}

	tests := []struct {
		name        string
		upkeepId    *big.Int
		payload     ocr2keepers.UpkeepPayload
		blocks      map[int64]string
		makeEthCall bool
		reason      uint8
		retryable   bool
		ethCallErr  error
		receipt     *types.Receipt
	}{
		{
			name:     "log block number invalid",
			upkeepId: big.NewInt(12345),
			payload: ocr2keepers.UpkeepPayload{
				UpkeepID: upkeepId,
				Trigger:  ocr2keepers.NewLogTrigger(550, common.HexToHash("0x5bff03de234fe771ac0d685f9ee0fb0b757ea02ec9e6f10e8e2ee806db1b6b83"), extension),
				WorkID:   "work",
			},
			reason:      UpkeepFailureReasonLogBlockInvalid,
			retryable:   true,
			makeEthCall: true,
			ethCallErr:  fmt.Errorf("error"),
		},
		{
			name:     "log block no longer exists",
			upkeepId: big.NewInt(12345),
			payload: ocr2keepers.UpkeepPayload{
				UpkeepID: upkeepId,
				Trigger:  ocr2keepers.NewLogTrigger(550, common.HexToHash("0x5bff03de234fe771ac0d685f9ee0fb0b757ea02ec9e6f10e8e2ee806db1b6b83"), extension),
				WorkID:   "work",
			},
			reason:      UpkeepFailureReasonLogBlockNoLongerExists,
			retryable:   false,
			makeEthCall: true,
			blocks: map[int64]string{
				500: "0xb2173b4b75f23f56b7b2b6b2cc5fa9ed1079b9d1655b12b40fdb4dbf59006419",
			},
			receipt: &types.Receipt{},
		},
		{
			name:     "eth client returns a matching block",
			upkeepId: big.NewInt(12345),
			payload: ocr2keepers.UpkeepPayload{
				UpkeepID: upkeepId,
				Trigger:  ocr2keepers.NewLogTrigger(550, common.HexToHash("0x5bff03de234fe771ac0d685f9ee0fb0b757ea02ec9e6f10e8e2ee806db1b6b83"), extension1),
				WorkID:   "work",
			},
			reason:    UpkeepFailureReasonNone,
			retryable: false,
			blocks: map[int64]string{
				500: "0xa518faeadcc423338c62572da84dda35fe44b34f521ce88f6081b703b250cca4",
			},
			makeEthCall: true,
			receipt: &types.Receipt{
				BlockNumber: big.NewInt(550),
				BlockHash:   common.HexToHash("0x5bff03de234fe771ac0d685f9ee0fb0b757ea02ec9e6f10e8e2ee806db1b6b83"),
			},
		},
		{
			name:     "log block is valid",
			upkeepId: big.NewInt(12345),
			payload: ocr2keepers.UpkeepPayload{
				UpkeepID: upkeepId,
				Trigger:  ocr2keepers.NewLogTrigger(550, common.HexToHash("0x5bff03de234fe771ac0d685f9ee0fb0b757ea02ec9e6f10e8e2ee806db1b6b83"), extension),
				WorkID:   "work",
			},
			reason:    UpkeepFailureReasonNone,
			retryable: false,
			blocks: map[int64]string{
				500: "0x3df0e926f3e21ec1195ffe007a2899214905eb02e768aa89ce0b94accd7f3d71",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bs := &BlockSubscriber{
				blocks: tc.blocks,
			}
			e := &EvmRegistry{
				lggr: lggr,
				bs:   bs,
				ctx:  context.Background(),
			}

			if tc.makeEthCall {
				client := new(evmClientMocks.Client)
				client.On("TransactionReceipt", mock.Anything, common.HexToHash("0xc8def8abdcf3a4eaaf6cc13bff3e4e2a7168d86ea41dbbf97451235aa76c3651")).
					Return(tc.receipt, tc.ethCallErr)
				e.client = client
			}

			reason, retryable := e.verifyLogExists(tc.upkeepId, tc.payload)
			assert.Equal(t, tc.reason, reason)
			assert.Equal(t, tc.retryable, retryable)
		})
	}
}
