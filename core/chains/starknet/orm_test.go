package starknet_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/db"

	"github.com/smartcontractkit/chainlink/core/chains/starknet"
	"github.com/smartcontractkit/chainlink/core/chains/starknet/types"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
)

func setupORM(t *testing.T) (*sqlx.DB, types.ORM) {
	t.Helper()

	db := pgtest.NewSqlxDB(t)
	orm := starknet.NewORM(db, logger.TestLogger(t), pgtest.NewQConfig(true))

	return db, orm
}

func Test_ORM(t *testing.T) {
	_, orm := setupORM(t)

	dbcs, err := orm.EnabledChains()
	require.NoError(t, err)
	require.Empty(t, dbcs)

	chainIDA := fmt.Sprintf("Chainlinktest-%d", rand.Int31n(999999))
	_, err = orm.CreateChain(chainIDA, &db.ChainCfg{})
	require.NoError(t, err)
	chainIDB := fmt.Sprintf("Chainlinktest-%d", rand.Int31n(999999))
	_, err = orm.CreateChain(chainIDB, &db.ChainCfg{})
	require.NoError(t, err)

	dbcs, err = orm.EnabledChains()
	require.NoError(t, err)
	require.Len(t, dbcs, 2)

	newNode := db.Node{
		Name:    "first",
		ChainID: chainIDA,
		URL:     "http://starknet.devnet.test/first",
	}
	gotNode, err := orm.CreateNode(newNode)
	require.NoError(t, err)
	assertEqual(t, newNode, gotNode)

	gotNode, err = orm.NodeNamed(gotNode.Name)
	require.NoError(t, err)
	assertEqual(t, newNode, gotNode)

	newNode2 := db.Node{
		Name:    "second",
		ChainID: chainIDB,
		URL:     "http://starknet.devnet.test/second",
	}
	gotNode2, err := orm.CreateNode(newNode2)
	require.NoError(t, err)
	assertEqual(t, newNode2, gotNode2)

	gotNodes, count, err := orm.Nodes(0, 3)
	require.NoError(t, err)
	require.Equal(t, 2, count)
	if assert.Len(t, gotNodes, 2) {
		assertEqual(t, newNode, gotNodes[0])
		assertEqual(t, newNode2, gotNodes[1])
	}

	gotNodes, count, err = orm.NodesForChain(newNode2.ChainID, 0, 3)
	require.NoError(t, err)
	require.Equal(t, 1, count)
	if assert.Len(t, gotNodes, 1) {
		assertEqual(t, newNode2, gotNodes[0])
	}

	err = orm.DeleteNode(gotNode.ID)
	require.NoError(t, err)

	gotNodes, count, err = orm.Nodes(0, 3)
	require.NoError(t, err)
	require.Equal(t, 1, count)
	if assert.Len(t, gotNodes, 1) {
		assertEqual(t, newNode2, gotNodes[0])
	}

	newNode3 := db.Node{
		Name:    "third",
		ChainID: chainIDB,
		URL:     "http://starknet.devnet.test/third",
	}
	gotNode3, err := orm.CreateNode(newNode3)
	require.NoError(t, err)
	assertEqual(t, newNode3, gotNode3)

	gotNamed, err := orm.NodeNamed("third")
	require.NoError(t, err)
	assertEqual(t, newNode3, gotNamed)
}

func assertEqual(t *testing.T, newNode db.Node, gotNode db.Node) {
	t.Helper()

	assert.Equal(t, newNode.Name, gotNode.Name)
	assert.Equal(t, newNode.ChainID, gotNode.ChainID)
	assert.Equal(t, newNode.URL, gotNode.URL)
}
