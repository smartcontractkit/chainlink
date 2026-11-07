package localauth_test

import (
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	"github.com/smartcontractkit/chainlink/v2/core/sessions"
	"github.com/smartcontractkit/chainlink/v2/core/sessions/localauth"
)

type sessionReaperConfig struct{}

func (c sessionReaperConfig) SessionTimeout() commonconfig.Duration {
	return *commonconfig.MustNewDuration(42 * time.Second)
}

func (c sessionReaperConfig) SessionReaperExpiration() commonconfig.Duration {
	return *commonconfig.MustNewDuration(142 * time.Second)
}

func TestSessionReaper_ReapSessions(t *testing.T) {
	t.Parallel()

	config := sessionReaperConfig{}
	lggr := logger.TestSugared(t)

	tests := []struct {
		name     string
		lastUsed time.Time
		wantReap bool
	}{
		{"current", time.Now(), false},
		{"expired", time.Now().Add(-config.SessionTimeout().Duration()), false},
		{"almost stale", time.Now().Add(-config.SessionReaperExpiration().Duration()), false},
		{"stale", time.Now().Add(-config.SessionReaperExpiration().Duration()).
			Add(-config.SessionTimeout().Duration()), true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			// Each subtest gets its own database: the assertions below read the whole
			// sessions table, so siblings must not share rows.
			db := pgtest.NewSqlxDB(t)
			orm := localauth.NewORM(db, config.SessionTimeout().Duration(), lggr, audit.NoopLogger)
			r := localauth.NewSessionReaper(db, config, lggr)
			t.Cleanup(func() {
				assert.NoError(t, r.Stop())
			})

			_, err := db.Exec("INSERT INTO sessions (last_used, email, id, created_at) VALUES ($1, $2, $3, now())", test.lastUsed, cltest.APIEmailAdmin, test.name)
			require.NoError(t, err)

			r.WakeUp()

			if test.wantReap {
				gomega.NewWithT(t).Eventually(func() []sessions.Session {
					sessions, err := orm.Sessions(ctx, 0, 10)
					assert.NoError(t, err)
					return sessions
				}).Should(gomega.BeEmpty())
			} else {
				gomega.NewWithT(t).Consistently(func() []sessions.Session {
					sessions, err := orm.Sessions(ctx, 0, 10)
					assert.NoError(t, err)
					return sessions
				}).Should(gomega.HaveLen(1))
			}
		})
	}
}
