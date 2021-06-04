package web_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/feeds"
	"github.com/smartcontractkit/chainlink/core/web"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_FeedsManagersController_Create(t *testing.T) {
	t.Parallel()

	app, client := setupFeedsManagerTest(t)

	pubKey, err := feeds.PublicKeyFromHex("3b0f149627adb7b6fafe1497a9dfc357f22295a5440786c3bc566dfdb0176808")
	require.NoError(t, err)
	body, err := json.Marshal(web.CreateFeedsManagerRequest{
		Name:      "Chainlink FM",
		URI:       "wss://127.0.0.1:2000",
		JobTypes:  []string{"fluxmonitor"},
		PublicKey: *pubKey,
		Network:   "mainnet",
	})
	require.NoError(t, err)

	fsvc := app.GetFeedsService()

	count, err := fsvc.CountManagerServices()
	require.NoError(t, err)
	require.Equal(t, int64(0), count)

	resp, cleanup := client.Post("/v2/feeds_managers", bytes.NewReader(body))
	t.Cleanup(cleanup)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	mss, err := fsvc.ListManagerServices()
	require.NoError(t, err)
	require.Len(t, mss, 1)
	ms := mss[0]

	resource := presenters.FeedsManagerResource{}
	err = web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, resp), &resource)
	require.NoError(t, err)

	assert.Equal(t, resource.ID, strconv.Itoa(int(ms.ID)))
	assert.Equal(t, resource.Name, ms.Name)
	assert.Equal(t, resource.URI, ms.URI)
	assert.Equal(t, resource.JobTypes, []string(ms.JobTypes))
	assert.Equal(t, resource.PublicKey, ms.PublicKey)
	assert.Equal(t, resource.Network, ms.Network)
}

func Test_FeedsManagersController_List(t *testing.T) {
	t.Parallel()

	app, client := setupFeedsManagerTest(t)

	pubKey, err := feeds.PublicKeyFromHex("3b0f149627adb7b6fafe1497a9dfc357f22295a5440786c3bc566dfdb0176808")
	require.NoError(t, err)

	// Seed feed managers
	fsvc := app.GetFeedsService()
	ms1 := feeds.ManagerService{
		Name:      "Chainlink FM",
		URI:       "wss://127.0.0.1:2000",
		JobTypes:  []string{"fluxmonitor"},
		PublicKey: *pubKey,
		Network:   "mainnet",
	}
	ms1ID, err := fsvc.RegisterManagerService(&ms1)
	require.NoError(t, err)

	resp, cleanup := client.Get("/v2/feeds_managers")
	t.Cleanup(cleanup)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	resources := []presenters.FeedsManagerResource{}
	err = web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, resp), &resources)
	require.NoError(t, err)

	assert.Equal(t, resources[0].ID, strconv.Itoa(int(ms1ID)))
	assert.Equal(t, resources[0].Name, ms1.Name)
	assert.Equal(t, resources[0].URI, ms1.URI)
	assert.Equal(t, resources[0].JobTypes, []string(ms1.JobTypes))
	assert.Equal(t, resources[0].PublicKey, ms1.PublicKey)
	assert.Equal(t, resources[0].Network, ms1.Network)
}

func Test_FeedsManagersController_Show(t *testing.T) {
	pubKey, err := feeds.PublicKeyFromHex("3b0f149627adb7b6fafe1497a9dfc357f22295a5440786c3bc566dfdb0176808")
	require.NoError(t, err)
	var (
		ms1 = feeds.ManagerService{
			Name:      "Chainlink FM",
			URI:       "wss://127.0.0.1:2000",
			JobTypes:  []string{"fluxmonitor"},
			PublicKey: *pubKey,
			Network:   "mainnet",
		}
	)

	testCases := []struct {
		name           string
		before         func(t *testing.T, app *cltest.TestApplication, id *string)
		want           *feeds.ManagerService
		wantStatusCode int
	}{
		{
			name: "success",
			before: func(t *testing.T, app *cltest.TestApplication, id *string) {
				// Seed feed managers
				fsvc := app.GetFeedsService()

				ms1ID, err := fsvc.RegisterManagerService(&ms1)
				require.NoError(t, err)

				*id = strconv.Itoa(int(ms1ID))
			},
			wantStatusCode: http.StatusOK,
			want:           &ms1,
		},
		{
			name: "invalid id",
			before: func(t *testing.T, app *cltest.TestApplication, id *string) {
				*id = "notanint"
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "not found",
			before: func(t *testing.T, app *cltest.TestApplication, id *string) {
				*id = "999999999"
			},
			wantStatusCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			app, client := setupFeedsManagerTest(t)

			var id string
			if tc.before != nil {
				tc.before(t, app, &id)
			}

			resp, cleanup := client.Get(fmt.Sprintf("/v2/feeds_managers/%s", id))
			t.Cleanup(cleanup)
			require.Equal(t, tc.wantStatusCode, resp.StatusCode)

			if tc.want != nil {
				resource := presenters.FeedsManagerResource{}
				err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, resp), &resource)
				require.NoError(t, err)

				assert.Equal(t, resource.ID, id)
				assert.Equal(t, resource.Name, tc.want.Name)
				assert.Equal(t, resource.URI, tc.want.URI)
				assert.Equal(t, resource.JobTypes, []string(tc.want.JobTypes))
				assert.Equal(t, resource.PublicKey, tc.want.PublicKey)
				assert.Equal(t, resource.Network, tc.want.Network)
			}
		})
	}
}

func setupFeedsManagerTest(t *testing.T) (*cltest.TestApplication, cltest.HTTPClientCleaner) {
	app, cleanup := cltest.NewApplication(t)
	t.Cleanup(cleanup)
	app.Start()

	client := app.NewHTTPClient()

	return app, client
}
