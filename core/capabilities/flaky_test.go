package capabilities

import (
	"math/rand"
	"testing"
)

func TestFlakyOriginallyTenPercent(t *testing.T) {
	t.Parallel()

	if rand.Intn(10) == 0 {
		t.Log("I flake 10% of the time")
		t.FailNow()
	}
}

func TestFlakyOriginallyTwentyFivePercent(t *testing.T) {
	t.Parallel()

	if rand.Intn(4) == 0 {
		t.Log("I flake 25% of the time")
		t.FailNow()
	}
}

func TestFlakyOriginallyFiftyPercent(t *testing.T) {
	t.Parallel()

	if rand.Intn(2) == 0 {
		t.Log("I flake 50% of the time")
		t.FailNow()
	}
}

func TestFlakyOriginallySeventyFivePercent(t *testing.T) {
	t.Parallel()

	if rand.Intn(4) != 0 {
		t.Log("I flake 75% of the time")
		t.FailNow()
	}
}

func TestPassOriginally(t *testing.T) {
	t.Parallel()

	t.Log("I pass")
}
