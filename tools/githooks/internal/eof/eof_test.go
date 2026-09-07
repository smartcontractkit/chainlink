package eof_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/eof"
)

func TestFixContent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		input         []byte
		expectedFixed []byte
		expectedChg   bool
	}{
		{
			name:          "empty content stays empty",
			input:         []byte(""),
			expectedFixed: []byte(""),
			expectedChg:   false,
		},
		{
			name:          "single line with newline unchanged",
			input:         []byte("hello world\n"),
			expectedFixed: []byte("hello world\n"),
			expectedChg:   false,
		},
		{
			name:          "single line without newline fixed",
			input:         []byte("hello world"),
			expectedFixed: []byte("hello world\n"),
			expectedChg:   true,
		},
		{
			name:          "multiple trailing newlines collapsed to single newline",
			input:         []byte("hello world\n\n\n"),
			expectedFixed: []byte("hello world\n"),
			expectedChg:   true,
		},
		{
			name:          "file with only newlines collapsed to single newline",
			input:         []byte("\n\n\n"),
			expectedFixed: []byte("\n"),
			expectedChg:   true,
		},
		{
			name:          "file with single newline unchanged",
			input:         []byte("\n"),
			expectedFixed: []byte("\n"),
			expectedChg:   false,
		},
		{
			name:          "CRLF endings normalized to LF with single newline",
			input:         []byte("hello world\r\n\r\n"),
			expectedFixed: []byte("hello world\n"),
			expectedChg:   true,
		},
		{
			name:          "multiline code with single trailing newline unchanged",
			input:         []byte("package main\n\nfunc main() {}\n"),
			expectedFixed: []byte("package main\n\nfunc main() {}\n"),
			expectedChg:   false,
		},
		{
			name:          "multiline code missing trailing newline",
			input:         []byte("package main\n\nfunc main() {}"),
			expectedFixed: []byte("package main\n\nfunc main() {}\n"),
			expectedChg:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fixed, changed := eof.FixContent(tt.input)
			assert.Equal(t, tt.expectedChg, changed)
			assert.Equal(t, string(tt.expectedFixed), string(fixed))
		})
	}
}

func TestFixFile(t *testing.T) {
	t.Parallel()

	t.Run("fixes file on disk when checkOnly is false", func(t *testing.T) {
		t.Parallel()
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "test.go")
		require.NoError(t, os.WriteFile(filePath, []byte("package main"), 0o600))

		changed, err := eof.FixFile(filePath, false)
		require.NoError(t, err)
		assert.True(t, changed)

		content, err := os.ReadFile(filePath)
		require.NoError(t, err)
		assert.Equal(t, "package main\n", string(content))
	})

	t.Run("checkOnly true does not write file but reports change", func(t *testing.T) {
		t.Parallel()
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "test.go")
		require.NoError(t, os.WriteFile(filePath, []byte("package main"), 0o600))

		changed, err := eof.FixFile(filePath, true)
		require.NoError(t, err)
		assert.True(t, changed)

		content, err := os.ReadFile(filePath)
		require.NoError(t, err)
		assert.Equal(t, "package main", string(content))
	})

	t.Run("skips non-eligible binary file", func(t *testing.T) {
		t.Parallel()
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "image.png")
		require.NoError(t, os.WriteFile(filePath, []byte{0x89, 'P', 'N', 'G'}, 0o600))

		changed, err := eof.FixFile(filePath, false)
		require.NoError(t, err)
		assert.False(t, changed)
	})
}
