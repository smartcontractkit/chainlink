package services_test

import (
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
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
		{"expired", time.Now().Add(-store.Config.SessionTimeout.Duration), false},
		{"almost stale", time.Now().Add(-store.Config.ReaperExpiration.Duration), false},
		{"stale", time.Now().Add(-store.Config.ReaperExpiration.Duration).Add(-store.Config.SessionTimeout.Duration), true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer cltest.ResetBucket(store, &models.Session{})

			session := cltest.NewSession(test.name)
			session.LastUsed = models.Time{test.lastUsed}
			require.NoError(t, store.Save(&session))

			r.ReapSessions()

			if test.wantReap {
				gomega.NewGomegaWithT(t).Eventually(func() []models.Session {
					sessions := []models.Session{}
					assert.Nil(t, store.All(&sessions))
					return sessions
				}).Should(gomega.HaveLen(0))
			} else {
				gomega.NewGomegaWithT(t).Consistently(func() []models.Session {
					sessions := []models.Session{}
					assert.Nil(t, store.All(&sessions))
					return sessions
				}).Should(gomega.HaveLen(1))
			}
		})
	}
}
