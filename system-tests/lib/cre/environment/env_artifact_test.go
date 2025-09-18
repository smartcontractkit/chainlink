// system-tests/lib/cre/environment/env_artifact_test.go
package environment

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

// set to true if marshal/unmarshal golden files should be updated after code changes
var update = flag.Bool("update", false, "update golden files")

func TestDONCapabilityConfigRoundTrip(t *testing.T) {
	cases := []struct{ dir string }{
		{dir: "000-empty-config"},
		{dir: "001-empty-remote-config-null"},
		{dir: "002-remote-trigger"},
		{dir: "003-remote-target"},
		{dir: "004-remote-executable"},
		{dir: "005-method-configs-trigger"},
		{dir: "006-method-configs-executable-all-at-once"},   // transmission_schedule=0, no deltaStage
		{dir: "007-method-configs-executable-one-at-a-time"}, // transmission_schedule=1, with deltaStage
		{dir: "008-method-configs-all-options"},
	}
	for _, tc := range cases {
		t.Run(tc.dir, func(t *testing.T) {
			in := readFile(t, filepath.Join("testdata/roundtrip", tc.dir, "input.json"))

			// struct under test. Adjust if top-level shape differs.
			var obj struct {
				Config DONCapabilityConfig `json:"config"`
			}

			if err := json.Unmarshal(in, &obj); err != nil { // triggers custom UnmarshalJSON
				t.Fatalf("failed when unmarshal: %v", err)
			}

			out, err := json.Marshal(obj) // triggers custom MarshalJSON
			if err != nil {
				t.Fatalf("fail when marshal: %v", err)
			}

			got := canonJSON(t, out)
			expPath := filepath.Join("testdata/roundtrip", tc.dir, "expected.json")

			if *update {
				// only used when updating testdata or changing the env_artifact code
				writeFile(t, expPath, got)
			}

			want := canonJSON(t, readFile(t, expPath))
			diff := cmp.Diff(string(want), string(got))
			assert.Empty(t, diff, "mismatch (-want +got):\n%s", diff)
		})
	}
}
func TestDONCapabilityConfigUnknownRemoteConfigType(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		errContains string
	}{
		{
			name: "unknown remote config type",
			input: []byte(`{
			"config": {
				"RemoteConfig": {
					"UnknownConfigType": {
						"foo": "bar"
					}
				}
			}
		}`),
			errContains: "unknown remote config type in capability config, keys: [UnknownConfigType]",
		},
		{
			name: "unknown method_configs config type",
			input: []byte(`{
			"config": {
				"method_configs": {
					"SomeRandomMethod": {
						"RemoteConfig": {
							"UnknownConfigType2": {
								"foo": "bar"
							}
						}
					}
				}
			}
		}`),
			errContains: "unknown method config type for method SomeRandomMethod, unknown config value keys: UnknownConfigType2",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var obj struct {
				Config DONCapabilityConfig `json:"config"`
			}
			err := json.Unmarshal(tc.input, &obj)
			assert.ErrorContains(t, err, tc.errContains)
		})
	}
}

func readFile(t *testing.T, file string) []byte {
	t.Helper()
	b, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("couldn't read %s: %v", file, err)
	}
	return b
}

func writeFile(t *testing.T, file string, bytes []byte) {
	t.Helper()
	if err := os.WriteFile(file, bytes, 0600); err != nil {
		t.Fatalf("couldn't write %s: %v", file, err)
	}
}

func canonJSON(t *testing.T, bytes []byte) []byte {
	t.Helper()
	var v any
	if err := json.Unmarshal(bytes, &v); err != nil {
		t.Fatalf("couldn't read bad json: %v\n%s", err, string(bytes))
	}
	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatalf("couldn't marshal indent: %v", err)
	}
	return append(out, '\n')
}
