package webhook_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	uuid "github.com/satori/go.uuid"

	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/webhook"
	"github.com/smartcontractkit/chainlink/core/static"
	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func JSONFromString(t *testing.T, arg string) *models.JSON {
	if arg == "" {
		return nil
	}
	ret := cltest.JSONFromString(t, arg)
	return &ret
}

func TestNotifyExternalInitiator_Notified(t *testing.T) {
	tests := []struct {
		Name          string
		ExInitr       models.ExternalInitiatorRequest
		JobSpec       models.JobSpec
		JobSpecNotice webhook.JobSpecNotice
	}{
		{
			"Job Spec w/ External Initiator",
			models.ExternalInitiatorRequest{
				Name: "somecoin",
			},
			models.JobSpec{
				ID: models.NewJobID(),
				Initiators: []models.Initiator{
					{
						Type: models.InitiatorExternal,
						InitiatorParams: models.InitiatorParams{
							Name: "somecoin",
							Body: JSONFromString(t, `{"foo":"bar"}`),
						},
					},
				},
			},
			webhook.JobSpecNotice{
				Type:   models.InitiatorExternal,
				Params: cltest.JSONFromString(t, `{"foo":"bar"}`),
			},
		},
		{
			"Job Spec w/ multiple initiators",
			models.ExternalInitiatorRequest{
				Name: "somecoin",
			},
			models.JobSpec{
				ID: models.NewJobID(),
				Initiators: []models.Initiator{
					{
						Type: models.InitiatorCron,
					},
					{
						Type: models.InitiatorWeb,
					},
					{
						Type: models.InitiatorExternal,
						InitiatorParams: models.InitiatorParams{
							Name: "somecoin",
							Body: JSONFromString(t, `{"foo":"bar"}`),
						},
					},
				},
			},
			webhook.JobSpecNotice{
				Type:   models.InitiatorExternal,
				Params: *JSONFromString(t, `{"foo":"bar"}`),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			store, cleanup := cltest.NewStore(t)
			defer cleanup()

			exInitr := struct {
				Header http.Header
				Body   webhook.JobSpecNotice
			}{}
			eiMockServer, assertCalled := cltest.NewHTTPMockServer(t, http.StatusOK, "POST", "",
				func(header http.Header, body string) {
					exInitr.Header = header
					err := json.Unmarshal([]byte(body), &exInitr.Body)
					require.NoError(t, err)
				},
			)
			defer assertCalled()

			url := cltest.WebURL(t, eiMockServer.URL)
			test.ExInitr.URL = &url
			eia := auth.NewToken()
			ei, err := models.NewExternalInitiator(eia, &test.ExInitr)
			require.NoError(t, err)
			err = store.CreateExternalInitiator(ei)
			require.NoError(t, err)

			err = store.CreateJob(&test.JobSpec)
			require.NoError(t, err)

			manager := webhook.NewExternalInitiatorManager(store.DB)
			err = manager.Notify(test.JobSpec)
			require.NoError(t, err)
			assert.Equal(t,
				ei.OutgoingToken,
				exInitr.Header.Get(static.ExternalInitiatorAccessKeyHeader),
			)
			assert.Equal(t,
				ei.OutgoingSecret,
				exInitr.Header.Get(static.ExternalInitiatorSecretHeader),
			)
			test.JobSpecNotice.JobID = test.JobSpec.ID
			assert.Equal(t, test.JobSpecNotice, exInitr.Body)
		})
	}
}

