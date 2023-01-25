package cmd_test

import (
	"flag"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
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
