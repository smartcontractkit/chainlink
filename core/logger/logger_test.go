package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	// no sampling
	assert.Nil(t, newBaseConfig().Sampling)
	assert.Nil(t, newTestConfig().Sampling)
	assert.Nil(t, newProductionConfig("", false, true, false, false).Sampling)

	// not development, which would trigger panics for Critical level
	assert.False(t, newBaseConfig().Development)
	assert.False(t, newTestConfig().Development)
	assert.False(t, newProductionConfig("", false, true, false, false).Development)
}
