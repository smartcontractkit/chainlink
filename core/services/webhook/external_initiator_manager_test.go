package webhook_test

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"net/http"
	"net/url"
	"strings"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/tidwall/gjson"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/services/webhook"
	webhookmocks "github.com/smartcontractkit/chainlink/core/services/webhook/mocks"
	"github.com/smartcontractkit/chainlink/core/static"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

			manager := webhook.NewExternalInitiatorManager(store.DB, utils.UnrestrictedClient)
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

			manager := webhook.NewExternalInitiatorManager(store.DB, utils.UnrestrictedClient)
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

			manager := webhook.NewExternalInitiatorManager(store.DB, utils.UnrestrictedClient)
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

func Test_ExternalInitiatorManager_Load(t *testing.T) {
	db := pgtest.NewGormDB(t)

	eiFoo := cltest.MustInsertExternalInitiator(t, db)
	eiBar := cltest.MustInsertExternalInitiator(t, db)

	jb1, webhookSpecOneEI := cltest.MustInsertWebhookSpec(t, db)
	jb2, webhookSpecTwoEIs := cltest.MustInsertWebhookSpec(t, db)
	jb3, webhookSpecNoEIs := cltest.MustInsertWebhookSpec(t, db)

	err := multierr.Combine(
		db.Exec(`INSERT INTO external_initiator_webhook_specs (external_initiator_id, webhook_spec_id, spec) VALUES (?,?,?)`, eiFoo.ID, webhookSpecTwoEIs.ID, `{"ei": "foo", "name": "webhookSpecTwoEIs"}`).Error,
		db.Exec(`INSERT INTO external_initiator_webhook_specs (external_initiator_id, webhook_spec_id, spec) VALUES (?,?,?)`, eiBar.ID, webhookSpecTwoEIs.ID, `{"ei": "bar", "name": "webhookSpecTwoEIs"}`).Error,
		db.Exec(`INSERT INTO external_initiator_webhook_specs (external_initiator_id, webhook_spec_id, spec) VALUES (?,?,?)`, eiFoo.ID, webhookSpecOneEI.ID, `{"ei": "foo", "name": "webhookSpecOneEI"}`).Error,
	)
	require.NoError(t, err)

	eim := webhook.NewExternalInitiatorManager(db, nil)

	eiWebhookSpecs, jobID, err := eim.Load(webhookSpecNoEIs.ID)
	require.NoError(t, err)
	assert.Len(t, eiWebhookSpecs, 0)
	assert.Equal(t, jb3.ExternalJobID, jobID)

	eiWebhookSpecs, jobID, err = eim.Load(webhookSpecOneEI.ID)
	require.NoError(t, err)
	assert.Len(t, eiWebhookSpecs, 1)
	assert.Equal(t, `{"ei": "foo", "name": "webhookSpecOneEI"}`, eiWebhookSpecs[0].Spec.Raw)
	assert.Equal(t, eiFoo.ID, eiWebhookSpecs[0].ExternalInitiator.ID)
	assert.Equal(t, jb1.ExternalJobID, jobID)

	eiWebhookSpecs, jobID, err = eim.Load(webhookSpecTwoEIs.ID)
	require.NoError(t, err)
	assert.Len(t, eiWebhookSpecs, 2)
	assert.Equal(t, jb2.ExternalJobID, jobID)
}

