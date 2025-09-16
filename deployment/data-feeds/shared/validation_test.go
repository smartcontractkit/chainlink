package shared

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateFeedID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		feedID string
		valid  bool
	}{
		{"0x011e22d6bf0003320000000000000000", true},  // valid feed ID
		{"0x026d06ebb60700020000000000000000", true},  // valid feed ID
		{"0x126d06ebb60700020000000000000000", false}, // invalid feed format
		{"0x126d06ebb6000020000000000000000", false},  // invalid attribute bucket
		{"126d06ebb6000020000000000000000", false},    // no 0x prefix
		{"0x026d06ebb6070002000000000000", false},     // invalid length
		{"0x011e22d6bf0003870000000000000000", false}, // invalid data type
		{"0x01zz22d6bf0003320000000000000000", false}, // invalid random bytes
	}

	for _, test := range tests {
		err := ValidateFeedID(test.feedID)
		if test.valid {
			require.NoError(t, err, "expected feedID %s to be valid", test.feedID)
		} else {
			require.Error(t, err, "expected feedID %s to be invalid", test.feedID)
		}
	}
}
