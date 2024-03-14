package webhook_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/webhook"
	"github.com/smartcontractkit/chainlink/v2/core/sessions"
)

type eiEnabledCfg struct{}

func (eiEnabledCfg) ExternalInitiatorsEnabled() bool { return true }

type eiDisabledCfg struct{}

func (eiDisabledCfg) ExternalInitiatorsEnabled() bool { return false }

func Test_Authorizer(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	borm := bridges.NewORM(db)

	eiFoo := cltest.MustInsertExternalInitiator(t, borm)
	eiBar := cltest.MustInsertExternalInitiator(t, borm)

	jobWithFooAndBarEI, webhookSpecWithFooAndBarEI := cltest.MustInsertWebhookSpec(t, db)
	jobWithBarEI, webhookSpecWithBarEI := cltest.MustInsertWebhookSpec(t, db)
	jobWithNoEI, _ := cltest.MustInsertWebhookSpec(t, db)

	_, err := db.Exec(`INSERT INTO external_initiator_webhook_specs (external_initiator_id, webhook_spec_id, spec) VALUES ($1,$2,$3)`, eiFoo.ID, webhookSpecWithFooAndBarEI.ID, `{"ei": "foo", "name": "webhookSpecWithFooAndBarEI"}`)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO external_initiator_webhook_specs (external_initiator_id, webhook_spec_id, spec) VALUES ($1,$2,$3)`, eiBar.ID, webhookSpecWithFooAndBarEI.ID, `{"ei": "bar", "name": "webhookSpecWithFooAndBarEI"}`)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO external_initiator_webhook_specs (external_initiator_id, webhook_spec_id, spec) VALUES ($1,$2,$3)`, eiBar.ID, webhookSpecWithBarEI.ID, `{"ei": "bar", "name": "webhookSpecTwoEIs"}`)
	require.NoError(t, err)

	t.Run("no user no ei never authorizes", func(t *testing.T) {
		a := webhook.NewAuthorizer(db, nil, nil)

		can, err := a.CanRun(testutils.Context(t), nil, jobWithFooAndBarEI.ExternalJobID)
		require.NoError(t, err)
		assert.False(t, can)
		can, err = a.CanRun(testutils.Context(t), nil, jobWithNoEI.ExternalJobID)
		require.NoError(t, err)
		assert.False(t, can)
		can, err = a.CanRun(testutils.Context(t), nil, uuid.New())
		require.NoError(t, err)
		assert.False(t, can)
	})

	t.Run("with user no ei always authorizes", func(t *testing.T) {
		a := webhook.NewAuthorizer(db, &sessions.User{}, nil)

		can, err := a.CanRun(testutils.Context(t), nil, jobWithFooAndBarEI.ExternalJobID)
		require.NoError(t, err)
		assert.True(t, can)
		can, err = a.CanRun(testutils.Context(t), nil, jobWithNoEI.ExternalJobID)
		require.NoError(t, err)
		assert.True(t, can)
		can, err = a.CanRun(testutils.Context(t), nil, uuid.New())
		require.NoError(t, err)
		assert.True(t, can)
	})

	t.Run("no user with ei authorizes conditionally", func(t *testing.T) {
		a := webhook.NewAuthorizer(db, nil, &eiFoo)

		can, err := a.CanRun(testutils.Context(t), eiEnabledCfg{}, jobWithFooAndBarEI.ExternalJobID)
		require.NoError(t, err)
		assert.True(t, can)
		can, err = a.CanRun(testutils.Context(t), eiDisabledCfg{}, jobWithFooAndBarEI.ExternalJobID)
		require.NoError(t, err)
		assert.False(t, can)
		can, err = a.CanRun(testutils.Context(t), eiEnabledCfg{}, jobWithBarEI.ExternalJobID)
		require.NoError(t, err)
		assert.False(t, can)
		can, err = a.CanRun(testutils.Context(t), eiEnabledCfg{}, jobWithNoEI.ExternalJobID)
		require.NoError(t, err)
		assert.False(t, can)
		can, err = a.CanRun(testutils.Context(t), eiEnabledCfg{}, uuid.New())
		require.NoError(t, err)
		assert.False(t, can)
	})
}
