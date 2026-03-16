package attestation

import (
	"crypto/sha256"
	"testing"
)

func TestDomainHash(t *testing.T) {
	tag := "TestTag"
	data := []byte(`{"key":"value"}`)

	got := DomainHash(tag, data)

	// Manually compute the expected hash.
	h := sha256.New()
	h.Write([]byte(DomainSeparator))
	h.Write([]byte("\n" + tag + "\n"))
	h.Write(data)
	want := h.Sum(nil)

	if len(got) != sha256.Size {
		t.Fatalf("expected %d bytes, got %d", sha256.Size, len(got))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("hash mismatch at byte %d: want %x, got %x", i, want, got)
		}
	}
}

func TestDomainHash_DifferentTags(t *testing.T) {
	data := []byte("same-data")
	h1 := DomainHash("Tag1", data)
	h2 := DomainHash("Tag2", data)

	for i := range h1 {
		if h1[i] != h2[i] {
			return // different, as expected
		}
	}
	t.Fatal("different tags should produce different hashes")
}

func TestDomainHash_DifferentData(t *testing.T) {
	tag := "SameTag"
	h1 := DomainHash(tag, []byte("data-a"))
	h2 := DomainHash(tag, []byte("data-b"))

	for i := range h1 {
		if h1[i] != h2[i] {
			return // different, as expected
		}
	}
	t.Fatal("different data should produce different hashes")
}
