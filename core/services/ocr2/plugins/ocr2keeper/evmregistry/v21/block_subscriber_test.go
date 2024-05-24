package evm

import (
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	ocr2keepers "github.com/smartcontractkit/chainlink-common/pkg/types/automation"

	htmocks "github.com/smartcontractkit/chainlink/v2/common/headtracker/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

const historySize = 4
const blockSize = int64(4)
const finality = uint32(4)

func TestBlockSubscriber_Subscribe(t *testing.T) {
	lggr := logger.TestLogger(t)
	var hb types.HeadBroadcaster
	var lp logpoller.LogPoller

	bs := NewBlockSubscriber(hb, lp, finality, lggr)
	bs.blockHistorySize = historySize
	bs.blockSize = blockSize
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

	bs := NewBlockSubscriber(hb, lp, finality, lggr)
	bs.blockHistorySize = historySize
	bs.blockSize = blockSize
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

	bs := NewBlockSubscriber(hb, lp, finality, lggr)
	bs.blockHistorySize = historySize
	bs.blockSize = blockSize
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
			ExpectedBlocks: []uint64{97, 98, 99, 100},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			lp := new(mocks.LogPoller)
			lp.On("LatestBlock", mock.Anything).Return(logpoller.LogPollerBlock{BlockNumber: tc.LatestBlock}, tc.LatestBlockErr)
			bs := NewBlockSubscriber(hb, lp, finality, lggr)
			bs.blockHistorySize = historySize
			bs.blockSize = blockSize
			blocks, err := bs.getBlockRange(testutils.Context(t))

			if tc.LatestBlockErr != nil {
				assert.Equal(t, tc.LatestBlockErr.Error(), err.Error())
			} else {
				assert.Equal(t, tc.ExpectedBlocks, blocks)
			}
		})
	}
}

func TestBlockSubscriber_InitializeBlocks(t *testing.T) {
	lggr := logger.TestLogger(t)
	var hb types.HeadBroadcaster

	tests := []struct {
		Name             string
		Blocks           []uint64
		PollerBlocks     []logpoller.LogPollerBlock
		LastClearedBlock int64
		Error            error
	}{
		{
			Name:  "failed to get latest block",
			Error: fmt.Errorf("failed to get log poller blocks"),
		},
		{
			Name:   "get block range",
			Blocks: []uint64{97, 98, 99, 100},
			PollerBlocks: []logpoller.LogPollerBlock{
				{
					BlockNumber: 97,
					BlockHash:   common.HexToHash("0x5e7fadfc14e1cfa9c05a91128c16a20c6cbc3be38b4723c3d482d44bf9c0e07b"),
				},
				{
					BlockNumber: 98,
					BlockHash:   common.HexToHash("0xaf3f8b36a27837e9f1ea3b4da7cdbf2ce0bdf7ef4e87d23add83b19438a2fcba"),
				},
				{
					BlockNumber: 99,
					BlockHash:   common.HexToHash("0xa7ac5bbc905b81f3a2ad9fb8ef1fe45f4a95768df456736952e4ec6c21296abe"),
				},
				{
					BlockNumber: 100,
					BlockHash:   common.HexToHash("0xa7ac5bbc905b81f3a2ad9fb8ef1fe45f4a95768df456736952e4ec6c21296abe"),
				},
			},
			LastClearedBlock: 96,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			lp := new(mocks.LogPoller)
			lp.On("GetBlocksRange", mock.Anything, tc.Blocks).Return(tc.PollerBlocks, tc.Error)
			bs := NewBlockSubscriber(hb, lp, finality, lggr)
			bs.blockHistorySize = historySize
			bs.blockSize = blockSize
			err := bs.initializeBlocks(testutils.Context(t), tc.Blocks)

			if tc.Error != nil {
				assert.Equal(t, tc.Error.Error(), err.Error())
			} else {
				for _, b := range tc.PollerBlocks {
					h, ok := bs.blocks[b.BlockNumber]
					assert.True(t, ok)
					assert.Equal(t, b.BlockHash.Hex(), h)
				}
				assert.Equal(t, tc.LastClearedBlock, bs.lastClearedBlock)
			}
		})
	}
}

