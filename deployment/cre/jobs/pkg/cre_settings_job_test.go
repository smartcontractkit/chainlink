package pkg

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/cresettings"
)

func TestCRESettingsJob_ResolveJob(t *testing.T) {
	tests := []struct {
		name     string
		settings string
	}{
		{"empty", ""},
		{"simple", "Foo = 42"},
		{"multi-line", `Foo = 42
[Bar]
Baz = "test"
`},
		{"single-quotes", `Foo = 'bar'
`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := CRESettingsJob{Settings: tt.settings}
			spec, err := j.ResolveJob()
			require.NoError(t, err)

			t.Log("Spec:", spec)

			got, err := cresettings.ValidatedCRESettingsSpec(spec)
			require.NoError(t, err)
			require.Equal(t, tt.settings, got.CRESettingsSpec.Settings)
		})
	}
}
