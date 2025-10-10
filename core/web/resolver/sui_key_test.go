package resolver

import (
	"context"
	"errors"
	"fmt"
	"testing"

	gqlerrors "github.com/graph-gophers/graphql-go/errors"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/keystest"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/suikey"
)

func TestResolver_SuiKeys(t *testing.T) {
	t.Parallel()

	query := `
		query GetSuiKeys {
			suiKeys {
				results {
					id
					account
				}
			}
		}`
	k := suikey.MustNewInsecure(keystest.NewRandReaderFromSeed(1))
	result := fmt.Sprintf(`
	{
		"suiKeys": {
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
		unauthorizedTestCase(GQLTestCase{query: query}, "suiKeys"),
		{
			name:          "success",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.Mocks.sui.On("GetAll").Return([]suikey.Key{k}, nil)
				f.Mocks.keystore.On("Sui").Return(f.Mocks.sui)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
			},
			query:  query,
			result: result,
		},
		{
			name:          "no keys returned by GetAll",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.Mocks.sui.On("GetAll").Return([]suikey.Key{}, gError)
				f.Mocks.keystore.On("Sui").Return(f.Mocks.sui)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
			},
			query:  query,
			result: `null`,
			errors: []*gqlerrors.QueryError{
				{
					Extensions:    nil,
					ResolverError: gError,
					Path:          []any{"suiKeys"},
					Message:       gError.Error(),
				},
			},
		},
	}

	RunGQLTests(t, testCases)
}
