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
	"github.com/ethereum/go-ethereum/rpc"

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
		state       PipelineExecutionState
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
			state: CheckBlockTooOld,
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
			state:       RpcFlakyFailure,
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
			state: CheckBlockInvalid,
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
			state: NoPipelineError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bs := &BlockSubscriber{
				latestBlock: atomic.Int64{},
				blocks:      tc.blocks,
			}
			bs.latestBlock.Store(tc.latestBlock.Int64())
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
		reason      UpkeepFailureReason
		state       PipelineExecutionState
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
			reason:      UpkeepFailureReasonNone,
			state:       RpcFlakyFailure,
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

			reason, state, retryable := e.verifyLogExists(tc.upkeepId, tc.payload)
			assert.Equal(t, tc.reason, reason)
			assert.Equal(t, tc.state, state)
			assert.Equal(t, tc.retryable, retryable)
		})
	}
}

func TestRegistry_CheckUpkeeps(t *testing.T) {
	lggr := logger.TestLogger(t)
	uid0 := genUpkeepID(ocr2keepers.UpkeepType(0), "p0")
	uid1 := genUpkeepID(ocr2keepers.UpkeepType(1), "p1")
	uid2 := genUpkeepID(ocr2keepers.UpkeepType(1), "p2")

	extension1 := &ocr2keepers.LogTriggerExtension{
		TxHash:      common.HexToHash("0xc8def8abdcf3a4eaaf6cc13bff3e4e2a7168d86ea41dbbf97451235aa76c3651"),
		Index:       0,
		BlockHash:   common.HexToHash("0x0919c83363b439ea634ce2b576cf3e30db26b340fb7a12058c2fcc401bd04ba0"),
		BlockNumber: 550,
	}
	extension2 := &ocr2keepers.LogTriggerExtension{
		TxHash:      common.HexToHash("0xc8def8abdcf3a4eaaf6cc13bff3e4e2a7168d86ea41dbbf97451235aa76c3651"),
		Index:       0,
		BlockHash:   common.HexToHash("0x9840e5b709bfccf6a1b44f34c884bc39403f57923f3f5ead6243cc090546b857"),
		BlockNumber: 550,
	}

	trigger0 := ocr2keepers.NewTrigger(150, common.HexToHash("0x1c77db0abe32327cf3ea9de2aadf79876f9e6b6dfcee9d4719a8a2dc8ca289d0"))
	trigger1 := ocr2keepers.NewLogTrigger(560, common.HexToHash("0x9840e5b709bfccf6a1b44f34c884bc39403f57923f3f5ead6243cc090546b857"), extension1)
	trigger2 := ocr2keepers.NewLogTrigger(570, common.HexToHash("0x1222d75217e2dd461cc77e4091c37abe76277430d97f1963a822b4e94ebb83fc"), extension2)

	tests := []struct {
		name          string
		inputs        []ocr2keepers.UpkeepPayload
		blocks        map[int64]string
		latestBlock   *big.Int
		results       []ocr2keepers.CheckResult
		err           error
		ethCalls      map[string]bool
		receipts      map[string]*types.Receipt
		ethCallErrors map[string]error
	}{
		{
			name: "check upkeeps with different upkeep types",
			inputs: []ocr2keepers.UpkeepPayload{
				{
					UpkeepID: uid0,
					Trigger:  trigger0,
					WorkID:   "work0",
				},
				{
					UpkeepID: uid1,
					Trigger:  trigger1,
					WorkID:   "work1",
				},
				{
					UpkeepID: uid2,
					Trigger:  trigger2,
					WorkID:   "work2",
					// check data byte slice length cannot be odd number, abi pack error
					CheckData: []byte{0, 0, 0, 0, 1},
				},
			},
			blocks: map[int64]string{
				550: "0x9840e5b709bfccf6a1b44f34c884bc39403f57923f3f5ead6243cc090546b857",
				560: "0x9840e5b709bfccf6a1b44f34c884bc39403f57923f3f5ead6243cc090546b857",
				570: "0x1222d75217e2dd461cc77e4091c37abe76277430d97f1963a822b4e94ebb83fc",
			},
			latestBlock: big.NewInt(580),
			results: []ocr2keepers.CheckResult{
				{
					PipelineExecutionState: uint8(CheckBlockTooOld),
					Retryable:              false,
					Eligible:               false,
					IneligibilityReason:    0,
					UpkeepID:               uid0,
					Trigger:                trigger0,
					WorkID:                 "work0",
					GasAllocated:           0,
					PerformData:            nil,
					FastGasWei:             big.NewInt(0),
					LinkNative:             big.NewInt(0),
				},
				{
					PipelineExecutionState: uint8(RpcFlakyFailure),
					Retryable:              true,
					Eligible:               false,
					IneligibilityReason:    0,
					UpkeepID:               uid1,
					Trigger:                trigger1,
					WorkID:                 "work1",
					GasAllocated:           0,
					PerformData:            nil,
					FastGasWei:             big.NewInt(0),
					LinkNative:             big.NewInt(0),
				},
				{
					PipelineExecutionState: uint8(PackUnpackDecodeFailed),
					Retryable:              false,
					Eligible:               false,
					IneligibilityReason:    0,
					UpkeepID:               uid2,
					Trigger:                trigger2,
					WorkID:                 "work2",
					GasAllocated:           0,
					PerformData:            nil,
					FastGasWei:             big.NewInt(0),
					LinkNative:             big.NewInt(0),
				},
			},
			ethCalls: map[string]bool{
				uid1.String(): true,
			},
			receipts: map[string]*types.Receipt{
				//uid1.String(): {
				//	BlockNumber: big.NewInt(550),
				//	BlockHash:   common.HexToHash("0x5bff03de234fe771ac0d685f9ee0fb0b757ea02ec9e6f10e8e2ee806db1b6b83"),
				//},
			},
			ethCallErrors: map[string]error{
				uid1.String(): fmt.Errorf("error"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bs := &BlockSubscriber{
				latestBlock: atomic.Int64{},
				blocks:      tc.blocks,
			}
			bs.latestBlock.Store(tc.latestBlock.Int64())
			e := &EvmRegistry{
				lggr: lggr,
				bs:   bs,
			}
			client := new(evmClientMocks.Client)
			for _, i := range tc.inputs {
				uid := i.UpkeepID.String()
				if tc.ethCalls[uid] {
					client.On("TransactionReceipt", mock.Anything, common.HexToHash("0xc8def8abdcf3a4eaaf6cc13bff3e4e2a7168d86ea41dbbf97451235aa76c3651")).
						Return(tc.receipts[uid], tc.ethCallErrors[uid])
				}
			}
			e.client = client

			results, err := e.checkUpkeeps(context.Background(), tc.inputs)
			assert.Equal(t, tc.results, results)
			assert.Equal(t, tc.err, err)
		})
	}
}

func TestRegistry_SimulatePerformUpkeeps(t *testing.T) {
	uid0 := genUpkeepID(ocr2keepers.UpkeepType(0), "p0")
	uid1 := genUpkeepID(ocr2keepers.UpkeepType(1), "p1")
	uid2 := genUpkeepID(ocr2keepers.UpkeepType(1), "p2")

	extension1 := &ocr2keepers.LogTriggerExtension{
		TxHash:      common.HexToHash("0xc8def8abdcf3a4eaaf6cc13bff3e4e2a7168d86ea41dbbf97451235aa76c3651"),
		Index:       0,
		BlockHash:   common.HexToHash("0x9840e5b709bfccf6a1b44f34c884bc39403f57923f3f5ead6243cc090546b857"),
		BlockNumber: 550,
	}

	trigger0 := ocr2keepers.NewTrigger(150, common.HexToHash("0x1c77db0abe32327cf3ea9de2aadf79876f9e6b6dfcee9d4719a8a2dc8ca289d0"))
	trigger1 := ocr2keepers.NewLogTrigger(570, common.HexToHash("0x1222d75217e2dd461cc77e4091c37abe76277430d97f1963a822b4e94ebb83fc"), extension1)
	trigger2 := ocr2keepers.NewLogTrigger(570, common.HexToHash("0x1222d75217e2dd461cc77e4091c37abe76277430d97f1963a822b4e94ebb83fc"), extension1)

	cr0 := ocr2keepers.CheckResult{
		PipelineExecutionState: uint8(CheckBlockTooOld),
		Retryable:              false,
		Eligible:               false,
		IneligibilityReason:    0,
		UpkeepID:               uid0,
		Trigger:                trigger0,
		WorkID:                 "work0",
		GasAllocated:           0,
		PerformData:            nil,
		FastGasWei:             big.NewInt(0),
		LinkNative:             big.NewInt(0),
	}

	tests := []struct {
		name    string
		inputs  []ocr2keepers.CheckResult
		results []ocr2keepers.CheckResult
		err     error
	}{
		{
			name: "simulate multiple upkeeps",
			inputs: []ocr2keepers.CheckResult{
				cr0,
				{
					PipelineExecutionState: 0,
					Retryable:              false,
					Eligible:               true,
					IneligibilityReason:    0,
					UpkeepID:               uid1,
					Trigger:                trigger1,
					WorkID:                 "work1",
					GasAllocated:           20000,
					PerformData:            []byte{0, 0, 0, 1, 2, 3},
					FastGasWei:             big.NewInt(20000),
					LinkNative:             big.NewInt(20000),
				},
				{
					PipelineExecutionState: 0,
					Retryable:              false,
					Eligible:               true,
					IneligibilityReason:    0,
					UpkeepID:               uid2,
					Trigger:                trigger2,
					WorkID:                 "work2",
					GasAllocated:           20000,
					PerformData:            []byte{0, 0, 0, 1, 2, 3},
					FastGasWei:             big.NewInt(20000),
					LinkNative:             big.NewInt(20000),
				},
			},
			results: []ocr2keepers.CheckResult{
				cr0,
				{
					PipelineExecutionState: uint8(RpcFlakyFailure),
					Retryable:              true,
					Eligible:               false,
					IneligibilityReason:    0,
					UpkeepID:               uid1,
					Trigger:                trigger1,
					WorkID:                 "work1",
					GasAllocated:           20000,
					PerformData:            []byte{0, 0, 0, 1, 2, 3},
					FastGasWei:             big.NewInt(20000),
					LinkNative:             big.NewInt(20000),
				},
				{
					PipelineExecutionState: uint8(PackUnpackDecodeFailed),
					Retryable:              false,
					Eligible:               false,
					IneligibilityReason:    0,
					UpkeepID:               uid2,
					Trigger:                trigger2,
					WorkID:                 "work2",
					GasAllocated:           20000,
					PerformData:            []byte{0, 0, 0, 1, 2, 3},
					FastGasWei:             big.NewInt(20000),
					LinkNative:             big.NewInt(20000),
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := setupEVMRegistry(t)
			client := new(evmClientMocks.Client)
			client.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
				return len(b) == 2 && b[0].Method == "eth_call" && b[1].Method == "eth_call"
			})).Return(nil).
				Run(func(args mock.Arguments) {
					be := args.Get(1).([]rpc.BatchElem)
					be[0].Error = fmt.Errorf("error")
					res := "0x0001"
					be[1].Result = res
				}).Once()
			e.client = client

			results, err := e.simulatePerformUpkeeps(context.Background(), tc.inputs)
			assert.Equal(t, tc.results, results)
			assert.Equal(t, tc.err, err)
		})
	}

}
