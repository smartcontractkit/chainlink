package web_test

import (
	"bytes"
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/stretchr/testify/assert"
)

func TestUserController_UpdatePassword(t *testing.T) {
	appWithUser, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()
	client := appWithUser.NewHTTPClient()

	// Invalid request
	resp, cleanup := client.Patch("/v2/user/password", bytes.NewBufferString(""))
	defer cleanup()
	errors := cltest.ParseJSONAPIErrors(resp.Body)
	assert.Equal(t, 422, resp.StatusCode)
	assert.Len(t, errors.Errors, 1)

	// Old password is wrong
	resp, cleanup = client.Patch(
		"/v2/user/password",
		bytes.NewBufferString(`{"oldPassword": "wrong password"}`))
	defer cleanup()
	errors = cltest.ParseJSONAPIErrors(resp.Body)
	assert.Equal(t, 401, resp.StatusCode)
	assert.Len(t, errors.Errors, 1)
	assert.Equal(t, "Old password does not match", errors.Errors[0].Detail)

	// Success
	resp, cleanup = client.Patch(
		"/v2/user/password",
		bytes.NewBufferString(`{"newPassword": "password", "oldPassword": "password"}`))
	defer cleanup()
	errors = cltest.ParseJSONAPIErrors(resp.Body)
	assert.Equal(t, 200, resp.StatusCode)
}
