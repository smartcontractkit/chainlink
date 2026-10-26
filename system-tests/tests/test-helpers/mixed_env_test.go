package helpers

import "testing"

func TestIsMixedEnvTopology(t *testing.T) {
	cases := []struct {
		name         string
		topologyName string
		want         bool
	}{
		{name: "empty", topologyName: "", want: false},
		{name: "single-image topology", topologyName: "workflow-gateway-capabilities-don", want: false},
		{name: "exact mixed-env", topologyName: "mixed-env", want: true},
		{name: "mixed-env confidential variant", topologyName: "mixed-env-confidential-workflows", want: true},
		{name: "case-insensitive", topologyName: "MIXED-ENV", want: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("TOPOLOGY_NAME", tc.topologyName)
			if got := IsMixedEnvTopology(); got != tc.want {
				t.Fatalf("IsMixedEnvTopology() with TOPOLOGY_NAME=%q = %v, want %v", tc.topologyName, got, tc.want)
			}
		})
	}
}
