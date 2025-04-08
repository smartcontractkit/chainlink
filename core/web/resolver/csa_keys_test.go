package resolver

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
)

type expectedKey struct {
	ID      string `json:"id"`
	PubKey  string `json:"publicKey"`
	Version int    `json:"version"`
}

func Test_CSAKeysQuery(t *testing.T) {
	query := `
	query GetCSAKeys {
		csaKeys {
			results {
				id
				publicKey
				version
			}
		}
	}`

	var fakeKeys []csakey.KeyV2
	var expectedKeys []expectedKey

	for i := 0; i < 5; i++ {
		k, err := csakey.NewV2()
		assert.NoError(t, err)

		fakeKeys = append(fakeKeys, k)
		expectedKeys = append(expectedKeys, expectedKey{
			ID:      k.ID(),
			Version: k.Version,
			PubKey:  "csa_" + k.PublicKeyString(),
		})
	}

	d, err := json.Marshal(map[string]interface{}{
		"csaKeys": map[string]interface{}{
			"results": expectedKeys,
		},
	})
	assert.NoError(t, err)
	expectedResult := string(d)

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "csaKeys"),
		{
			name:          "success",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.Mocks.csa.On("GetAll").Return(fakeKeys, nil)
				f.Mocks.keystore.On("CSA").Return(f.Mocks.csa)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
			},
			query:  query,
			result: expectedResult,
		},
	}

	RunGQLTests(t, testCases)
}

func Test_CreateCSAKey(t *testing.T) {
	query := `
	mutation CreateCSAKey {
		createCSAKey {
			... on CreateCSAKeySuccess {
				csaKey {
					id
					version
					publicKey
				}
			}
			... on CSAKeyExistsError {
				message
				code
			}
		}
	}`

	fakeKey, err := csakey.NewV2()
	assert.NoError(t, err)

	expected := `
		{
			"createCSAKey": {
				"csaKey": {
					"id": "%s",
					"version": %d,
					"publicKey": "csa_%s"
				}
			}
		}
	`

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "createCSAKey"),
		{
			name:          "success",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.Mocks.csa.On("Create", mock.Anything).Return(fakeKey, nil)
				f.Mocks.keystore.On("CSA").Return(f.Mocks.csa)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
			},
			query:  query,
			result: fmt.Sprintf(expected, fakeKey.ID(), fakeKey.Version, fakeKey.PublicKeyString()),
		},
		{
			name:          "csa key exists error",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.Mocks.csa.On("Create", mock.Anything).Return(csakey.KeyV2{}, keystore.ErrCSAKeyExists)
				f.Mocks.keystore.On("CSA").Return(f.Mocks.csa)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
			},
			query: query,
			result: fmt.Sprintf(`
			{
				"createCSAKey": {
					"message": "%s",
					"code": "UNPROCESSABLE"
				}
			}`, keystore.ErrCSAKeyExists.Error()),
		},
	}

	RunGQLTests(t, testCases)
}

func Test_DeleteCSAKey(t *testing.T) {
	query := `
	mutation DeleteCSAKey($id: ID!) {
		deleteCSAKey(id: $id) {
			... on DeleteCSAKeySuccess {
				csaKey {
					id
					version
					publicKey
				}
			}
			... on NotFoundError {
				message
				code
			}
		}
	}`

	fakeKey, err := csakey.NewV2()
	assert.NoError(t, err)

	expected := `
		{
			"deleteCSAKey": {
				"csaKey": {
					"id": "%s",
					"version": %d,
					"publicKey": "csa_%s"
				}
			}
		}
	`
	variables := map[string]interface{}{"id": fakeKey.ID()}

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query, variables: variables}, "deleteCSAKey"),
		{
			name:          "success",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
				f.Mocks.keystore.On("CSA").Return(f.Mocks.csa)
				f.Mocks.csa.On("Delete", mock.Anything, fakeKey.ID()).Return(fakeKey, nil)
			},
			query:     query,
			variables: variables,
			result:    fmt.Sprintf(expected, fakeKey.ID(), fakeKey.Version, fakeKey.PublicKeyString()),
		},
		{
			name:          "not found error",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
				f.Mocks.keystore.On("CSA").Return(f.Mocks.csa)
				f.Mocks.csa.
					On("Delete", mock.Anything, fakeKey.ID()).
					Return(csakey.KeyV2{}, keystore.KeyNotFoundError{ID: fakeKey.ID(), KeyType: "CSA"})
			},
			query:     query,
			variables: variables,
			result: fmt.Sprintf(`
			{
				"deleteCSAKey": {
					"message": "unable to find CSA key with id %s",
					"code": "NOT_FOUND"
				}
			}`, fakeKey.ID()),
		},
	}

	RunGQLTests(t, testCases)
}
