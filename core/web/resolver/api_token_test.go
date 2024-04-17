package resolver

import (
	"testing"

	gqlerrors "github.com/graph-gophers/graphql-go/errors"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/auth"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	webauth "github.com/smartcontractkit/chainlink/v2/core/web/auth"
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
	variablesIncorrect := map[string]interface{}{
		"input": map[string]interface{}{
			"password": "wrong-password",
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

				f.Mocks.authProvider.On("FindUser", mock.Anything, session.User.Email).Return(*session.User, nil)
				f.Mocks.authProvider.On("TestPassword", mock.Anything, session.User.Email, defaultPassword).Return(nil)
				f.Mocks.authProvider.On("CreateAndSetAuthToken", mock.Anything, session.User).Return(&auth.Token{
					Secret:    "new-secret",
					AccessKey: "new-access-key",
				}, nil)
				f.App.On("AuthenticationProvider").Return(f.Mocks.authProvider)
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

				f.Mocks.authProvider.On("FindUser", mock.Anything, session.User.Email).Return(*session.User, nil)
				f.Mocks.authProvider.On("TestPassword", mock.Anything, session.User.Email, "wrong-password").Return(gError)
				f.App.On("AuthenticationProvider").Return(f.Mocks.authProvider)
			},
			query:     mutation,
			variables: variablesIncorrect,
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

				f.Mocks.authProvider.On("FindUser", mock.Anything, session.User.Email).Return(*session.User, gError)
				f.App.On("AuthenticationProvider").Return(f.Mocks.authProvider)
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

				f.Mocks.authProvider.On("FindUser", mock.Anything, session.User.Email).Return(*session.User, nil)
				f.Mocks.authProvider.On("TestPassword", mock.Anything, session.User.Email, defaultPassword).Return(nil)
				f.Mocks.authProvider.On("CreateAndSetAuthToken", mock.Anything, session.User).Return(nil, gError)
				f.App.On("AuthenticationProvider").Return(f.Mocks.authProvider)
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

func TestResolver_DeleteAPIToken(t *testing.T) {
	t.Parallel()

	defaultPassword := "my-password"
	mutation := `
		mutation DeleteAPIToken($input: DeleteAPITokenInput!) {
			deleteAPIToken(input: $input) {
				... on DeleteAPITokenSuccess {
					token {
						accessKey
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
	variablesIncorrect := map[string]interface{}{
		"input": map[string]interface{}{
			"password": "wrong-password",
		},
	}
	gError := errors.New("error")

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: mutation, variables: variables}, "deleteAPIToken"),
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
				err = session.User.TokenKey.UnmarshalText([]byte("new-access-key"))
				require.NoError(t, err)

				f.Mocks.authProvider.On("FindUser", mock.Anything, session.User.Email).Return(*session.User, nil)
				f.Mocks.authProvider.On("TestPassword", mock.Anything, session.User.Email, defaultPassword).Return(nil)
				f.Mocks.authProvider.On("DeleteAuthToken", mock.Anything, session.User).Return(nil)
				f.App.On("AuthenticationProvider").Return(f.Mocks.authProvider)
			},
			query:     mutation,
			variables: variables,
			result: `
				{
					"deleteAPIToken": {
						"token": {
							"accessKey": "new-access-key"
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

				f.Mocks.authProvider.On("FindUser", mock.Anything, session.User.Email).Return(*session.User, nil)
				f.Mocks.authProvider.On("TestPassword", mock.Anything, session.User.Email, "wrong-password").Return(gError)
				f.App.On("AuthenticationProvider").Return(f.Mocks.authProvider)
			},
			query:     mutation,
			variables: variablesIncorrect,
			result: `
				{
					"deleteAPIToken": {
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

				f.Mocks.authProvider.On("FindUser", mock.Anything, session.User.Email).Return(*session.User, gError)
				f.App.On("AuthenticationProvider").Return(f.Mocks.authProvider)
			},
			query:     mutation,
			variables: variables,
			result:    `null`,
			errors: []*gqlerrors.QueryError{
				{
					Extensions:    nil,
					ResolverError: gError,
					Path:          []interface{}{"deleteAPIToken"},
					Message:       "error",
				},
			},
		},
		{
			name:          "failed to delete token",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				session, ok := webauth.GetGQLAuthenticatedSession(f.Ctx)
				require.True(t, ok)
				require.NotNil(t, session)

				pwd, err := utils.HashPassword(defaultPassword)
				require.NoError(t, err)

				session.User.HashedPassword = pwd

				f.Mocks.authProvider.On("FindUser", mock.Anything, session.User.Email).Return(*session.User, nil)
				f.Mocks.authProvider.On("TestPassword", mock.Anything, session.User.Email, defaultPassword).Return(nil)
				f.Mocks.authProvider.On("DeleteAuthToken", mock.Anything, session.User).Return(gError)
				f.App.On("AuthenticationProvider").Return(f.Mocks.authProvider)
			},
			query:     mutation,
			variables: variables,
			result:    `null`,
			errors: []*gqlerrors.QueryError{
				{
					Extensions:    nil,
					ResolverError: gError,
					Path:          []interface{}{"deleteAPIToken"},
					Message:       "error",
				},
			},
		},
	}

	RunGQLTests(t, testCases)
}
