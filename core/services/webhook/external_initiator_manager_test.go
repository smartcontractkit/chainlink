package webhook_test

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	_ "github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/webhook"
	webhookmocks "github.com/smartcontractkit/chainlink/v2/core/services/webhook/mocks"
)

func Test_ExternalInitiatorManager_Load(t *testing.T) {
	ctx := testutils.Context(t)
	db := pgtest.NewSqlxDB(t)
	borm := bridges.NewORM(db)

	eiFoo := cltest.MustInsertExternalInitiator(t, borm)
	eiBar := cltest.MustInsertExternalInitiator(t, borm)

	jb1, webhookSpecOneEI := cltest.MustInsertWebhookSpec(t, db)
	jb2, webhookSpecTwoEIs := cltest.MustInsertWebhookSpec(t, db)
	jb3, webhookSpecNoEIs := cltest.MustInsertWebhookSpec(t, db)

	pgtest.MustExec(t, db, `INSERT INTO external_initiator_webhook_specs (external_initiator_id, webhook_spec_id, spec) VALUES ($1,$2,$3)`, eiFoo.ID, webhookSpecTwoEIs.ID, `{"ei": "foo", "name": "webhookSpecTwoEIs"}`)
	pgtest.MustExec(t, db, `INSERT INTO external_initiator_webhook_specs (external_initiator_id, webhook_spec_id, spec) VALUES ($1,$2,$3)`, eiBar.ID, webhookSpecTwoEIs.ID, `{"ei": "bar", "name": "webhookSpecTwoEIs"}`)
	pgtest.MustExec(t, db, `INSERT INTO external_initiator_webhook_specs (external_initiator_id, webhook_spec_id, spec) VALUES ($1,$2,$3)`, eiFoo.ID, webhookSpecOneEI.ID, `{"ei": "foo", "name": "webhookSpecOneEI"}`)

	eim := webhook.NewExternalInitiatorManager(db, nil)

	eiWebhookSpecs, jobID, err := eim.Load(ctx, webhookSpecNoEIs.ID)
	require.NoError(t, err)
	assert.Len(t, eiWebhookSpecs, 0)
	assert.Equal(t, jb3.ExternalJobID, jobID)

	eiWebhookSpecs, jobID, err = eim.Load(ctx, webhookSpecOneEI.ID)
	require.NoError(t, err)
	assert.Len(t, eiWebhookSpecs, 1)
	assert.Equal(t, `{"ei": "foo", "name": "webhookSpecOneEI"}`, eiWebhookSpecs[0].Spec.Raw)
	assert.Equal(t, eiFoo.ID, eiWebhookSpecs[0].ExternalInitiator.ID)
	assert.Equal(t, jb1.ExternalJobID, jobID)

	eiWebhookSpecs, jobID, err = eim.Load(ctx, webhookSpecTwoEIs.ID)
	require.NoError(t, err)
	assert.Len(t, eiWebhookSpecs, 2)
	assert.Equal(t, jb2.ExternalJobID, jobID)
}

