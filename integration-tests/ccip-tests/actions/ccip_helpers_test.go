package actions

import (
	"errors"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testreporters"
)

func TestIsPhaseValid(t *testing.T) {
	// isPhaseValid has some complex logic that could lead to false negatives
	t.Parallel()
	logger := zerolog.New(zerolog.Nop())

	testCases := []struct {
		name         string
		currentPhase testreporters.Phase
		opts         validationOptions
		phaseErr     error

		expectedShouldReturn bool
		expectedErr          error
	}{
		{
			name:         "should return error immediately if phase error is present and no phase is expected to fail",
			currentPhase: testreporters.CCIPSendRe,
			opts:         validationOptions{},
			phaseErr:     errors.New("some error"),

			expectedShouldReturn: true,
			expectedErr:          errors.New("some error"),
		},
		{
			name:         "should return with no error if phase is expected to fail and phase error present",
			currentPhase: testreporters.CCIPSendRe,
			opts: validationOptions{
				phaseExpectedToFail: testreporters.CCIPSendRe,
			},
			phaseErr: errors.New("some error"),

			expectedShouldReturn: true,
			expectedErr:          nil,
		},
		{
			name:         "should return with error if phase is expected to fail and no phase error present",
			currentPhase: testreporters.CCIPSendRe,
			opts: validationOptions{
				phaseExpectedToFail: testreporters.CCIPSendRe,
			},
			phaseErr: nil,

			expectedShouldReturn: true,
			expectedErr:          errors.New("expected phase 'CCIPSendRequested' to fail, but it passed"),
		},
		{
			name:         "should not return if phase is not expected to fail and no phase error present",
			currentPhase: testreporters.CCIPSendRe,
			opts: validationOptions{
				phaseExpectedToFail: testreporters.ExecStateChanged,
			},
			phaseErr: nil,

			expectedShouldReturn: false,
			expectedErr:          nil,
		},
		{
			name:         "should return with no error if phase is expected to fail with specific error message and that error message is present",
			currentPhase: testreporters.CCIPSendRe,
			opts: validationOptions{
				phaseExpectedToFail:  testreporters.CCIPSendRe,
				expectedErrorMessage: "some error",
			},
			phaseErr: errors.New("some error"),

			expectedShouldReturn: true,
			expectedErr:          nil,
		},
		{
			name:         "should return with error if phase is expected to fail with specific error message and that error message is not present",
			currentPhase: testreporters.CCIPSendRe,
			opts: validationOptions{
				phaseExpectedToFail:  testreporters.CCIPSendRe,
				expectedErrorMessage: "some error",
			},
			phaseErr: errors.New("some other error"),

			expectedShouldReturn: true,
			expectedErr:          errors.New("expected phase 'CCIPSendRequested' to fail with error message 'some error' but got error 'some other error'"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			shouldReturn, err := isPhaseValid(&logger, tc.currentPhase, tc.opts, tc.phaseErr)
			require.Equal(t, tc.expectedShouldReturn, shouldReturn, "shouldReturn not as expected")
			require.Equal(t, tc.expectedErr, err, "err not as expected")
		})
	}
}
