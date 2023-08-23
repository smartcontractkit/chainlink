package logprovider

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/core"
)

func TestLogEventProvider_LifeCycle(t *testing.T) {
	tests := []struct {
		name           string
		errored        bool
		upkeepID       *big.Int
		upkeepCfg      LogTriggerConfig
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
			uint64(1),
			true,
			false,
		},
		{
			"empty config",
			true,
			big.NewInt(111),
			LogTriggerConfig{},
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
			uint64(2),
			false,
			false,
		},
		{
			"existing config",
			true,
			big.NewInt(111),
			LogTriggerConfig{
				ContractAddress: common.BytesToAddress(common.LeftPadBytes([]byte{1, 2, 3, 4}, 20)),
				Topic0:          common.BytesToHash(common.LeftPadBytes([]byte{1, 2, 3, 4}, 32)),
			},
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
			uint64(2),
			true,
			true,
		},
	}

	mp := new(mocks.LogPoller)
	mp.On("RegisterFilter", mock.Anything).Return(nil)
	mp.On("UnregisterFilter", mock.Anything).Return(nil)
	mp.On("ReplayAsync", mock.Anything).Return(nil)
	p := NewLogProvider(logger.TestLogger(t), mp, &mockedPacker{}, NewUpkeepFilterStore(), nil)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := p.RegisterFilter(FilterOptions{
				UpkeepID:      tc.upkeepID,
				TriggerConfig: tc.upkeepCfg,
				UpdateBlock:   tc.cfgUpdateBlock,
			})
			if tc.errored {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				if tc.unregister {
					require.NoError(t, p.UnregisterFilter(tc.upkeepID))
				}
			}
		})
	}
}

func TestEventLogProvider_RefreshActiveUpkeeps(t *testing.T) {
	mp := new(mocks.LogPoller)
	mp.On("RegisterFilter", mock.Anything).Return(nil)
	mp.On("UnregisterFilter", mock.Anything).Return(nil)
	mp.On("ReplayAsync", mock.Anything).Return(nil)

	p := NewLogProvider(logger.TestLogger(t), mp, &mockedPacker{}, NewUpkeepFilterStore(), nil)

	require.NoError(t, p.RegisterFilter(FilterOptions{
		UpkeepID: core.GenUpkeepID(ocr2keepers.LogTrigger, "1111").BigInt(),
		TriggerConfig: LogTriggerConfig{
			ContractAddress: common.BytesToAddress(common.LeftPadBytes([]byte{1, 2, 3, 4}, 20)),
			Topic0:          common.BytesToHash(common.LeftPadBytes([]byte{1, 2, 3, 4}, 32)),
		},
		UpdateBlock: uint64(0),
	}))
	require.NoError(t, p.RegisterFilter(FilterOptions{
		UpkeepID: core.GenUpkeepID(ocr2keepers.LogTrigger, "2222").BigInt(),
		TriggerConfig: LogTriggerConfig{
			ContractAddress: common.BytesToAddress(common.LeftPadBytes([]byte{1, 2, 3, 4}, 20)),
			Topic0:          common.BytesToHash(common.LeftPadBytes([]byte{1, 2, 3, 4}, 32)),
		},
		UpdateBlock: uint64(0),
	}))
	require.Equal(t, 2, p.filterStore.Size())

	newIds, err := p.RefreshActiveUpkeeps()
	require.NoError(t, err)
	require.Len(t, newIds, 0)
	newIds, err = p.RefreshActiveUpkeeps(
		core.GenUpkeepID(ocr2keepers.LogTrigger, "2222").BigInt(),
		core.GenUpkeepID(ocr2keepers.LogTrigger, "1234").BigInt(),
		core.GenUpkeepID(ocr2keepers.LogTrigger, "123").BigInt())
	require.NoError(t, err)
	require.Len(t, newIds, 2)
	require.Equal(t, 1, p.filterStore.Size())
}

func TestLogEventProvider_GetFiltersBySelector(t *testing.T) {
	var zeroBytes [32]byte
	tests := []struct {
		name           string
		filterSelector uint8
		filters        [][]byte
		expectedSigs   []common.Hash
	}{
		{
			"invalid filters",
			1,
			[][]byte{
				zeroBytes[:],
			},
			[]common.Hash{},
		},
		{
			"selector 000",
			0,
			[][]byte{
				{1},
			},
			[]common.Hash{},
		},
		{
			"selector 001",
			1,
			[][]byte{
				{1},
				{2},
				{3},
			},
			[]common.Hash{
				common.BytesToHash(common.LeftPadBytes([]byte{1}, 32)),
			},
		},
		{
			"selector 010",
			2,
			[][]byte{
				{1},
				{2},
				{3},
			},
			[]common.Hash{
				common.BytesToHash(common.LeftPadBytes([]byte{2}, 32)),
			},
		},
		{
			"selector 011",
			3,
			[][]byte{
				{1},
				{2},
				{3},
			},
			[]common.Hash{
				common.BytesToHash(common.LeftPadBytes([]byte{1}, 32)),
				common.BytesToHash(common.LeftPadBytes([]byte{2}, 32)),
			},
		},
		{
			"selector 100",
			4,
			[][]byte{
				{1},
				{2},
				{3},
			},
			[]common.Hash{
				common.BytesToHash(common.LeftPadBytes([]byte{3}, 32)),
			},
		},
		{
			"selector 101",
			5,
			[][]byte{
				{1},
				{2},
				{3},
			},
			[]common.Hash{
				common.BytesToHash(common.LeftPadBytes([]byte{1}, 32)),
				common.BytesToHash(common.LeftPadBytes([]byte{3}, 32)),
			},
		},
		{
			"selector 110",
			6,
			[][]byte{
				{1},
				{2},
				{3},
			},
			[]common.Hash{
				common.BytesToHash(common.LeftPadBytes([]byte{2}, 32)),
				common.BytesToHash(common.LeftPadBytes([]byte{3}, 32)),
			},
		},
		{
			"selector 111",
			7,
			[][]byte{
				{1},
				{2},
				{3},
			},
			[]common.Hash{
				common.BytesToHash(common.LeftPadBytes([]byte{1}, 32)),
				common.BytesToHash(common.LeftPadBytes([]byte{2}, 32)),
				common.BytesToHash(common.LeftPadBytes([]byte{3}, 32)),
			},
		},
	}

	p := NewLogProvider(logger.TestLogger(t), nil, &mockedPacker{}, NewUpkeepFilterStore(), nil)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sigs := p.getFiltersBySelector(tc.filterSelector, tc.filters...)
			if len(sigs) != len(tc.expectedSigs) {
				t.Fatalf("expected %v, got %v", len(tc.expectedSigs), len(sigs))
			}
			for i := range sigs {
				if sigs[i] != tc.expectedSigs[i] {
					t.Fatalf("expected %v, got %v", tc.expectedSigs[i], sigs[i])
				}
			}
		})
	}
}
