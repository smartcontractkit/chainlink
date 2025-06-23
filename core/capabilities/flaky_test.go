package capabilities

import (
	"testing"
)

func TestFlakyOriginallyTenPercent(t *testing.T) {
	t.Parallel()
	t.Log("I'm no longer flaky")
}

func TestFlakyOriginallyTwentyFivePercent(t *testing.T) {
	t.Parallel()

	t.Log("I'm no longer flaky")
}

func TestFlakyOriginallyFiftyPercent(t *testing.T) {
	t.Parallel()
	t.Log("I'm no longer flaky")
}

func TestFlakyOriginallySeventyFivePercent(t *testing.T) {
	t.Parallel()
	t.Log("I'm no longer flaky")
}

func TestFlakyOriginallyNinetyPercent(t *testing.T) {
	t.Parallel()
	t.Log("I'm no longer flaky")
}
