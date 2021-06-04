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
	store, cleanup := cltest.NewStore(t)
	t.Cleanup(cleanup)

	orm := feeds.NewORM(store.DB)

	ms := &feeds.ManagerService{
		URI:       uri,
		Name:      name,
		PublicKey: publicKey,
		JobTypes:  jobTypes,
		Network:   network,
	}

	count, err := orm.Count()
	require.NoError(t, err)
	require.Equal(t, int64(0), count)

	id, err := orm.CreateManagerService(context.Background(), ms)
	require.NoError(t, err)

	count, err = orm.Count()
	require.NoError(t, err)
	require.Equal(t, int64(1), count)

	assert.NotZero(t, id)
}

func Test_ORM_ListManagerServices(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	t.Cleanup(cleanup)

	orm := feeds.NewORM(store.DB)

	ms := &feeds.ManagerService{
		URI:       uri,
		Name:      name,
		PublicKey: publicKey,
		JobTypes:  jobTypes,
		Network:   network,
	}

	id, err := orm.CreateManagerService(context.Background(), ms)
	require.NoError(t, err)

	mss, err := orm.ListManagerServices(context.Background())
	require.NoError(t, err)
	require.Len(t, mss, 1)

	actual := mss[0]
	assert.Equal(t, id, actual.ID)
	assert.Equal(t, uri, actual.URI)
	assert.Equal(t, name, actual.Name)
	assert.Equal(t, publicKey, actual.PublicKey)
	assert.Equal(t, jobTypes, actual.JobTypes)
	assert.Equal(t, network, actual.Network)
}

func Test_ORM_GetManagerService(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	t.Cleanup(cleanup)

	orm := feeds.NewORM(store.DB)

	ms := &feeds.ManagerService{
		URI:       uri,
		Name:      name,
		PublicKey: publicKey,
		JobTypes:  jobTypes,
		Network:   network,
	}

	id, err := orm.CreateManagerService(context.Background(), ms)
	require.NoError(t, err)

	actual, err := orm.GetManagerService(context.Background(), id)
	require.NoError(t, err)

	assert.Equal(t, id, actual.ID)
	assert.Equal(t, uri, actual.URI)
	assert.Equal(t, name, actual.Name)
	assert.Equal(t, publicKey, actual.PublicKey)
	assert.Equal(t, jobTypes, actual.JobTypes)
	assert.Equal(t, network, actual.Network)
}
