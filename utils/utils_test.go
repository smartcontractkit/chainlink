package utils_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/utils"
	"github.com/stretchr/testify/assert"
)

func TestNewBytes32ID(t *testing.T) {
	t.Parallel()

	id := utils.NewBytes32ID()
	assert.NotContains(t, id, "-")
}
