package llo

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo/channeldefinitions"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo/types"
)

func Test_ORM(t *testing.T) {
	const ETHMainnetChainSelector uint64 = 5009297550715157269
	const OtherChainSelector uint64 = 1234567890

	db := pgtest.NewSqlxDB(t)
	orm := NewChainScopedORM(db, ETHMainnetChainSelector)
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
			sourceID := uint32(1)

			channelDefsJSON := fmt.Sprintf(`
{
	"%d": {
		"trigger": {
			"source": %d,
			"url": "",
			"sha": [0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],
			"block_num": 142,
			"log_index": 0,
			"version": 1,
			"tx_hash": [0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]
		},
		"definitions": {
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
	}
}
			`, sourceID, sourceID, cid1, cid2)
			pgtest.MustExec(t, db, `
			INSERT INTO channel_definitions(addr, chain_selector, don_id, definitions, block_num, version, updated_at, format)
			VALUES ($1, $2, $3, $4, $5, $6, NOW(), $7)
			`, addr1, ETHMainnetChainSelector, 1, channelDefsJSON, expectedBlockNum, 1, 1)

			pgtest.MustExec(t, db, `
			INSERT INTO channel_definitions(addr, chain_selector, don_id, definitions, block_num, version, updated_at, format)
			VALUES ($1, $2, $3, $4, $5, $6, NOW(), $7)
			`, addr2, ETHMainnetChainSelector, 1, `{}`, expectedBlockNum2, 1, 1)

			{
				// alternative chain selector; we expect these ones to be ignored
				pgtest.MustExec(t, db, `
			INSERT INTO channel_definitions(addr, chain_selector, don_id, definitions, block_num, version, updated_at, format)
			VALUES ($1, $2, $3, $4, $5, $6, NOW(), $7)
			`, addr1, OtherChainSelector, 1, channelDefsJSON, expectedBlockNum, 1, 1)
				pgtest.MustExec(t, db, `
			INSERT INTO channel_definitions(addr, chain_selector, don_id, definitions, block_num, version, updated_at, format)
			VALUES ($1, $2, $3, $4, $5, $6, NOW(), $7)
			`, addr3, OtherChainSelector, 1, channelDefsJSON, expectedBlockNum, 1, 1)
			}

			pd, err := orm.LoadChannelDefinitions(ctx, addr1, donID1)
			require.NoError(t, err)

			assert.Equal(t, ETHMainnetChainSelector, pd.ChainSelector)
			assert.Equal(t, addr1, pd.Address)
			assert.Equal(t, expectedBlockNum, pd.BlockNum)
			assert.Equal(t, donID1, pd.DonID)
			assert.Equal(t, uint32(1), pd.Version)
			assert.Equal(t, uint32(1), pd.Format)

			// Unmarshal the definitions from json.RawMessage
			var loadedDefs map[uint32]types.SourceDefinition
			err = json.Unmarshal(pd.Definitions, &loadedDefs)
			require.NoError(t, err)
			require.Len(t, loadedDefs, 1)
			assert.Equal(t, sourceID, loadedDefs[sourceID].Trigger.Source)
			assert.Equal(t, int64(142), loadedDefs[sourceID].Trigger.BlockNum)
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
			}, loadedDefs[sourceID].Definitions)

			// does not load erroneously for a different address
			pd, err = orm.LoadChannelDefinitions(ctx, addr2, donID1)
			require.NoError(t, err)

			// Unmarshal empty definitions
			var emptyDefs map[uint32]types.SourceDefinition
			err = json.Unmarshal(pd.Definitions, &emptyDefs)
			require.NoError(t, err)
			assert.Empty(t, emptyDefs)
			assert.Equal(t, expectedBlockNum2, pd.BlockNum)

			// does not load erroneously for a different don ID
			pd, err = orm.LoadChannelDefinitions(ctx, addr1, donID2)
			require.NoError(t, err)

			assert.Equal(t, (*types.PersistedDefinitions)(nil), pd)
		})
	})

	t.Run("StoreChannelDefinitions", func(t *testing.T) {
		expectedBlockNum := rand.Int63()
		cid1 := rand.Uint32()
		cid2 := rand.Uint32()
		cid3 := rand.Uint32()
		cid4 := rand.Uint32()
		defs := map[uint32]types.SourceDefinition{
			1: {
				Trigger: types.Trigger{
					Source:   1,
					BlockNum: 142,
					Version:  42,
				},
				Definitions: llotypes.ChannelDefinitions{
					cid1: llotypes.ChannelDefinition{
						ReportFormat: llotypes.ReportFormatJSON,
						Streams:      []llotypes.Stream{{StreamID: 1, Aggregator: llotypes.AggregatorMedian}, {StreamID: 2, Aggregator: llotypes.AggregatorMode}},
						Opts:         []byte(`{"foo":"bar"}`),
					},
					cid2: llotypes.ChannelDefinition{
						ReportFormat: llotypes.ReportFormatEVMPremiumLegacy,
						Streams:      []llotypes.Stream{{StreamID: 1, Aggregator: llotypes.AggregatorMedian}, {StreamID: 3, Aggregator: llotypes.AggregatorQuote}},
					},
				},
			},
			2: {
				Trigger: types.Trigger{
					Source:   2,
					BlockNum: 142,
					Version:  42,
				},
				Definitions: llotypes.ChannelDefinitions{
					cid3: llotypes.ChannelDefinition{
						ReportFormat: llotypes.ReportFormatJSON,
						Streams:      []llotypes.Stream{{StreamID: 1, Aggregator: llotypes.AggregatorMedian}, {StreamID: 2, Aggregator: llotypes.AggregatorMode}},
						Opts:         []byte(`{"foo":"bar"}`),
					},
					cid4: llotypes.ChannelDefinition{
						ReportFormat: llotypes.ReportFormatEVMPremiumLegacy,
						Streams:      []llotypes.Stream{{StreamID: 1, Aggregator: llotypes.AggregatorMedian}, {StreamID: 3, Aggregator: llotypes.AggregatorQuote}},
					},
				},
			},
		}

		// Marshal definitions to json.RawMessage
		defsJSON, err := json.Marshal(defs)
		require.NoError(t, err)

		t.Run("stores channel definitions in the database", func(t *testing.T) {
			err := orm.StoreChannelDefinitions(ctx, addr1, donID1, 42, defsJSON, expectedBlockNum, channeldefinitions.MultiChannelDefinitionsFormat)
			require.NoError(t, err)

			pd, err := orm.LoadChannelDefinitions(ctx, addr1, donID1)
			require.NoError(t, err)
			assert.Equal(t, ETHMainnetChainSelector, pd.ChainSelector)
			assert.Equal(t, addr1, pd.Address)
			assert.Equal(t, expectedBlockNum, pd.BlockNum)
			assert.Equal(t, donID1, pd.DonID)
			assert.Equal(t, uint32(42), pd.Version)
			assert.Equal(t, channeldefinitions.MultiChannelDefinitionsFormat, pd.Format)

			// Unmarshal and compare
			var loadedDefs map[uint32]types.SourceDefinition
			err = json.Unmarshal(pd.Definitions, &loadedDefs)
			require.NoError(t, err)
			assert.Equal(t, defs, loadedDefs)
		})
		t.Run("does not update if version is older than the database persisted version", func(t *testing.T) {
			// try to update with an older version
			emptyDefsJSON, err := json.Marshal(map[uint32]types.SourceDefinition{})
			require.NoError(t, err)
			err = orm.StoreChannelDefinitions(ctx, addr1, donID1, 41, emptyDefsJSON, expectedBlockNum-1, channeldefinitions.MultiChannelDefinitionsFormat)
			require.NoError(t, err)

			pd, err := orm.LoadChannelDefinitions(ctx, addr1, donID1)
			require.NoError(t, err)
			assert.Equal(t, uint32(42), pd.Version)

			// Unmarshal and verify original definitions are still there
			var loadedDefs map[uint32]types.SourceDefinition
			err = json.Unmarshal(pd.Definitions, &loadedDefs)
			require.NoError(t, err)
			assert.Equal(t, defs, loadedDefs)
		})
	})
}
