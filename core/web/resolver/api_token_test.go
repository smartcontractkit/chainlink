package resolver

import (
	"testing"

	gqlerrors "github.com/graph-gophers/graphql-go/errors"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/chainlink/core/utils"
	webauth "github.com/smartcontractkit/chainlink/core/web/auth"
)

func TestResolver_CreateAPIToken(t *testing.T) {
	t.Parallel()

	defaultPassword := "my-password"
	mutation := `
		mutation CreateAPIToken($input: CreateAPITokenInput!) {
			createAPIToken(input: $input) {
				... on CreateAPITokenSuccess {
					token {
						accessKey
						secret
					}
				}
				... on InputErrors {
					errors {
						path
						message
						code
					}
				}
			}
		}`
	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"password": defaultPassword,
		},
	}
	gError := errors.New("error")

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: mutation, variables: variables}, "createAPIToken"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				session, ok := webauth.GetGQLAuthenticatedSession(f.Ctx)
				require.True(t, ok)
				require.NotNil(t, session)

				pwd, err := utils.HashPassword(defaultPassword)
				require.NoError(t, err)

				session.User.HashedPassword = pwd

				f.Mocks.sessionsORM.On("FindUser").Return(*session.User, nil)
				f.Mocks.sessionsORM.On("GenerateAuthToken", session.User).Return(&auth.Token{
					Secret:    "new-secret",
					AccessKey: "new-access-key",
				}, nil)
				f.App.On("SessionORM").Return(f.Mocks.sessionsORM)
			},
			query:     mutation,
			variables: variables,
			result: `
				{
					"createAPIToken": {
						"token": {
							"accessKey": "new-access-key",
							"secret": "new-secret"
						}
					}
				}`,
		},
		{
			name:          "input errors",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				session, ok := webauth.GetGQLAuthenticatedSession(f.Ctx)
				require.True(t, ok)
				require.NotNil(t, session)

				session.User.HashedPassword = "wrong-password"

				f.Mocks.sessionsORM.On("FindUser").Return(*session.User, nil)
				f.App.On("SessionORM").Return(f.Mocks.sessionsORM)
			},
			query:     mutation,
			variables: variables,
			result: `
				{
					"createAPIToken": {
						"errors": [{
							"path": "password",
							"message": "incorrect password",
							"code": "INVALID_INPUT"
						}]
					}
				}`,
		},
		{
			name:          "failed to find user",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				session, ok := webauth.GetGQLAuthenticatedSession(f.Ctx)
				require.True(t, ok)
				require.NotNil(t, session)

				pwd, err := utils.HashPassword(defaultPassword)
				require.NoError(t, err)

				session.User.HashedPassword = pwd

				f.Mocks.sessionsORM.On("FindUser").Return(*session.User, gError)
				f.App.On("SessionORM").Return(f.Mocks.sessionsORM)
			},
			query:     mutation,
			variables: variables,
			result:    `null`,
			errors: []*gqlerrors.QueryError{
				{
					Extensions:    nil,
					ResolverError: gError,
					Path:          []interface{}{"createAPIToken"},
					Message:       "error",
				},
			},
		},
		{
			name:          "failed to generate token",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				session, ok := webauth.GetGQLAuthenticatedSession(f.Ctx)
				require.True(t, ok)
				require.NotNil(t, session)

				pwd, err := utils.HashPassword(defaultPassword)
				require.NoError(t, err)

				session.User.HashedPassword = pwd

				f.Mocks.sessionsORM.On("FindUser").Return(*session.User, nil)
				f.Mocks.sessionsORM.On("GenerateAuthToken", session.User).Return(nil, gError)
				f.App.On("SessionORM").Return(f.Mocks.sessionsORM)
			},
			query:     mutation,
			variables: variables,
			result:    `null`,
			errors: []*gqlerrors.QueryError{
				{
					Extensions:    nil,
					ResolverError: gError,
					Path:          []interface{}{"createAPIToken"},
					Message:       "error",
				},
			},
		},
	}

	RunGQLTests(t, testCases)
}
