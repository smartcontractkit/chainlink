package evm

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync/atomic"
	"testing"

	types3 "github.com/smartcontractkit/chainlink-automation/pkg/v3/types"

	types2 "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	ocr2keepers "github.com/smartcontractkit/chainlink-common/pkg/types/automation"

	evmClientMocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	gasMocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	ac "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_automation_v21_plus_common"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/streams_lookup_compatible_interface"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/core"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/encoding"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/mocks"
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
		latestBlock *ocr2keepers.BlockKey
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
			name:        "for an invalid check block number, if hash does not match the check hash, return CheckBlockInvalid",
			checkBlock:  big.NewInt(500),
			latestBlock: &ocr2keepers.BlockKey{Number: 560},
			upkeepId:    big.NewInt(12345),
			checkHash:   common.HexToHash("0x5bff03de234fe771ac0d685f9ee0fb0b757ea02ec9e6f10e8e2ee806db1b6b83"),
			payload: ocr2keepers.UpkeepPayload{
				UpkeepID: upkeepId,
				Trigger:  ocr2keepers.NewTrigger(500, common.HexToHash("0x5bff03de234fe771ac0d685f9ee0fb0b757ea02ec9e6f10e8e2ee806db1b6b83")),
				WorkID:   "work",
			},
			poller: &mockLogPoller{
				GetBlocksRangeFn: func(ctx context.Context, numbers []uint64) ([]logpoller.LogPollerBlock, error) {
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
			latestBlock: &ocr2keepers.BlockKey{Number: 560},
			upkeepId:    big.NewInt(12345),
			checkHash:   common.HexToHash("0x5bff03de234fe771ac0d685f9ee0fb0b757ea02ec9e6f10e8e2ee806db1b6b83"),
			payload: ocr2keepers.UpkeepPayload{
				UpkeepID: upkeepId,
				Trigger:  ocr2keepers.NewTrigger(500, common.HexToHash("0x5bff03de234fe771ac0d685f9ee0fb0b757ea02ec9e6f10e8e2ee806db1b6b83")),
				WorkID:   "work",
			},
			poller: &mockLogPoller{
				GetBlocksRangeFn: func(ctx context.Context, numbers []uint64) ([]logpoller.LogPollerBlock, error) {
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
			latestBlock: &ocr2keepers.BlockKey{Number: 560},
			upkeepId:    big.NewInt(12345),
			checkHash:   common.HexToHash("0x5bff03de234fe771ac0d685f9ee0fb0b757ea02ec9e6f10e8e2ee806db1b6b83"),
			payload: ocr2keepers.UpkeepPayload{
				UpkeepID: upkeepId,
				Trigger:  ocr2keepers.NewTrigger(500, common.HexToHash("0x5bff03de234fe771ac0d685f9ee0fb0b757ea02ec9e6f10e8e2ee806db1b6b83")),
				WorkID:   "work",
			},
			poller: &mockLogPoller{
				GetBlocksRangeFn: func(ctx context.Context, numbers []uint64) ([]logpoller.LogPollerBlock, error) {
					return []logpoller.LogPollerBlock{
						{
							BlockHash: common.HexToHash("0xcba5cf9e2bb32373c76015384e1098912d9510a72481c78057fcb088209167de"),
						},
					}, nil
				},
			},
			blocks: map[int64]string{
				500: "0xa518faeadcc423338c62572da84dda35fe44b34f521ce88f6081b703b250cca4",
			},
			state: encoding.CheckBlockInvalid,
		},
		{
			name:        "check block is valid",
			checkBlock:  big.NewInt(500),
			latestBlock: &ocr2keepers.BlockKey{Number: 560},
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
			bs.latestBlock.Store(tc.latestBlock)
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

			state, retryable := e.verifyCheckBlock(testutils.Context(t), tc.checkBlock, tc.upkeepId, tc.checkHash)
			assert.Equal(t, tc.state, state)
			assert.Equal(t, tc.retryable, retryable)
		})
	}
}

type mockLogPoller struct {
	logpoller.LogPoller
	GetBlocksRangeFn func(ctx context.Context, numbers []uint64) ([]logpoller.LogPollerBlock, error)
	IndexedLogsFn    func(ctx context.Context, eventSig common.Hash, address common.Address, topicIndex int, topicValues []common.Hash, confs evmtypes.Confirmations) ([]logpoller.Log, error)
}

func (p *mockLogPoller) GetBlocksRange(ctx context.Context, numbers []uint64) ([]logpoller.LogPollerBlock, error) {
	return p.GetBlocksRangeFn(ctx, numbers)
}

