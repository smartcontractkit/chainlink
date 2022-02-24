package bridges_test

import (
	"testing"

	"github.com/smartcontractkit/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/chainlink/core/bridges"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
)

func setupORM(t *testing.T) (*sqlx.DB, bridges.ORM) {
	t.Helper()

	cfg := cltest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)
	orm := bridges.NewORM(db, logger.TestLogger(t), cfg)

	return db, orm
}

func TestORM_FindBridges(t *testing.T) {
	t.Parallel()
	_, orm := setupORM(t)

	bt := bridges.BridgeType{
		Name: "bridge1",
		URL:  cltest.WebURL(t, "https://bridge1.com"),
	}
	assert.NoError(t, orm.CreateBridgeType(&bt))
	bt2 := bridges.BridgeType{
		Name: "bridge2",
		URL:  cltest.WebURL(t, "https://bridge2.com"),
	}
	assert.NoError(t, orm.CreateBridgeType(&bt2))
	bts, err := orm.FindBridges([]bridges.BridgeName{"bridge2", "bridge1"})
	require.NoError(t, err)
	require.Equal(t, 2, len(bts))

	bts, err = orm.FindBridges([]bridges.BridgeName{"bridge1"})
	require.NoError(t, err)
	require.Equal(t, 1, len(bts))
	require.Equal(t, "bridge1", bts[0].Name.String())

	// One invalid bridge errors
	bts, err = orm.FindBridges([]bridges.BridgeName{"bridge1", "bridgeX"})
	require.Error(t, err, bts)

	// All invalid bridges error
	bts, err = orm.FindBridges([]bridges.BridgeName{"bridgeY", "bridgeX"})
	require.Error(t, err, bts)

	// Requires at least one bridge
	bts, err = orm.FindBridges([]bridges.BridgeName{})
	require.Error(t, err, bts)
}

func TestORM_FindBridge(t *testing.T) {
	t.Parallel()

	_, orm := setupORM(t)

	bt := bridges.BridgeType{}
	bt.Name = bridges.MustParseBridgeName("solargridreporting")
	bt.URL = cltest.WebURL(t, "https://denergy.eth")
	assert.NoError(t, orm.CreateBridgeType(&bt))

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
			tt, err := orm.FindBridge(test.name)
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
	_, orm := setupORM(t)

	firstBridge := &bridges.BridgeType{
		Name: "UniqueName",
		URL:  cltest.WebURL(t, "http:/oneurl.com"),
	}

	require.NoError(t, orm.CreateBridgeType(firstBridge))

	updateBridge := &bridges.BridgeTypeRequest{
		URL: cltest.WebURL(t, "http:/updatedurl.com"),
	}

	require.NoError(t, orm.UpdateBridgeType(firstBridge, updateBridge))

	foundbridge, err := orm.FindBridge("UniqueName")
	require.NoError(t, err)
	require.Equal(t, updateBridge.URL, foundbridge.URL)
}

func TestORM_CreateExternalInitiator(t *testing.T) {
	_, orm := setupORM(t)

	token := auth.NewToken()
	req := bridges.ExternalInitiatorRequest{
		Name: "externalinitiator",
	}
	exi, err := bridges.NewExternalInitiator(token, &req)
	require.NoError(t, err)
	require.NoError(t, orm.CreateExternalInitiator(exi))

	exi2, err := bridges.NewExternalInitiator(token, &req)
	require.NoError(t, err)
	require.Contains(t, orm.CreateExternalInitiator(exi2).Error(), `ERROR: duplicate key value violates unique constraint "external_initiators_name_key" (SQLSTATE 23505)`)
}

func TestORM_DeleteExternalInitiator(t *testing.T) {
	_, orm := setupORM(t)

	token := auth.NewToken()
	req := bridges.ExternalInitiatorRequest{
		Name: "externalinitiator",
	}
	exi, err := bridges.NewExternalInitiator(token, &req)
	require.NoError(t, err)
	require.NoError(t, orm.CreateExternalInitiator(exi))

	_, err = orm.FindExternalInitiator(token)
	require.NoError(t, err)

	err = orm.DeleteExternalInitiator(exi.Name)
	require.NoError(t, err)

	_, err = orm.FindExternalInitiator(token)
	require.Error(t, err)

	require.NoError(t, orm.CreateExternalInitiator(exi))
}
