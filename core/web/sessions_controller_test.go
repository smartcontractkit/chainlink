package web_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/web"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionsController_Create(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t, cltest.LenientEthMock)
	app.Start()
	defer cleanup()

	config := app.Store.Config
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
				user, err := app.Store.AuthorizedUserWithSession(decrypted)
				assert.NoError(t, err)
				assert.Equal(t, test.email, user.Email)

				b, err := ioutil.ReadAll(resp.Body)
				assert.NoError(t, err)
				assert.Contains(t, string(b), `"attributes":{"authenticated":true}`)
			} else {
				require.True(t, resp.StatusCode >= 400, "Should not be able to create session")
				// Ignore fixture session
				sessions, err := app.Store.Sessions(1, 2)
				assert.NoError(t, err)
				assert.Empty(t, sessions)
			}
		})
	}
}

func TestSessionsController_Create_ReapSessions(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t, cltest.LenientEthMock)
	app.Start()
	defer cleanup()

	staleSession := cltest.NewSession()
	staleSession.LastUsed = time.Now().Add(-cltest.MustParseDuration(t, "241h"))
	require.NoError(t, app.Store.SaveSession(&staleSession))

	body := fmt.Sprintf(`{"email":"%s","password":"%s"}`, cltest.APIEmail, cltest.Password)
	resp, err := http.Post(app.Config.ClientNodeURL()+"/sessions", "application/json", bytes.NewBufferString(body))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	gomega.NewGomegaWithT(t).Eventually(func() []models.Session {
		sessions, err := app.Store.Sessions(0, 10)
		assert.NoError(t, err)
		for _, session := range sessions {
			assert.NotEqual(t, session.ID, staleSession.ID)
		}
		return sessions
	}).Should(gomega.HaveLen(1))
}

func TestSessionsController_Destroy(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t, cltest.LenientEthMock)
	require.NoError(t, app.Start())

	correctSession := models.NewSession()
	require.NoError(t, app.Store.SaveSession(&correctSession))
	defer cleanup()

	config := app.Store.Config
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
			cookie := cltest.MustGenerateSessionCookie(test.sessionID)
			request, err := http.NewRequest("DELETE", config.ClientNodeURL()+"/sessions", nil)
			assert.NoError(t, err)
			request.AddCookie(cookie)

			resp, err := client.Do(request)
			assert.NoError(t, err)

			_, err = app.Store.AuthorizedUserWithSession(test.sessionID)
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
	app, cleanup := cltest.NewApplication(t, cltest.LenientEthMock)
	defer cleanup()
	require.NoError(t, app.Start())

	correctSession := models.NewSession()
	require.NoError(t, app.Store.SaveSession(&correctSession))
	cookie := cltest.MustGenerateSessionCookie(correctSession.ID)

	staleSession := cltest.NewSession()
	staleSession.LastUsed = time.Now().Add(-cltest.MustParseDuration(t, "241h"))
	require.NoError(t, app.Store.SaveSession(&staleSession))

	request, err := http.NewRequest("DELETE", app.Config.ClientNodeURL()+"/sessions", nil)
	assert.NoError(t, err)
	request.AddCookie(cookie)

	resp, err := client.Do(request)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	gomega.NewGomegaWithT(t).Eventually(func() []models.Session {
		sessions, err := app.Store.Sessions(0, 10)
		assert.NoError(t, err)
		return sessions
	}).Should(gomega.HaveLen(0))
}
