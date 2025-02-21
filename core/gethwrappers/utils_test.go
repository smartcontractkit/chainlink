package gethwrappers_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers"

	"github.com/stretchr/testify/assert"
)

func TestBoxOutput(t *testing.T) {
	t.Parallel()

	output := gethwrappers.BoxOutput("some error %d %s", 123, "foo")
	const expected = "\n" +
		"↘↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↙\n" +
		"→                      ←\n" +
		"→  README README       ←\n" +
		"→                      ←\n" +
		"→  some error 123 foo  ←\n" +
		"→                      ←\n" +
		"→  README README       ←\n" +
		"→                      ←\n" +
		"↗↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↖\n" +
		"\n"
	assert.Equal(t, expected, output)
}
