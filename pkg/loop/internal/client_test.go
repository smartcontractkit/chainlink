package internal

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

func TestWrappedErrorError(t *testing.T) {
	t.Parallel()
	t.Run("Is returns false for different error code", func(t *testing.T) {
		// silly, but to verify that it's only looking at error code here, we need to make the message the same
		err := types.NotFoundError(types.ErrInvalidType.Error())
		assert.False(t, errors.Is(wrapRPCErr(types.ErrInvalidType), err))
	})

	t.Run("Is returns false for different message", func(t *testing.T) {
		// Both are InvalidArgumentError
		assert.False(t, errors.Is(wrapRPCErr(types.ErrInvalidType), types.ErrInvalidEncoding))
	})

	t.Run("Is returns true if the message and code are the same", func(t *testing.T) {
		assert.True(t, errors.Is(wrapRPCErr(types.ErrInvalidType), types.ErrInvalidType))
	})

	t.Run("Is returns true if the message is contained and the code is the same", func(t *testing.T) {
		wrapped := wrapRPCErr(fmt.Errorf("%w: %w", types.ErrInvalidType, errors.New("some other error")))
		assert.True(t, errors.Is(wrapped, types.ErrInvalidType))
	})
}
