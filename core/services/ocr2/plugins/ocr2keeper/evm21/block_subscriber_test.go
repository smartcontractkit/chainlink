package evm

import (
	"context"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

const blockHistorySize = 4

func TestBlockSubscriber_Subscribe(t *testing.T) {
	lggr := logger.TestLogger(t)
	var hb types.HeadBroadcaster
	var lp logpoller.LogPoller

	bs := NewBlockSubscriber(hb, lp, blockHistorySize, lggr)
	subId, _, err := bs.Subscribe()
	assert.Nil(t, err)
	assert.Equal(t, subId, 1)
	subId, _, err = bs.Subscribe()
	assert.Nil(t, err)
	assert.Equal(t, subId, 2)
	subId, _, err = bs.Subscribe()
	assert.Nil(t, err)
	assert.Equal(t, subId, 3)
}

func TestBlockSubscriber_Unsubscribe(t *testing.T) {
	lggr := logger.TestLogger(t)
	var hb types.HeadBroadcaster
	var lp logpoller.LogPoller

	bs := NewBlockSubscriber(hb, lp, blockHistorySize, lggr)
	subId, _, err := bs.Subscribe()
	assert.Nil(t, err)
	assert.Equal(t, subId, 1)
	subId, _, err = bs.Subscribe()
	assert.Nil(t, err)
	assert.Equal(t, subId, 2)
	err = bs.Unsubscribe(1)
	assert.Nil(t, err)
}

func TestBlockSubscriber_Unsubscribe_Failure(t *testing.T) {
	lggr := logger.TestLogger(t)
	var hb types.HeadBroadcaster
	var lp logpoller.LogPoller

	bs := NewBlockSubscriber(hb, lp, blockHistorySize, lggr)
	err := bs.Unsubscribe(2)
	assert.Equal(t, err.Error(), "subscriber 2 does not exist")
}

func TestBlockSubscriber_GetBlockRange(t *testing.T) {
	lggr := logger.TestLogger(t)
	var hb types.HeadBroadcaster

	tests := []struct {
		Name           string
		LatestBlock    int64
		LatestBlockErr error
		ExpectedBlocks []uint64
	}{
		{
			Name:           "failed to get latest block",
			LatestBlockErr: fmt.Errorf("failed to get latest block"),
		},
		{
			Name:           "get block range",
			LatestBlock:    100,
			ExpectedBlocks: []uint64{100, 99, 98, 97},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			lp := new(mocks.LogPoller)
			lp.On("LatestBlock", mock.Anything).Return(tc.LatestBlock, tc.LatestBlockErr)
			bs := NewBlockSubscriber(hb, lp, blockHistorySize, lggr)
			blocks, err := bs.getBlockRange(context.Background())

			if tc.LatestBlockErr != nil {
				assert.Equal(t, tc.LatestBlockErr.Error(), err.Error())
			} else {
				assert.Equal(t, tc.ExpectedBlocks, blocks)
			}
		})
	}
}

