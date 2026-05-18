//go:build go1.18

package pipeline

import (
	"testing"
)

func FuzzParseETHABIArgsString(f *testing.F) {
	for _, tt := range testsABIDecode {
		f.Add(tt.abi, false)
	}
	f.Fuzz(func(t *testing.T, theABI string, isLog bool) {
		_, _, err := ParseETHABIArgsString([]byte(theABI), isLog)
		if err != nil {
			t.Skip()
		}
	})
}
