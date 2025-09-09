package contracts

import "testing"

func TestDonsOrderedByID(t *testing.T) {
	// Test donsOrderedByID sorts by id ascending
	d := dons{
		c: make(map[string]donConfig),
	}

	d.c["don3"] = donConfig{id: 3}
	d.c["don1"] = donConfig{id: 1}
	d.c["don2"] = donConfig{id: 2}

	ordered := d.donsOrderedByID()
	if len(ordered) != 3 {
		t.Fatalf("expected 3 dons, got %d", len(ordered))
	}

	if ordered[0].id != 1 || ordered[1].id != 2 || ordered[2].id != 3 {
		t.Fatalf("expected dons ordered by id 1,2,3 got %d,%d,%d", ordered[0].id, ordered[1].id, ordered[2].id)
	}
}
