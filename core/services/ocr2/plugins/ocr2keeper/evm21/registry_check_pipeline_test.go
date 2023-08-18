package evm

import (
	"context"
	"fmt"
	"math/big"
	"sync/atomic"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"

	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	evmClientMocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/core"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/encoding"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

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
				UpkeepID: core.UpkeepIDFromInt("10"),
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
				UpkeepID: core.UpkeepIDFromInt("10"),
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
		latestBlock ocr2keepers.BlockKey
		upkeepId    *big.Int
		checkHash   common.Hash
		payload     ocr2keepers.UpkeepPayload
		blocks      map[int64]string
		poller      logpoller.LogPoller
		state       encoding.PipelineExecutionState
		retryable   bool
		makeEthCall bool
	}{
		{
			name:        "check block number too told",
			checkBlock:  big.NewInt(500),
			latestBlock: ocr2keepers.BlockKey{Number: 800},
			upkeepId:    big.NewInt(12345),
			checkHash:   common.HexToHash("0x5bff03de234fe771ac0d685f9ee0fb0b757ea02ec9e6f10e8e2ee806db1b6b83"),
			payload: ocr2keepers.UpkeepPayload{
				UpkeepID: upkeepId,
				Trigger:  ocr2keepers.NewTrigger(500, common.HexToHash("0x5bff03de234fe771ac0d685f9ee0fb0b757ea02ec9e6f10e8e2ee806db1b6b83")),
				WorkID:   "work",
			},
			state: encoding.CheckBlockTooOld,
		},
		{
			name:        "for an invalid check block number, if hash does not match the check hash, return CheckBlockInvalid",
			checkBlock:  big.NewInt(500),
			latestBlock: ocr2keepers.BlockKey{Number: 560},
			upkeepId:    big.NewInt(12345),
			checkHash:   common.HexToHash("0x5bff03de234fe771ac0d685f9ee0fb0b757ea02ec9e6f10e8e2ee806db1b6b83"),
			payload: ocr2keepers.UpkeepPayload{
				UpkeepID: upkeepId,
				Trigger:  ocr2keepers.NewTrigger(500, common.HexToHash("0x5bff03de234fe771ac0d685f9ee0fb0b757ea02ec9e6f10e8e2ee806db1b6b83")),
				WorkID:   "work",
			},
			poller: &mockLogPoller{
				GetBlocksRangeFn: func(ctx context.Context, numbers []uint64, qopts ...pg.QOpt) ([]logpoller.LogPollerBlock, error) {
					return []logpoller.LogPollerBlock{
						{
							BlockHash: common.HexToHash("abcdef"),
						},
					}, nil
				},
			},
			state:       encoding.CheckBlockInvalid,
			retryable:   false,
			makeEthCall: true,
		},
		{
			name:        "for an invalid check block number, if hash does match the check hash, return NoPipelineError",
			checkBlock:  big.NewInt(500),
			latestBlock: ocr2keepers.BlockKey{Number: 560},
			upkeepId:    big.NewInt(12345),
			checkHash:   common.HexToHash("0x5bff03de234fe771ac0d685f9ee0fb0b757ea02ec9e6f10e8e2ee806db1b6b83"),
			payload: ocr2keepers.UpkeepPayload{
				UpkeepID: upkeepId,
				Trigger:  ocr2keepers.NewTrigger(500, common.HexToHash("0x5bff03de234fe771ac0d685f9ee0fb0b757ea02ec9e6f10e8e2ee806db1b6b83")),
				WorkID:   "work",
			},
			poller: &mockLogPoller{
				GetBlocksRangeFn: func(ctx context.Context, numbers []uint64, qopts ...pg.QOpt) ([]logpoller.LogPollerBlock, error) {
					return []logpoller.LogPollerBlock{
						{
							BlockHash: common.HexToHash("0x5bff03de234fe771ac0d685f9ee0fb0b757ea02ec9e6f10e8e2ee806db1b6b83"),
						},
					}, nil
				},
			},
			state:       encoding.NoPipelineError,
			retryable:   false,
			makeEthCall: true,
		},
		{
			name:        "check block hash does not match",
			checkBlock:  big.NewInt(500),
			latestBlock: ocr2keepers.BlockKey{Number: 560},
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
			state: encoding.CheckBlockInvalid,
		},
		{
			name:        "check block is valid",
			checkBlock:  big.NewInt(500),
			latestBlock: ocr2keepers.BlockKey{Number: 560},
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
			state: encoding.NoPipelineError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bs := &BlockSubscriber{
				latestBlock: atomic.Pointer[ocr2keepers.BlockKey]{},
				blocks:      tc.blocks,
			}
			bs.latestBlock.Store(&tc.latestBlock)
			e := &EvmRegistry{
				lggr:   lggr,
				bs:     bs,
				poller: tc.poller,
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

type mockLogPoller struct {
	logpoller.LogPoller
	GetBlocksRangeFn func(ctx context.Context, numbers []uint64, qopts ...pg.QOpt) ([]logpoller.LogPollerBlock, error)
}

func (p *mockLogPoller) GetBlocksRange(ctx context.Context, numbers []uint64, qopts ...pg.QOpt) ([]logpoller.LogPollerBlock, error) {
	return p.GetBlocksRangeFn(ctx, numbers, qopts...)
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
		reason      encoding.UpkeepFailureReason
		state       encoding.PipelineExecutionState
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
			reason:      encoding.UpkeepFailureReasonNone,
			state:       encoding.RpcFlakyFailure,
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
			reason:      encoding.UpkeepFailureReasonTxHashNoLongerExists,
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
			reason:    encoding.UpkeepFailureReasonNone,
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
			reason:    encoding.UpkeepFailureReasonNone,
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
	uid0 := core.GenUpkeepID(ocr2keepers.UpkeepType(0), "p0")
	uid1 := core.GenUpkeepID(ocr2keepers.UpkeepType(1), "p1")
	uid2 := core.GenUpkeepID(ocr2keepers.UpkeepType(1), "p2")

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
		latestBlock   ocr2keepers.BlockKey
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
			latestBlock: ocr2keepers.BlockKey{Number: 580},
			results: []ocr2keepers.CheckResult{
				{
					PipelineExecutionState: uint8(encoding.CheckBlockTooOld),
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
					PipelineExecutionState: uint8(encoding.RpcFlakyFailure),
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
					PipelineExecutionState: uint8(encoding.PackUnpackDecodeFailed),
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
				latestBlock: atomic.Pointer[ocr2keepers.BlockKey]{},
				blocks:      tc.blocks,
			}
			bs.latestBlock.Store(&tc.latestBlock)
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
	uid0 := core.GenUpkeepID(ocr2keepers.UpkeepType(0), "p0")
	uid1 := core.GenUpkeepID(ocr2keepers.UpkeepType(1), "p1")
	uid2 := core.GenUpkeepID(ocr2keepers.UpkeepType(1), "p2")

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
		PipelineExecutionState: uint8(encoding.CheckBlockTooOld),
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
					PipelineExecutionState: uint8(encoding.RpcFlakyFailure),
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
					PipelineExecutionState: uint8(encoding.PackUnpackDecodeFailed),
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
