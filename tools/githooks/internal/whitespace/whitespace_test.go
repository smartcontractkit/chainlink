package whitespace_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/whitespace"
)

func TestFixContent_Go(t *testing.T) {
	t.Parallel()

	t.Run("trims trailing whitespace on regular go code lines", func(t *testing.T) {
		t.Parallel()
		input := []byte("package main    \n\nfunc main() {  \n\tx := 1    \n\t_ = x\n}\n")
		expected := []byte("package main\n\nfunc main() {\n\tx := 1\n\t_ = x\n}\n")

		fixed, changed, err := whitespace.FixContent("main.go", input)
		require.NoError(t, err)
		assert.True(t, changed)
		assert.Equal(t, string(expected), string(fixed))
	})

	t.Run("trims blank lines with whitespace", func(t *testing.T) {
		t.Parallel()
		input := []byte("package main\n    \nfunc main() {}\n")
		expected := []byte("package main\n\nfunc main() {}\n")

		fixed, changed, err := whitespace.FixContent("main.go", input)
		require.NoError(t, err)
		assert.True(t, changed)
		assert.Equal(t, string(expected), string(fixed))
	})

	t.Run("trims trailing spaces on comments", func(t *testing.T) {
		t.Parallel()
		input := []byte("package main\n\n// line comment with spaces   \n/* block comment with spaces   \nanother line   \n*/\n")
		expected := []byte("package main\n\n// line comment with spaces\n/* block comment with spaces\nanother line\n*/\n")

		fixed, changed, err := whitespace.FixContent("main.go", input)
		require.NoError(t, err)
		assert.True(t, changed)
		assert.Equal(t, string(expected), string(fixed))
	})

	t.Run("STRICT INVARIANT: preserves trailing whitespace inside Go raw backtick strings", func(t *testing.T) {
		t.Parallel()
		// Inside raw string literal, trailing whitespace is intentional and semantically significant!
		input := []byte("package main\n\nconst query = `\nSELECT *    \nFROM users   \nWHERE id = 1  \n`\n\nfunc Run() {   \n\t_ = query\n}\n")
		expected := []byte("package main\n\nconst query = `\nSELECT *    \nFROM users   \nWHERE id = 1  \n`\n\nfunc Run() {\n\t_ = query\n}\n")

		fixed, changed, err := whitespace.FixContent("main.go", input)
		require.NoError(t, err)
		assert.True(t, changed)
		assert.Equal(t, string(expected), string(fixed))
	})

	t.Run("interpreted string literal on line with trailing spaces outside string", func(t *testing.T) {
		t.Parallel()
		input := []byte("package main\n\nfunc main() {\n\tmsg := \"hello world\"   \n\t_ = msg\n}\n")
		expected := []byte("package main\n\nfunc main() {\n\tmsg := \"hello world\"\n\t_ = msg\n}\n")

		fixed, changed, err := whitespace.FixContent("main.go", input)
		require.NoError(t, err)
		assert.True(t, changed)
		assert.Equal(t, string(expected), string(fixed))
	})
}

func TestFixContent_Markdown(t *testing.T) {
	t.Parallel()

	t.Run("preserves two trailing spaces for hard line breaks", func(t *testing.T) {
		t.Parallel()
		input := []byte("# Header\n\nThis line has a hard break.  \nNext line here.\n")
		expected := []byte("# Header\n\nThis line has a hard break.  \nNext line here.\n")

		fixed, changed, err := whitespace.FixContent("README.md", input)
		require.NoError(t, err)
		assert.False(t, changed)
		assert.Equal(t, string(expected), string(fixed))
	})

	t.Run("trims single trailing space in markdown", func(t *testing.T) {
		t.Parallel()
		input := []byte("# Header \n\nLine with one trailing space. \nNext line.\n")
		expected := []byte("# Header\n\nLine with one trailing space.\nNext line.\n")

		fixed, changed, err := whitespace.FixContent("README.md", input)
		require.NoError(t, err)
		assert.True(t, changed)
		assert.Equal(t, string(expected), string(fixed))
	})

	t.Run("trims 3+ trailing spaces down to 2 in markdown for hard line breaks", func(t *testing.T) {
		t.Parallel()
		input := []byte("This line has four trailing spaces.    \nNext line.\n")
		expected := []byte("This line has four trailing spaces.  \nNext line.\n")

		fixed, changed, err := whitespace.FixContent("README.md", input)
		require.NoError(t, err)
		assert.True(t, changed)
		assert.Equal(t, string(expected), string(fixed))
	})

	t.Run("trims blank lines with whitespace in markdown", func(t *testing.T) {
		t.Parallel()
		input := []byte("# Header\n  \nParagraph\n")
		expected := []byte("# Header\n\nParagraph\n")

		fixed, changed, err := whitespace.FixContent("README.md", input)
		require.NoError(t, err)
		assert.True(t, changed)
		assert.Equal(t, string(expected), string(fixed))
	})
}

func TestFixContent_Generic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		filename string
		input    []byte
		expected []byte
	}{
		{
			name:     "yaml file",
			filename: "config.yaml",
			input:    []byte("key: value   \nlist:  \n  - item 1   \n"),
			expected: []byte("key: value\nlist:\n  - item 1\n"),
		},
		{
			name:     "json file",
			filename: "data.json",
			input:    []byte("{\n  \"key\": \"value\"   \n}\n"),
			expected: []byte("{\n  \"key\": \"value\"\n}\n"),
		},
		{
			name:     "sql file",
			filename: "query.sql",
			input:    []byte("SELECT *   \nFROM users;   \n"),
			expected: []byte("SELECT *\nFROM users;\n"),
		},
		{
			name:     "shell file",
			filename: "run.sh",
			input:    []byte("#!/bin/bash   \necho \"hello\"   \n"),
			expected: []byte("#!/bin/bash\necho \"hello\"\n"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fixed, changed, err := whitespace.FixContent(tt.filename, tt.input)
			require.NoError(t, err)
			assert.True(t, changed)
			assert.Equal(t, string(tt.expected), string(fixed))
		})
	}
}

func TestFixFile(t *testing.T) {
	t.Parallel()

	t.Run("fixes trailing whitespace on disk when checkOnly is false", func(t *testing.T) {
		t.Parallel()
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "app.go")
		require.NoError(t, os.WriteFile(filePath, []byte("package main   \n"), 0o600))

		changed, err := whitespace.FixFile(filePath, false)
		require.NoError(t, err)
		assert.True(t, changed)

		content, err := os.ReadFile(filePath)
		require.NoError(t, err)
		assert.Equal(t, "package main\n", string(content))
	})

	t.Run("checkOnly true reports changes without modifying disk", func(t *testing.T) {
		t.Parallel()
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "app.go")
		require.NoError(t, os.WriteFile(filePath, []byte("package main   \n"), 0o600))

		changed, err := whitespace.FixFile(filePath, true)
		require.NoError(t, err)
		assert.True(t, changed)

		content, err := os.ReadFile(filePath)
		require.NoError(t, err)
		assert.Equal(t, "package main   \n", string(content))
	})
}
