package resolver

import (
	"context"
	"errors"
	"fmt"
	"testing"

	gqlerrors "github.com/graph-gophers/graphql-go/errors"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/stellarkey"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/keystest"
)

func TestResolver_StellarKeys(t *testing.T) {
	t.Parallel()

	query := `
		query GetStellarKeys {
			stellarKeys {
				results {
					id
					account
				}
			}
		}`
	k := stellarkey.MustNewInsecure(keystest.NewRandReaderFromSeed(1))
	result := fmt.Sprintf(`
	{
		"stellarKeys": {
			"results": [
				{
					"id": "%s",
					"account": "%s"
				}
			]
		}
	}`, k.ID(), k.Account())
	gError := errors.New("error")

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "stellarKeys"),
		{
			name:          "success",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.Mocks.stellar.On("GetAll").Return([]stellarkey.Key{k}, nil)
				f.Mocks.keystore.On("Stellar").Return(f.Mocks.stellar)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
			},
			query:  query,
			result: result,
		},
		{
			name:          "no keys returned by GetAll",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.Mocks.stellar.On("GetAll").Return([]stellarkey.Key{}, gError)
				f.Mocks.keystore.On("Stellar").Return(f.Mocks.stellar)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
			},
			query:  query,
			result: `null`,
			errors: []*gqlerrors.QueryError{
				{
					Extensions:    nil,
					ResolverError: gError,
					Path:          []any{"stellarKeys"},
					Message:       gError.Error(),
				},
			},
		},
	}

	RunGQLTests(t, testCases)
}