func TestBlockSubscriber_GetLogPollerBlocks(t *testing.T) {
	lggr := logger.TestLogger(t)
	var hb types.HeadBroadcaster

	tests := []struct {
		Name         string
		Blocks       []uint64
		PollerBlocks []logpoller.LogPollerBlock
		Error        error
	}{
		{
			Name:  "failed to get latest block",
			Error: fmt.Errorf("failed to get log poller blocks"),
		},
		{
			Name:   "get block range",
			Blocks: []uint64{100, 99, 98, 97},
			PollerBlocks: []logpoller.LogPollerBlock{
				{
					BlockNumber: 100,
					BlockHash:   common.HexToHash("0x5e7fadfc14e1cfa9c05a91128c16a20c6cbc3be38b4723c3d482d44bf9c0e07b"),
				},
				{
					BlockNumber: 99,
					BlockHash:   common.HexToHash("0xaf3f8b36a27837e9f1ea3b4da7cdbf2ce0bdf7ef4e87d23add83b19438a2fcba"),
				},
				{
					BlockNumber: 98,
					BlockHash:   common.HexToHash("0xa7ac5bbc905b81f3a2ad9fb8ef1fe45f4a95768df456736952e4ec6c21296abe"),
				},
				{
					BlockNumber: 97,
					BlockHash:   common.HexToHash("0xa7ac5bbc905b81f3a2ad9fb8ef1fe45f4a95768df456736952e4ec6c21296abe"),
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			lp := new(mocks.LogPoller)
			lp.On("GetBlocksRange", mock.Anything, tc.Blocks, mock.Anything).Return(tc.PollerBlocks, tc.Error)
			bs := NewBlockSubscriber(hb, lp, blockHistorySize, lggr)
			err := bs.getLogPollerBlocks(context.Background(), tc.Blocks)

			if tc.Error != nil {
				assert.Equal(t, tc.Error.Error(), err.Error())
			} else {
				for _, b := range tc.PollerBlocks {
					h, ok := bs.blocksFromPoller[b.BlockNumber]
					assert.True(t, ok)
					assert.Equal(t, b.BlockHash, h)
				}
			}
		})
	}
}

func TestBlockSubscriber_BuildHistory(t *testing.T) {
	lggr := logger.TestLogger(t)
	var hb types.HeadBroadcaster
	lp := new(mocks.LogPoller)

	tests := []struct {
		Name                  string
		BlocksFromLogPoller   map[int64]common.Hash
		BlocksFromBroadcaster map[int64]common.Hash
		Block                 int64
		ExpectedHistory       ocr2keepers.BlockHistory
	}{
		{
			Name: "build history",
			BlocksFromLogPoller: map[int64]common.Hash{
				100: common.HexToHash("0x5e7fadfc14e1cfa9c05a91128c16a20c6cbc3be38b4723c3d482d44bf9c0e07b"),
				97:  common.HexToHash("0xa7ac5bbc905b81f3a2ad9fb8ef1fe45f4a95768df456736952e4ec6c21296abe"),
				96:  common.HexToHash("0x44f23c588193695abd92697ddc1ba032376d0a784818eddd2d159eee4c41f03f"),
			},
			BlocksFromBroadcaster: map[int64]common.Hash{
				100: common.HexToHash("0xaf3f8b36a27837e9f1ea3b4da7cdbf2ce0bdf7ef4e87d23add83b19438a2fcba"),
				98:  common.HexToHash("0xc20c7b47466c081a44a3b168994e89affe85cb894547845d938f923b67c633c0"),
			},
			Block: 100,
			ExpectedHistory: ocr2keepers.BlockHistory{
				"100|0x5e7fadfc14e1cfa9c05a91128c16a20c6cbc3be38b4723c3d482d44bf9c0e07b",
				"98|0xc20c7b47466c081a44a3b168994e89affe85cb894547845d938f923b67c633c0",
				"97|0xa7ac5bbc905b81f3a2ad9fb8ef1fe45f4a95768df456736952e4ec6c21296abe",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			bs := NewBlockSubscriber(hb, lp, blockHistorySize, lggr)
			bs.blocksFromPoller = tc.BlocksFromLogPoller
			bs.blocksFromBroadcaster = tc.BlocksFromBroadcaster

			history := bs.buildHistory(tc.Block)
			assert.Equal(t, history, tc.ExpectedHistory)
		})
	}
}

func TestBlockSubscriber_Cleanup(t *testing.T) {
	lggr := logger.TestLogger(t)
	var hb types.HeadBroadcaster
	lp := new(mocks.LogPoller)

	tests := []struct {
		Name                      string
		BlocksFromLogPoller       map[int64]common.Hash
		BlocksFromBroadcaster     map[int64]common.Hash
		LastClearedBlock          int64
		LastSentBlock             int64
		ExpectedLastClearedBlock  int64
		ExpectedLogPollerBlocks   map[int64]common.Hash
		ExpectedBroadcasterBlocks map[int64]common.Hash
	}{
		{
			Name: "build history",
			BlocksFromLogPoller: map[int64]common.Hash{
				101: common.HexToHash("0x5e7fadfc14e1cfa9c05a91128c16a20c6cbc3be38b4723c3d482d44bf9c0e07b"),
				100: common.HexToHash("0x5e7fadfc14e1cfa9c05a91128c16a20c6cbc3be38b4723c3d482d44bf9c0e07b"),
				97:  common.HexToHash("0xa7ac5bbc905b81f3a2ad9fb8ef1fe45f4a95768df456736952e4ec6c21296abe"),
				96:  common.HexToHash("0x44f23c588193695abd92697ddc1ba032376d0a784818eddd2d159eee4c41f03f"),
				95:  common.HexToHash("0x44f23c588193695abd92697ddc1ba032376d0a784818eddd2d159eee4c41f03f"),
			},
			BlocksFromBroadcaster: map[int64]common.Hash{
				102: common.HexToHash("0xaf3f8b36a27837e9f1ea3b4da7cdbf2ce0bdf7ef4e87d23add83b19438a2fcba"),
				100: common.HexToHash("0xaf3f8b36a27837e9f1ea3b4da7cdbf2ce0bdf7ef4e87d23add83b19438a2fcba"),
				98:  common.HexToHash("0xc20c7b47466c081a44a3b168994e89affe85cb894547845d938f923b67c633c0"),
				95:  common.HexToHash("0xc20c7b47466c081a44a3b168994e89affe85cb894547845d938f923b67c633c0"),
			},
			LastClearedBlock:         94,
			LastSentBlock:            101,
			ExpectedLastClearedBlock: 97,
			ExpectedLogPollerBlocks: map[int64]common.Hash{
				101: common.HexToHash("0x5e7fadfc14e1cfa9c05a91128c16a20c6cbc3be38b4723c3d482d44bf9c0e07b"),
				100: common.HexToHash("0x5e7fadfc14e1cfa9c05a91128c16a20c6cbc3be38b4723c3d482d44bf9c0e07b"),
			},
			ExpectedBroadcasterBlocks: map[int64]common.Hash{
				102: common.HexToHash("0xaf3f8b36a27837e9f1ea3b4da7cdbf2ce0bdf7ef4e87d23add83b19438a2fcba"),
				100: common.HexToHash("0xaf3f8b36a27837e9f1ea3b4da7cdbf2ce0bdf7ef4e87d23add83b19438a2fcba"),
				98:  common.HexToHash("0xc20c7b47466c081a44a3b168994e89affe85cb894547845d938f923b67c633c0"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			bs := NewBlockSubscriber(hb, lp, blockHistorySize, lggr)
			bs.blocksFromPoller = tc.BlocksFromLogPoller
			bs.blocksFromBroadcaster = tc.BlocksFromBroadcaster
			bs.lastClearedBlock = tc.LastClearedBlock
			bs.lastSentBlock = tc.LastSentBlock
			bs.cleanup()

			assert.Equal(t, tc.ExpectedLastClearedBlock, bs.lastClearedBlock)
			assert.Equal(t, tc.ExpectedBroadcasterBlocks, bs.blocksFromBroadcaster)
			assert.Equal(t, tc.ExpectedLogPollerBlocks, bs.blocksFromPoller)
		})
	}
}
