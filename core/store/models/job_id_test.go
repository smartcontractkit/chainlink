package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestID_UnmarshalText(t *testing.T) {
	t.Parallel()

	i := &JobID{}

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"slim uuid", `3d7af8e1-ede5-4350-864c-663cfe0ad8e5`, "3d7af8e1ede54350864c663cfe0ad8e5"},
		{"uuid with dashes", `3d7af8e1-ede5-4350-864c-663cfe0ad8e5`, "3d7af8e1ede54350864c663cfe0ad8e5"},
		{"uppercase uuid", `3D7AF8E1-EDE5-4350-864C-663CFE0AD8E5`, "3d7af8e1ede54350864c663cfe0ad8e5"},
		{"wrapped uuid", `"3d7af8e1-ede5-4350-864c-663cfe0ad8e5"`, "3d7af8e1ede54350864c663cfe0ad8e5"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := i.UnmarshalText([]byte(test.input))
			require.NoError(t, err)
			assert.Equal(t, test.want, i.String())
		})
	}
}

func TestID_UnmarshalString(t *testing.T) {
	t.Parallel()

	i := &JobID{}

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"slim uuid", `3d7af8e1-ede5-4350-864c-663cfe0ad8e5`, "3d7af8e1ede54350864c663cfe0ad8e5"},
		{"uuid with dashes", `3d7af8e1-ede5-4350-864c-663cfe0ad8e5`, "3d7af8e1ede54350864c663cfe0ad8e5"},
		{"uppercase uuid", `3D7AF8E1-EDE5-4350-864C-663CFE0AD8E5`, "3d7af8e1ede54350864c663cfe0ad8e5"},
		{"wrapped uuid", `"3d7af8e1-ede5-4350-864c-663cfe0ad8e5"`, "3d7af8e1ede54350864c663cfe0ad8e5"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := i.UnmarshalString(test.input)
			require.NoError(t, err)
			assert.Equal(t, test.want, i.String())
		})
	}
}
