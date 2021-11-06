package web_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/sessions"
	"github.com/smartcontractkit/chainlink/core/web"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionsController_Create(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplicationEVMDisabled(t)
	require.NoError(t, app.Start())

	config := app.GetConfig()
	client := http.Client{}
	tests := []struct {
		name        string
		email       string
		password    string
		wantSession bool
	}{
		{"incorrect pwd", cltest.APIEmail, "incorrect", false},
		{"incorrect email", "incorrect@test.net", cltest.Password, false},
		{"correct", cltest.APIEmail, cltest.Password, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body := fmt.Sprintf(`{"email":"%s","password":"%s"}`, test.email, test.password)
			request, err := http.NewRequest("POST", config.ClientNodeURL()+"/sessions", bytes.NewBufferString(body))
			assert.NoError(t, err)
			resp, err := client.Do(request)
			assert.NoError(t, err)
			defer resp.Body.Close()

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

				b, err := ioutil.ReadAll(resp.Body)
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

func TestSessionsController_Create_ReapSessions(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplicationEVMDisabled(t)
	require.NoError(t, app.Start())

	staleSession := cltest.NewSession()
	staleSession.LastUsed = time.Now().Add(-cltest.MustParseDuration(t, "241h"))
	require.NoError(t, app.GetDB().Save(&staleSession).Error)

	body := fmt.Sprintf(`{"email":"%s","password":"%s"}`, cltest.APIEmail, cltest.Password)
	resp, err := http.Post(app.Config.ClientNodeURL()+"/sessions", "application/json", bytes.NewBufferString(body))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var s []sessions.Session
	cltest.NewGomegaWithT(t).Eventually(func() []sessions.Session {
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
	require.NoError(t, app.Start())

	correctSession := sessions.NewSession()
	require.NoError(t, app.GetDB().Save(&correctSession).Error)

	config := app.GetConfig()
	client := http.Client{}
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
			request, err := http.NewRequest("DELETE", config.ClientNodeURL()+"/sessions", nil)
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

	client := http.Client{}
	app := cltest.NewApplicationEVMDisabled(t)
	require.NoError(t, app.Start())

	correctSession := sessions.NewSession()
	require.NoError(t, app.GetDB().Save(&correctSession).Error)
	cookie := cltest.MustGenerateSessionCookie(t, correctSession.ID)

	staleSession := cltest.NewSession()
	staleSession.LastUsed = time.Now().Add(-cltest.MustParseDuration(t, "241h"))
	require.NoError(t, app.GetDB().Save(&staleSession).Error)

	request, err := http.NewRequest("DELETE", app.Config.ClientNodeURL()+"/sessions", nil)
	assert.NoError(t, err)
	request.AddCookie(cookie)

	resp, err := client.Do(request)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	cltest.NewGomegaWithT(t).Eventually(func() []sessions.Session {
		sessions, err := app.SessionORM().Sessions(0, 10)
		assert.NoError(t, err)
		return sessions
	}).Should(gomega.HaveLen(0))
}
