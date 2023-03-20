package smoke

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups/mercury"
)

func genUuid() string {
	id, _ := uuid.NewV4()
	return id.String()
}

func TestMercuryServerAdminAPI(t *testing.T) {
	testEnv, err := mercury.NewEnv(t.Name(), "smoke")
	testEnv.InitEnv()

	t.Cleanup(func() {
		testEnv.Cleanup(t)
	})
	require.NoError(t, err)

	admin := mercury.User{
		Id:       testEnv.MSInfo.AdminId,
		Key:      "admintestkey",
		Secret:   "mz1I4AgYtvo3Wumrgtlyh9VWkCf/IzZ6JROnuw==",
		Role:     "admin",
		Disabled: false,
	}
	user := mercury.User{
		Id:       genUuid(),
		Key:      "admintestkey",
		Secret:   "mz1I4AgYtvo3Wumrgtlyh9VWkCf/IzZ6JROnuw==",
		Role:     "user",
		Disabled: false,
	}
	initUsers := []mercury.User{admin, user}
	err = testEnv.AddMercuryServer(nil, nil, &initUsers)
	require.NoError(t, err)
	msUrl := testEnv.MSInfo.LocalUrl

	t.Run("GET /admin/user as admin role", func(t *testing.T) {
		c := client.NewMercuryServerClient(msUrl, admin.Id, admin.Key)
		users, resp, err := c.GetUsers()
		require.NoError(t, err)
		require.Equal(t, 200, resp.StatusCode)
		require.Equal(t, len(initUsers), len(users))
	})

	t.Run("GET /admin/user as user role", func(t *testing.T) {
		c := client.NewMercuryServerClient(msUrl, user.Id, user.Key)
		users, resp, err := c.GetUsers()
		require.NoError(t, err)
		require.Equal(t, 401, resp.StatusCode)
		require.Equal(t, 0, len(users))
	})
}
