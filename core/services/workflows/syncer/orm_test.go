package syncer

import (
	"encoding/hex"
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils/crypto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkflowArtifactsORM_GetAndUpdate(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	ctx := testutils.Context(t)
	lggr := logger.TestLogger(t)
	orm := &orm{ds: db, lggr: lggr}

	giveURL := "https://example.com"
	giveBytes, err := crypto.Keccak256([]byte(giveURL))
	require.NoError(t, err)
	giveHash := hex.EncodeToString(giveBytes)
	giveContent := "some contents"

	gotID, err := orm.Create(ctx, giveURL, giveHash, giveContent)
	require.NoError(t, err)

	url, err := orm.GetSecretsURLByID(ctx, gotID)
	require.NoError(t, err)
	assert.Equal(t, giveURL, url)

	contents, err := orm.GetContents(ctx, giveURL)
	require.NoError(t, err)
	assert.Equal(t, "some contents", contents)

	contents, err = orm.GetContentsByHash(ctx, giveHash)
	require.NoError(t, err)
	assert.Equal(t, "some contents", contents)

	_, err = orm.Update(ctx, giveHash, "new contents")
	require.NoError(t, err)

	contents, err = orm.GetContents(ctx, giveURL)
	require.NoError(t, err)
	assert.Equal(t, "new contents", contents)

	contents, err = orm.GetContentsByHash(ctx, giveHash)
	require.NoError(t, err)
	assert.Equal(t, "new contents", contents)
}
