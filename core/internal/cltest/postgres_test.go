package cltest

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractDB(t *testing.T) {
	tests := []struct {
		url, expectation string
	}{
		{"postgres://localhost:5432/chainlink_test?sslmode=disable", "chainlink_test"},
		{"postgres://circleci_postgres@localhost:5432/circleci_test?sslmode=disable", "circleci_test"},
		{"postgres://localhost:5432/random_db", "random_db"},
		{"postgres://username:password@localhost:5432/authenticated_db", "authenticated_db"},
	}

	for _, test := range tests {
		t.Run(test.url, func(t *testing.T) {
			require.Equal(t, test.expectation, extractDB(t, test.url))
		})
	}
}
