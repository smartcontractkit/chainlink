package logprovider

import (
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-automation/pkg/v3/types"

	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/core"
)

func TestLogEventProvider_LifeCycle(t *testing.T) {
	tests := []struct {
		name           string
		errored        bool
		upkeepID       *big.Int
		upkeepCfg      LogTriggerConfig
		hasFilter      bool
		replyed        bool
		cfgUpdateBlock uint64
		mockPoller     bool
		unregister     bool
	}{
		{
			"new upkeep",
			false,
			big.NewInt(111),
			LogTriggerConfig{
				ContractAddress: common.BytesToAddress(common.LeftPadBytes([]byte{1, 2, 3, 4}, 20)),
				Topic0:          common.BytesToHash(common.LeftPadBytes([]byte{1, 2, 3, 4}, 32)),
			},
			false,
			true,
			uint64(1),
			true,
			false,
		},
		{
			"empty config",
			true,
			big.NewInt(111),
			LogTriggerConfig{},
			false,
			false,
			uint64(0),
			false,
			false,
		},
		{
			"invalid config",
			true,
			big.NewInt(111),
			LogTriggerConfig{
				ContractAddress: common.BytesToAddress(common.LeftPadBytes([]byte{}, 20)),
				Topic0:          common.BytesToHash(common.LeftPadBytes([]byte{}, 32)),
			},
			false,
			false,
			uint64(2),
			false,
			false,
		},
		{
			"existing config with old block",
			true,
			big.NewInt(111),
			LogTriggerConfig{
				ContractAddress: common.BytesToAddress(common.LeftPadBytes([]byte{1, 2, 3, 4}, 20)),
				Topic0:          common.BytesToHash(common.LeftPadBytes([]byte{1, 2, 3, 4}, 32)),
			},
			true,
			false,
			uint64(0),
			true,
			false,
		},
		{
			"existing config with newer block",
			false,
			big.NewInt(111),
			LogTriggerConfig{
				ContractAddress: common.BytesToAddress(common.LeftPadBytes([]byte{1, 2, 3, 4}, 20)),
				Topic0:          common.BytesToHash(common.LeftPadBytes([]byte{1, 2, 3, 4}, 32)),
			},
			true,
			false,
			uint64(2),
			true,
			true,
		},
	}

	p := NewLogProvider(logger.TestLogger(t), nil, big.NewInt(1), &mockedPacker{}, NewUpkeepFilterStore(), NewOptions(200, big.NewInt(1)))

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := testutils.Context(t)

			if tc.mockPoller {
				lp := new(mocks.LogPoller)
				lp.On("RegisterFilter", mock.Anything, mock.Anything).Return(nil)
				lp.On("UnregisterFilter", mock.Anything, mock.Anything).Return(nil)
				lp.On("LatestBlock", mock.Anything).Return(logpoller.Block{}, nil)
				hasFitlerTimes := 1
				if tc.unregister {
					hasFitlerTimes = 2
				}
				lp.On("HasFilter", p.filterName(tc.upkeepID)).Return(tc.hasFilter).Times(hasFitlerTimes)
				if tc.replyed {
					lp.On("ReplayAsync", mock.Anything).Return(nil).Times(1)
				} else {
					lp.On("ReplayAsync", mock.Anything).Return(nil).Times(0)
				}
				p.lock.Lock()
				p.poller = lp
				p.lock.Unlock()
			}

			err := p.RegisterFilter(ctx, FilterOptions{
				UpkeepID:      tc.upkeepID,
				TriggerConfig: tc.upkeepCfg,
				UpdateBlock:   tc.cfgUpdateBlock,
			})
			if tc.errored {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				if tc.unregister {
					require.NoError(t, p.UnregisterFilter(ctx, tc.upkeepID))
				}
			}
		})
	}
}

