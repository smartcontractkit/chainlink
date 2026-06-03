package dbdetect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractGoListFlags(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args []string
		want []string
	}{
		{
			name: "tags as separate args",
			args: []string{"-tags", "integration", "./core/..."},
			want: []string{"-tags", "integration"},
		},
		{
			name: "tags with equals",
			args: []string{"-tags=integration", "./core/..."},
			want: []string{"-tags=integration"},
		},
		{
			name: "mod and modfile",
			args: []string{"-mod", "readonly", "-modfile", "go.local.mod", "./core/..."},
			want: []string{"-mod", "readonly", "-modfile", "go.local.mod"},
		},
		{
			name: "ignores go test only flags",
			args: []string{"-run", "TestFoo", "-count", "2", "-tags", "integration", "./core/..."},
			want: []string{"-tags", "integration"},
		},
		{
			name: "ignores harness flags",
			args: []string{"--database-url", "postgres://example", "./core/..."},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, extractGoListFlags(tt.args))
		})
	}
}

func TestBuildGoListArgs(t *testing.T) {
	t.Parallel()

	got := buildGoListArgs(
		[]string{"./core/internal/testutils/dbdetectfixture"},
		[]string{"-tags", "dbdetecttag", "-run", "TestFoo"},
	)
	assert.Equal(t, []string{
		"list", "-deps", "-test",
		"-tags", "dbdetecttag",
		"./core/internal/testutils/dbdetectfixture",
	}, got)
}