func TestBlockSubscriber_BuildHistory(t *testing.T) {
	lggr := logger.TestLogger(t)
	var hb types.HeadBroadcaster
	lp := new(mocks.LogPoller)

	tests := []struct {
		Name            string
		Blocks          map[int64]string
		Block           int64
		ExpectedHistory ocr2keepers.BlockHistory
	}{
		{
			Name: "build history",
			Blocks: map[int64]string{
				100: "0xaf3f8b36a27837e9f1ea3b4da7cdbf2ce0bdf7ef4e87d23add83b19438a2fcba",
				98:  "0xc20c7b47466c081a44a3b168994e89affe85cb894547845d938f923b67c633c0",
				97:  "0xa7ac5bbc905b81f3a2ad9fb8ef1fe45f4a95768df456736952e4ec6c21296abe",
				95:  "0xc20c7b47466c081a44a3b168994e89affe85cb894547845d938f923b67c633c0",
			},
			Block: 100,
			ExpectedHistory: ocr2keepers.BlockHistory{
				ocr2keepers.BlockKey{
					Number: 100,
					Hash:   common.HexToHash("0xaf3f8b36a27837e9f1ea3b4da7cdbf2ce0bdf7ef4e87d23add83b19438a2fcba"),
				},
				ocr2keepers.BlockKey{
					Number: 98,
					Hash:   common.HexToHash("0xc20c7b47466c081a44a3b168994e89affe85cb894547845d938f923b67c633c0"),
				},
				ocr2keepers.BlockKey{
					Number: 97,
					Hash:   common.HexToHash("0xa7ac5bbc905b81f3a2ad9fb8ef1fe45f4a95768df456736952e4ec6c21296abe"),
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			bs := NewBlockSubscriber(hb, lp, finality, lggr)
			bs.blockHistorySize = historySize
			bs.blockSize = blockSize
			bs.blocks = tc.Blocks

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
		Name                     string
		Blocks                   map[int64]string
		LastClearedBlock         int64
		LastSentBlock            int64
		ExpectedLastClearedBlock int64
		ExpectedBlocks           map[int64]string
	}{
		{
			Name: "build history",
			Blocks: map[int64]string{
				102: "0xaf3f8b36a27837e9f1ea3b4da7cdbf2ce0bdf7ef4e87d23add83b19438a2fcba",
				100: "0xaf3f8b36a27837e9f1ea3b4da7cdbf2ce0bdf7ef4e87d23add83b19438a2fcba",
				98:  "0xc20c7b47466c081a44a3b168994e89affe85cb894547845d938f923b67c633c0",
				95:  "0xc20c7b47466c081a44a3b168994e89affe85cb894547845d938f923b67c633c0",
			},
			LastClearedBlock:         94,
			LastSentBlock:            101,
			ExpectedLastClearedBlock: 97,
			ExpectedBlocks: map[int64]string{
				102: "0xaf3f8b36a27837e9f1ea3b4da7cdbf2ce0bdf7ef4e87d23add83b19438a2fcba",
				100: "0xaf3f8b36a27837e9f1ea3b4da7cdbf2ce0bdf7ef4e87d23add83b19438a2fcba",
				98:  "0xc20c7b47466c081a44a3b168994e89affe85cb894547845d938f923b67c633c0",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			bs := NewBlockSubscriber(hb, lp, finality, lggr)
			bs.blockHistorySize = historySize
			bs.blockSize = blockSize
			bs.blocks = tc.Blocks
			bs.lastClearedBlock = tc.LastClearedBlock
			bs.lastSentBlock = tc.LastSentBlock
			bs.cleanup()

			assert.Equal(t, tc.ExpectedLastClearedBlock, bs.lastClearedBlock)
			assert.Equal(t, tc.ExpectedBlocks, bs.blocks)
		})
	}
}

func TestBlockSubscriber_Start(t *testing.T) {
	lggr := logger.TestLogger(t)
	hb := htmocks.NewHeadBroadcaster[*evmtypes.Head, common.Hash](t)
	hb.On("Subscribe", mock.Anything).Return(&evmtypes.Head{Number: 42}, func() {})
	lp := new(mocks.LogPoller)
	lp.On("LatestBlock", mock.Anything).Return(logpoller.LogPollerBlock{BlockNumber: 100}, nil)
	blocks := []uint64{97, 98, 99, 100}
	pollerBlocks := []logpoller.LogPollerBlock{
		{
			BlockNumber: 97,
			BlockHash:   common.HexToHash("0xda2f9d1359eadd7b93338703adc07d942021a78195564038321ef53f23f87333"),
		},
		{
			BlockNumber: 98,
			BlockHash:   common.HexToHash("0xc20c7b47466c081a44a3b168994e89affe85cb894547845d938f923b67c633c0"),
		},
		{
			BlockNumber: 99,
			BlockHash:   common.HexToHash("0x9bc2b51e147f9cad05f1614b7f1d8181cb24c544cbcf841f3155e54e752a3b44"),
		},
		{
			BlockNumber: 100,
			BlockHash:   common.HexToHash("0x5e7fadfc14e1cfa9c05a91128c16a20c6cbc3be38b4723c3d482d44bf9c0e07b"),
		},
	}

	lp.On("GetBlocksRange", mock.Anything, blocks).Return(pollerBlocks, nil)

	bs := NewBlockSubscriber(hb, lp, finality, lggr)
	bs.blockHistorySize = historySize
	bs.blockSize = blockSize
	err := bs.Start(testutils.Context(t))
	assert.Nil(t, err)

	h97 := evmtypes.Head{
		Number: 97,
		Hash:   common.HexToHash("0xda2f9d1359eadd7b93338703adc07d942021a78195564038321ef53f23f87333"),
		Parent: nil,
	}
	h98 := evmtypes.Head{
		Number: 98,
		Hash:   common.HexToHash("0xc20c7b47466c081a44a3b168994e89affe85cb894547845d938f923b67c633c0"),
		Parent: &h97,
	}
	h99 := evmtypes.Head{
		Number: 99,
		Hash:   common.HexToHash("0x9bc2b51e147f9cad05f1614b7f1d8181cb24c544cbcf841f3155e54e752a3b44"),
		Parent: &h98,
	}
	h100 := evmtypes.Head{
		Number: 100,
		Hash:   common.HexToHash("0x5e7fadfc14e1cfa9c05a91128c16a20c6cbc3be38b4723c3d482d44bf9c0e07b"),
		Parent: &h99,
	}

	// no subscribers yet
	bs.headC <- &h100

	expectedBlocks := map[int64]string{
		97:  "0xda2f9d1359eadd7b93338703adc07d942021a78195564038321ef53f23f87333",
		98:  "0xc20c7b47466c081a44a3b168994e89affe85cb894547845d938f923b67c633c0",
		99:  "0x9bc2b51e147f9cad05f1614b7f1d8181cb24c544cbcf841f3155e54e752a3b44",
		100: "0x5e7fadfc14e1cfa9c05a91128c16a20c6cbc3be38b4723c3d482d44bf9c0e07b",
	}

	// sleep 100 milli to wait for the go routine to finish
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, int64(historySize), bs.blockHistorySize)
	assert.Equal(t, int64(96), bs.lastClearedBlock)
	assert.Equal(t, int64(100), bs.lastSentBlock)
	assert.Equal(t, expectedBlocks, bs.blocks)

	// add 1 subscriber
	subId1, c1, err := bs.Subscribe()
	assert.Nil(t, err)
	assert.Equal(t, 1, subId1)

	h101 := &evmtypes.Head{
		Number: 101,
		Hash:   common.HexToHash("0xc20c7b47466c081a44a3b168994e89affe85cb894547845d938f923b67c633c0"),
		Parent: &h100,
	}
	bs.headC <- h101

	time.Sleep(100 * time.Millisecond)
	bk1 := <-c1
	assert.Equal(t, ocr2keepers.BlockHistory{
		ocr2keepers.BlockKey{
			Number: 101,
			Hash:   common.HexToHash("0xc20c7b47466c081a44a3b168994e89affe85cb894547845d938f923b67c633c0"),
		},
		ocr2keepers.BlockKey{
			Number: 100,
			Hash:   common.HexToHash("0x5e7fadfc14e1cfa9c05a91128c16a20c6cbc3be38b4723c3d482d44bf9c0e07b"),
		},
		ocr2keepers.BlockKey{
			Number: 99,
			Hash:   common.HexToHash("0x9bc2b51e147f9cad05f1614b7f1d8181cb24c544cbcf841f3155e54e752a3b44"),
		},
		ocr2keepers.BlockKey{
			Number: 98,
			Hash:   common.HexToHash("0xc20c7b47466c081a44a3b168994e89affe85cb894547845d938f923b67c633c0"),
		},
	}, bk1)

	// add 2nd subscriber
	subId2, c2, err := bs.Subscribe()
	assert.Nil(t, err)
	assert.Equal(t, 2, subId2)

	// re-org happens
	new99 := &evmtypes.Head{
		Number: 99,
		Hash:   common.HexToHash("0x70c03acc4ddbfb253ba41a25dc13fb21b25da8b63bcd1aa7fb55713d33a36c71"),
		Parent: &h98,
	}
	new100 := &evmtypes.Head{
		Number: 100,
		Hash:   common.HexToHash("0x8a876b62d252e63e16cf3487db3486c0a7c0a8e06bc3792a3b116c5ca480503f"),
		Parent: new99,
	}
	new101 := &evmtypes.Head{
		Number: 101,
		Hash:   common.HexToHash("0x41b5842b8847dcf834e39556d2ac51cc7d960a7de9471ec504673d0038fd6c8e"),
		Parent: new100,
	}

	new102 := &evmtypes.Head{
		Number: 102,
		Hash:   common.HexToHash("0x9ac1ebc307554cf1bcfcc2a49462278e89d6878d613a33df38a64d0aeac971b5"),
		Parent: new101,
	}

	bs.headC <- new102

	time.Sleep(100 * time.Millisecond)
	bk1 = <-c1
	assert.Equal(t,
		ocr2keepers.BlockHistory{
			ocr2keepers.BlockKey{
				Number: 102,
				Hash:   common.HexToHash("0x9ac1ebc307554cf1bcfcc2a49462278e89d6878d613a33df38a64d0aeac971b5"),
			},
			ocr2keepers.BlockKey{
				Number: 101,
				Hash:   common.HexToHash("0x41b5842b8847dcf834e39556d2ac51cc7d960a7de9471ec504673d0038fd6c8e"),
			},
			ocr2keepers.BlockKey{
				Number: 100,
				Hash:   common.HexToHash("0x8a876b62d252e63e16cf3487db3486c0a7c0a8e06bc3792a3b116c5ca480503f"),
			},
			ocr2keepers.BlockKey{
				Number: 99,
				Hash:   common.HexToHash("0x70c03acc4ddbfb253ba41a25dc13fb21b25da8b63bcd1aa7fb55713d33a36c71"),
			},
		},
		bk1,
	)

	bk2 := <-c2
	assert.Equal(t,
		ocr2keepers.BlockHistory{
			ocr2keepers.BlockKey{
				Number: 102,
				Hash:   common.HexToHash("0x9ac1ebc307554cf1bcfcc2a49462278e89d6878d613a33df38a64d0aeac971b5"),
			},
			ocr2keepers.BlockKey{
				Number: 101,
				Hash:   common.HexToHash("0x41b5842b8847dcf834e39556d2ac51cc7d960a7de9471ec504673d0038fd6c8e"),
			},
			ocr2keepers.BlockKey{
				Number: 100,
				Hash:   common.HexToHash("0x8a876b62d252e63e16cf3487db3486c0a7c0a8e06bc3792a3b116c5ca480503f"),
			},
			ocr2keepers.BlockKey{
				Number: 99,
				Hash:   common.HexToHash("0x70c03acc4ddbfb253ba41a25dc13fb21b25da8b63bcd1aa7fb55713d33a36c71"),
			},
		},
		bk2,
	)

	assert.Equal(t, int64(102), bs.lastSentBlock)
	assert.Equal(t, int64(96), bs.lastClearedBlock)
}
