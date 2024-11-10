package syncer

import (
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/capabilities/workflows/common"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	url, err := orm.GetSecretsURL(ctx, common.Keccak256Hash([]byte(giveURL)))
	require.NoError(t, err)
	assert.Equal(t, giveURL, url)

	artifact, err := orm.GetArtifactByHash(ctx, common.Keccak256Hash([]byte(giveURL)))
	require.NoError(t, err)
	assert.Equal(t, giveURL, artifact.SecretsURL)
	assert.Equal(t, "some contents", artifact.Contents)

	_, err = orm.Update(ctx, giveURL, "new contents")
	require.NoError(t, err)

	artifact, err = orm.GetArtifactByHash(ctx, common.Keccak256Hash([]byte(giveURL)))
	require.NoError(t, err)
	assert.Equal(t, giveURL, artifact.SecretsURL)
	assert.Equal(t, "new contents", artifact.Contents)
}