func Test_ExternalInitiatorManager_NotifyV2(t *testing.T) {
	db := pgtest.NewGormDB(t)

	eiWithURL := cltest.MustInsertExternalInitiatorWithOpts(t, db, cltest.ExternalInitiatorOpts{
		URL:            cltest.MustWebURL(t, "http://example.com/foo"),
		OutgoingSecret: "secret",
		OutgoingToken:  "token",
	})
	eiNoURL := cltest.MustInsertExternalInitiator(t, db)

	jb, webhookSpecTwoEIs := cltest.MustInsertWebhookSpec(t, db)
	_, webhookSpecNoEIs := cltest.MustInsertWebhookSpec(t, db)

	err := multierr.Combine(
		db.Exec(`INSERT INTO external_initiator_webhook_specs (external_initiator_id, webhook_spec_id, spec) VALUES (?,?,?)`, eiWithURL.ID, webhookSpecTwoEIs.ID, `{"ei": "foo", "name": "webhookSpecTwoEIs"}`).Error,
		db.Exec(`INSERT INTO external_initiator_webhook_specs (external_initiator_id, webhook_spec_id, spec) VALUES (?,?,?)`, eiNoURL.ID, webhookSpecTwoEIs.ID, `{"ei": "bar", "name": "webhookSpecTwoEIs"}`).Error,
	)
	require.NoError(t, err)

	client := new(webhookmocks.HTTPClient)
	eim := webhook.NewExternalInitiatorManager(db, client)

	// Does nothing with no EI
	eim.NotifyV2(webhookSpecNoEIs.ID)

	client.On("Do", mock.MatchedBy(func(r *http.Request) bool {
		body, err := r.GetBody()
		require.NoError(t, err)
		b, err := ioutil.ReadAll(body)
		require.NoError(t, err)

		assert.Equal(t, jb.ExternalJobID.String(), gjson.GetBytes(b, "jobId").Str)
		assert.Equal(t, eiWithURL.Name, gjson.GetBytes(b, "type").Str)
		assert.Equal(t, `{"ei":"foo","name":"webhookSpecTwoEIs"}`, gjson.GetBytes(b, "params").Raw)

		return r.Method == "POST" && r.URL.String() == eiWithURL.URL.String() && r.Header["Content-Type"][0] == "application/json" && r.Header["X-Chainlink-Ea-Accesskey"][0] == "token" && r.Header["X-Chainlink-Ea-Secret"][0] == "secret"
	})).Once().Return(&http.Response{Body: io.NopCloser(strings.NewReader(""))}, nil)
	eim.NotifyV2(webhookSpecTwoEIs.ID)

	client.AssertExpectations(t)
}

func Test_ExternalInitiatorManager_DeleteJobV2(t *testing.T) {
	db := pgtest.NewGormDB(t)

	eiWithURL := cltest.MustInsertExternalInitiatorWithOpts(t, db, cltest.ExternalInitiatorOpts{
		URL:            cltest.MustWebURL(t, "http://example.com/foo"),
		OutgoingSecret: "secret",
		OutgoingToken:  "token",
	})
	eiNoURL := cltest.MustInsertExternalInitiator(t, db)

	jb, webhookSpecTwoEIs := cltest.MustInsertWebhookSpec(t, db)
	_, webhookSpecNoEIs := cltest.MustInsertWebhookSpec(t, db)

	err := multierr.Combine(
		db.Exec(`INSERT INTO external_initiator_webhook_specs (external_initiator_id, webhook_spec_id, spec) VALUES (?,?,?)`, eiWithURL.ID, webhookSpecTwoEIs.ID, `{"ei": "foo", "name": "webhookSpecTwoEIs"}`).Error,
		db.Exec(`INSERT INTO external_initiator_webhook_specs (external_initiator_id, webhook_spec_id, spec) VALUES (?,?,?)`, eiNoURL.ID, webhookSpecTwoEIs.ID, `{"ei": "bar", "name": "webhookSpecTwoEIs"}`).Error,
	)
	require.NoError(t, err)

	client := new(webhookmocks.HTTPClient)
	eim := webhook.NewExternalInitiatorManager(db, client)

	// Does nothing with no EI
	eim.DeleteJobV2(webhookSpecNoEIs.ID)

	client.On("Do", mock.MatchedBy(func(r *http.Request) bool {
		expectedURL := fmt.Sprintf("%s/%s", eiWithURL.URL.String(), jb.ExternalJobID.String())
		return r.Method == "DELETE" && r.URL.String() == expectedURL && r.Header["Content-Type"][0] == "application/json" && r.Header["X-Chainlink-Ea-Accesskey"][0] == "token" && r.Header["X-Chainlink-Ea-Secret"][0] == "secret"
	})).Once().Return(&http.Response{Body: io.NopCloser(strings.NewReader(""))}, nil)
	eim.DeleteJobV2(webhookSpecTwoEIs.ID)

	client.AssertExpectations(t)
}
