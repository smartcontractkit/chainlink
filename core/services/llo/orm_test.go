package llo

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
)

func Test_ORM(t *testing.T) {
	const ETHMainnetChainSelector uint64 = 5009297550715157269
	const OtherChainSelector uint64 = 1234567890

	db := pgtest.NewSqlxDB(t)
	orm := NewORM(db, ETHMainnetChainSelector)
	ctx := testutils.Context(t)

	addr1 := testutils.NewAddress()
	addr2 := testutils.NewAddress()
	addr3 := testutils.NewAddress()

	donID1 := uint32(1)
	donID2 := uint32(2)

	t.Run("LoadChannelDefinitions", func(t *testing.T) {
		t.Run("returns zero values if nothing in database", func(t *testing.T) {
			pd, err := orm.LoadChannelDefinitions(ctx, addr1, donID1)
			assert.NoError(t, err)
			assert.Nil(t, pd)
		})
		t.Run("loads channel definitions from database for the given don ID", func(t *testing.T) {
			expectedBlockNum := rand.Int63()
			expectedBlockNum2 := rand.Int63()
			cid1 := rand.Uint32()
			cid2 := rand.Uint32()

			channelDefsJSON := fmt.Sprintf(`
{
	"%d": {
		"reportFormat": 42,
		"chainSelector": 142,
		"streams": [{"streamId": 1, "aggregator": "median"}, {"streamId": 2, "aggregator": "mode"}],
		"opts": {"foo":"bar"}
	},
	"%d": {
		"reportFormat": 43,
		"chainSelector": 142,
		"streams": [{"streamId": 1, "aggregator": "median"}, {"streamId": 3, "aggregator": "quote"}]
	}
}
			`, cid1, cid2)
			pgtest.MustExec(t, db, `
			INSERT INTO channel_definitions(addr, chain_selector, don_id, definitions, block_num, version, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, NOW())
			`, addr1, ETHMainnetChainSelector, 1, channelDefsJSON, expectedBlockNum, 1)

			pgtest.MustExec(t, db, `
			INSERT INTO channel_definitions(addr, chain_selector, don_id, definitions, block_num, version, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, NOW())
			`, addr2, ETHMainnetChainSelector, 1, `{}`, expectedBlockNum2, 1)

			{
				// alternative chain selector; we expect these ones to be ignored
				pgtest.MustExec(t, db, `
			INSERT INTO channel_definitions(addr, chain_selector, don_id, definitions, block_num, version, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, NOW())
			`, addr1, OtherChainSelector, 1, channelDefsJSON, expectedBlockNum, 1)
				pgtest.MustExec(t, db, `
			INSERT INTO channel_definitions(addr, chain_selector, don_id, definitions, block_num, version, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, NOW())
			`, addr3, OtherChainSelector, 1, channelDefsJSON, expectedBlockNum, 1)
			}

			pd, err := orm.LoadChannelDefinitions(ctx, addr1, donID1)
			require.NoError(t, err)

			assert.Equal(t, ETHMainnetChainSelector, pd.ChainSelector)
			assert.Equal(t, addr1, pd.Address)
			assert.Equal(t, expectedBlockNum, pd.BlockNum)
			assert.Equal(t, donID1, pd.DonID)
			assert.Equal(t, uint32(1), pd.Version)
			assert.Equal(t, llotypes.ChannelDefinitions{
				cid1: llotypes.ChannelDefinition{
					ReportFormat: 42,
					Streams:      []llotypes.Stream{{StreamID: 1, Aggregator: llotypes.AggregatorMedian}, {StreamID: 2, Aggregator: llotypes.AggregatorMode}},
					Opts:         []byte(`{"foo":"bar"}`),
				},
				cid2: llotypes.ChannelDefinition{
					ReportFormat: 43,
					Streams:      []llotypes.Stream{{StreamID: 1, Aggregator: llotypes.AggregatorMedian}, {StreamID: 3, Aggregator: llotypes.AggregatorQuote}},
				},
			}, pd.Definitions)

			// does not load erroneously for a different address
			pd, err = orm.LoadChannelDefinitions(ctx, addr2, donID1)
			require.NoError(t, err)

			assert.Equal(t, llotypes.ChannelDefinitions{}, pd.Definitions)
			assert.Equal(t, expectedBlockNum2, pd.BlockNum)

			// does not load erroneously for a different don ID
			pd, err = orm.LoadChannelDefinitions(ctx, addr1, donID2)
			require.NoError(t, err)

			assert.Equal(t, (*PersistedDefinitions)(nil), pd)
		})
	})

	t.Run("StoreChannelDefinitions", func(t *testing.T) {
		expectedBlockNum := rand.Int63()
		cid1 := rand.Uint32()
		cid2 := rand.Uint32()
		defs := llotypes.ChannelDefinitions{
			cid1: llotypes.ChannelDefinition{
				ReportFormat: llotypes.ReportFormatJSON,
				Streams:      []llotypes.Stream{{StreamID: 1, Aggregator: llotypes.AggregatorMedian}, {StreamID: 2, Aggregator: llotypes.AggregatorMode}},
				Opts:         []byte(`{"foo":"bar"}`),
			},
			cid2: llotypes.ChannelDefinition{
				ReportFormat: llotypes.ReportFormatEVMPremiumLegacy,
				Streams:      []llotypes.Stream{{StreamID: 1, Aggregator: llotypes.AggregatorMedian}, {StreamID: 3, Aggregator: llotypes.AggregatorQuote}},
			},
		}

		t.Run("stores channel definitions in the database", func(t *testing.T) {
			err := orm.StoreChannelDefinitions(ctx, addr1, donID1, 42, defs, expectedBlockNum)
			require.NoError(t, err)

			pd, err := orm.LoadChannelDefinitions(ctx, addr1, donID1)
			require.NoError(t, err)
			assert.Equal(t, ETHMainnetChainSelector, pd.ChainSelector)
			assert.Equal(t, addr1, pd.Address)
			assert.Equal(t, expectedBlockNum, pd.BlockNum)
			assert.Equal(t, donID1, pd.DonID)
			assert.Equal(t, uint32(42), pd.Version)
			assert.Equal(t, defs, pd.Definitions)
		})
		t.Run("does not update if version is older than the database persisted version", func(t *testing.T) {
			// try to update with an older version
			err := orm.StoreChannelDefinitions(ctx, addr1, donID1, 41, llotypes.ChannelDefinitions{}, expectedBlockNum)
			require.NoError(t, err)

			pd, err := orm.LoadChannelDefinitions(ctx, addr1, donID1)
			require.NoError(t, err)
			assert.Equal(t, uint32(42), pd.Version)
			assert.Equal(t, defs, pd.Definitions)
		})
	})
}
