package cmd_test

import (
	"bytes"
	"flag"
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/sessions"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func TestShell_CreateUser(t *testing.T) {
	app := startNewApplicationV2(t, nil)
	client, _ := app.NewShellAndRenderer()
	client.PasswordPrompter = cltest.MockPasswordPrompter{
		Password: cltest.Password,
	}

	tests := []struct {
		name  string
		email string
		role  string
		err   string
	}{
		{"Invalid email", "//", "", "mail: missing '@' or angle-addr"},
		{"No params", "", "", "Must enter an email"},
		{"No email", "", "view", "Must enter an email"},
		{"User exists", cltest.APIEmailAdmin, "admin", fmt.Sprintf(`user with email %s already exists`, cltest.APIEmailAdmin)},
		{"Valid params", cltest.MustRandomUser(t).Email, "view", ""},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			set := flag.NewFlagSet("test", 0)
			flagSetApplyFromAction(client.CreateUser, set, "")

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

func TestShell_ChangeRole(t *testing.T) {
	ctx := testutils.Context(t)
	app := startNewApplicationV2(t, nil)
	client, _ := app.NewShellAndRenderer()
	user := cltest.MustRandomUser(t)
	require.NoError(t, app.AuthenticationProvider().CreateUser(ctx, &user))

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
			flagSetApplyFromAction(client.ChangeRole, set, "")

			require.NoError(t, set.Set("email", test.email))
			require.NoError(t, set.Set("new-role", test.role))
			c := cli.NewContext(nil, set, nil)
			if test.err != "" {
				assert.ErrorContains(t, client.ChangeRole(c), test.err)
			} else {
				assert.NoError(t, client.ChangeRole(c))
			}
		})
	}
}

func TestShell_DeleteUser(t *testing.T) {
	ctx := testutils.Context(t)
	app := startNewApplicationV2(t, nil)
	client, _ := app.NewShellAndRenderer()
	user := cltest.MustRandomUser(t)
	require.NoError(t, app.BasicAdminUsersORM().CreateUser(ctx, &user))

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
			flagSetApplyFromAction(client.DeleteUser, set, "")

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

func TestShell_ListUsers(t *testing.T) {
	ctx := testutils.Context(t)
	app := startNewApplicationV2(t, nil)
	client, _ := app.NewShellAndRenderer()
	user := cltest.MustRandomUser(t)
	require.NoError(t, app.AuthenticationProvider().CreateUser(ctx, &user))

	set := flag.NewFlagSet("test", 0)
	flagSetApplyFromAction(client.ListUsers, set, "")
	c := cli.NewContext(nil, set, nil)

	testRenderer := &testRenderer{}
	client.Renderer = testRenderer
	assert.NoError(t, client.ListUsers(c), user.Email)

	userPresenterFound := false
	for _, presenter := range testRenderer.presenters {
		if presenter.Email == user.Email {
			userPresenterFound = true
			assert.Equal(t, presenter.Role, user.Role)
			userHasActiveApiToken, err := strconv.ParseBool(presenter.HasActiveApiToken)
			assert.NoError(t, err)
			assert.Equal(t, userHasActiveApiToken, user.TokenKey.String != "")
			assert.True(t, presenter.CreatedAt.Equal(user.CreatedAt))
			assert.True(t, presenter.CreatedAt.Equal(user.UpdatedAt))
		}
	}
	assert.Truef(t, userPresenterFound, "expected to find user %s in presenter list", user.Email)
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

type testRenderer struct {
	presenters []cmd.AdminUsersPresenter
}

func (t *testRenderer) Render(i interface{}, s ...string) error {
	adminPresenters := i.(*cmd.AdminUsersPresenters)
	t.presenters = *adminPresenters
	return nil
}
