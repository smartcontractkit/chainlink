package monitoring

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// tests CustomMessageAgent does not share state across new instances created by `With`
func Test_CustomMessageAgent(t *testing.T) {
	cma := NewCustomMessageLabeler()
	cma1 := cma.With("key1", "value1")
	cma2 := cma1.With("key2", "value2")

	assert.NotEqual(t, cma1.labels, cma2.labels)
}
