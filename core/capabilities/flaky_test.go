package capabilities

import (
	"math/rand"
	"testing"

	"github.com/smartcontractkit/flakeguard"
)

func TestFlakyOriginallyTenPercent(t *testing.T) {
	t.Parallel()

	flakeguard.Quarantine(t, "flaky originally ten percent")

	if rand.Intn(10) == 0 {
		t.Log("I flake 10% of the time")
		t.FailNow()
	}
}

func TestFlakyOriginallyTwentyFivePercent(t *testing.T) {
	t.Parallel()
	flakeguard.Quarantine(t, "flaky originally twenty five percent")

	if rand.Intn(4) == 0 {
		t.Log("I flake 25% of the time")
		t.FailNow()
	}
}

func TestFlakyOriginallyFiftyPercent(t *testing.T) {
	t.Parallel()
	flakeguard.Quarantine(t, "flaky originally fifty percent")

	if rand.Intn(2) == 0 {
		t.Log("I flake 50% of the time")
		t.FailNow()
	}
}

func TestFlakyOriginallySeventyFivePercent(t *testing.T) {
	t.Parallel()
	flakeguard.Quarantine(t, "flaky originally seventy five percent")
	if rand.Intn(4) != 0 {
		t.Log("I flake 75% of the time")
		t.FailNow()
	}
}

func TestPassOriginally(t *testing.T) {
	t.Parallel()

	flakeguard.Quarantine(t, "pass originally")
	t.Log("I pass")
}
