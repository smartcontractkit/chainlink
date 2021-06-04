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
	var (
		orm = &mocks.ORM{}
		svc = feeds.NewService(orm)
		id  = int32(1)
		ms  = feeds.ManagerService{}
	)

	orm.On("CreateManagerService", context.Background(), &ms).
		Return(id, nil)

	actual, err := svc.RegisterManagerService(&ms)
	require.NoError(t, err)

	assert.Equal(t, actual, id)
}

func Test_Service_ListManagerServices(t *testing.T) {
	var (
		orm = &mocks.ORM{}
		svc = feeds.NewService(orm)
		ms  = feeds.ManagerService{}
		mss = []feeds.ManagerService{ms}
	)

	orm.On("ListManagerServices", context.Background()).
		Return(mss, nil)

	actual, err := svc.ListManagerServices()
	require.NoError(t, err)

	assert.Equal(t, actual, mss)
}

func Test_Service_GetManagerServices(t *testing.T) {
	var (
		orm = &mocks.ORM{}
		svc = feeds.NewService(orm)
		id  = int32(1)
		ms  = feeds.ManagerService{ID: id}
	)

	orm.On("GetManagerService", context.Background(), id).
		Return(&ms, nil)

	actual, err := svc.GetManagerService(id)
	require.NoError(t, err)

	assert.Equal(t, actual, &ms)
}
