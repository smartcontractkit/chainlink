package resolver

import (
	"context"
	"errors"
	"fmt"
	"testing"

	gqlerrors "github.com/graph-gophers/graphql-go/errors"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/keystest"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/tonkey"
)

func TestResolver_TONKeys(t *testing.T) {
	t.Parallel()

	query := `
		query GetTONKeys {
			tonKeys {
				results {
					id
					addressBase64
					rawAddress
				}
			}
		}`
	k := tonkey.MustNewInsecure(keystest.NewRandReaderFromSeed(1))
	result := fmt.Sprintf(`
	{
		"tonKeys": {
			"results": [
				{
					"id": "%s",
					"addressBase64": "%s",
					"rawAddress": "%s"
				}
			]
		}
	}`, k.ID(), k.AddressBase64(), k.RawAddress())
	gError := errors.New("error")

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "tonKeys"),
		{
			name:          "success",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.Mocks.ton.On("GetAll").Return([]tonkey.Key{k}, nil)
				f.Mocks.keystore.On("TON").Return(f.Mocks.ton)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
			},
			query:  query,
			result: result,
		},
		{
			name:          "no keys returned by GetAll",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.Mocks.ton.On("GetAll").Return([]tonkey.Key{}, gError)
				f.Mocks.keystore.On("TON").Return(f.Mocks.ton)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
			},
			query:  query,
			result: `null`,
			errors: []*gqlerrors.QueryError{
				{
					Extensions:    nil,
					ResolverError: gError,
					Path:          []interface{}{"tonKeys"},
					Message:       gError.Error(),
				},
			},
		},
	}

	RunGQLTests(t, testCases)
}
