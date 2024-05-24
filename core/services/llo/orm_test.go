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
	db := pgtest.NewSqlxDB(t)
	orm := NewORM(db, testutils.FixtureChainID)
	ctx := testutils.Context(t)

	addr1 := testutils.NewAddress()
	addr2 := testutils.NewAddress()
	addr3 := testutils.NewAddress()

	t.Run("LoadChannelDefinitions", func(t *testing.T) {
		t.Run("returns zero values if nothing in database", func(t *testing.T) {
			cd, blockNum, err := orm.LoadChannelDefinitions(ctx, addr1)
			require.NoError(t, err)

			assert.Zero(t, cd)
			assert.Zero(t, blockNum)
		})
		t.Run("loads channel definitions from database", func(t *testing.T) {
			expectedBlockNum := rand.Int63()
			expectedBlockNum2 := rand.Int63()
			cid1 := rand.Uint32()
			cid2 := rand.Uint32()

			channelDefsJSON := fmt.Sprintf(`
{
	"%d": {
		"reportFormat": 42,
		"chainSelector": 142,
		"streamIds": [1, 2]
	},
	"%d": {
		"reportFormat": 42,
		"chainSelector": 142,
		"streamIds": [1, 3]
	}
}
			`, cid1, cid2)
			pgtest.MustExec(t, db, `
			INSERT INTO channel_definitions(addr, evm_chain_id, definitions, block_num,  updated_at)
			VALUES ( $1, $2, $3, $4, NOW())
			`, addr1, testutils.FixtureChainID.String(), channelDefsJSON, expectedBlockNum)

			pgtest.MustExec(t, db, `
			INSERT INTO channel_definitions(addr, evm_chain_id, definitions, block_num, updated_at)
			VALUES ( $1, $2, $3, $4, NOW())
			`, addr2, testutils.FixtureChainID.String(), `{}`, expectedBlockNum2)

			{
				// alternative chain ID; we expect these ones to be ignored
				pgtest.MustExec(t, db, `
			INSERT INTO channel_definitions(addr, evm_chain_id, definitions, block_num, updated_at)
			VALUES ( $1, $2, $3, $4, NOW())
			`, addr1, testutils.SimulatedChainID.String(), channelDefsJSON, expectedBlockNum)
				pgtest.MustExec(t, db, `
			INSERT INTO channel_definitions(addr, evm_chain_id, definitions, block_num, updated_at)
			VALUES ( $1, $2, $3, $4, NOW())
			`, addr3, testutils.SimulatedChainID.String(), channelDefsJSON, expectedBlockNum)
			}

			cd, blockNum, err := orm.LoadChannelDefinitions(ctx, addr1)
			require.NoError(t, err)

			assert.Equal(t, llotypes.ChannelDefinitions{
				cid1: llotypes.ChannelDefinition{
					ReportFormat:  42,
					ChainSelector: 142,
					StreamIDs:     []llotypes.StreamID{1, 2},
				},
				cid2: llotypes.ChannelDefinition{
					ReportFormat:  42,
					ChainSelector: 142,
					StreamIDs:     []llotypes.StreamID{1, 3},
				},
			}, cd)
			assert.Equal(t, expectedBlockNum, blockNum)

			cd, blockNum, err = orm.LoadChannelDefinitions(ctx, addr2)
			require.NoError(t, err)

			assert.Equal(t, llotypes.ChannelDefinitions{}, cd)
			assert.Equal(t, expectedBlockNum2, blockNum)
		})
	})
}
