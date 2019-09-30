package web_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNotifyExternalInitiator_Notified(t *testing.T) {
	tests := []struct {
		Name          string
		ExInitr       models.ExternalInitiatorRequest
		JobSpec       models.JobSpec
		JobSpecNotice web.JobSpecNotice
	}{
		{
			"Job Spec w/ External Initiator",
			models.ExternalInitiatorRequest{
				Name: "somecoin",
			},
			models.JobSpec{
				ID: models.NewID(),
				Initiators: []models.Initiator{
					models.Initiator{
						Type: models.InitiatorExternal,
						InitiatorParams: models.InitiatorParams{
							Name:   "somecoin",
							Params: `{"foo":"bar"}`,
						},
					},
				},
			},
			web.JobSpecNotice{
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
				ID: models.NewID(),
				Initiators: []models.Initiator{
					models.Initiator{
						Type: models.InitiatorCron,
					},
					models.Initiator{
						Type: models.InitiatorWeb,
					},
					models.Initiator{
						Type: models.InitiatorExternal,
						InitiatorParams: models.InitiatorParams{
							Name:   "somecoin",
							Params: `{"foo":"bar"}`,
						},
					},
				},
			},
			web.JobSpecNotice{
				Type:   models.InitiatorExternal,
				Params: cltest.JSONFromString(t, `{"foo":"bar"}`),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			app, cleanup := cltest.NewApplicationWithKey(t)
			defer cleanup()
			exInitr := struct {
				Header http.Header
				Body   web.JobSpecNotice
			}{}
			eiMockServer, assertCalled := cltest.NewHTTPMockServer(t, http.StatusOK, "POST", "",
				func(header http.Header, body string) {
					exInitr.Header = header
					err := json.Unmarshal([]byte(body), &exInitr.Body)
					require.NoError(t, err)
				},
			)
			defer assertCalled()

			test.ExInitr.URL = cltest.WebURL(t, eiMockServer.URL)
			eia := models.NewExternalInitiatorAuthentication()
			ei, err := models.NewExternalInitiator(eia, &test.ExInitr)
			require.NoError(t, err)
			err = app.GetStore().CreateExternalInitiator(ei)
			require.NoError(t, err)

			err = app.GetStore().CreateJob(&test.JobSpec)
			require.NoError(t, err)

			err = web.NotifyExternalInitiator(test.JobSpec, app.GetStore())
			require.NoError(t, err)
			assert.Equal(t,
				ei.OutgoingToken,
				exInitr.Header.Get(web.ExternalInitiatorAccessKeyHeader),
			)
			assert.Equal(t,
				ei.OutgoingSecret,
				exInitr.Header.Get(web.ExternalInitiatorSecretHeader),
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
				ID:         models.NewID(),
				Initiators: []models.Initiator{},
			},
		},
		{
			"Job Spec w/ multiple initiators",
			models.ExternalInitiatorRequest{
				Name: "somecoin",
			},
			models.JobSpec{
				ID: models.NewID(),
				Initiators: []models.Initiator{
					models.Initiator{
						Type: models.InitiatorCron,
					},
					models.Initiator{
						Type: models.InitiatorWeb,
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			app, cleanup := cltest.NewApplicationWithKey(t)
			defer cleanup()

			var remoteNotified bool
			eiMockServer, _ := cltest.NewHTTPMockServer(t, http.StatusOK, "POST", "",
				func(header http.Header, body string) {
					remoteNotified = true
				},
			)
			defer eiMockServer.Close()

			test.ExInitr.URL = cltest.WebURL(t, eiMockServer.URL)
			eia := models.NewExternalInitiatorAuthentication()
			ei, err := models.NewExternalInitiator(eia, &test.ExInitr)
			require.NoError(t, err)
			err = app.GetStore().CreateExternalInitiator(ei)
			require.NoError(t, err)

			err = app.GetStore().CreateJob(&test.JobSpec)
			require.NoError(t, err)

			err = web.NotifyExternalInitiator(test.JobSpec, app.GetStore())
			require.NoError(t, err)

			require.False(t, remoteNotified)
		})
	}
}
