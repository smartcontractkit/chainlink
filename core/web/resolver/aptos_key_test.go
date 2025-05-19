package resolver

import (
	"context"
	"errors"
	"fmt"
	"testing"

	gqlerrors "github.com/graph-gophers/graphql-go/errors"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/keystest"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/aptoskey"
)

func TestResolver_AptosKeys(t *testing.T) {
	t.Parallel()

	query := `
		query GetAptosKeys {
			aptosKeys {
				results {
					id
					account
				}
			}
		}`
	k := aptoskey.MustNewInsecure(keystest.NewRandReaderFromSeed(1))
	result := fmt.Sprintf(`
	{
		"aptosKeys": {
			"results": [
				{
					"id": "%s",
					"account": "%s"
				}
			]
		}
	}`, k.PublicKeyStr(), k.Account())
	gError := errors.New("error")

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "aptosKeys"),
		{
			name:          "success",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.Mocks.aptos.On("GetAll").Return([]aptoskey.Key{k}, nil)
				f.Mocks.keystore.On("Aptos").Return(f.Mocks.aptos)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
			},
			query:  query,
			result: result,
		},
		{
			name:          "no keys returned by GetAll",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.Mocks.aptos.On("GetAll").Return([]aptoskey.Key{}, gError)
				f.Mocks.keystore.On("Aptos").Return(f.Mocks.aptos)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
			},
			query:  query,
			result: `null`,
			errors: []*gqlerrors.QueryError{
				{
					Extensions:    nil,
					ResolverError: gError,
					Path:          []interface{}{"aptosKeys"},
					Message:       gError.Error(),
				},
			},
		},
	}

	RunGQLTests(t, testCases)
}
