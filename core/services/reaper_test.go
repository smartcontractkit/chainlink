package services_test

import (
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"gorm.io/gorm"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionReaper_ReapSessions(t *testing.T) {
	t.Parallel()

	db := pgtest.NewGormDB(t)
	config := cltest.NewTestConfig(t)

	r := services.NewSessionReaper(db, config)
	defer r.Stop()

	tests := []struct {
		name     string
		lastUsed time.Time
		wantReap bool
	}{
		{"current", time.Now(), false},
		{"expired", time.Now().Add(-config.SessionTimeout().Duration()), false},
		{"almost stale", time.Now().Add(-config.ReaperExpiration().Duration()), false},
		{"stale", time.Now().Add(-config.ReaperExpiration().Duration()).
			Add(-config.SessionTimeout().Duration()), true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Cleanup(func() {
				clearSessions(t, db)
			})

			session := cltest.NewSession(test.name)
			session.LastUsed = test.lastUsed
			require.NoError(t, db.Save(&session).Error)

			r.WakeUp()

			if test.wantReap {
				gomega.NewGomegaWithT(t).Eventually(func() []models.Session {
					sessions, err := postgres.Sessions(db, 0, 10)
					assert.NoError(t, err)
					return sessions
				}).Should(gomega.HaveLen(0))
			} else {
				gomega.NewGomegaWithT(t).Consistently(func() []models.Session {
					sessions, err := postgres.Sessions(db, 0, 10)
					assert.NoError(t, err)
					return sessions
				}).Should(gomega.HaveLen(1))
			}
		})
	}
}

// clearSessions removes all sessions.
func clearSessions(t *testing.T, db *gorm.DB) {
	require.NoError(t, db.Exec("DELETE FROM sessions").Error)
}