func TestEventLogProvider_RefreshActiveUpkeeps(t *testing.T) {
	ctx := testutils.Context(t)
	mp := new(mocks.LogPoller)
	mp.On("RegisterFilter", mock.Anything, mock.Anything).Return(nil)
	mp.On("UnregisterFilter", mock.Anything, mock.Anything).Return(nil)
	mp.On("HasFilter", mock.Anything).Return(false)
	mp.On("LatestBlock", mock.Anything).Return(logpoller.Block{}, nil)
	mp.On("ReplayAsync", mock.Anything).Return(nil)

	p := NewLogProvider(logger.TestLogger(t), mp, big.NewInt(1), &mockedPacker{}, NewUpkeepFilterStore(), NewOptions(200, big.NewInt(1)))

	require.NoError(t, p.RegisterFilter(ctx, FilterOptions{
		UpkeepID: core.GenUpkeepID(types.LogTrigger, "1111").BigInt(),
		TriggerConfig: LogTriggerConfig{
			ContractAddress: common.BytesToAddress(common.LeftPadBytes([]byte{1, 2, 3, 4}, 20)),
			Topic0:          common.BytesToHash(common.LeftPadBytes([]byte{1, 2, 3, 4}, 32)),
		},
		UpdateBlock: uint64(0),
	}))
	require.NoError(t, p.RegisterFilter(ctx, FilterOptions{
		UpkeepID: core.GenUpkeepID(types.LogTrigger, "2222").BigInt(),
		TriggerConfig: LogTriggerConfig{
			ContractAddress: common.BytesToAddress(common.LeftPadBytes([]byte{1, 2, 3, 4}, 20)),
			Topic0:          common.BytesToHash(common.LeftPadBytes([]byte{1, 2, 3, 4}, 32)),
		},
		UpdateBlock: uint64(0),
	}))
	require.Equal(t, 2, p.filterStore.Size())

	newIds, err := p.RefreshActiveUpkeeps(ctx)
	require.NoError(t, err)
	require.Empty(t, newIds)
	mp.On("HasFilter", p.filterName(core.GenUpkeepID(types.LogTrigger, "2222").BigInt())).Return(true)
	newIds, err = p.RefreshActiveUpkeeps(
		ctx,
		core.GenUpkeepID(types.LogTrigger, "2222").BigInt(),
		core.GenUpkeepID(types.LogTrigger, "1234").BigInt(),
		core.GenUpkeepID(types.LogTrigger, "123").BigInt())
	require.NoError(t, err)
	require.Len(t, newIds, 2)
	require.Equal(t, 1, p.filterStore.Size())
}

func TestLogEventProvider_ValidateLogTriggerConfig(t *testing.T) {
	contractAddress := common.HexToAddress("0xB9F3af0c2CbfE108efd0E23F7b0a151Ea42f764E")
	eventSig := common.HexToHash("0x3bdab8bffae631cfee411525ebae27f3fb61b10c662c09ec2a7dbb5854c87e8c")
	tests := []struct {
		name        string
		cfg         LogTriggerConfig
		expectedErr error
	}{
		{
			"success",
			LogTriggerConfig{
				ContractAddress: contractAddress,
				FilterSelector:  0,
				Topic0:          eventSig,
			},
			nil,
		},
		{
			"invalid contract address",
			LogTriggerConfig{
				ContractAddress: common.Address{},
				FilterSelector:  0,
				Topic0:          eventSig,
			},
			errors.New("invalid contract address: zeroed"),
		},
		{
			"invalid topic0",
			LogTriggerConfig{
				ContractAddress: contractAddress,
				FilterSelector:  0,
			},
			errors.New("invalid topic0: zeroed"),
		},
		{
			"success",
			LogTriggerConfig{
				ContractAddress: contractAddress,
				FilterSelector:  8,
				Topic0:          eventSig,
			},
			errors.New("invalid filter selector: larger or equal to 8"),
		},
	}

	p := NewLogProvider(logger.TestLogger(t), nil, big.NewInt(1), &mockedPacker{}, NewUpkeepFilterStore(), NewOptions(200, big.NewInt(1)))
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := p.validateLogTriggerConfig(tc.cfg)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