func (p *mockLogPoller) IndexedLogs(ctx context.Context, eventSig common.Hash, address common.Address, topicIndex int, topicValues []common.Hash, confs evmtypes.Confirmations) ([]logpoller.Log, error) {
	return p.IndexedLogsFn(ctx, eventSig, address, topicIndex, topicValues, confs)
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
			receipt: &types.Receipt{Status: 0},
		},
		{
			name:     "eth client returns a matching block but different hash",
			upkeepId: big.NewInt(12345),
			payload: ocr2keepers.UpkeepPayload{
				UpkeepID: upkeepId,
				Trigger:  ocr2keepers.NewLogTrigger(550, common.HexToHash("0x5bff03de234fe771ac0d685f9ee0fb0b757ea02ec9e6f10e8e2ee806db1b6b83"), extension1),
				WorkID:   "work",
			},
			reason:    encoding.UpkeepFailureReasonTxHashReorged,
			retryable: false,
			blocks: map[int64]string{
				500: "0xa518faeadcc423338c62572da84dda35fe44b34f521ce88f6081b703b250cca4",
			},
			makeEthCall: true,
			receipt: &types.Receipt{
				Status:      1,
				BlockNumber: big.NewInt(550),
				BlockHash:   common.HexToHash("0x5bff03de234fe771ac0d685f9ee0fb0b757ea02ec9e6f10e8e2ee806db1b6b83"),
			},
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
				Status:      1,
				BlockNumber: big.NewInt(550),
				BlockHash:   common.HexToHash("0x3df0e926f3e21ec1195ffe007a2899214905eb02e768aa89ce0b94accd7f3d71"),
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
			ctx := testutils.Context(t)
			bs := &BlockSubscriber{
				blocks: tc.blocks,
			}
			e := &EvmRegistry{
				lggr: lggr,
				bs:   bs,
			}

			if tc.makeEthCall {
				client := new(evmClientMocks.Client)
				client.On("CallContext", mock.Anything, mock.Anything, "eth_getTransactionReceipt", common.BytesToHash(tc.payload.Trigger.LogTriggerExtension.TxHash[:])).
					Return(tc.ethCallErr).Run(func(args mock.Arguments) {
					if tc.receipt != nil {
						res := args.Get(1).(*types.Receipt)
						res.Status = tc.receipt.Status
						res.TxHash = tc.receipt.TxHash
						res.BlockNumber = tc.receipt.BlockNumber
						res.BlockHash = tc.receipt.BlockHash
					}
				})
				e.client = client
			}

			reason, state, retryable := e.verifyLogExists(ctx, tc.upkeepId, tc.payload)
			assert.Equal(t, tc.reason, reason)
			assert.Equal(t, tc.state, state)
			assert.Equal(t, tc.retryable, retryable)
		})
	}
}

