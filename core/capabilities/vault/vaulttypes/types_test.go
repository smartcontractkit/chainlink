package vaulttypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeNamespace(t *testing.T) {
	assert.Equal(t, DefaultNamespace, NormalizeNamespace(""))
	assert.Equal(t, "custom", NormalizeNamespace("custom"))
}
