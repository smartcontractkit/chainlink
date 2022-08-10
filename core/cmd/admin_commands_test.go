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
	app := startNewApplication(t)
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
			set.String("email", "", "")
			set.String("role", "", "")
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
