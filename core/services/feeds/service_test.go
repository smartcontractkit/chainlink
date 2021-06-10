package feeds_test

import (
	"context"
	"testing"

	"github.com/smartcontractkit/chainlink/core/services/feeds"
	"github.com/smartcontractkit/chainlink/core/services/feeds/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Service_RegisterManagerService(t *testing.T) {
	t.Parallel()

	var (
		orm = &mocks.ORM{}
		svc = feeds.NewService(orm)
		id  = int32(1)
		ms  = feeds.FeedsManager{}
	)

	orm.On("CreateManager", context.Background(), &ms).
		Return(id, nil)

	actual, err := svc.RegisterManager(&ms)
	require.NoError(t, err)

	assert.Equal(t, actual, id)
}

func Test_Service_ListManagerServices(t *testing.T) {
	t.Parallel()

	var (
		orm = &mocks.ORM{}
		svc = feeds.NewService(orm)
		ms  = feeds.FeedsManager{}
		mss = []feeds.FeedsManager{ms}
	)

	orm.On("ListManagers", context.Background()).
		Return(mss, nil)

	actual, err := svc.ListManagers()
	require.NoError(t, err)

	assert.Equal(t, actual, mss)
}

func Test_Service_GetManagerServices(t *testing.T) {
	t.Parallel()

	var (
		orm = &mocks.ORM{}
		svc = feeds.NewService(orm)
		id  = int32(1)
		ms  = feeds.FeedsManager{ID: id}
	)

	orm.On("GetManager", context.Background(), id).
		Return(&ms, nil)

	actual, err := svc.GetManager(id)
	require.NoError(t, err)

	assert.Equal(t, actual, &ms)
}