func Test_ExternalInitiatorManager_Notify(t *testing.T) {
	ctx := tests.Context(t)
	db := pgtest.NewSqlxDB(t)
	borm := bridges.NewORM(db)

	eiWithURL := cltest.MustInsertExternalInitiatorWithOpts(t, borm, cltest.ExternalInitiatorOpts{
		URL:            cltest.MustWebURL(t, "http://example.com/foo"),
		OutgoingSecret: "secret",
		OutgoingToken:  "token",
	})
	eiNoURL := cltest.MustInsertExternalInitiator(t, borm)

	jb, webhookSpecTwoEIs := cltest.MustInsertWebhookSpec(t, db)
	_, webhookSpecNoEIs := cltest.MustInsertWebhookSpec(t, db)

	pgtest.MustExec(t, db, `INSERT INTO external_initiator_webhook_specs (external_initiator_id, webhook_spec_id, spec) VALUES ($1,$2,$3)`, eiWithURL.ID, webhookSpecTwoEIs.ID, `{"ei": "foo", "name": "webhookSpecTwoEIs"}`)
	pgtest.MustExec(t, db, `INSERT INTO external_initiator_webhook_specs (external_initiator_id, webhook_spec_id, spec) VALUES ($1,$2,$3)`, eiNoURL.ID, webhookSpecTwoEIs.ID, `{"ei": "bar", "name": "webhookSpecTwoEIs"}`)

	client := webhookmocks.NewHTTPClient(t)
	eim := webhook.NewExternalInitiatorManager(db, client)

	// Does nothing with no EI
	require.NoError(t, eim.Notify(ctx, webhookSpecNoEIs.ID))

	client.On("Do", mock.MatchedBy(func(r *http.Request) bool {
		body, err := r.GetBody()
		require.NoError(t, err)
		b, err := io.ReadAll(body)
		require.NoError(t, err)

		assert.Equal(t, jb.ExternalJobID.String(), gjson.GetBytes(b, "jobId").Str)
		assert.Equal(t, eiWithURL.Name, gjson.GetBytes(b, "type").Str)
		assert.Equal(t, `{"ei":"foo","name":"webhookSpecTwoEIs"}`, gjson.GetBytes(b, "params").Raw)

		return r.Method == "POST" && r.URL.String() == eiWithURL.URL.String() && r.Header["Content-Type"][0] == "application/json" && r.Header["X-Chainlink-Ea-Accesskey"][0] == "token" && r.Header["X-Chainlink-Ea-Secret"][0] == "secret"
	})).Once().Return(&http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(""))}, nil)
	require.NoError(t, eim.Notify(ctx, webhookSpecTwoEIs.ID))
}

func Test_ExternalInitiatorManager_DeleteJob(t *testing.T) {
	ctx := testutils.Context(t)
	db := pgtest.NewSqlxDB(t)
	borm := bridges.NewORM(db)

	eiWithURL := cltest.MustInsertExternalInitiatorWithOpts(t, borm, cltest.ExternalInitiatorOpts{
		URL:            cltest.MustWebURL(t, "http://example.com/foo"),
		OutgoingSecret: "secret",
		OutgoingToken:  "token",
	})
	eiNoURL := cltest.MustInsertExternalInitiator(t, borm)

	jb, webhookSpecTwoEIs := cltest.MustInsertWebhookSpec(t, db)
	_, webhookSpecNoEIs := cltest.MustInsertWebhookSpec(t, db)

	pgtest.MustExec(t, db, `INSERT INTO external_initiator_webhook_specs (external_initiator_id, webhook_spec_id, spec) VALUES ($1,$2,$3)`, eiWithURL.ID, webhookSpecTwoEIs.ID, `{"ei": "foo", "name": "webhookSpecTwoEIs"}`)
	pgtest.MustExec(t, db, `INSERT INTO external_initiator_webhook_specs (external_initiator_id, webhook_spec_id, spec) VALUES ($1,$2,$3)`, eiNoURL.ID, webhookSpecTwoEIs.ID, `{"ei": "bar", "name": "webhookSpecTwoEIs"}`)

	client := webhookmocks.NewHTTPClient(t)
	eim := webhook.NewExternalInitiatorManager(db, client)

	// Does nothing with no EI
	require.NoError(t, eim.DeleteJob(ctx, webhookSpecNoEIs.ID))

	client.On("Do", mock.MatchedBy(func(r *http.Request) bool {
		expectedURL := fmt.Sprintf("%s/%s", eiWithURL.URL.String(), jb.ExternalJobID.String())
		return r.Method == "DELETE" && r.URL.String() == expectedURL && r.Header["Content-Type"][0] == "application/json" && r.Header["X-Chainlink-Ea-Accesskey"][0] == "token" && r.Header["X-Chainlink-Ea-Secret"][0] == "secret"
	})).Once().Return(&http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(""))}, nil)
	require.NoError(t, eim.DeleteJob(ctx, webhookSpecTwoEIs.ID))
}
