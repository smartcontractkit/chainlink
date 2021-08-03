package headtracker_test

import (
	"context"
	"testing"

	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/services/headtracker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestORM_Heads_Chain(t *testing.T) {
	t.Parallel()

	db := pgtest.NewGormDB(t)
	orm := headtracker.NewORM(db, cltest.FixtureChainID)

	// A competing chain existed from block num 3 to 4
	var baseOfForkHash common.Hash
	var longestChainHeadHash common.Hash
	var parentHash *common.Hash
	for idx := 0; idx < 8; idx++ {
		h := *cltest.Head(idx)
		if parentHash != nil {
			h.ParentHash = *parentHash
		}
		parentHash = &h.Hash
		if idx == 2 {
			baseOfForkHash = h.Hash
		} else if idx == 7 {
			longestChainHeadHash = h.Hash
		}
		assert.Nil(t, orm.IdempotentInsertHead(context.TODO(), h))
	}

	competingHead1 := *cltest.Head(3)
	competingHead1.ParentHash = baseOfForkHash
	assert.Nil(t, orm.IdempotentInsertHead(context.TODO(), competingHead1))
	competingHead2 := *cltest.Head(4)
	competingHead2.ParentHash = competingHead1.Hash
	assert.Nil(t, orm.IdempotentInsertHead(context.TODO(), competingHead2))

	// Query for the top of the longer chain does not include the competing chain
	h, err := orm.Chain(context.TODO(), longestChainHeadHash, 12)
	require.NoError(t, err)
	assert.Equal(t, longestChainHeadHash, h.Hash)
	count := 1
	for {
		if h.Parent == nil {
			break
		}
		require.NotEqual(t, competingHead1.Hash, h.Hash)
		require.NotEqual(t, competingHead2.Hash, h.Hash)
		h = *h.Parent
		count++
	}
	assert.Equal(t, 8, count)

	// If we set the limit lower we get fewer heads in chain
	h, err = orm.Chain(context.TODO(), longestChainHeadHash, 2)
	require.NoError(t, err)
	assert.Equal(t, longestChainHeadHash, h.Hash)
	count = 1
	for {
		if h.Parent == nil {
			break
		}
		h = *h.Parent
		count++
	}
	assert.Equal(t, 2, count)

	// If we query for the top of the competing chain we get its parents
	head, err := orm.Chain(context.TODO(), competingHead2.Hash, 12)
	require.NoError(t, err)
	assert.Equal(t, competingHead2.Hash, head.Hash)
	require.NotNil(t, head.Parent)
	assert.Equal(t, competingHead1.Hash, head.Parent.Hash)
	require.NotNil(t, head.Parent.Parent)
	assert.Equal(t, baseOfForkHash, head.Parent.Parent.Hash)
	assert.NotNil(t, head.Parent.Parent.Parent) // etc...

	// Returns error if hash has no matches
	_, err = orm.Chain(context.TODO(), utils.NewHash(), 12)
	require.Error(t, err)

	t.Run("depth of 0 returns error", func(t *testing.T) {
		_, err = orm.Chain(context.TODO(), longestChainHeadHash, 0)
		require.EqualError(t, err, "record not found")
	})
}

func TestORM_Heads_IdempotentInsertHead(t *testing.T) {
	t.Parallel()

	db := pgtest.NewGormDB(t)
	orm := headtracker.NewORM(db, cltest.FixtureChainID)

	// Returns nil when inserting first head
	head := *cltest.Head(0)
	require.NoError(t, orm.IdempotentInsertHead(context.TODO(), head))

	// Head is inserted
	foundHead, err := orm.LastHead(context.TODO())
	require.NoError(t, err)
	assert.Equal(t, head.Hash, foundHead.Hash)

	// Returns nil when inserting same head again
	require.NoError(t, orm.IdempotentInsertHead(context.TODO(), head))

	// Head is still inserted
	foundHead, err = orm.LastHead(context.TODO())
	require.NoError(t, err)
	assert.Equal(t, head.Hash, foundHead.Hash)
}
