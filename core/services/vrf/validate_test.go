package vrf

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/stretchr/testify/require"
)

func TestValidateVRFJobSpec(t *testing.T) {
	var tt = []struct {
		name      string
		toml      string
		assertion func(t *testing.T, os job.Job, err error)
	}{
		{
			name: "valid spec",
			toml: `
type            = "vrf"
schemaVersion   = 1
confirmations = 10
publicKey = "0x79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F8179800"
coordinatorAddress = "0xB3b7874F13387D44a3398D298B075B7A3505D8d4"
observationSource   = """
getrandomvalue [type=vrf];
"""
`,
			assertion: func(t *testing.T, s job.Job, err error) {
				require.NoError(t, err)
				require.NotNil(t, s.VRFSpec)
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s, err := ValidateVRFSpec(tc.toml)
			tc.assertion(t, s, err)
		})
	}
}
