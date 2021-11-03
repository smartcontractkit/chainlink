package resolver

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocrkey"
)

func TestOCR(t *testing.T) {
	t.Parallel()

	query := `
		query GetOCRKeys {
			ocrKeys {
				results {
					id
					configPublicKey
					offChainPublicKey
					onChainSigningAddress
				}
			}
		}
	`

	fakeKeys := []ocrkey.KeyV2{}
	expectedKeys := []map[string]string{}
	for i := 0; i < 2; i++ {
		k, err := ocrkey.NewV2()
		assert.NoError(t, err)
		fakeKeys = append(fakeKeys, k)
		expectedKeys = append(expectedKeys, map[string]string{
			"id":                    k.ID(),
			"configPublicKey":       ocrkey.ConfigPublicKey(k.PublicKeyConfig()).String(),
			"offChainPublicKey":     k.OffChainSigning.PublicKey().String(),
			"onChainSigningAddress": k.OnChainSigning.Address().String(),
		})
	}

	d, err := json.Marshal(map[string]interface{}{
		"ocrKeys": map[string]interface{}{
			"results": expectedKeys,
		},
	})
	assert.NoError(t, err)
	expected := string(d)

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "ocrKeys"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.ocr.On("GetAll").Return(fakeKeys, nil)
				f.Mocks.keystore.On("OCR").Return(f.Mocks.ocr)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
			},
			query:  query,
			result: expected,
		},
	}

	RunGQLTests(t, testCases)
}
