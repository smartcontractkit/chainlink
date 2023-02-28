package src

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_writeLines(t *testing.T) {
	type args struct {
		lines []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "write read lines",
			args: args{
				lines: []string{"a", "b"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pth := filepath.Join(t.TempDir(), strings.ReplaceAll(tt.name, " ", "_"))
			err := writeLines(tt.args.lines, pth)
			assert.NoError(t, err)
			got, err := readLines(pth)
			assert.NoError(t, err)
			assert.Equal(t, tt.args.lines, got)

		})
	}
}
