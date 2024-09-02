package standardcapabilities_test

import (
	"testing"

	"github.com/test-go/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/standardcapabilities"
)

func Test_ValidatedStandardCapabilitiesSpec(t *testing.T) {

	type testCase struct {
		name          string
		expectedError string
	}

	testCases := []testCase{
		{
			name:          "invalid TOML string",
			expectedError: "toml error on load standard capabilities",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := standardcapabilities.ValidatedStandardCapabilitiesSpec("invalid TOML string")

			if tc.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err, tc.expectedError)
			}
		})
	}

}
