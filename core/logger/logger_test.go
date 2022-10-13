package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	// no sampling
	assert.Nil(t, newZapConfigBase().Sampling)
	assert.Nil(t, newZapConfigProd(false, false).Sampling)

	// not development, which would trigger panics for Critical level
	assert.False(t, newZapConfigBase().Development)
	assert.False(t, newZapConfigProd(false, false).Development)
}
