package beholder

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gotoml "github.com/pelletier/go-toml/v2"
)

func TestSanitizeURLString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		in       string
		expected string
	}{
		{
			name:     "url with userinfo query and fragment",
			in:       "wss://user:pass@rpc.example.com/path?key=k#frag",
			expected: "wss://userxxx:passwordxxx@rpc.example.com/path?<QUERY>#<FRAGMENT>",
		},
		{
			name:     "url with userinfo only",
			in:       "wss://user:pass@rpc.example.com/path",
			expected: "wss://userxxx:passwordxxx@rpc.example.com/path",
		},
		{
			name:     "url with username only",
			in:       "wss://user@rpc.example.com/path",
			expected: "wss://userxxx@rpc.example.com/path",
		},
		{
			name:     "url with password only",
			in:       "wss://:pass@rpc.example.com/path",
			expected: "wss://:passwordxxx@rpc.example.com/path",
		},
		{
			name:     "url with empty password",
			in:       "wss://user:@rpc.example.com/path",
			expected: "wss://userxxx:passwordxxx@rpc.example.com/path",
		},
		{
			name:     "url with query only",
			in:       "https://rpc.example.com/v1?key=k",
			expected: "https://rpc.example.com/v1?<QUERY>",
		},
		{
			name:     "url with fragment only",
			in:       "https://rpc.example.com/v1#section",
			expected: "https://rpc.example.com/v1#<FRAGMENT>",
		},
		{
			name:     "url with query and fragment no userinfo",
			in:       "https://rpc.example.com/v1?key=k#section",
			expected: "https://rpc.example.com/v1?<QUERY>#<FRAGMENT>",
		},
		{
			name:     "clean url unchanged",
			in:       "https://rpc.example.com/v1",
			expected: "https://rpc.example.com/v1",
		},
		{
			name:     "url with port unchanged",
			in:       "wss://user:pass@rpc.example.com:8546/path",
			expected: "wss://userxxx:passwordxxx@rpc.example.com:8546/path",
		},
		{
			name:     "non-url string unchanged",
			in:       "prom.test",
			expected: "prom.test",
		},
		{
			name:     "string with host:port but no scheme unchanged",
			in:       "localhost:50051",
			expected: "localhost:50051",
		},
		{
			name:     "string with hash but no scheme-host unchanged",
			in:       "section#2",
			expected: "section#2",
		},
		{
			name:     "string with question but no scheme-host unchanged",
			in:       "color?blue",
			expected: "color?blue",
		},
		{
			name:     "mailto url without host unchanged",
			in:       "mailto:user@example.com",
			expected: "mailto:user@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, sanitizeURLString(tt.in))
		})
	}
}

func TestSanitizeConfigTOML(t *testing.T) {
	t.Parallel()

	t.Run("valid toml with urls", func(t *testing.T) {
		t.Parallel()
		in := `
URL = 'wss://user:pass@rpc.example.com/path?key=k#frag'
QueryOnly = 'https://rpc.example.com/v1?key=k'
Clean = 'https://rpc.example.com/v1'
Header = 'Authorization: token'
Duration = '10s'
`

		out, err := SanitizeConfigTOML(in)
		require.NoError(t, err)

		var m map[string]any
		require.NoError(t, gotoml.Unmarshal([]byte(out), &m))

		assert.Equal(t, "wss://userxxx:passwordxxx@rpc.example.com/path?<QUERY>#<FRAGMENT>", m["URL"])
		assert.Equal(t, "https://rpc.example.com/v1?<QUERY>", m["QueryOnly"])
		assert.Equal(t, "https://rpc.example.com/v1", m["Clean"])
		assert.Equal(t, "Authorization: token", m["Header"])
		assert.Equal(t, "10s", m["Duration"])
	})

	t.Run("urls in array", func(t *testing.T) {
		t.Parallel()
		in := `URLs = ['https://u:p@h1/q?k1=v1', 'https://h2/q']`

		out, err := SanitizeConfigTOML(in)
		require.NoError(t, err)

		var m map[string]any
		require.NoError(t, gotoml.Unmarshal([]byte(out), &m))

		urls, ok := m["URLs"].([]any)
		require.True(t, ok)
		require.Len(t, urls, 2)
		assert.Equal(t, "https://userxxx:passwordxxx@h1/q?<QUERY>", urls[0])
		assert.Equal(t, "https://h2/q", urls[1])
	})

	t.Run("nested table urls", func(t *testing.T) {
		t.Parallel()
		in := `
[Node]
URL = 'https://user:pass@rpc.example.com'
Name = 'primary'
`

		out, err := SanitizeConfigTOML(in)
		require.NoError(t, err)

		var m map[string]any
		require.NoError(t, gotoml.Unmarshal([]byte(out), &m))

		node, ok := m["Node"].(map[string]any)
		require.True(t, ok)
		assert.Equal(t, "https://userxxx:passwordxxx@rpc.example.com", node["URL"])
		assert.Equal(t, "primary", node["Name"])
	})

	t.Run("array of tables urls", func(t *testing.T) {
		t.Parallel()
		in := `
[[Nodes]]
URL = 'https://user:pass@rpc.example.com'

[[Nodes]]
URL = 'https://rpc.example.com?key=k'
`

		out, err := SanitizeConfigTOML(in)
		require.NoError(t, err)

		var m map[string]any
		require.NoError(t, gotoml.Unmarshal([]byte(out), &m))

		nodes, ok := m["Nodes"].([]any)
		require.True(t, ok)
		require.Len(t, nodes, 2)

		node0 := nodes[0].(map[string]any)
		assert.Equal(t, "https://userxxx:passwordxxx@rpc.example.com", node0["URL"])

		node1 := nodes[1].(map[string]any)
		assert.Equal(t, "https://rpc.example.com?<QUERY>", node1["URL"])
	})

	t.Run("empty input", func(t *testing.T) {
		t.Parallel()
		out, err := SanitizeConfigTOML("")
		require.NoError(t, err)
		assert.Empty(t, out)
	})

	t.Run("invalid toml returns error", func(t *testing.T) {
		t.Parallel()
		_, err := SanitizeConfigTOML("[unclosed")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "decode config TOML")
	})
}
