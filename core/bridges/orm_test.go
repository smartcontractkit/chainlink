package bridges_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jmoiron/sqlx"

	"github.com/smartcontractkit/chainlink/v2/core/auth"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

func setupORM(t *testing.T) (*sqlx.DB, bridges.ORM) {
	t.Helper()

	db := pgtest.NewSqlxDB(t)
	orm := bridges.NewORM(db)

	return db, orm
}

func TestORM_FindBridges(t *testing.T) {
	t.Parallel()
	_, orm := setupORM(t)

	bt := bridges.BridgeType{
		Name: "bridge1",
		URL:  cltest.WebURL(t, "https://bridge1.com"),
	}
	ctx := testutils.Context(t)
	assert.NoError(t, orm.CreateBridgeType(ctx, &bt))
	bt2 := bridges.BridgeType{
		Name: "bridge2",
		URL:  cltest.WebURL(t, "https://bridge2.com"),
	}
	assert.NoError(t, orm.CreateBridgeType(ctx, &bt2))
	bts, err := orm.FindBridges(ctx, []bridges.BridgeName{"bridge2", "bridge1"})
	require.NoError(t, err)
	require.Len(t, bts, 2)

	bts, err = orm.FindBridges(ctx, []bridges.BridgeName{"bridge1"})
	require.NoError(t, err)
	require.Len(t, bts, 1)
	require.Equal(t, "bridge1", bts[0].Name.String())

	// One invalid bridge errors
	bts, err = orm.FindBridges(ctx, []bridges.BridgeName{"bridge1", "bridgeX"})
	require.Error(t, err, bts)

	// All invalid bridges error
	bts, err = orm.FindBridges(ctx, []bridges.BridgeName{"bridgeY", "bridgeX"})
	require.Error(t, err, bts)

	// Requires at least one bridge
	bts, err = orm.FindBridges(ctx, []bridges.BridgeName{})
	require.Error(t, err, bts)
}

func TestORM_FindBridge(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	_, orm := setupORM(t)

	bt := bridges.BridgeType{}
	bt.Name = bridges.MustParseBridgeName("solargridreporting")
	bt.URL = cltest.WebURL(t, "https://denergy.eth")
	assert.NoError(t, orm.CreateBridgeType(ctx, &bt))

	cases := []struct {
		description string
		name        bridges.BridgeName
		want        bridges.BridgeType
		errored     bool
	}{
		{"actual external adapter", bt.Name, bt, false},
		{"core adapter", "ethtx", bridges.BridgeType{}, true},
		{"non-existent adapter", "nonExistent", bridges.BridgeType{}, true},
	}

	for _, test := range cases {
		t.Run(test.description, func(t *testing.T) {
			tt, err := orm.FindBridge(ctx, test.name)
			tt.CreatedAt = test.want.CreatedAt
			tt.UpdatedAt = test.want.UpdatedAt
			if test.errored {
				require.Error(t, err)
			} else {
				// we can't make any assumptions about the return type if scanning failed
				require.Equal(t, test.want, tt)
			}
		})
	}
}
func TestORM_UpdateBridgeType(t *testing.T) {
	ctx := testutils.Context(t)
	_, orm := setupORM(t)

	firstBridge := &bridges.BridgeType{
		Name: "UniqueName",
		URL:  cltest.WebURL(t, "http:/oneurl.com"),
	}

	require.NoError(t, orm.CreateBridgeType(ctx, firstBridge))

	updateBridge := &bridges.BridgeTypeRequest{
		URL: cltest.WebURL(t, "http:/updatedurl.com"),
	}

	require.NoError(t, orm.UpdateBridgeType(ctx, firstBridge, updateBridge))

	foundbridge, err := orm.FindBridge(ctx, "UniqueName")
	require.NoError(t, err)
	require.Equal(t, updateBridge.URL, foundbridge.URL)

	bs, count, err := orm.BridgeTypes(ctx, 0, 10)
	require.NoError(t, err)
	require.Equal(t, 1, count)
	require.Len(t, bs, 1)

	require.NoError(t, orm.DeleteBridgeType(ctx, &foundbridge))

	bs, count, err = orm.BridgeTypes(ctx, 0, 10)
	require.NoError(t, err)
	require.Equal(t, 0, count)
	require.Empty(t, bs)
}

func TestORM_TestCachedResponse(t *testing.T) {
	ctx := testutils.Context(t)
	cfg := configtest.NewGeneralConfig(t, nil)
	db := pgtest.NewSqlxDB(t)
	orm := bridges.NewORM(db)

	trORM := pipeline.NewORM(db, logger.TestLogger(t), cfg.JobPipeline().MaxSuccessfulRuns())
	specID, err := trORM.CreateSpec(ctx, pipeline.Pipeline{}, *models.NewInterval(5 * time.Minute))
	require.NoError(t, err)

	_, err = orm.GetCachedResponse(ctx, "dot", specID, 1*time.Second)
	require.Error(t, err)
	require.Contains(t, err.Error(), "no rows in result set")

	err = orm.UpsertBridgeResponse(ctx, "dot", specID, []byte{111, 222, 2})
	require.NoError(t, err)

	val, err := orm.GetCachedResponse(ctx, "dot", specID, 1*time.Second)
	require.NoError(t, err)
	require.Equal(t, []byte{111, 222, 2}, val)
}

func TestORM_CreateExternalInitiator(t *testing.T) {
	ctx := testutils.Context(t)
	_, orm := setupORM(t)

	token := auth.NewToken()
	name := uuid.New().String()
	req := bridges.ExternalInitiatorRequest{
		Name: name,
	}
	exi, err := bridges.NewExternalInitiator(token, &req)
	require.NoError(t, err)
	require.NoError(t, orm.CreateExternalInitiator(ctx, exi))

	exi2, err := bridges.NewExternalInitiator(token, &req)
	require.NoError(t, err)
	require.Contains(t, orm.CreateExternalInitiator(ctx, exi2).Error(), `ERROR: duplicate key value violates unique constraint "external_initiators_name_key" (SQLSTATE 23505)`)
}

func TestORM_DeleteExternalInitiator(t *testing.T) {
	ctx := testutils.Context(t)
	_, orm := setupORM(t)

	token := auth.NewToken()
	name := uuid.New().String()
	req := bridges.ExternalInitiatorRequest{
		Name: name,
	}
	exi, err := bridges.NewExternalInitiator(token, &req)
	require.NoError(t, err)
	require.NoError(t, orm.CreateExternalInitiator(ctx, exi))

	_, err = orm.FindExternalInitiator(ctx, token)
	require.NoError(t, err)
	_, err = orm.FindExternalInitiatorByName(ctx, exi.Name)
	require.NoError(t, err)

	err = orm.DeleteExternalInitiator(ctx, exi.Name)
	require.NoError(t, err)

	_, err = orm.FindExternalInitiator(ctx, token)
	require.Error(t, err)
	_, err = orm.FindExternalInitiatorByName(ctx, exi.Name)
	require.Error(t, err)

	require.NoError(t, orm.CreateExternalInitiator(ctx, exi))
}
