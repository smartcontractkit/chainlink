package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
)

func Test_Callback(t *testing.T) {
	cb := NewCallback()
	payload := handlers.UserCallbackPayload{RawResponse: []byte("test")}

	err := cb.SendResponse(payload)
	require.NoError(t, err)

	err = cb.SendResponse(payload)
	require.ErrorContains(t, err, "response already sent")

	resp, err := cb.Wait(t.Context())
	require.NoError(t, err)

	assert.Equal(t, payload, resp)

	_, err = cb.Wait(t.Context())
	require.ErrorContains(t, err, "Wait can only be called once per Callback instance")
}
