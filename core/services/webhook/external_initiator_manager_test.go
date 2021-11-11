package webhook_test

import (
	"fmt"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	_ "github.com/smartcontractkit/chainlink/core/services/postgres"
	"io"
	"io/ioutil"

	"net/http"
	"strings"
	"testing"

	"github.com/tidwall/gjson"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/services/webhook"
	webhookmocks "github.com/smartcontractkit/chainlink/core/services/webhook/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_ExternalInitiatorManager_Load(t *testing.T) {
	gdb := pgtest.NewGormDB(t)
	db := postgres.UnwrapGormDB(gdb)

	eiFoo := cltest.MustInsertExternalInitiator(t, gdb)
	eiBar := cltest.MustInsertExternalInitiator(t, gdb)

	jb1, webhookSpecOneEI := cltest.MustInsertWebhookSpec(t, db)
	jb2, webhookSpecTwoEIs := cltest.MustInsertWebhookSpec(t, db)
	jb3, webhookSpecNoEIs := cltest.MustInsertWebhookSpec(t, db)

	err := multierr.Combine(
		gdb.Exec(`INSERT INTO external_initiator_webhook_specs (external_initiator_id, webhook_spec_id, spec) VALUES (?,?,?)`, eiFoo.ID, webhookSpecTwoEIs.ID, `{"ei": "foo", "name": "webhookSpecTwoEIs"}`).Error,
		gdb.Exec(`INSERT INTO external_initiator_webhook_specs (external_initiator_id, webhook_spec_id, spec) VALUES (?,?,?)`, eiBar.ID, webhookSpecTwoEIs.ID, `{"ei": "bar", "name": "webhookSpecTwoEIs"}`).Error,
		gdb.Exec(`INSERT INTO external_initiator_webhook_specs (external_initiator_id, webhook_spec_id, spec) VALUES (?,?,?)`, eiFoo.ID, webhookSpecOneEI.ID, `{"ei": "foo", "name": "webhookSpecOneEI"}`).Error,
	)
	require.NoError(t, err)

	eim := webhook.NewExternalInitiatorManager(gdb, nil)

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

func Test_ExternalInitiatorManager_Notify(t *testing.T) {
	gdb := pgtest.NewGormDB(t)
	db := postgres.UnwrapGormDB(gdb)

	eiWithURL := cltest.MustInsertExternalInitiatorWithOpts(t, gdb, cltest.ExternalInitiatorOpts{
		URL:            cltest.MustWebURL(t, "http://example.com/foo"),
		OutgoingSecret: "secret",
		OutgoingToken:  "token",
	})
	eiNoURL := cltest.MustInsertExternalInitiator(t, gdb)

	jb, webhookSpecTwoEIs := cltest.MustInsertWebhookSpec(t, db)
	_, webhookSpecNoEIs := cltest.MustInsertWebhookSpec(t, db)

	err := multierr.Combine(
		gdb.Exec(`INSERT INTO external_initiator_webhook_specs (external_initiator_id, webhook_spec_id, spec) VALUES (?,?,?)`, eiWithURL.ID, webhookSpecTwoEIs.ID, `{"ei": "foo", "name": "webhookSpecTwoEIs"}`).Error,
		gdb.Exec(`INSERT INTO external_initiator_webhook_specs (external_initiator_id, webhook_spec_id, spec) VALUES (?,?,?)`, eiNoURL.ID, webhookSpecTwoEIs.ID, `{"ei": "bar", "name": "webhookSpecTwoEIs"}`).Error,
	)
	require.NoError(t, err)

	client := new(webhookmocks.HTTPClient)
	eim := webhook.NewExternalInitiatorManager(gdb, client)

	// Does nothing with no EI
	eim.Notify(webhookSpecNoEIs.ID)

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
	eim.Notify(webhookSpecTwoEIs.ID)

	client.AssertExpectations(t)
}

func Test_ExternalInitiatorManager_DeleteJob(t *testing.T) {
	gdb := pgtest.NewGormDB(t)
	db := postgres.UnwrapGormDB(gdb)

	eiWithURL := cltest.MustInsertExternalInitiatorWithOpts(t, gdb, cltest.ExternalInitiatorOpts{
		URL:            cltest.MustWebURL(t, "http://example.com/foo"),
		OutgoingSecret: "secret",
		OutgoingToken:  "token",
	})
	eiNoURL := cltest.MustInsertExternalInitiator(t, gdb)

	jb, webhookSpecTwoEIs := cltest.MustInsertWebhookSpec(t, db)
	_, webhookSpecNoEIs := cltest.MustInsertWebhookSpec(t, db)

	err := multierr.Combine(
		gdb.Exec(`INSERT INTO external_initiator_webhook_specs (external_initiator_id, webhook_spec_id, spec) VALUES (?,?,?)`, eiWithURL.ID, webhookSpecTwoEIs.ID, `{"ei": "foo", "name": "webhookSpecTwoEIs"}`).Error,
		gdb.Exec(`INSERT INTO external_initiator_webhook_specs (external_initiator_id, webhook_spec_id, spec) VALUES (?,?,?)`, eiNoURL.ID, webhookSpecTwoEIs.ID, `{"ei": "bar", "name": "webhookSpecTwoEIs"}`).Error,
	)
	require.NoError(t, err)

	client := new(webhookmocks.HTTPClient)
	eim := webhook.NewExternalInitiatorManager(gdb, client)

	// Does nothing with no EI
	eim.DeleteJob(webhookSpecNoEIs.ID)

	client.On("Do", mock.MatchedBy(func(r *http.Request) bool {
		expectedURL := fmt.Sprintf("%s/%s", eiWithURL.URL.String(), jb.ExternalJobID.String())
		return r.Method == "DELETE" && r.URL.String() == expectedURL && r.Header["Content-Type"][0] == "application/json" && r.Header["X-Chainlink-Ea-Accesskey"][0] == "token" && r.Header["X-Chainlink-Ea-Secret"][0] == "secret"
	})).Once().Return(&http.Response{Body: io.NopCloser(strings.NewReader(""))}, nil)
	eim.DeleteJob(webhookSpecTwoEIs.ID)

	client.AssertExpectations(t)
}
