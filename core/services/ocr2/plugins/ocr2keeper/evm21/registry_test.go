package evm

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
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
		{Name: "AvailableActiveIDs", LatestHead: 1, ActiveIDs: []string{"8", "9", "3", "1"}, ExpectedKeys: []ocr2keepers.UpkeepIdentifier{
			ocr2keepers.UpkeepIdentifier("8"),
			ocr2keepers.UpkeepIdentifier("9"),
			ocr2keepers.UpkeepIdentifier("3"),
			ocr2keepers.UpkeepIdentifier("1"),
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
				ocr2keepers.UpkeepIdentifier("32329108151019397958065800113404894502874153543356521479058624064899121404671"),
			},
			Triggers: []uint8{uint8(logTrigger)},
		},
		{
			Name:       "get conditional upkeeps",
			LatestHead: 1,
			ActiveIDs:  []string{"8", "32329108151019397958065800113404894502874153543356521479058624064899121404671"},
			ExpectedKeys: []ocr2keepers.UpkeepIdentifier{
				ocr2keepers.UpkeepIdentifier("8"),
			},
			Triggers: []uint8{uint8(conditionTrigger)},
		},
		{
			Name:       "get multiple types of upkeeps",
			LatestHead: 1,
			ActiveIDs:  []string{"8", "32329108151019397958065800113404894502874153543356521479058624064899121404671"},
			ExpectedKeys: []ocr2keepers.UpkeepIdentifier{
				ocr2keepers.UpkeepIdentifier("8"),
				ocr2keepers.UpkeepIdentifier("32329108151019397958065800113404894502874153543356521479058624064899121404671"),
			},
			Triggers: []uint8{uint8(logTrigger), uint8(conditionTrigger)},
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

	maxBigInt, _ := new(big.Int).SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 10)
	overMaxBigInt, _ := new(big.Int).SetString("115792089237316195423570985008687907853269984665640564039457584007913129639936", 10)

	tests := []struct {
		name       string
		input      ocr2keepers.UpkeepPayload
		latest     int64
		wantBlock  *big.Int
		wantUpkeep *big.Int
		wantErr    bool
	}{
		{
			"happy flow",
			ocr2keepers.UpkeepPayload{
				Upkeep: ocr2keepers.ConfiguredUpkeep{
					ID: ocr2keepers.UpkeepIdentifier(big.NewInt(10).Bytes()),
				},
				Trigger: ocr2keepers.Trigger{
					BlockNumber: 1,
				},
			},
			10,
			big.NewInt(1),
			big.NewInt(10),
			true,
		},
		{
			"maximum size integer ok",
			ocr2keepers.UpkeepPayload{
				Upkeep: ocr2keepers.ConfiguredUpkeep{
					ID: ocr2keepers.UpkeepIdentifier(maxBigInt.Bytes()),
				},
				Trigger: ocr2keepers.Trigger{
					BlockNumber: 1,
				},
			},
			10,
			big.NewInt(1),
			maxBigInt,
			true,
		},
		{
			"block number too high",
			ocr2keepers.UpkeepPayload{
				Upkeep: ocr2keepers.ConfiguredUpkeep{
					ID: ocr2keepers.UpkeepIdentifier(big.NewInt(10).Bytes()),
				},
				Trigger: ocr2keepers.Trigger{
					BlockNumber: 1000,
				},
			},
			999,
			nil,
			nil,
			false,
		},
		{
			"block number too low",
			ocr2keepers.UpkeepPayload{
				Upkeep: ocr2keepers.ConfiguredUpkeep{
					ID: ocr2keepers.UpkeepIdentifier(big.NewInt(10).Bytes()),
				},
				Trigger: ocr2keepers.Trigger{
					BlockNumber: 100,
				},
			},
			357,
			nil,
			nil,
			false,
		},
		{
			"empty block number",
			ocr2keepers.UpkeepPayload{
				Upkeep: ocr2keepers.ConfiguredUpkeep{
					ID: ocr2keepers.UpkeepIdentifier(big.NewInt(10).Bytes()),
				},
			},
			5000,
			nil,
			nil,
			false,
		},
		{
			"empty payload",
			ocr2keepers.UpkeepPayload{},
			10,
			nil,
			nil,
			false,
		},
		{
			"upkeep id as text can be parsed as a number because big int reads the raw bytes and interprets them as a number",
			ocr2keepers.UpkeepPayload{
				Upkeep: ocr2keepers.ConfiguredUpkeep{
					ID: ocr2keepers.UpkeepIdentifier("test"),
				},
				Trigger: ocr2keepers.Trigger{
					BlockNumber: 1,
				},
			},
			100,
			big.NewInt(1),
			big.NewInt(1952805748),
			true,
		},
		{
			"upkeep id larger than largest value should fail",
			ocr2keepers.UpkeepPayload{
				Upkeep: ocr2keepers.ConfiguredUpkeep{
					ID: ocr2keepers.UpkeepIdentifier(overMaxBigInt.Bytes()),
				},
				Trigger: ocr2keepers.Trigger{
					BlockNumber: 1,
				},
			},
			100,
			nil,
			nil,
			false,
		},
		{
			"upkeep id parsing with bytes is interpreted as absolute value",
			ocr2keepers.UpkeepPayload{
				Upkeep: ocr2keepers.ConfiguredUpkeep{
					ID: ocr2keepers.UpkeepIdentifier(big.NewInt(-12).Bytes()),
				},
				Trigger: ocr2keepers.Trigger{
					BlockNumber: 1,
				},
			},
			100,
			big.NewInt(1),
			big.NewInt(12),
			true,
		},
		{
			"upkeep id should not be zero",
			ocr2keepers.UpkeepPayload{
				Upkeep: ocr2keepers.ConfiguredUpkeep{
					ID: ocr2keepers.UpkeepIdentifier(big.NewInt(0).Bytes()),
				},
				Trigger: ocr2keepers.Trigger{
					BlockNumber: 1,
				},
			},
			100,
			nil,
			nil,
			false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			block, upkeep, err := r.getBlockAndUpkeepId(tc.input, tc.latest)

			if tc.wantBlock != nil {
				assert.Equal(t, tc.wantBlock.String(), block.String(), "block number should match expected")
			} else {
				assert.Nil(t, block, "block number should be nil")
			}

			if tc.wantUpkeep != nil {
				assert.Equal(t, tc.wantUpkeep.String(), upkeep.String(), "upkeep id should match expected")
			} else {
				assert.Nil(t, upkeep, "upkeep id should be nil")
			}

			assert.Equal(t, err == nil, tc.wantErr, "err nil should match expected")
		})
	}
}
