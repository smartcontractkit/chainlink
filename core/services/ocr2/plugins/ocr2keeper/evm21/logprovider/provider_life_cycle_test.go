package logprovider

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestLogEventProvider_LifeCycle(t *testing.T) {
	tests := []struct {
		name       string
		errored    bool
		upkeepID   *big.Int
		upkeepCfg  LogTriggerConfig
		mockPoller bool
	}{
		{
			"new upkeep",
			false,
			big.NewInt(111),
			LogTriggerConfig{
				ContractAddress: common.BytesToAddress(common.LeftPadBytes([]byte{1, 2, 3, 4}, 20)),
				Topic0:          common.BytesToHash(common.LeftPadBytes([]byte{1, 2, 3, 4}, 32)),
			},
			true,
		},
		{
			"empty config",
			true,
			big.NewInt(111),
			LogTriggerConfig{},
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
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mp := new(mocks.LogPoller)
			if tc.mockPoller {
				mp.On("RegisterFilter", mock.Anything).Return(nil)
				mp.On("UnregisterFilter", mock.Anything, mock.Anything).Return(nil)
			}
			p := New(logger.TestLogger(t), mp, &mockedPacker{}, NewUpkeepFilterStore(), nil)
			err := p.RegisterFilter(tc.upkeepID, tc.upkeepCfg)
			if tc.errored {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NoError(t, p.UnregisterFilter(tc.upkeepID))
			}
		})
	}
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

	p := New(logger.TestLogger(t), nil, &mockedPacker{}, NewUpkeepFilterStore(), nil)

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
