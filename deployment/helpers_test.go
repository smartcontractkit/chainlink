package deployment

import (
	"testing"

	"github.com/stretchr/testify/require"

	forwarder "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/forwarder_1_0_0"
)

func TestParseErrorFromABI(t *testing.T) {
	testCases := []struct {
		Name               string
		RevertReason       string
		ABI                string
		ParsedRevertReason string
	}{
		{
			Name:               "Generic error with string msg",
			RevertReason:       "0x08c379a0000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000164f6e6c792063616c6c61626c65206279206f776e657200000000000000000000",
			ABI:                "", // ABI is not required for this case
			ParsedRevertReason: "error - `Only callable by owner`",
		},
		{
			Name:               "Custom typed error",
			RevertReason:       "0xdf3b81ea0000000000000000000000000000000000000000000000000000000100000001",
			ABI:                forwarder.KeystoneForwarderABI,
			ParsedRevertReason: "error -`InvalidConfig` args [4294967297]",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			revertReason, err := parseErrorFromABI(tc.RevertReason, tc.ABI)
			require.NoError(t, err)
			require.Equal(t, tc.ParsedRevertReason, revertReason)
		})
	}
}
