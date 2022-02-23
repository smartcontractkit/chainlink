package web_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/sessions"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserController_UpdatePassword(t *testing.T) {
	app := cltest.NewApplicationEVMDisabled(t)
	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient()

	testCases := []struct {
		name           string
		reqBody        string
		wantStatusCode int
		wantErrCount   int
		wantErrMessage string
	}{
		{
			name:           "Invalid request",
			reqBody:        "",
			wantStatusCode: http.StatusUnprocessableEntity,
			wantErrCount:   1,
		},
		{
			name:           "Incorrect old password",
			reqBody:        `{"oldPassword": "wrong password"}`,
			wantStatusCode: http.StatusConflict,
			wantErrCount:   1,
			wantErrMessage: "old password does not match",
		},
		{
			name:           "Success",
			reqBody:        fmt.Sprintf(`{"newPassword": "%v", "oldPassword": "%v"}`, cltest.Password, cltest.Password),
			wantStatusCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			resp, cleanup := client.Patch("/v2/user/password", bytes.NewBufferString(tc.reqBody))
			t.Cleanup(cleanup)
			errors := cltest.ParseJSONAPIErrors(t, resp.Body)

			require.Equal(t, tc.wantStatusCode, resp.StatusCode)
			assert.Len(t, errors.Errors, tc.wantErrCount)
			if tc.wantErrMessage != "" {
				assert.Equal(t, tc.wantErrMessage, errors.Errors[0].Detail)
			}
		})
	}
}

func TestUserController_NewAPIToken(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplicationEVMDisabled(t)
	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient()
	req, err := json.Marshal(sessions.ChangeAuthTokenRequest{
		Password: cltest.Password,
	})
	require.NoError(t, err)
	resp, cleanup := client.Post("/v2/user/token", bytes.NewBuffer(req))
	defer cleanup()

	require.Equal(t, http.StatusCreated, resp.StatusCode)
	var authToken auth.Token
	err = cltest.ParseJSONAPIResponse(t, resp, &authToken)
	require.NoError(t, err)
	assert.NotEmpty(t, authToken.AccessKey)
	assert.NotEmpty(t, authToken.Secret)
}

func TestUserController_NewAPIToken_unauthorized(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplicationEVMDisabled(t)
	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient()
	req, err := json.Marshal(sessions.ChangeAuthTokenRequest{
		Password: "wrong-password",
	})
	require.NoError(t, err)
	resp, cleanup := client.Post("/v2/user/token", bytes.NewBuffer(req))
	defer cleanup()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestUserController_DeleteAPIKey(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplicationEVMDisabled(t)
	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient()
	req, err := json.Marshal(sessions.ChangeAuthTokenRequest{
		Password: cltest.Password,
	})
	require.NoError(t, err)
	resp, cleanup := client.Post("/v2/user/token/delete", bytes.NewBuffer(req))
	defer cleanup()

	require.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestUserController_DeleteAPIKey_unauthorized(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplicationEVMDisabled(t)
	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient()
	req, err := json.Marshal(sessions.ChangeAuthTokenRequest{
		Password: "wrong-password",
	})
	require.NoError(t, err)
	resp, cleanup := client.Post("/v2/user/token/delete", bytes.NewBuffer(req))
	defer cleanup()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}
