package resolver

import (
	"context"
	"errors"
	"fmt"
	"testing"

	gqlerrors "github.com/graph-gophers/graphql-go/errors"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/keystest"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/starkkey"
)

func TestResolver_StarkNetKeys(t *testing.T) {
	t.Parallel()

	query := `
		query GetStarkNetKeys {
			starknetKeys {
				results {
					id
				}
			}
		}`
	k := starkkey.MustNewInsecure(keystest.NewRandReaderFromSeed(1))
	result := fmt.Sprintf(`
	{
		"starknetKeys": {
			"results": [
				{
					"id": "%s"
				}
			]
		}
	}`, k.StarkKeyStr())
	gError := errors.New("error")

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "starknetKeys"),
		{
			name:          "success",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.Mocks.starknet.On("GetAll").Return([]starkkey.Key{k}, nil)
				f.Mocks.keystore.On("StarkNet").Return(f.Mocks.starknet)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
			},
			query:  query,
			result: result,
		},
		{
			name:          "no keys returned by GetAll",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.Mocks.starknet.On("GetAll").Return([]starkkey.Key{}, gError)
				f.Mocks.keystore.On("StarkNet").Return(f.Mocks.starknet)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
			},
			query:  query,
			result: `null`,
			errors: []*gqlerrors.QueryError{
				{
					Extensions:    nil,
					ResolverError: gError,
					Path:          []interface{}{"starknetKeys"},
					Message:       gError.Error(),
				},
			},
		},
	}

	RunGQLTests(t, testCases)
}
