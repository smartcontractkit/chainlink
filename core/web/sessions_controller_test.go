package web_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	clhttptest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/httptest"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/sessions"
	"github.com/smartcontractkit/chainlink/v2/core/web"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionsController_Create(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplicationEVMDisabled(t)
	require.NoError(t, app.Start(testutils.Context(t)))

	client := clhttptest.NewTestLocalOnlyHTTPClient()
	tests := []struct {
		name        string
		email       string
		password    string
		wantSession bool
	}{
		{"incorrect pwd", cltest.APIEmailAdmin, "incorrect", false},
		{"incorrect email", "incorrect@test.net", cltest.Password, false},
		{"correct", cltest.APIEmailAdmin, cltest.Password, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body := fmt.Sprintf(`{"email":"%s","password":"%s"}`, test.email, test.password)
			request, err := http.NewRequest("POST", app.Server.URL+"/sessions", bytes.NewBufferString(body))
			assert.NoError(t, err)
			resp, err := client.Do(request)
			assert.NoError(t, err)
			defer func() { assert.NoError(t, resp.Body.Close()) }()

			if test.wantSession {
				require.Equal(t, http.StatusOK, resp.StatusCode)

				cookies := resp.Cookies()
				sessionCookie := web.FindSessionCookie(cookies)
				require.NotNil(t, sessionCookie)

				decrypted, err := cltest.DecodeSessionCookie(sessionCookie.Value)
				require.NoError(t, err)
				user, err := app.SessionORM().AuthorizedUserWithSession(decrypted)
				assert.NoError(t, err)
				assert.Equal(t, test.email, user.Email)

				b, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)
				assert.Contains(t, string(b), `"attributes":{"authenticated":true}`)
			} else {
				require.True(t, resp.StatusCode >= 400, "Should not be able to create session")
				// Ignore fixture session
				sessions, err := app.SessionORM().Sessions(1, 2)
				assert.NoError(t, err)
				assert.Empty(t, sessions)
			}
		})
	}
}

func mustInsertSession(t *testing.T, q pg.Q, session *sessions.Session) {
	sql := "INSERT INTO sessions (id, email, last_used, created_at) VALUES ($1, $2, $3, $4) RETURNING *"
	_, err := q.Exec(sql, session.ID, cltest.APIEmailAdmin, session.LastUsed, session.CreatedAt)
	require.NoError(t, err)
}

func TestSessionsController_Create_ReapSessions(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplicationEVMDisabled(t)
	require.NoError(t, app.Start(testutils.Context(t)))

	staleSession := cltest.NewSession()
	staleSession.LastUsed = time.Now().Add(-cltest.MustParseDuration(t, "241h"))
	q := pg.NewQ(app.GetSqlxDB(), app.GetLogger(), app.GetConfig())
	mustInsertSession(t, q, &staleSession)

	body := fmt.Sprintf(`{"email":"%s","password":"%s"}`, cltest.APIEmailAdmin, cltest.Password)
	resp, err := http.Post(app.Server.URL+"/sessions", "application/json", bytes.NewBufferString(body))
	assert.NoError(t, err)
	defer func() { assert.NoError(t, resp.Body.Close()) }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var s []sessions.Session
	gomega.NewWithT(t).Eventually(func() []sessions.Session {
		s, err = app.SessionORM().Sessions(0, 10)
		assert.NoError(t, err)
		return s
	}).Should(gomega.HaveLen(1))

	for _, session := range s {
		assert.NotEqual(t, session.ID, staleSession.ID)
	}
}

func TestSessionsController_Destroy(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplicationEVMDisabled(t)
	require.NoError(t, app.Start(testutils.Context(t)))

	correctSession := sessions.NewSession()
	q := pg.NewQ(app.GetSqlxDB(), app.GetLogger(), app.GetConfig())
	mustInsertSession(t, q, &correctSession)

	client := clhttptest.NewTestLocalOnlyHTTPClient()
	tests := []struct {
		name, sessionID string
		success         bool
	}{
		{"correct cookie", correctSession.ID, true},
		{"incorrect cookie", "wrongsessionid", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cookie := cltest.MustGenerateSessionCookie(t, test.sessionID)
			request, err := http.NewRequest("DELETE", app.Server.URL+"/sessions", nil)
			assert.NoError(t, err)
			request.AddCookie(cookie)

			resp, err := client.Do(request)
			assert.NoError(t, err)

			_, err = app.SessionORM().AuthorizedUserWithSession(test.sessionID)
			assert.Error(t, err)
			if test.success {
				assert.Equal(t, http.StatusOK, resp.StatusCode)
			} else {
				assert.True(t, resp.StatusCode >= 400, "Should get an erroneous status code for deleting a nonexistent session id")
			}
		})
	}
}

func TestSessionsController_Destroy_ReapSessions(t *testing.T) {
	t.Parallel()

	client := clhttptest.NewTestLocalOnlyHTTPClient()
	app := cltest.NewApplicationEVMDisabled(t)
	q := pg.NewQ(app.GetSqlxDB(), app.GetLogger(), app.GetConfig())
	require.NoError(t, app.Start(testutils.Context(t)))

	correctSession := sessions.NewSession()
	mustInsertSession(t, q, &correctSession)
	cookie := cltest.MustGenerateSessionCookie(t, correctSession.ID)

	staleSession := cltest.NewSession()
	staleSession.LastUsed = time.Now().Add(-cltest.MustParseDuration(t, "241h"))
	mustInsertSession(t, q, &staleSession)

	request, err := http.NewRequest("DELETE", app.Server.URL+"/sessions", nil)
	assert.NoError(t, err)
	request.AddCookie(cookie)

	resp, err := client.Do(request)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	gomega.NewWithT(t).Eventually(func() []sessions.Session {
		sessions, err := app.SessionORM().Sessions(0, 10)
		assert.NoError(t, err)
		return sessions
	}).Should(gomega.HaveLen(0))
}
