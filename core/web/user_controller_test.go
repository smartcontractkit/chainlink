package web_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/sessions"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func TestUserController_UpdatePassword(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplicationEVMDisabled(t)
	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(cltest.APIEmailAdmin)

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
			name:           "Insufficient length of new password",
			reqBody:        fmt.Sprintf(`{"newPassword": "%v", "oldPassword": "%v"}`, "foo", cltest.Password),
			wantStatusCode: http.StatusUnprocessableEntity,
			wantErrCount:   1,
			wantErrMessage: fmt.Sprintf("%s	%s\n", utils.ErrMsgHeader, "password is less than 16 characters long"),
		},
		{
			name:           "New password includes api email",
			reqBody:        fmt.Sprintf(`{"newPassword": "%slonglonglonglong", "oldPassword": "%s"}`, cltest.APIEmailAdmin, cltest.Password),
			wantStatusCode: http.StatusUnprocessableEntity,
			wantErrCount:   1,
			wantErrMessage: fmt.Sprintf("%s	%s\n", utils.ErrMsgHeader, "password may not contain: \"apiuser@chainlink.test\""),
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

func TestUserController_CreateUser(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplicationEVMDisabled(t)
	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(cltest.APIEmailAdmin)

	longPassword := strings.Repeat("x", sessions.MaxBcryptPasswordLength+1)

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
			name:           "Wrong email format",
			reqBody:        fmt.Sprintf(`{"email": "12345678", "role": "view", "password": "%v"}`, cltest.Password),
			wantStatusCode: http.StatusBadRequest,
			wantErrCount:   1,
			wantErrMessage: "mail: missing '@' or angle-addr",
		},
		{
			name:           "Empty email format",
			reqBody:        fmt.Sprintf(`{"email": "", "role": "view", "password": "%v"}`, cltest.Password),
			wantStatusCode: http.StatusBadRequest,
			wantErrCount:   1,
			wantErrMessage: "Must enter an email",
		},
		{
			name:           "Empty role",
			reqBody:        fmt.Sprintf(`{"email": "abc@email.com", "role": "", "password": "%v"}`, cltest.Password),
			wantStatusCode: http.StatusBadRequest,
			wantErrCount:   1,
			wantErrMessage: "Invalid role",
		},
		{
			name:           "Too long password",
			reqBody:        fmt.Sprintf(`{"email": "abc@email.com", "role": "view", "password": "%v"}`, longPassword),
			wantStatusCode: http.StatusBadRequest,
			wantErrCount:   1,
			wantErrMessage: "must enter a password less than 50 characters",
		},
		{
			name:           "Too short password",
			reqBody:        `{"email": "abc@email.com", "role": "view", "password": "short"}`,
			wantStatusCode: http.StatusBadRequest,
			wantErrCount:   1,
			wantErrMessage: "Must be at least 16 characters long",
		},
		{
			name:           "Empty password",
			reqBody:        `{"email": "abc@email.com", "role": "view", "password": ""}`,
			wantStatusCode: http.StatusBadRequest,
			wantErrCount:   1,
			wantErrMessage: "Must be at least 16 characters long",
		},
		{
			name:           "Password contains email",
			reqBody:        `{"email": "asd@email.com", "role": "view", "password": "asd@email.comasd@email.comasd@email.com"}`,
			wantStatusCode: http.StatusBadRequest,
			wantErrCount:   1,
			wantErrMessage: `password may not contain: "asd@email.com"`,
		},
		{
			name:           "Success",
			reqBody:        fmt.Sprintf(`{"email": "%s", "role": "edit", "password": "%v"}`, cltest.MustRandomUser(t).Email, cltest.Password),
			wantStatusCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			resp, cleanup := client.Post("/v2/users", bytes.NewBufferString(tc.reqBody))
			t.Cleanup(cleanup)
			errors := cltest.ParseJSONAPIErrors(t, resp.Body)

			require.Equal(t, tc.wantStatusCode, resp.StatusCode)
			assert.Len(t, errors.Errors, tc.wantErrCount)
			if tc.wantErrMessage != "" {
				assert.Contains(t, errors.Errors[0].Detail, tc.wantErrMessage)
			}
		})
	}
}

func TestUserController_UpdateRole(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplicationEVMDisabled(t)
	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(cltest.APIEmailAdmin)
	user := cltest.MustRandomUser(t)
	err := app.SessionORM().CreateUser(&user)
	require.NoError(t, err)

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
			name:           "Success",
			reqBody:        fmt.Sprintf(`{"email": "%s", "newRole": "edit"}`, user.Email),
			wantStatusCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			resp, cleanup := client.Patch("/v2/users", bytes.NewBufferString(tc.reqBody))
			t.Cleanup(cleanup)
			errors := cltest.ParseJSONAPIErrors(t, resp.Body)

			require.Equal(t, tc.wantStatusCode, resp.StatusCode)
			assert.Len(t, errors.Errors, tc.wantErrCount)
			if tc.wantErrMessage != "" {
				assert.Contains(t, errors.Errors[0].Detail, tc.wantErrMessage)
			}
		})
	}
}

func TestUserController_DeleteUser(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplicationEVMDisabled(t)
	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(cltest.APIEmailAdmin)
	user := cltest.MustRandomUser(t)
	err := app.SessionORM().CreateUser(&user)
	require.NoError(t, err)

	resp, cleanup := client.Delete(fmt.Sprintf("/v2/users/%s", url.QueryEscape(user.Email)))
	t.Cleanup(cleanup)
	errors := cltest.ParseJSONAPIErrors(t, resp.Body)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Empty(t, errors.Errors)

	// second attempt would fail
	resp, cleanup = client.Delete(fmt.Sprintf("/v2/users/%s", url.QueryEscape(user.Email)))
	t.Cleanup(cleanup)
	errors = cltest.ParseJSONAPIErrors(t, resp.Body)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Len(t, errors.Errors, 1)
	assert.Contains(t, errors.Errors[0].Detail, "specified user not found")
}

func TestUserController_NewAPIToken(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplicationEVMDisabled(t)
	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(cltest.APIEmailAdmin)
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

	client := app.NewHTTPClient(cltest.APIEmailAdmin)
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

	client := app.NewHTTPClient(cltest.APIEmailAdmin)
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

	client := app.NewHTTPClient(cltest.APIEmailAdmin)
	req, err := json.Marshal(sessions.ChangeAuthTokenRequest{
		Password: "wrong-password",
	})
	require.NoError(t, err)
	resp, cleanup := client.Post("/v2/user/token/delete", bytes.NewBuffer(req))
	defer cleanup()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}
