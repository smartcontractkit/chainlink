package tunnel

import "testing"

func TestCanonicalComponentID(t *testing.T) {
	if got := CanonicalComponentID(KindBlockchain, 0, "Anvil-Main"); got != "blockchain:0:anvil-main" {
		t.Fatalf("unexpected canonical id: %s", got)
	}

	if got := CanonicalComponentID(KindJD, 2, ""); got != "jd:2" {
		t.Fatalf("unexpected canonical id for empty name: %s", got)
	}

	if got := CanonicalComponentID(KindNodeSet, 1, "   "); got != "nodeset:1" {
		t.Fatalf("unexpected canonical id for whitespace name: %s", got)
	}
}
