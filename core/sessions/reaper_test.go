package sessions_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/sessions"
	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type sessionReaperConfig struct{}

func (c sessionReaperConfig) SessionTimeout() models.Duration {
	return models.MustMakeDuration(42 * time.Second)
}

func (c sessionReaperConfig) ReaperExpiration() models.Duration {
	return models.MustMakeDuration(142 * time.Second)
}

func TestSessionReaper_ReapSessions(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	config := sessionReaperConfig{}
	lggr := logger.TestLogger(t)
	orm := sessions.NewORM(db, config.SessionTimeout().Duration(), lggr)

	r := sessions.NewSessionReaper(db.DB, config, lggr)
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
				clearSessions(t, db.DB)
			})

			_, err := db.Exec("INSERT INTO sessions (last_used, id, created_at) VALUES ($1, $2, now())", test.lastUsed, test.name)
			require.NoError(t, err)

			r.WakeUp()

			if test.wantReap {
				gomega.NewWithT(t).Eventually(func() []sessions.Session {
					sessions, err := orm.Sessions(0, 10)
					assert.NoError(t, err)
					return sessions
				}).Should(gomega.HaveLen(0))
			} else {
				gomega.NewWithT(t).Consistently(func() []sessions.Session {
					sessions, err := orm.Sessions(0, 10)
					assert.NoError(t, err)
					return sessions
				}).Should(gomega.HaveLen(1))
			}
		})
	}
}

// clearSessions removes all sessions.
func clearSessions(t *testing.T, db *sql.DB) {
	_, err := db.Exec("DELETE FROM sessions")
	require.NoError(t, err)
}
