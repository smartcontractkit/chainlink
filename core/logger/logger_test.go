package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoSampling(t *testing.T) {
	assert.Nil(t, newBaseConfig().Sampling)
	assert.Nil(t, newTestConfig().Sampling)
	assert.Nil(t, newProductionConfig("", false, true, false).Sampling)
}
