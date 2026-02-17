package environment

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsTopologyConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		content string
		want    bool
	}{
		{
			name: "valid topology without fake sections",
			content: `
[[blockchains]]
chain_id = "1337"

[[nodesets]]
name = "workflow"

[jd]
csa_encryption_key = "dummy"

[infra]
type = "docker"
`,
			want: true,
		},
		{
			name: "missing infra section",
			content: `
[[blockchains]]
chain_id = "1337"

[[nodesets]]
name = "workflow"

[jd]
`,
			want: false,
		},
		{
			name: "missing jd section",
			content: `
[[blockchains]]
chain_id = "1337"

[[nodesets]]
name = "workflow"

[infra]
type = "docker"
`,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			path := writeTempTopologyFile(t, tt.content)
			got, err := isTopologyConfig(path)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestIsTopologyConfig_InvalidTOMLReturnsError(t *testing.T) {
	t.Parallel()

	path := writeTempTopologyFile(t, "not-valid=[")
	_, err := isTopologyConfig(path)
	require.Error(t, err)
}

func writeTempTopologyFile(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "topology.toml")
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	return path
}
