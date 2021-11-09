package resolver

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocrkey"
)

func TestOCR(t *testing.T) {
	t.Parallel()

	query := `
		query GetOCRKeyBundles {
			ocrKeyBundles {
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
	expectedBundles := []map[string]string{}
	for i := 0; i < 2; i++ {
		k, err := ocrkey.NewV2()
		assert.NoError(t, err)
		fakeKeys = append(fakeKeys, k)
		expectedBundles = append(expectedBundles, map[string]string{
			"id":                    k.ID(),
			"configPublicKey":       ocrkey.ConfigPublicKey(k.PublicKeyConfig()).String(),
			"offChainPublicKey":     k.OffChainSigning.PublicKey().String(),
			"onChainSigningAddress": k.OnChainSigning.Address().String(),
		})
	}

	d, err := json.Marshal(map[string]interface{}{
		"ocrKeyBundles": map[string]interface{}{
			"results": expectedBundles,
		},
	})
	assert.NoError(t, err)
	expected := string(d)

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "ocrKeyBundles"),
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

func TestOCRCreateBundle(t *testing.T) {
	t.Parallel()

	mutation := `
		mutation CreateOCRKeyBundle {
			createOCRKeyBundle {
				bundle {
					id
					configPublicKey
					offChainPublicKey
					onChainSigningAddress
				}
			}
		}
	`

	fakeKey, err := ocrkey.NewV2()
	assert.NoError(t, err)

	d, err := json.Marshal(map[string]interface{}{
		"createOCRKeyBundle": map[string]interface{}{
			"bundle": map[string]interface{}{
				"id":                    fakeKey.ID(),
				"configPublicKey":       ocrkey.ConfigPublicKey(fakeKey.PublicKeyConfig()).String(),
				"offChainPublicKey":     fakeKey.OffChainSigning.PublicKey().String(),
				"onChainSigningAddress": fakeKey.OnChainSigning.Address().String(),
			},
		},
	})
	assert.NoError(t, err)
	expected := string(d)

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: mutation}, "createOCRKeyBundle"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.ocr.On("Create").Return(fakeKey, nil)
				f.Mocks.keystore.On("OCR").Return(f.Mocks.ocr)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
			},
			query:  mutation,
			result: expected,
		},
	}

	RunGQLTests(t, testCases)
}

func TestOCRDeleteBundle(t *testing.T) {
	t.Parallel()

	fakeKey, err := ocrkey.NewV2()
	assert.NoError(t, err)

	mutation := `
		mutation DeleteOCRKeyBundle($id: ID!) {
			deleteOCRKeyBundle(id: $id) {
				... on DeleteOCRKeyBundleSuccess {
					bundle {
						id
						configPublicKey
						offChainPublicKey
						onChainSigningAddress
					}
				}
				... on NotFoundError {
					message
					code
				}
			}
		}
	`
	variables := map[string]interface{}{
		"id": fakeKey.ID(),
	}

	d, err := json.Marshal(map[string]interface{}{
		"deleteOCRKeyBundle": map[string]interface{}{
			"bundle": map[string]interface{}{
				"id":                    fakeKey.ID(),
				"configPublicKey":       ocrkey.ConfigPublicKey(fakeKey.PublicKeyConfig()).String(),
				"offChainPublicKey":     fakeKey.OffChainSigning.PublicKey().String(),
				"onChainSigningAddress": fakeKey.OnChainSigning.Address().String(),
			},
		},
	})
	assert.NoError(t, err)
	expected := string(d)

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: mutation, variables: variables}, "deleteOCRKeyBundle"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.ocr.On("Delete", fakeKey.ID()).Return(fakeKey, nil)
				f.Mocks.keystore.On("OCR").Return(f.Mocks.ocr)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
			},
			query:     mutation,
			variables: variables,
			result:    expected,
		},
		{
			name:          "not found error",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.ocr.On("Delete", fakeKey.ID()).Return(ocrkey.KeyV2{}, keystore.KeyNotFoundError{ID: "helloWorld"})
				f.Mocks.keystore.On("OCR").Return(f.Mocks.ocr)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
			},
			query:     mutation,
			variables: variables,
			result: `{
				"deleteOCRKeyBundle": {
					"code":"NOT_FOUND",
					"message":"unable to find OCR key with id helloWorld"
				}
			}`,
		},
	}

	RunGQLTests(t, testCases)
}