func TestNotifyExternalInitiator_NotNotified(t *testing.T) {
	tests := []struct {
		Name    string
		ExInitr models.ExternalInitiatorRequest
		JobSpec models.JobSpec
	}{
		{
			"Job Spec w/ no Initiators",
			models.ExternalInitiatorRequest{
				Name: "somecoin",
			},
			models.JobSpec{
				ID:         models.NewJobID(),
				Initiators: []models.Initiator{},
			},
		},
		{
			"Job Spec w/ multiple initiators",
			models.ExternalInitiatorRequest{
				Name: "somecoin",
			},
			models.JobSpec{
				ID: models.NewJobID(),
				Initiators: []models.Initiator{
					{
						Type: models.InitiatorCron,
					},
					{
						Type: models.InitiatorWeb,
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			store, cleanup := cltest.NewStore(t)
			defer cleanup()

			var remoteNotified bool
			eiMockServer, _ := cltest.NewHTTPMockServer(t, http.StatusOK, "POST", "",
				func(header http.Header, body string) {
					remoteNotified = true
				},
			)
			defer eiMockServer.Close()

			url := cltest.WebURL(t, eiMockServer.URL)
			test.ExInitr.URL = &url
			eia := auth.NewToken()
			ei, err := models.NewExternalInitiator(eia, &test.ExInitr)
			require.NoError(t, err)
			err = store.CreateExternalInitiator(ei)
			require.NoError(t, err)

			err = store.CreateJob(&test.JobSpec)
			require.NoError(t, err)

			manager := webhook.NewExternalInitiatorManager(store.DB)
			err = manager.Notify(test.JobSpec)
			require.NoError(t, err)

			require.False(t, remoteNotified)
		})
	}
}

func Test_ExternalInitiatorManager_DeleteJob(t *testing.T) {
	tests := []struct {
		Name    string
		ExInitr models.ExternalInitiatorRequest
		JobSpec models.JobSpec
	}{
		{
			"Job Spec w/ External Initiator",
			models.ExternalInitiatorRequest{
				Name: "somecoin",
			},
			models.JobSpec{
				ID: models.NewJobID(),
				Initiators: []models.Initiator{
					{
						Type: models.InitiatorExternal,
						InitiatorParams: models.InitiatorParams{
							Name: "somecoin",
							Body: JSONFromString(t, `{"foo":"bar"}`),
						},
					},
				},
			},
		},
		{
			"Job Spec w/ multiple initiators",
			models.ExternalInitiatorRequest{
				Name: "somecoin",
			},
			models.JobSpec{
				ID: models.NewJobID(),
				Initiators: []models.Initiator{
					{
						Type: models.InitiatorCron,
					},
					{
						Type: models.InitiatorWeb,
					},
					{
						Type: models.InitiatorExternal,
						InitiatorParams: models.InitiatorParams{
							Name: "somecoin",
							Body: JSONFromString(t, `{"foo":"bar"}`),
						},
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			store, cleanup := cltest.NewStore(t)
			defer cleanup()

			var deleteURL *url.URL
			var outgoingToken string
			var outgoingSecret string
			var method string
			var body []byte
			eiMockServer, assertCalled := cltest.NewHTTPMockServerWithRequest(t, http.StatusOK, "",
				func(r *http.Request) {
					b, err := ioutil.ReadAll(r.Body)
					require.NoError(t, err)
					body = b
					deleteURL = r.URL
					outgoingToken = r.Header.Get(static.ExternalInitiatorAccessKeyHeader)
					outgoingSecret = r.Header.Get(static.ExternalInitiatorSecretHeader)
					method = r.Method
				},
			)
			defer assertCalled()

			url := cltest.WebURL(t, eiMockServer.URL)
			test.ExInitr.URL = &url
			eia := auth.NewToken()
			ei, err := models.NewExternalInitiator(eia, &test.ExInitr)
			require.NoError(t, err)
			err = store.CreateExternalInitiator(ei)
			require.NoError(t, err)

			err = store.CreateJob(&test.JobSpec)
			require.NoError(t, err)

			manager := webhook.NewExternalInitiatorManager(store.DB)
			err = manager.DeleteJob(test.JobSpec.ID)
			require.NoError(t, err)

			assert.Equal(t, fmt.Sprintf("/%s", uuid.UUID(test.JobSpec.ID).String()), deleteURL.String())
			assert.Equal(t, ei.OutgoingToken, outgoingToken)
			assert.Equal(t, ei.OutgoingSecret, outgoingSecret)
			assert.Equal(t, "DELETE", method)
			assert.Len(t, body, 0)
		})
	}
}
