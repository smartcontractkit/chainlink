package localauth_test

import (
	"testing"
	"time"

	"github.com/onsi/gomega"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
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

	db := pgtest.NewSqlxDB(t)
	config := sessionReaperConfig{}
	lggr := logger.TestLogger(t)
	orm := localauth.NewORM(db, config.SessionTimeout().Duration(), lggr, audit.NoopLogger)

	r := localauth.NewSessionReaper(db, config, lggr)
	t.Cleanup(func() {
		assert.NoError(t, r.Stop())
	})

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
			ctx := testutils.Context(t)
			t.Cleanup(func() {
				_, err2 := db.Exec("DELETE FROM sessions where email = $1", cltest.APIEmailAdmin)
				require.NoError(t, err2)
			})

			_, err := db.Exec("INSERT INTO sessions (last_used, email, id, created_at) VALUES ($1, $2, $3, now())", test.lastUsed, cltest.APIEmailAdmin, test.name)
			require.NoError(t, err)

			r.WakeUp()

			if test.wantReap {
				gomega.NewWithT(t).Eventually(func() []sessions.Session {
					sessions, err := orm.Sessions(ctx, 0, 10)
					assert.NoError(t, err)
					return sessions
				}).Should(gomega.HaveLen(0))
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
