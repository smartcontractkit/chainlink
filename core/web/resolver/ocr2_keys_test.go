package resolver

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	gqlerrors "github.com/graph-gophers/graphql-go/errors"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/keystest"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocr2key"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolver_GetOCR2KeyBundles(t *testing.T) {
	t.Parallel()

	query := `
		query GetOCR2KeyBundles {
			ocr2KeyBundles {
				results {
					chainType
					configPublicKey
					offChainPublicKey
					onChainPublicKey
				}
			}
		}
	`

	gError := errors.New("error")
	fakeKeys := []ocr2key.KeyBundle{
		ocr2key.MustNewInsecure(keystest.NewRandReaderFromSeed(1), "evm"),
		ocr2key.MustNewInsecure(keystest.NewRandReaderFromSeed(1), "terra"),
		ocr2key.MustNewInsecure(keystest.NewRandReaderFromSeed(1), "solana"),
	}
	expectedBundles := []map[string]interface{}{}
	for _, k := range fakeKeys {
		configPublic := k.ConfigEncryptionPublicKey()
		ct, err := ToOCR2ChainType(string(k.ChainType()))
		require.NoError(t, err)

		expectedBundles = append(expectedBundles, map[string]interface{}{
			"chainType":         ct,
			"onChainPublicKey":  fmt.Sprintf("ocr2on_%s_%s", k.ChainType(), k.OnChainPublicKey()),
			"configPublicKey":   fmt.Sprintf("ocr2cfg_%s_%s", k.ChainType(), hex.EncodeToString(configPublic[:])),
			"offChainPublicKey": fmt.Sprintf("ocr2off_%s_%s", k.ChainType(), hex.EncodeToString(k.OffchainPublicKey())),
		})
	}

	d, err := json.Marshal(map[string]interface{}{
		"ocr2KeyBundles": map[string]interface{}{
			"results": expectedBundles,
		},
	})
	assert.NoError(t, err)
	expected := string(d)

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "ocr2KeyBundles"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.ocr2.On("GetAll").Return(fakeKeys, nil)
				f.Mocks.keystore.On("OCR2").Return(f.Mocks.ocr2)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
			},
			query:  query,
			result: expected,
		},
		{
			name:          "generic error on GetAll()",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.ocr2.On("GetAll").Return(nil, gError)
				f.Mocks.keystore.On("OCR2").Return(f.Mocks.ocr2)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
			},
			query:  query,
			result: `null`,
			errors: []*gqlerrors.QueryError{
				{
					Extensions:    nil,
					ResolverError: gError,
					Path:          []interface{}{"ocr2KeyBundles"},
					Message:       gError.Error(),
				},
			},
		},
	}

	RunGQLTests(t, testCases)
}
