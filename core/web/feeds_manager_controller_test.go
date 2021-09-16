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
	"github.com/smartcontractkit/chainlink/core/utils/crypto"
	"github.com/smartcontractkit/chainlink/core/web"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
)

func Test_FeedsManagersController_Create(t *testing.T) {
	t.Parallel()

	app, client := setupFeedsManagerTest(t)

	pubKey, err := crypto.PublicKeyFromHex("3b0f149627adb7b6fafe1497a9dfc357f22295a5440786c3bc566dfdb0176808")
	require.NoError(t, err)
	body, err := json.Marshal(web.CreateFeedsManagerRequest{
		Name:                   "Chainlink FM",
		URI:                    "127.0.0.1:2000",
		JobTypes:               []string{"fluxmonitor"},
		PublicKey:              *pubKey,
		IsBootstrapPeer:        true,
		BootstrapPeerMultiaddr: null.StringFrom("/dns4/ocr-bootstrap.chain.link/tcp/0000/p2p/7777777"),
	})
	require.NoError(t, err)

	fsvc := app.GetFeedsService()

	count, err := fsvc.CountManagers()
	require.NoError(t, err)
	require.Equal(t, int64(0), count)

	resp, cleanup := client.Post("/v2/feeds_managers", bytes.NewReader(body))
	t.Cleanup(cleanup)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	mss, err := fsvc.ListManagers()
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
	assert.True(t, ms.IsOCRBootstrapPeer)
	assert.Equal(t, resource.BootstrapPeerMultiaddr, resource.BootstrapPeerMultiaddr)
}

func Test_FeedsManagersController_List(t *testing.T) {
	t.Parallel()

	app, client := setupFeedsManagerTest(t)

	pubKey, err := crypto.PublicKeyFromHex("3b0f149627adb7b6fafe1497a9dfc357f22295a5440786c3bc566dfdb0176808")
	require.NoError(t, err)

	// Seed feed managers
	fsvc := app.GetFeedsService()
	ms1 := feeds.FeedsManager{
		Name:                      "Chainlink FM",
		URI:                       "wss://127.0.0.1:2000",
		JobTypes:                  []string{"fluxmonitor"},
		PublicKey:                 *pubKey,
		IsOCRBootstrapPeer:        true,
		OCRBootstrapPeerMultiaddr: null.StringFrom("/dns4/ocr-bootstrap.chain.link/tcp/0000/p2p/7777777"),
	}
	ms1ID, err := fsvc.RegisterManager(&ms1)
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
	assert.True(t, resources[0].IsBootstrapPeer)
	assert.Equal(t, resources[0].BootstrapPeerMultiaddr, ms1.OCRBootstrapPeerMultiaddr)
}

func Test_FeedsManagersController_Show(t *testing.T) {
	t.Parallel()

	pubKey, err := crypto.PublicKeyFromHex("3b0f149627adb7b6fafe1497a9dfc357f22295a5440786c3bc566dfdb0176808")
	require.NoError(t, err)
	var (
		ms1 = feeds.FeedsManager{
			Name:                      "Chainlink FM",
			URI:                       "wss://127.0.0.1:2000",
			JobTypes:                  []string{"fluxmonitor"},
			PublicKey:                 *pubKey,
			IsOCRBootstrapPeer:        true,
			OCRBootstrapPeerMultiaddr: null.StringFrom("/dns4/ocr-bootstrap.chain.link/tcp/0000/p2p/7777777"),
		}
	)

	testCases := []struct {
		name           string
		before         func(t *testing.T, app *cltest.TestApplication, id *string)
		want           *feeds.FeedsManager
		wantStatusCode int
	}{
		{
			name: "success",
			before: func(t *testing.T, app *cltest.TestApplication, id *string) {
				// Seed feed managers
				fsvc := app.GetFeedsService()

				ms1ID, err := fsvc.RegisterManager(&ms1)
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
				assert.True(t, resource.IsBootstrapPeer)
				assert.Equal(t, resource.BootstrapPeerMultiaddr, ms1.OCRBootstrapPeerMultiaddr)
			}
		})
	}
}

func setupFeedsManagerTest(t *testing.T) (*cltest.TestApplication, cltest.HTTPClientCleaner) {
	app := cltest.NewApplication(t)
	require.NoError(t, app.Start())
	// We need a CSA key to establish a connection to the FMS
	app.KeyStore.CSA().Create()

	client := app.NewHTTPClient()

	return app, client
}
