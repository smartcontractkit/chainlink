package capabilities

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResponseValidate(t *testing.T) {
	tt := []struct {
		name        string
		response    Response
		expectError string
	}{
		{
			name:     "valid Response with ExecutionError",
			response: Response{ExecutionError: true, ErrorMessage: "Some error"},
		},
		{
			name:        "invalid Response with ExecutionError but no ErrorMessage",
			response:    Response{ExecutionError: true},
			expectError: "executionError is true but errorMessage is empty",
		},
		{
			name:     "valid HTTP Response",
			response: Response{StatusCode: 200},
		},
		{
			name: "invalid status code",
			response: Response{
				Body: []byte("body"),
			},
			expectError: "statusCode must be set when executionError is false",
		},
		{
			name:        "invalid HTTP Response with bad StatusCode",
			response:    Response{StatusCode: 700},
			expectError: "statusCode must be a valid HTTP status code (100-599)",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.response.Validate()

			if tc.expectError != "" {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
