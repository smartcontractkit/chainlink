package web_test

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionsController_create(t *testing.T) {
	t.Parallel()

	user := cltest.MustUser("email@test.net", "password123")
	app, cleanup := cltest.NewApplication()
	app.Start()
	err := app.Store.Save(&user)
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
			request, err := http.NewRequest("POST", config.ClientNodeURL+"/sessions", bytes.NewBufferString(body))
			assert.NoError(t, err)
			resp, err := client.Do(request)
			assert.NoError(t, err)

			if test.wantSession {
				assert.Equal(t, 200, resp.StatusCode)
				cookies := resp.Cookies()
				assert.Equal(t, 1, len(cookies))
				decrypted, err := cltest.DecodeSessionCookie(cookies[0].Value)
				require.NoError(t, err)
				user, err := app.Store.FindUserBySession(decrypted)
				assert.NoError(t, err)
				assert.Equal(t, test.email, user.Email)
			} else {
				assert.True(t, resp.StatusCode >= 400, "Should not be able to create session")
				user, err := app.Store.FindUser()
				assert.NoError(t, err)
				assert.Empty(t, user.SessionID)
			}
		})
	}
}

func TestSessionsController_destroy(t *testing.T) {
	t.Parallel()

	seedUser := cltest.MustUser("email@test.net", "password123", "ShouldBeDeleted")
	app, cleanup := cltest.NewApplication()
	app.Start()
	err := app.Store.Save(&seedUser)
	assert.NoError(t, err)
	defer cleanup()

	config := app.Store.Config
	client := http.Client{}
	tests := []struct {
		name    string
		cookie  *http.Cookie
		success bool
	}{
		{"incorrect cookie", cltest.MustGenerateSessionCookie("deadbeef"), false},
		{"correct cookie", cltest.MustGenerateSessionCookie(seedUser.SessionID), true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request, err := http.NewRequest("DELETE", config.ClientNodeURL+"/sessions", nil)
			request.AddCookie(test.cookie)
			assert.NoError(t, err)
			resp, err := client.Do(request)
			assert.NoError(t, err)

			if test.success {
				assert.Equal(t, 200, resp.StatusCode)
				user, err := app.Store.FindUser()
				assert.NoError(t, err)
				assert.Empty(t, user.SessionID)
			} else {
				assert.True(t, resp.StatusCode >= 400, "Should not be able to destroy session")
				user, err := app.Store.FindUser()
				assert.NoError(t, err)
				assert.Equal(t, seedUser.SessionID, user.SessionID)
			}
		})
	}
}
