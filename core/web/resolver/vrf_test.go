package resolver

import (
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/vrfkey"
)

func TestResolver_GetVRFKey(t *testing.T) {
	t.Parallel()

	query := `
		query GetVRFKey($id: ID!) {
			vrfKey(id: $id) {
				... on VRFKeySuccess {
					key {
						id
						compressed
						uncompressed
						hash
					}
				}

				... on NotFoundError {
					message
					code
				}
			}
		}
	`

	fakeKey := vrfkey.MustNewV2XXXTestingOnly(big.NewInt(1))

	uncompressed, err := fakeKey.PublicKey.StringUncompressed()
	assert.NoError(t, err)

	d, err := json.Marshal(map[string]interface{}{
		"vrfKey": map[string]interface{}{
			"key": map[string]interface{}{
				"id":           fakeKey.PublicKey.String(),
				"compressed":   fakeKey.PublicKey.String(),
				"uncompressed": uncompressed,
				"hash":         fakeKey.PublicKey.MustHash().String(),
			},
		},
	})
	assert.NoError(t, err)
	expected := string(d)

	variables := map[string]interface{}{
		"id": fakeKey.PublicKey.String(),
	}

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query, variables: variables}, "vrfKey"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.vrf.On("Get", fakeKey.PublicKey.String()).Return(fakeKey, nil)
				f.Mocks.keystore.On("VRF").Return(f.Mocks.vrf)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
			},
			query:     query,
			variables: variables,
			result:    expected,
		},
		{
			name:          "not found error",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.vrf.
					On("Get", fakeKey.PublicKey.String()).
					Return(vrfkey.KeyV2{}, errors.Wrapf(
						keystore.ErrMissingVRFKey,
						"unable to find VRF key with id %s",
						fakeKey.PublicKey.String(),
					))
				f.Mocks.keystore.On("VRF").Return(f.Mocks.vrf)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
			},
			query:     query,
			variables: variables,
			result: fmt.Sprintf(`{
				"vrfKey": {
					"code": "NOT_FOUND",
					"message": "unable to find VRF key with id %s: unable to find VRF key"
				}
			}`, fakeKey.PublicKey.String()),
		},
	}

	RunGQLTests(t, testCases)
}

func TestResolver_GetVRFKeys(t *testing.T) {
	t.Parallel()

	query := `
		query GetVRFKeys {
			vrfKeys {
				results {
					id
					compressed
					uncompressed
					hash
				}
			}
		}
	`

	fakeKeys := []vrfkey.KeyV2{}
	expectedKeys := []map[string]string{}
	for i := 0; i < 2; i++ {
		fakeKey := vrfkey.MustNewV2XXXTestingOnly(big.NewInt(1))
		uncompressed, err := fakeKey.PublicKey.StringUncompressed()
		assert.NoError(t, err)

		fakeKeys = append(fakeKeys, fakeKey)
		expectedKeys = append(expectedKeys, map[string]string{
			"id":           fakeKey.PublicKey.String(),
			"compressed":   fakeKey.PublicKey.String(),
			"uncompressed": uncompressed,
			"hash":         fakeKey.PublicKey.MustHash().String(),
		})
	}

	d, err := json.Marshal(map[string]interface{}{
		"vrfKeys": map[string]interface{}{
			"results": expectedKeys,
		},
	})
	assert.NoError(t, err)
	expected := string(d)

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "vrfKeys"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.vrf.On("GetAll").Return(fakeKeys, nil)
				f.Mocks.keystore.On("VRF").Return(f.Mocks.vrf)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
			},
			query:  query,
			result: expected,
		},
	}

	RunGQLTests(t, testCases)
}

func TestResolver_CreateVRFKey(t *testing.T) {
	t.Parallel()

	mutation := `
		mutation CreateVRFKey {
			createVRFKey {
				key {
					id
					compressed
					uncompressed
					hash
				}
			}
		}
	`

	fakeKey := vrfkey.MustNewV2XXXTestingOnly(big.NewInt(1))

	uncompressed, err := fakeKey.PublicKey.StringUncompressed()
	assert.NoError(t, err)

	d, err := json.Marshal(map[string]interface{}{
		"createVRFKey": map[string]interface{}{
			"key": map[string]interface{}{
				"id":           fakeKey.PublicKey.String(),
				"compressed":   fakeKey.PublicKey.String(),
				"uncompressed": uncompressed,
				"hash":         fakeKey.PublicKey.MustHash().String(),
			},
		},
	})
	assert.NoError(t, err)
	expected := string(d)

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: mutation}, "createVRFKey"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.vrf.On("Create").Return(fakeKey, nil)
				f.Mocks.keystore.On("VRF").Return(f.Mocks.vrf)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
			},
			query:  mutation,
			result: expected,
		},
	}

	RunGQLTests(t, testCases)
}

func TestResolver_DeleteVRFKey(t *testing.T) {
	t.Parallel()

	mutation := `
		mutation DeleteVRFKey($id: ID!) {
			deleteVRFKey(id: $id) {
				... on DeleteVRFKeySuccess {
					key {
						id
						compressed
						uncompressed
						hash
					}
				}

				... on NotFoundError {
					message
					code
				}
			}
		}
	`

	fakeKey := vrfkey.MustNewV2XXXTestingOnly(big.NewInt(1))

	uncompressed, err := fakeKey.PublicKey.StringUncompressed()
	assert.NoError(t, err)

	d, err := json.Marshal(map[string]interface{}{
		"deleteVRFKey": map[string]interface{}{
			"key": map[string]interface{}{
				"id":           fakeKey.PublicKey.String(),
				"compressed":   fakeKey.PublicKey.String(),
				"uncompressed": uncompressed,
				"hash":         fakeKey.PublicKey.MustHash().String(),
			},
		},
	})
	assert.NoError(t, err)
	expected := string(d)

	variables := map[string]interface{}{
		"id": fakeKey.PublicKey.String(),
	}

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: mutation, variables: variables}, "deleteVRFKey"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.vrf.On("Delete", fakeKey.PublicKey.String()).Return(fakeKey, nil)
				f.Mocks.keystore.On("VRF").Return(f.Mocks.vrf)
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
				f.Mocks.vrf.
					On("Delete", fakeKey.PublicKey.String()).
					Return(vrfkey.KeyV2{}, errors.Wrapf(
						keystore.ErrMissingVRFKey,
						"unable to find VRF key with id %s",
						fakeKey.PublicKey.String(),
					))
				f.Mocks.keystore.On("VRF").Return(f.Mocks.vrf)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
			},
			query:     mutation,
			variables: variables,
			result: fmt.Sprintf(`{
				"deleteVRFKey": {
					"message": "unable to find VRF key with id %s: unable to find VRF key",
					"code": "NOT_FOUND"
				}
			}`, fakeKey.PublicKey.String()),
		},
	}

	RunGQLTests(t, testCases)
}
