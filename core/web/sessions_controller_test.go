package web_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionsController_Create(t *testing.T) {
	t.Parallel()

	user := cltest.MustUser("email@test.net", "password123")
	app, cleanup := cltest.NewApplication(t)
	app.Start()
	err := app.Store.SaveUser(&user)
	assert.NoError(t, err)
	defer cleanup()

	config := app.Store.Config
	client := http.Client{}
	tests := []struct {
		name        string
		email       string
		password    string
		wantSession bool
	}{
		{"incorrect pwd", "email@test.net", "incorrect", false},
		{"incorrect email", "incorrect@test.net", "password123", false},
		{"correct", "email@test.net", "password123", true},
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
				require.Equal(t, 200, resp.StatusCode)
				cookies := resp.Cookies()
				require.Equal(t, 1, len(cookies))
				decrypted, err := cltest.DecodeSessionCookie(cookies[0].Value)
				require.NoError(t, err)
				user, err := app.Store.AuthorizedUserWithSession(decrypted)
				assert.NoError(t, err)
				assert.Equal(t, test.email, user.Email)

				b, err := ioutil.ReadAll(resp.Body)
				assert.NoError(t, err)
				assert.Contains(t, string(b), `"attributes":{"authenticated":true}`)
			} else {
				require.True(t, resp.StatusCode >= 400, "Should not be able to create session")
				sessions, err := app.Store.Sessions(0, 1)
				assert.NoError(t, err)
				assert.Empty(t, sessions)
			}
		})
	}
}

func TestSessionsController_Create_ReapSessions(t *testing.T) {
	t.Parallel()

	user := cltest.MustUser("email@test.net", "password123")
	app, cleanup := cltest.NewApplication(t)
	app.Start()
	err := app.Store.SaveUser(&user)
	assert.NoError(t, err)
	defer cleanup()

	staleSession := cltest.NewSession()
	staleSession.LastUsed = time.Now().Add(-cltest.MustParseDuration("241h"))
	require.NoError(t, app.Store.SaveSession(&staleSession))

	body := fmt.Sprintf(`{"email":"%s","password":"%s"}`, "email@test.net", "password123")
	resp, err := http.Post(app.Config.ClientNodeURL()+"/sessions", "application/json", bytes.NewBufferString(body))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)
	gomega.NewGomegaWithT(t).Eventually(func() []models.Session {
		sessions, err := app.Store.Sessions(0, 10)
		assert.NoError(t, err)
		return sessions
	}).Should(gomega.HaveLen(1))
}

func TestSessionsController_Destroy(t *testing.T) {
	t.Parallel()

	seedUser := cltest.MustUser("email@test.net", "password123")
	app, cleanup := cltest.NewApplication(t)
	app.Start()
	err := app.Store.SaveUser(&seedUser)
	assert.NoError(t, err)

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
				assert.Equal(t, 200, resp.StatusCode)
			} else {
				assert.True(t, resp.StatusCode >= 400, "Should get an erroneous status code for deleting a nonexistent session id")
			}
		})
	}
}

func TestSessionsController_Destroy_ReapSessions(t *testing.T) {
	t.Parallel()

	client := http.Client{}
	user := cltest.MustUser("email@test.net", "password123")
	app, cleanup := cltest.NewApplication(t)
	defer cleanup()

	app.Start()
	err := app.Store.SaveUser(&user)
	assert.NoError(t, err)

	correctSession := models.NewSession()
	require.NoError(t, app.Store.SaveSession(&correctSession))
	cookie := cltest.MustGenerateSessionCookie(correctSession.ID)

	staleSession := cltest.NewSession()
	staleSession.LastUsed = time.Now().Add(-cltest.MustParseDuration("241h"))
	require.NoError(t, app.Store.SaveSession(&staleSession))

	request, err := http.NewRequest("DELETE", app.Config.ClientNodeURL()+"/sessions", nil)
	assert.NoError(t, err)
	request.AddCookie(cookie)

	resp, err := client.Do(request)
	assert.NoError(t, err)

	assert.Equal(t, 200, resp.StatusCode)
	gomega.NewGomegaWithT(t).Eventually(func() []models.Session {
		sessions, err := app.Store.Sessions(0, 10)
		assert.NoError(t, err)
		return sessions
	}).Should(gomega.HaveLen(0))
}
