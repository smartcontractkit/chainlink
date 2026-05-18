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

func TestStderrWriter(t *testing.T) {
	sw := stderrWriter{}

	// Test Write
	n, err := sw.Write([]byte("Hello, World!"))
	assert.NoError(t, err)
	assert.Equal(t, 13, n, "Expected 13 bytes written")

	// Test Close
	err = sw.Close()
	assert.NoError(t, err)
}
