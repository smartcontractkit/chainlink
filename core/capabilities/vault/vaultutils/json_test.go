package vaultutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
)

func TestToCanonicalJSON(t *testing.T) {
	t.Parallel()

	testData, err := structpb.NewValue(map[string]any{
		"field1": "value1",
		"field2": 42,
	})
	require.NoError(t, err)

	canonicalJSON, err := ToCanonicalJSON(testData, false)
	require.NoError(t, err)

	expectedJSON := `{"field1":"value1","field2":42}`
	assert.Equal(t, expectedJSON, string(canonicalJSON)) //nolint:testifylint // testifylint requires use of assert.JSONEq which doesn't do a string match of the JSON, but does a structural match.

	// same data, different order
	testData, err = structpb.NewValue(map[string]any{
		"field2": 42,
		"field1": "value1",
	})
	require.NoError(t, err)

	canonicalJSON, err = ToCanonicalJSON(testData, false)
	require.NoError(t, err)

	assert.Equal(t, expectedJSON, string(canonicalJSON)) //nolint:testifylint // testifylint requires use of assert.JSONEq which doesn't do a string match of the JSON, but does a structural match.
}

func TestToCanonicalJSON_OmitsEmptyFields(t *testing.T) {
	t.Parallel()

	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret1",
	}

	t.Run("nested empty string and empty repeated field", func(t *testing.T) {
		t.Parallel()

		msg := &vaultcommon.CreateSecretsResponse{
			RequestId: "owner::req-1",
			Responses: []*vaultcommon.CreateSecretResponse{
				{
					Id:      id,
					Success: true,
					Error:   "",
				},
			},
		}

		withEmptyFields, err := ToCanonicalJSON(msg, false)
		require.NoError(t, err)
		assert.JSONEq(t, `{
			"requestId":"owner::req-1",
			"responses":[{"id":{"owner":"owner","namespace":"main","key":"secret1"},"success":true,"error":""}]
		}`, string(withEmptyFields))

		canonicalJSON, err := ToCanonicalJSON(msg, true)
		require.NoError(t, err)
		assert.JSONEq(t, `{
			"requestId":"owner::req-1",
			"responses":[{"id":{"owner":"owner","namespace":"main","key":"secret1"},"success":true}]
		}`, string(canonicalJSON))
		assert.NotEqual(t, string(withEmptyFields), string(canonicalJSON))
	})

	t.Run("empty repeated field and empty top-level string", func(t *testing.T) {
		t.Parallel()

		msg := &vaultcommon.CreateSecretsResponse{
			Responses: []*vaultcommon.CreateSecretResponse{},
		}

		withEmptyFields, err := ToCanonicalJSON(msg, false)
		require.NoError(t, err)
		assert.JSONEq(t, `{"responses":[],"requestId":""}`, string(withEmptyFields))

		canonicalJSON, err := ToCanonicalJSON(msg, true)
		require.NoError(t, err)
		assert.JSONEq(t, `{}`, string(canonicalJSON))
		assert.NotEqual(t, string(withEmptyFields), string(canonicalJSON))
	})

	t.Run("empty repeated field only", func(t *testing.T) {
		t.Parallel()

		msg := &vaultcommon.CreateSecretsResponse{
			RequestId: "owner::req-1",
			Responses: []*vaultcommon.CreateSecretResponse{},
		}

		withEmptyFields, err := ToCanonicalJSON(msg, false)
		require.NoError(t, err)
		assert.JSONEq(t, `{"responses":[],"requestId":"owner::req-1"}`, string(withEmptyFields))

		canonicalJSON, err := ToCanonicalJSON(msg, true)
		require.NoError(t, err)
		assert.JSONEq(t, `{"requestId":"owner::req-1"}`, string(canonicalJSON))
		assert.NotEqual(t, string(withEmptyFields), string(canonicalJSON))
	})
}
