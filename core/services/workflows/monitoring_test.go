package workflows

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// tests CustomMessageAgent does not share state across new instances created by `with`
func Test_CustomMessageAgent(t *testing.T) {
	cma := NewCustomMessageAgent()
	cma1 := cma.with("key1", "value1")
	cma2 := cma1.with("key2", "value2")

	assert.NotEqual(t, cma1.labels, cma2.labels)
}