func TestRegistry_CheckUpkeeps(t *testing.T) {
	lggr := logger.TestLogger(t)
	uid0 := core.GenUpkeepID(types3.UpkeepType(0), "p0")
	uid1 := core.GenUpkeepID(types3.UpkeepType(1), "p1")
	uid2 := core.GenUpkeepID(types3.UpkeepType(1), "p2")

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

	trigger0 := ocr2keepers.NewTrigger(575, common.HexToHash("0x1c77db0abe32327cf3ea9de2aadf79876f9e6b6dfcee9d4719a8a2dc8ca289d0"))
	trigger1 := ocr2keepers.NewLogTrigger(560, common.HexToHash("0x9840e5b709bfccf6a1b44f34c884bc39403f57923f3f5ead6243cc090546b857"), extension1)
	trigger2 := ocr2keepers.NewLogTrigger(570, common.HexToHash("0x1222d75217e2dd461cc77e4091c37abe76277430d97f1963a822b4e94ebb83fc"), extension2)

	tests := []struct {
		name          string
		inputs        []ocr2keepers.UpkeepPayload
		blocks        map[int64]string
		latestBlock   *ocr2keepers.BlockKey
		results       []ocr2keepers.CheckResult
		err           error
		ethCalls      map[string]bool
		receipts      map[string]*types.Receipt
		poller        logpoller.LogPoller
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
				575: "0x9840e5b709bfccf6a1b44f34c884bc39403f57923f3f5ead6243cc090546b857",
			},
			latestBlock: &ocr2keepers.BlockKey{Number: 580},
			results: []ocr2keepers.CheckResult{
				{
					PipelineExecutionState: uint8(encoding.CheckBlockInvalid),
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
			receipts: map[string]*types.Receipt{},
			poller: &mockLogPoller{
				GetBlocksRangeFn: func(ctx context.Context, numbers []uint64) ([]logpoller.LogPollerBlock, error) {
					return []logpoller.LogPollerBlock{
						{
							BlockHash: common.HexToHash("0xcba5cf9e2bb32373c76015384e1098912d9510a72481c78057fcb088209167de"),
						},
					}, nil
				},
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
			bs.latestBlock.Store(tc.latestBlock)
			e := &EvmRegistry{
				lggr:   lggr,
				bs:     bs,
				poller: tc.poller,
			}
			client := new(evmClientMocks.Client)
			for _, i := range tc.inputs {
				uid := i.UpkeepID.String()
				if tc.ethCalls[uid] {
					client.On("CallContext", mock.Anything, mock.Anything, "eth_getTransactionReceipt", common.HexToHash("0xc8def8abdcf3a4eaaf6cc13bff3e4e2a7168d86ea41dbbf97451235aa76c3651")).
						Return(tc.ethCallErrors[uid]).Run(func(args mock.Arguments) {
						receipt := tc.receipts[uid]
						if receipt != nil {
							res := args.Get(1).(*types.Receipt)
							res.Status = receipt.Status
							res.TxHash = receipt.TxHash
							res.BlockNumber = receipt.BlockNumber
							res.BlockHash = receipt.BlockHash
						}
					})
				}
			}
			e.client = client

			results, err := e.checkUpkeeps(testutils.Context(t), tc.inputs)
			assert.Equal(t, tc.results, results)
			assert.Equal(t, tc.err, err)
		})
	}
}

func TestRegistry_SimulatePerformUpkeeps(t *testing.T) {
	uid0 := core.GenUpkeepID(types3.UpkeepType(0), "p0")
	uid1 := core.GenUpkeepID(types3.UpkeepType(1), "p1")
	uid2 := core.GenUpkeepID(types3.UpkeepType(1), "p2")

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

			mockReg := mocks.NewRegistry(t)
			mockReg.On("GetUpkeep", mock.Anything, mock.Anything).Return(
				encoding.UpkeepInfo{OffchainConfig: make([]byte, 0)},
				nil,
			).Times(2)
			e.registry = mockReg

			results, err := e.simulatePerformUpkeeps(testutils.Context(t), tc.inputs)
			assert.Equal(t, tc.results, results)
			assert.Equal(t, tc.err, err)
		})
	}
}

// setups up an evm registry for tests.
func setupEVMRegistry(t *testing.T) *EvmRegistry {
	lggr := logger.TestLogger(t)
	addr := common.HexToAddress("0x6cA639822c6C241Fa9A7A6b5032F6F7F1C513CAD")
	keeperRegistryABI, err := abi.JSON(strings.NewReader(ac.IAutomationV21PlusCommonABI))
	require.Nil(t, err, "need registry abi")
	streamsLookupCompatibleABI, err := abi.JSON(strings.NewReader(streams_lookup_compatible_interface.StreamsLookupCompatibleInterfaceABI))
	require.Nil(t, err, "need mercury abi")
	var logPoller logpoller.LogPoller
	mockReg := mocks.NewRegistry(t)
	mockHttpClient := mocks.NewHttpClient(t)
	client := evmClientMocks.NewClient(t)
	ge := gasMocks.NewEvmFeeEstimator(t)

	r := &EvmRegistry{
		lggr:         lggr,
		poller:       logPoller,
		addr:         addr,
		client:       client,
		logProcessed: make(map[string]bool),
		registry:     mockReg,
		abi:          keeperRegistryABI,
		active:       NewActiveUpkeepList(),
		packer:       encoding.NewAbiPacker(),
		headFunc:     func(ocr2keepers.BlockKey) {},
		chLog:        make(chan logpoller.Log, 1000),
		mercury: &MercuryConfig{
			cred: &types2.MercuryCredentials{
				LegacyURL: "https://google.old.com",
				URL:       "https://google.com",
				Username:  "FakeClientID",
				Password:  "FakeClientKey",
			},
			Abi:            streamsLookupCompatibleABI,
			AllowListCache: cache.New(defaultAllowListExpiration, cleanupInterval),
		},
		hc: mockHttpClient,
		bs: &BlockSubscriber{latestBlock: atomic.Pointer[ocr2keepers.BlockKey]{}},
		ge: ge,
	}
	return r
}
