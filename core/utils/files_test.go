package utils

import (
	"encoding/hex"
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func tempFileName() string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	return filepath.Join(os.TempDir(), hex.EncodeToString(randBytes))
}

func TestFileExists(t *testing.T) {
	t.Parallel()

	exists, err := FileExists(tempFileName())
	require.NoError(t, err)
	assert.False(t, exists)

	exists, err = FileExists(os.Args[0])
	require.NoError(t, err)
	assert.True(t, exists)
}

func TestTooPermissive(t *testing.T) {
	t.Parallel()

	res := TooPermissive(os.FileMode(0700), os.FileMode(0600))
	assert.True(t, res)

	res = TooPermissive(os.FileMode(0600), os.FileMode(0600))
	assert.False(t, res)

	res = TooPermissive(os.FileMode(0600), os.FileMode(0700))
	assert.False(t, res)
}

func TestFileSize_MarshalText_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    FileSize
		expected string
	}{
		{FileSize(0), "0b"},
		{FileSize(1), "1b"},
		{FileSize(MB), "1.00mb"},
		{FileSize(KB), "1.00kb"},
		{FileSize(MB), "1.00mb"},
		{FileSize(GB), "1.00gb"},
		{FileSize(TB), "1.00tb"},
		{FileSize(5 * GB), "5.00gb"},
		{FileSize(0.5 * GB), "500.00mb"},
	}

	for _, test := range tests {
		test := test

		t.Run(test.expected, func(t *testing.T) {
			t.Parallel()

			bstr, err := test.input.MarshalText()
			assert.NoError(t, err)
			assert.Equal(t, test.expected, string(bstr))
			assert.Equal(t, test.expected, test.input.String())
		})
	}
}

func TestFileSize_UnmarshalText(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    string
		expected FileSize
		valid    bool
	}{
		// valid
		{"0", FileSize(0), true},
		{"0.0", FileSize(0), true},
		{"1.12345", FileSize(1), true},
		{"123", FileSize(123), true},
		{"123", FileSize(123), true},
		{"123b", FileSize(123), true},
		{"123B", FileSize(123), true},
		{"123kb", FileSize(123 * KB), true},
		{"123KB", FileSize(123 * KB), true},
		{"123mb", FileSize(123 * MB), true},
		{"123gb", FileSize(123 * GB), true},
		{"123tb", FileSize(123 * TB), true},
		{"5.5mb", FileSize(5.5 * MB), true},
		{"0.5mb", FileSize(0.5 * MB), true},
		// invalid
		{"", FileSize(0), false},
		{"xyz", FileSize(0), false},
		{"-1g", FileSize(0), false},
		{"+1g", FileSize(0), false},
		{"1g", FileSize(0), false},
		{"1t", FileSize(0), false},
		{"1a", FileSize(0), false},
		{"1tbtb", FileSize(0), false},
		{"1tb1tb", FileSize(0), false},
	}

	for _, test := range tests {
		test := test

		t.Run(test.input, func(t *testing.T) {
			t.Parallel()

			var fs FileSize
			err := fs.UnmarshalText([]byte(test.input))
			if test.valid {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, fs)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
