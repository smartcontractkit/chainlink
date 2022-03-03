package resolver

import (
	"errors"
	"fmt"
	"testing"

	gqlerrors "github.com/graph-gophers/graphql-go/errors"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/keystest"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/solkey"
)

func TestResolver_SolanaKeys(t *testing.T) {
	t.Parallel()

	query := `
		query GetSolanaKeys {
			solanaKeys {
				results {
					id
				}
			}
		}`
	k := solkey.MustNewInsecure(keystest.NewRandReaderFromSeed(1))
	result := fmt.Sprintf(`
	{
		"solanaKeys": {
			"results": [
				{
					"id": "%s"
				}
			]
		}
	}`, k.PublicKeyStr())
	gError := errors.New("error")

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "solanaKeys"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.solana.On("GetAll").Return([]solkey.Key{k}, nil)
				f.Mocks.keystore.On("Solana").Return(f.Mocks.solana)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
			},
			query:  query,
			result: result,
		},
		{
			name:          "generic error on GetAll",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.solana.On("GetAll").Return([]solkey.Key{}, gError)
				f.Mocks.keystore.On("Solana").Return(f.Mocks.solana)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
			},
			query:  query,
			result: `null`,
			errors: []*gqlerrors.QueryError{
				{
					Extensions:    nil,
					ResolverError: gError,
					Path:          []interface{}{"solanaKeys"},
					Message:       gError.Error(),
				},
			},
		},
	}

	RunGQLTests(t, testCases)
}
