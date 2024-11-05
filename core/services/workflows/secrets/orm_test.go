package syncer

import (
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/stretchr/testify/assert"
	"github.com/test-go/testify/require"
)

func TestWorkflowArtifactsORM_GetAndUpdate(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	ctx := testutils.Context(t)
	lggr := logger.TestLogger(t)
	orm := &orm{ds: db, lggr: lggr}

	giveURL := "https://example.com"
	giveContent := "some contents"

	_, err := orm.Update(ctx, giveURL, giveContent)
	require.NoError(t, err)

	url, err := orm.GetSecretsURL(ctx, keccak256Hash(giveURL))
	require.NoError(t, err)
	assert.Equal(t, giveURL, url)

	artifact, err := orm.GetArtifactByHash(ctx, keccak256Hash(giveURL))
	require.NoError(t, err)
	assert.Equal(t, giveURL, artifact.SecretsURL)
	assert.Equal(t, "some contents", artifact.Contents)

	_, err = orm.Update(ctx, giveURL, "new contents")
	require.NoError(t, err)

	artifact, err = orm.GetArtifactByHash(ctx, keccak256Hash(giveURL))
	require.NoError(t, err)
	assert.Equal(t, giveURL, artifact.SecretsURL)
	assert.Equal(t, "new contents", artifact.Contents)
}
