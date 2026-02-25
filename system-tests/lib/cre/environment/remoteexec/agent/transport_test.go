package agent

import "testing"

type testNested struct {
	Value string `toml:"value"`
}

type testRuntimeStruct struct {
	Name   string      `toml:"name"`
	Nested *testNested `toml:"nested"`
	SkipMe string      `toml:"-"`
}

func TestTransportRoundtripDropsTomlIgnoredFields(t *testing.T) {
	input := &testRuntimeStruct{
		Name:   "abc",
		Nested: &testNested{Value: "x"},
		SkipMe: "should-not-travel",
	}

	encoded, err := EncodeForTransport(input)
	if err != nil {
		t.Fatalf("expected no error encoding transport payload, got %v", err)
	}

	decoded, err := DecodeFromTransport[testRuntimeStruct](encoded)
	if err != nil {
		t.Fatalf("expected no error decoding transport payload, got %v", err)
	}

	if decoded.Name != "abc" {
		t.Fatalf("expected name to roundtrip, got %q", decoded.Name)
	}
	if decoded.Nested == nil || decoded.Nested.Value != "x" {
		t.Fatalf("expected nested value to roundtrip, got %#v", decoded.Nested)
	}
	if decoded.SkipMe != "" {
		t.Fatalf("expected toml-ignored field to be dropped, got %q", decoded.SkipMe)
	}
}
