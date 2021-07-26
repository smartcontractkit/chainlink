package webhook_test

import (
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/services/webhook"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	null "gopkg.in/guregu/null.v4"
)

func Test_Authorizer(t *testing.T) {
	db := pgtest.NewGormDB(t)

	cltest.MustInsertExternalInitiator(t, db, "foo")
	cltest.MustInsertExternalInitiator(t, db, "bar")

	eiSpec := cltest.JSONFromString(t, `{"bar": 1}`)
	jobWithFooName, _ := cltest.MustInsertWebhookSpec(t, db, null.StringFrom("foo"), &eiSpec)
	jobWithBarName, _ := cltest.MustInsertWebhookSpec(t, db, null.StringFrom("bar"), &eiSpec)
	jobWithNoName, _ := cltest.MustInsertWebhookSpec(t, db, null.String{}, nil)

	t.Run("no user no ea never authorizes", func(t *testing.T) {
		a := webhook.NewAuthorizer(db, nil, nil)

		can, err := a.CanRun(jobWithFooName.ExternalJobID)
		require.NoError(t, err)
		assert.False(t, can)
		can, err = a.CanRun(jobWithNoName.ExternalJobID)
		require.NoError(t, err)
		assert.False(t, can)
		can, err = a.CanRun(uuid.NewV4())
		require.NoError(t, err)
		assert.False(t, can)
	})

	t.Run("with user no ei always authorizes", func(t *testing.T) {
		a := webhook.NewAuthorizer(db, &models.User{}, nil)

		can, err := a.CanRun(jobWithFooName.ExternalJobID)
		require.NoError(t, err)
		assert.True(t, can)
		can, err = a.CanRun(jobWithNoName.ExternalJobID)
		require.NoError(t, err)
		assert.True(t, can)
		can, err = a.CanRun(uuid.NewV4())
		require.NoError(t, err)
		assert.True(t, can)
	})

	t.Run("no user with ei authorizes conditionally", func(t *testing.T) {
		a := webhook.NewAuthorizer(db, nil, &models.ExternalInitiator{Name: "foo"})

		can, err := a.CanRun(jobWithFooName.ExternalJobID)
		require.NoError(t, err)
		assert.True(t, can)
		can, err = a.CanRun(jobWithBarName.ExternalJobID)
		require.NoError(t, err)
		assert.False(t, can)
		can, err = a.CanRun(jobWithNoName.ExternalJobID)
		require.NoError(t, err)
		assert.True(t, can)
		can, err = a.CanRun(uuid.NewV4())
		require.NoError(t, err)
		assert.False(t, can)
	})
}
