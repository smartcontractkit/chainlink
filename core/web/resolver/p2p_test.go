package resolver

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

func TestResolver_GetP2PKeys(t *testing.T) {
	t.Parallel()

	query := `
		query GetP2PKeys {
			p2pKeys {
				results {
					id
					peerID
					publicKey
				}
			}
		}
	`

	fakeKeys := []p2pkey.KeyV2{}
	expectedKeys := []map[string]string{}
	for i := 0; i < 2; i++ {
		k := p2pkey.MustNewV2XXXTestingOnly(big.NewInt(1))
		fakeKeys = append(fakeKeys, k)
		expectedKeys = append(expectedKeys, map[string]string{
			"id":        k.ID(),
			"peerID":    k.PeerID().String(),
			"publicKey": k.PublicKeyHex(),
		})
	}

	d, err := json.Marshal(map[string]interface{}{
		"p2pKeys": map[string]interface{}{
			"results": expectedKeys,
		},
	})
	assert.NoError(t, err)
	expected := string(d)

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "p2pKeys"),
		{
			name:          "success",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.Mocks.p2p.On("GetAll").Return(fakeKeys, nil)
				f.Mocks.keystore.On("P2P").Return(f.Mocks.p2p)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
			},
			query:  query,
			result: expected,
		},
	}

	RunGQLTests(t, testCases)
}

func TestResolver_CreateP2PKey(t *testing.T) {
	t.Parallel()

	query := `
		mutation CreateP2PKey {
			createP2PKey {
				... on CreateP2PKeySuccess {
					p2pKey {
						id
						peerID
						publicKey
					}
				}
			}
		}
	`

	fakeKey := p2pkey.MustNewV2XXXTestingOnly(big.NewInt(1))

	d, err := json.Marshal(map[string]interface{}{
		"createP2PKey": map[string]interface{}{
			"p2pKey": map[string]interface{}{
				"id":        fakeKey.ID(),
				"peerID":    fakeKey.PeerID().String(),
				"publicKey": fakeKey.PublicKeyHex(),
			},
		},
	})
	assert.NoError(t, err)
	expected := string(d)

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "createP2PKey"),
		{
			name:          "success",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.Mocks.p2p.On("Create", mock.Anything).Return(fakeKey, nil)
				f.Mocks.keystore.On("P2P").Return(f.Mocks.p2p)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
			},
			query:  query,
			result: expected,
		},
	}

	RunGQLTests(t, testCases)
}

func TestResolver_DeleteP2PKey(t *testing.T) {
	t.Parallel()

	fakeKey := p2pkey.MustNewV2XXXTestingOnly(big.NewInt(1))

	query := `
		mutation DeleteP2PKey($id: ID!) {
			deleteP2PKey(id: $id) {
				... on DeleteP2PKeySuccess {
					p2pKey {
						id
						peerID
						publicKey
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

	peerID, err := p2pkey.MakePeerID(fakeKey.ID())
	assert.NoError(t, err)

	d, err := json.Marshal(map[string]interface{}{
		"deleteP2PKey": map[string]interface{}{
			"p2pKey": map[string]interface{}{
				"id":        fakeKey.ID(),
				"peerID":    fakeKey.PeerID().String(),
				"publicKey": fakeKey.PublicKeyHex(),
			},
		},
	})
	assert.NoError(t, err)
	expected := string(d)

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query, variables: variables}, "deleteP2PKey"),
		{
			name:          "success",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.Mocks.p2p.On("Delete", mock.Anything, peerID).Return(fakeKey, nil)
				f.Mocks.keystore.On("P2P").Return(f.Mocks.p2p)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
			},
			query:     query,
			variables: variables,
			result:    expected,
		},
		{
			name:          "not found error",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.Mocks.p2p.
					On("Delete", mock.Anything, peerID).
					Return(
						p2pkey.KeyV2{},
						keystore.KeyNotFoundError{ID: peerID.String(), KeyType: "P2P"},
					)
				f.Mocks.keystore.On("P2P").Return(f.Mocks.p2p)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
			},
			query:     query,
			variables: variables,
			result: fmt.Sprintf(`{
				"deleteP2PKey": {
					"code":"NOT_FOUND",
					"message":"unable to find P2P key with id %s"
				}
			}`, peerID.String()),
		},
	}

	RunGQLTests(t, testCases)
}
