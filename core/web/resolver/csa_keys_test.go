package resolver

import (
	"encoding/json"
	"fmt"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/csakey"
	"github.com/stretchr/testify/assert"
	"testing"
)

type expectedKey struct {
	PubKey  string `json:"pubKey"`
	Version int    `json:"version"`
}

func Test_CSAKeysQuery(t *testing.T) {
	query := `
	query GetCSAKeys{
		csaKeys {
			results {
				pubKey
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
		fmt.Println(k.Version, k.PublicKeyString())
		expectedKeys = append(expectedKeys, expectedKey{
			Version: k.Version,
			PubKey:  k.PublicKeyString(),
		})
	}

	fmt.Println(expectedKeys)

	d, err := json.Marshal(map[string]interface{}{
		"csaKeys": map[string]interface{}{
			"results": expectedKeys,
		},
	})
	assert.NoError(t, err)
	expectedResult := string(d)

	fmt.Println(expectedResult)

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "csaKeys"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
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
