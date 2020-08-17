package web_test

import (
	"errors"
	"testing"

	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/web"

	"github.com/stretchr/testify/assert"
)

func TestHelpers_StatusCodeForError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		err        error
		statusCode int
	}{
		{"ValidationError", models.NewValidationError("test"), 400},
		{"DatabaseAccessError", models.NewDatabaseAccessError("test"), 500},
		{"DefaultError", errors.New("test"), 500},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.statusCode, web.StatusCodeForError(test.err))
		})
	}
}
