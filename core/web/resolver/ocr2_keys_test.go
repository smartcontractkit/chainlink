package resolver

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	gqlerrors "github.com/graph-gophers/graphql-go/errors"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/keystest"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
)

func TestResolver_GetOCR2KeyBundles(t *testing.T) {
	t.Parallel()

	query := `
		query GetOCR2KeyBundles {
			ocr2KeyBundles {
				results {
					id
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
		ocr2key.MustNewInsecure(keystest.NewRandReaderFromSeed(1), "cosmos"),
		ocr2key.MustNewInsecure(keystest.NewRandReaderFromSeed(1), "solana"),
		ocr2key.MustNewInsecure(keystest.NewRandReaderFromSeed(1), "starknet"),
		ocr2key.MustNewInsecure(keystest.NewRandReaderFromSeed(1), "aptos"),
		ocr2key.MustNewInsecure(keystest.NewRandReaderFromSeed(1), "tron"),
	}
	expectedBundles := []map[string]interface{}{}
	for _, k := range fakeKeys {
		configPublic := k.ConfigEncryptionPublicKey()
		ct, err := ToOCR2ChainType(string(k.ChainType()))
		require.NoError(t, err)
		pubKey := k.OffchainPublicKey()
		expectedBundles = append(expectedBundles, map[string]interface{}{
			"id":                k.ID(),
			"chainType":         ct,
			"onChainPublicKey":  fmt.Sprintf("ocr2on_%s_%s", k.ChainType(), k.OnChainPublicKey()),
			"configPublicKey":   fmt.Sprintf("ocr2cfg_%s_%s", k.ChainType(), hex.EncodeToString(configPublic[:])),
			"offChainPublicKey": fmt.Sprintf("ocr2off_%s_%s", k.ChainType(), hex.EncodeToString(pubKey[:])),
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
			before: func(ctx context.Context, f *gqlTestFramework) {
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
			before: func(ctx context.Context, f *gqlTestFramework) {
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

func TestResolver_CreateOCR2KeyBundle(t *testing.T) {
	t.Parallel()

	mutation := `
		mutation CreateOCR2KeyBundle($chainType: OCR2ChainType!) {
			createOCR2KeyBundle(chainType: $chainType) {
				... on CreateOCR2KeyBundleSuccess {
					bundle {
						id
						chainType
						configPublicKey
						offChainPublicKey
						onChainPublicKey
					}
				}
			}
		}
	`

	gError := errors.New("error")
	fakeKey := ocr2key.MustNewInsecure(keystest.NewRandReaderFromSeed(1), "evm")
	ct, err := ToOCR2ChainType(string(fakeKey.ChainType()))
	require.NoError(t, err)

	configPublic := fakeKey.ConfigEncryptionPublicKey()
	pubKey := fakeKey.OffchainPublicKey()
	d, err := json.Marshal(map[string]interface{}{
		"createOCR2KeyBundle": map[string]interface{}{
			"bundle": map[string]interface{}{
				"id":                fakeKey.ID(),
				"chainType":         ct,
				"onChainPublicKey":  fmt.Sprintf("ocr2on_%s_%s", fakeKey.ChainType(), fakeKey.OnChainPublicKey()),
				"configPublicKey":   fmt.Sprintf("ocr2cfg_%s_%s", fakeKey.ChainType(), hex.EncodeToString(configPublic[:])),
				"offChainPublicKey": fmt.Sprintf("ocr2off_%s_%s", fakeKey.ChainType(), hex.EncodeToString(pubKey[:])),
			},
		},
	})
	assert.NoError(t, err)
	expected := string(d)

	variables := map[string]interface{}{
		"chainType": OCR2ChainTypeEVM,
	}

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: mutation, variables: variables}, "createOCR2KeyBundle"),
		{
			name:          "success",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.Mocks.ocr2.On("Create", mock.Anything, chaintype.ChainType("evm")).Return(fakeKey, nil)
				f.Mocks.keystore.On("OCR2").Return(f.Mocks.ocr2)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
			},
			query:     mutation,
			variables: variables,
			result:    expected,
		},
		{
			name:          "generic error on Create()",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.Mocks.ocr2.On("Create", mock.Anything, chaintype.ChainType("evm")).Return(nil, gError)
				f.Mocks.keystore.On("OCR2").Return(f.Mocks.ocr2)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
			},
			query:     mutation,
			variables: variables,
			result:    `null`,
			errors: []*gqlerrors.QueryError{
				{
					Extensions:    nil,
					ResolverError: gError,
					Path:          []interface{}{"createOCR2KeyBundle"},
					Message:       gError.Error(),
				},
			},
		},
	}

	RunGQLTests(t, testCases)
}

func TestResolver_DeleteOCR2KeyBundle(t *testing.T) {
	t.Parallel()

	fakeKey := ocr2key.MustNewInsecure(keystest.NewRandReaderFromSeed(1), "evm")

	mutation := `
		mutation DeleteOCR2KeyBundle($id: ID!) {
			deleteOCR2KeyBundle(id: $id) {
				... on DeleteOCR2KeyBundleSuccess {
					bundle {
						id
						chainType
						configPublicKey
						offChainPublicKey
						onChainPublicKey
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

	ct, err := ToOCR2ChainType(string(fakeKey.ChainType()))
	require.NoError(t, err)

	configPublic := fakeKey.ConfigEncryptionPublicKey()
	pubKey := fakeKey.OffchainPublicKey()
	d, err := json.Marshal(map[string]interface{}{
		"deleteOCR2KeyBundle": map[string]interface{}{
			"bundle": map[string]interface{}{
				"id":                fakeKey.ID(),
				"chainType":         ct,
				"onChainPublicKey":  fmt.Sprintf("ocr2on_%s_%s", fakeKey.ChainType(), fakeKey.OnChainPublicKey()),
				"configPublicKey":   fmt.Sprintf("ocr2cfg_%s_%s", fakeKey.ChainType(), hex.EncodeToString(configPublic[:])),
				"offChainPublicKey": fmt.Sprintf("ocr2off_%s_%s", fakeKey.ChainType(), hex.EncodeToString(pubKey[:])),
			},
		},
	})
	assert.NoError(t, err)
	expected := string(d)

	gError := errors.New("error")

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: mutation, variables: variables}, "deleteOCR2KeyBundle"),
		{
			name:          "success",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.Mocks.ocr2.On("Delete", mock.Anything, fakeKey.ID()).Return(nil)
				f.Mocks.ocr2.On("Get", fakeKey.ID()).Return(fakeKey, nil)
				f.Mocks.keystore.On("OCR2").Return(f.Mocks.ocr2)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
			},
			query:     mutation,
			variables: variables,
			result:    expected,
		},
		{
			name:          "not found error",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.Mocks.ocr2.On("Get", fakeKey.ID()).Return(fakeKey, gError)
				f.Mocks.keystore.On("OCR2").Return(f.Mocks.ocr2)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
			},
			query:     mutation,
			variables: variables,
			result: `{
				"deleteOCR2KeyBundle": {
					"code": "NOT_FOUND",
					"message": "error"
				}
			}`,
		},
		{
			name:          "generic error on Delete()",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.Mocks.ocr2.On("Delete", mock.Anything, fakeKey.ID()).Return(gError)
				f.Mocks.ocr2.On("Get", fakeKey.ID()).Return(fakeKey, nil)
				f.Mocks.keystore.On("OCR2").Return(f.Mocks.ocr2)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
			},
			query:     mutation,
			variables: variables,
			result:    `null`,
			errors: []*gqlerrors.QueryError{
				{
					Extensions:    nil,
					ResolverError: gError,
					Path:          []interface{}{"deleteOCR2KeyBundle"},
					Message:       gError.Error(),
				},
			},
		},
	}

	RunGQLTests(t, testCases)
}
