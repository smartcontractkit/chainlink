package cmd_test

import (
	"bytes"
	"flag"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/sessions"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

func TestClient_CreateUser(t *testing.T) {
	app := startNewApplicationV2(t, nil)
	client, _ := app.NewClientAndRenderer()
	client.PasswordPrompter = cltest.MockPasswordPrompter{
		Password: cltest.Password,
	}

	tests := []struct {
		name  string
		email string
		role  string
		err   string
	}{
		{"Invalid request", "//", "", "parseResponse error"},
		{"No params", "", "", "Invalid role"},
		{"No email", "", "view", "Must enter an email"},
		{"User exists", cltest.APIEmailAdmin, "admin", fmt.Sprintf(`user with email %s already exists`, cltest.APIEmailAdmin)},
		{"Valid params", cltest.MustRandomUser(t).Email, "view", ""},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			set := flag.NewFlagSet("test", 0)
			cltest.FlagSetApplyFromAction(client.CreateUser, set, "")

			require.NoError(t, set.Set("email", test.email))
			require.NoError(t, set.Set("role", test.role))

			c := cli.NewContext(nil, set, nil)
			if test.err != "" {
				assert.ErrorContains(t, client.CreateUser(c), test.err)
			} else {
				assert.NoError(t, client.CreateUser(c))
			}
		})
	}
}

func TestClient_ChangeRole(t *testing.T) {
	app := startNewApplicationV2(t, nil)
	client, _ := app.NewClientAndRenderer()
	user := cltest.MustRandomUser(t)
	require.NoError(t, app.SessionORM().CreateUser(&user))

	tests := []struct {
		name  string
		email string
		role  string
		err   string
	}{
		{"Invalid request", "//", "", "parseResponse error"},
		{"No params", "", "", "must specify an email"},
		{"No email", "", "view", "must specify an email"},
		{"No role", user.Email, "", "must specify a new role"},
		{"Unknown role", user.Email, "foo", "new role does not exist"},
		{"Unknown user", cltest.MustRandomUser(t).Email, "admin", "error updating API user"},
		{"Valid params", user.Email, "view", ""},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			set := flag.NewFlagSet("test", 0)
			cltest.FlagSetApplyFromAction(client.ChangeRole, set, "")

			require.NoError(t, set.Set("email", test.email))
			require.NoError(t, set.Set("newrole", test.role))
			c := cli.NewContext(nil, set, nil)
			if test.err != "" {
				assert.ErrorContains(t, client.ChangeRole(c), test.err)
			} else {
				assert.NoError(t, client.ChangeRole(c))
			}
		})
	}
}

func TestClient_DeleteUser(t *testing.T) {
	app := startNewApplicationV2(t, nil)
	client, _ := app.NewClientAndRenderer()
	user := cltest.MustRandomUser(t)
	require.NoError(t, app.SessionORM().CreateUser(&user))

	tests := []struct {
		name  string
		email string
		err   string
	}{
		{"Invalid request", "//", "parseResponse error"},
		{"No email", "", "must specify an email"},
		{"Unknown email", "foo", "specified user not found"},
		{"Valid params", user.Email, ""},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			set := flag.NewFlagSet("test", 0)
			cltest.FlagSetApplyFromAction(client.DeleteUser, set, "")

			require.NoError(t, set.Set("email", test.email))
			c := cli.NewContext(nil, set, nil)
			if test.err != "" {
				assert.ErrorContains(t, client.DeleteUser(c), test.err)
			} else {
				assert.NoError(t, client.DeleteUser(c))
			}
		})
	}
}

func TestClient_ListUsers(t *testing.T) {
	app := startNewApplicationV2(t, nil)
	client, _ := app.NewClientAndRenderer()
	user := cltest.MustRandomUser(t)
	require.NoError(t, app.SessionORM().CreateUser(&user))

	set := flag.NewFlagSet("test", 0)
	cltest.FlagSetApplyFromAction(client.ListUsers, set, "")
	c := cli.NewContext(nil, set, nil)

	buffer := bytes.NewBufferString("")
	client.Renderer = cmd.RendererTable{Writer: buffer}

	assert.NoError(t, client.ListUsers(c), user.Email)

	output := buffer.String()
	assert.Contains(t, output, user.Email)
	assert.Contains(t, output, user.Role)
	assert.Contains(t, output, user.TokenKey.String)
	assert.Contains(t, output, user.CreatedAt.String())
	assert.Contains(t, output, user.UpdatedAt.String())
}

func TestAdminUsersPresenter_RenderTable(t *testing.T) {
	user := sessions.User{
		Email:     "foo@bar.com",
		Role:      "admin",
		CreatedAt: time.Now(),
		TokenKey:  null.StringFrom("tokenKey"),
		UpdatedAt: time.Now().Add(time.Duration(rand.Intn(10000)) * time.Second),
	}

	presenter := cmd.AdminUsersPresenter{
		JAID: cmd.JAID{ID: user.Email},
		UserResource: presenters.UserResource{
			JAID:              presenters.JAID{ID: user.Email},
			Email:             user.Email,
			Role:              user.Role,
			HasActiveApiToken: user.TokenKey.String,
			CreatedAt:         user.CreatedAt,
			UpdatedAt:         user.UpdatedAt,
		},
	}

	buffer := bytes.NewBufferString("")
	r := cmd.RendererTable{Writer: buffer}

	require.NoError(t, presenter.RenderTable(r))

	output := buffer.String()
	assert.Contains(t, output, user.Email)
	assert.Contains(t, output, user.Role)
	assert.Contains(t, output, user.TokenKey.String)
	assert.Contains(t, output, user.CreatedAt.String())
	assert.Contains(t, output, user.UpdatedAt.String())
}
