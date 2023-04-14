package gateway_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/gateway"
)

func TestAPI_Encode(t *testing.T) {
	t.Parallel()

	input := []byte(`{"don_id": "functions_local", "payload": {"request_id": 123, "source": "int main"}}`)
	output, err := gateway.Decode(input)
	require.NoError(t, err)
	require.Equal(t, "functions_local", output.DonId)
	require.NotNil(t, output.Payload)
}
