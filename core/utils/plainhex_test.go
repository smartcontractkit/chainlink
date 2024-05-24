package utils

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"testing"
)

type marshalTest struct {
	input interface{}
	want  string
}

type unmarshalTest struct {
	input   string
	want    interface{}
	wantErr string
}

var (
	unmarshalBytesTests = []unmarshalTest{
		// invalid encoding
		{input: "", wantErr: "unexpected end of JSON input"},
		{input: "null", wantErr: "json: cannot unmarshal non-string into Go value of type utils.PlainHexBytes"},
		{input: `"null"`, wantErr: "UnmarshalJSON failed: UnmarshalText failed: encoding/hex: invalid byte: U+006E 'n'"},
		{input: `"0x"`, wantErr: "UnmarshalJSON failed: UnmarshalText failed: encoding/hex: invalid byte: U+0078 'x'"},
		{input: `"0X"`, wantErr: "UnmarshalJSON failed: UnmarshalText failed: encoding/hex: invalid byte: U+0058 'X'"},
		{input: `"0"`, wantErr: "UnmarshalJSON failed: UnmarshalText failed: odd length"},
		{input: `"xx"`, wantErr: "UnmarshalJSON failed: UnmarshalText failed: encoding/hex: invalid byte: U+0078 'x'"},
		{input: `"01zz01"`, wantErr: "UnmarshalJSON failed: UnmarshalText failed: encoding/hex: invalid byte: U+007A 'z'"},

		// valid encoding
		{input: `""`, want: referenceBytes("")},
		{input: `"02"`, want: referenceBytes("02")},
		{input: `"ffffffffff"`, want: referenceBytes("ffffffffff")},
		{
			input: `"ffffffffffffffffffffffffffffffffffff"`,
			want:  referenceBytes("ffffffffffffffffffffffffffffffffffff"),
		},
	}

	encodeBytesTests = []marshalTest{
		{[]byte{}, ""},
		{[]byte{0}, "00"},
		{[]byte{0, 0, 1, 2}, "00000102"},
	}
)

func TestUnmarshalBytes(t *testing.T) {
	for _, test := range unmarshalBytesTests {
		var v PlainHexBytes
		err := json.Unmarshal([]byte(test.input), &v)
		if !checkError(t, test.input, err, test.wantErr) {
			continue
		}
		if !bytes.Equal(test.want.([]byte), v) {
			t.Errorf("input %s: value mismatch: got %x, want %x", test.input, &v, test.want)
			continue
		}
	}
}

func TestMarshalBytes(t *testing.T) {
	for _, test := range encodeBytesTests {
		in := test.input.([]byte)
		out, err := json.Marshal(PlainHexBytes(in))
		if err != nil {
			t.Errorf("%x: %v", in, err)
			continue
		}
		if want := `"` + test.want + `"`; string(out) != want {
			t.Errorf("%x: MarshalJSON output mismatch: got %q, want %q", in, out, want)
			continue
		}
		if out := PlainHexBytes(in).String(); out != test.want {
			t.Errorf("%x: String mismatch: got %q, want %q", in, out, test.want)
			continue
		}
	}
}

func referenceBytes(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return b
}

func checkError(t *testing.T, input string, got error, want string) bool {
	if got == nil {
		if want != "" {
			t.Errorf("input %s: got no error, want %q", input, want)
			return false
		}
		return true
	}
	if want == "" {
		t.Errorf("input %s: unexpected error %q", input, got)
	} else if got.Error() != want {
		t.Errorf("input %s: got error %q, want %q", input, got, want)
	}
	return false
}
