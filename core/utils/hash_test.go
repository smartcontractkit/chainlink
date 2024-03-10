package utils

import (
	"encoding/json"
	"strings"
	"testing"
)

func Test_Hash_UnmarshalText(t *testing.T) {
	var tests = []struct {
		Prefix string
		Size   int
		Error  string
	}{
		{"", 62, "hash: expected a hex string starting with '0x'"},
		{"0x", 66, "hash: expected 32-byte sequence, got 33 bytes"},
		{"0x", 63, "hash: UnmarshalText failed: odd length"},
		{"0x", 0, "hash: expected 32-byte sequence, got 0 bytes"},
		{"0x", 64, ""},
		{"0X", 64, "hash: expected a hex string starting with '0x'"},
	}
	for _, test := range tests {
		input := test.Prefix + strings.Repeat("0", test.Size)
		v := new(Hash)
		err := v.UnmarshalText([]byte(input))
		if err == nil {
			if test.Error != "" {
				t.Errorf("%s: error mismatch: have nil, want %q", input, test.Error)
			}
		} else {
			if err.Error() != test.Error {
				t.Errorf("%s: error mismatch: have %q, want %q", input, err, test.Error)
			}
		}
	}
}

func Test_Hash_UnmarshalJSON(t *testing.T) {
	var tests = []struct {
		Prefix string
		Size   int
		Error  string
	}{
		{"", 62, "hash: expected a hex string starting with '0x'"},
		{"0x", 66, "hash: expected 32-byte sequence, got 33 bytes"},
		{"0x", 63, "hash: UnmarshalText failed: odd length"},
		{"0x", 0, "hash: expected 32-byte sequence, got 0 bytes"},
		{"0x", 64, ""},
		{"0X", 64, "hash: expected a hex string starting with '0x'"},
	}
	for _, test := range tests {
		input := `"` + test.Prefix + strings.Repeat("0", test.Size) + `"`
		var v Hash
		err := json.Unmarshal([]byte(input), &v)
		if err == nil {
			if test.Error != "" {
				t.Errorf("%s: error mismatch: have nil, want %q", input, test.Error)
			}
		} else {
			if err.Error() != test.Error {
				t.Errorf("%s: error mismatch: have %q, want %q", input, err, test.Error)
			}
		}
	}
}
