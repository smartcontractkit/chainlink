package feeds_test

import (
	"context"
	"testing"

	"github.com/lib/pq"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/feeds"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	uri       = "http://192.168.0.1"
	name      = "Chainlink FMS"
	publicKey = feeds.PublicKey([]byte("11111111111111111111111111111111"))
	jobTypes  = pq.StringArray{feeds.JobTypeFluxMonitor, feeds.JobTypeOffchainReporting}
	network   = "mainnet"
)

func Test_ORM_CreateManagerService(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	t.Cleanup(cleanup)

	orm := feeds.NewORM(store.DB)

	mgr := &feeds.FeedsManager{
		URI:       uri,
		Name:      name,
		PublicKey: publicKey,
		JobTypes:  jobTypes,
		Network:   network,
	}

	count, err := orm.CountManagers()
	require.NoError(t, err)
	require.Equal(t, int64(0), count)

	id, err := orm.CreateManager(context.Background(), mgr)
	require.NoError(t, err)

	count, err = orm.CountManagers()
	require.NoError(t, err)
	require.Equal(t, int64(1), count)

	assert.NotZero(t, id)
}

func Test_ORM_ListManagerServices(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	t.Cleanup(cleanup)

	orm := feeds.NewORM(store.DB)

	mgr := &feeds.FeedsManager{
		URI:       uri,
		Name:      name,
		PublicKey: publicKey,
		JobTypes:  jobTypes,
		Network:   network,
	}

	id, err := orm.CreateManager(context.Background(), mgr)
	require.NoError(t, err)

	mgrs, err := orm.ListManagers(context.Background())
	require.NoError(t, err)
	require.Len(t, mgrs, 1)

	actual := mgrs[0]
	assert.Equal(t, id, actual.ID)
	assert.Equal(t, uri, actual.URI)
	assert.Equal(t, name, actual.Name)
	assert.Equal(t, publicKey, actual.PublicKey)
	assert.Equal(t, jobTypes, actual.JobTypes)
	assert.Equal(t, network, actual.Network)
}

func Test_ORM_GetManagerService(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	t.Cleanup(cleanup)

	orm := feeds.NewORM(store.DB)

	mgr := &feeds.FeedsManager{
		URI:       uri,
		Name:      name,
		PublicKey: publicKey,
		JobTypes:  jobTypes,
		Network:   network,
	}

	id, err := orm.CreateManager(context.Background(), mgr)
	require.NoError(t, err)

	actual, err := orm.GetManager(context.Background(), id)
	require.NoError(t, err)

	assert.Equal(t, id, actual.ID)
	assert.Equal(t, uri, actual.URI)
	assert.Equal(t, name, actual.Name)
	assert.Equal(t, publicKey, actual.PublicKey)
	assert.Equal(t, jobTypes, actual.JobTypes)
	assert.Equal(t, network, actual.Network)

	actual, err = orm.GetManager(context.Background(), -1)
	require.Nil(t, actual)
	require.Error(t, err)
}
