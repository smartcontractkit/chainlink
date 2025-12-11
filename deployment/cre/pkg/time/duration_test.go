package time_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"

	customtime "github.com/smartcontractkit/chainlink/deployment/cre/pkg/time"
)

func TestDuration_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		input       string
		expectedDur time.Duration
		expectErr   bool
	}{
		{
			name:        "unmarshal seconds",
			input:       `"1s"`,
			expectedDur: time.Second,
		},
		{
			name:        "unmarshal minutes",
			input:       `"1m"`,
			expectedDur: time.Minute,
		},
		{
			name:        "unmarshal hours",
			input:       `"1h"`,
			expectedDur: time.Hour,
		},
		{
			name:        "unmarshal days",
			input:       `"1d"`,
			expectedDur: 24 * time.Hour,
		},
		{
			name:        "unmarshal multiple days",
			input:       `"21d"`,
			expectedDur: 21 * 24 * time.Hour,
		},
		{
			name:        "unmarshal float days",
			input:       `"1.5d"`,
			expectedDur: 36 * time.Hour,
		},
		{
			name:      "unmarshal empty string",
			input:     `""`,
			expectErr: true,
		},
		{
			name:      "unmarshal invalid format",
			input:     `"1x"`,
			expectErr: true,
		},
		{
			name:      "unmarshal invalid day value",
			input:     `"xd"`,
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var d customtime.Duration
			err := json.Unmarshal([]byte(tc.input), &d)

			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedDur, commonconfig.Duration(d).Duration())
			}
		})
	}
}

func TestDuration_MarshalJSON(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		duration     customtime.Duration
		expectedJSON string
	}{
		{
			name:         "marshal seconds",
			duration:     customtime.Duration(*commonconfig.MustNewDuration(time.Second)),
			expectedJSON: `"1s"`,
		},
		{
			name:         "marshal minutes",
			duration:     customtime.Duration(*commonconfig.MustNewDuration(time.Minute)),
			expectedJSON: `"1m0s"`,
		},
		{
			name:         "marshal hours",
			duration:     customtime.Duration(*commonconfig.MustNewDuration(time.Hour)),
			expectedJSON: `"1h0m0s"`,
		},
		{
			name:         "marshal single day",
			duration:     customtime.Duration(*commonconfig.MustNewDuration(24 * time.Hour)),
			expectedJSON: `"1d"`,
		},
		{
			name:         "marshal multiple days",
			duration:     customtime.Duration(*commonconfig.MustNewDuration(21 * 24 * time.Hour)),
			expectedJSON: `"21d"`,
		},
		{
			name:         "marshal non-exact day",
			duration:     customtime.Duration(*commonconfig.MustNewDuration(25 * time.Hour)),
			expectedJSON: `"25h0m0s"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b, err := json.Marshal(tc.duration)
			require.NoError(t, err)
			require.JSONEq(t, tc.expectedJSON, string(b))
		})
	}
}
