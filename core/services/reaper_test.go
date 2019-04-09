package services_test

import (
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/tools/cltest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreReaper_ReapSessions(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()

	r := services.NewStoreReaper(store)
	r.Start()
	defer r.Stop()

	tests := []struct {
		name     string
		lastUsed time.Time
		wantReap bool
	}{
		{"current", time.Now(), false},
		{"expired", time.Now().Add(-store.Config.SessionTimeout()), false},
		{"almost stale", time.Now().Add(-store.Config.ReaperExpiration()), false},
		{"stale", time.Now().Add(-store.Config.ReaperExpiration()).Add(-store.Config.SessionTimeout()), true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer store.ORM.ClearSessions()

			session := cltest.NewSession(test.name)
			session.LastUsed = test.lastUsed
			require.NoError(t, store.SaveSession(&session))

			r.WakeUp()

			if test.wantReap {
				gomega.NewGomegaWithT(t).Eventually(func() []models.Session {
					sessions, err := store.Sessions(0, 10)
					assert.NoError(t, err)
					return sessions
				}).Should(gomega.HaveLen(0))
			} else {
				gomega.NewGomegaWithT(t).Consistently(func() []models.Session {
					sessions, err := store.Sessions(0, 10)
					assert.NoError(t, err)
					return sessions
				}).Should(gomega.HaveLen(1))
			}
		})
	}
}
