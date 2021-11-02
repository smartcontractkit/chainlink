package resolver

import (
	"encoding/json"
	"testing"

	coreMocks "github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocrkey"
	"github.com/stretchr/testify/assert"
)

func TestOCR(t *testing.T) {
	t.Parallel()

	query := `
		query GetOCRKeys {
			ocrKeys {
				keys {
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
			"configPublicKey":       ocrkey.ConfigPublicKey(k.PublicKeyConfig()).String(),
			"offChainPublicKey":     k.OffChainSigning.PublicKey().String(),
			"onChainSigningAddress": k.OnChainSigning.Address().String(),
		})
	}

	d, err := json.Marshal(map[string]interface{}{
		"ocrKeys": map[string]interface{}{
			"keys": expectedKeys,
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

				mockOCR := &coreMocks.OCR{}
				mockOCR.On("GetAll").Return(fakeKeys, nil)
				mockKeyStore := &coreMocks.Master{}
				mockKeyStore.On("OCR").Return(mockOCR)
				f.App.On("GetKeyStore").Return(mockKeyStore)
			},
			query:  query,
			result: expected,
		},
	}

	RunGQLTests(t, testCases)
}
