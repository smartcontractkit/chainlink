package helpers

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveCreEnvCommand(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name                 string
		createBinary         bool
		inputArgs            []string
		expectedPathSuffix   string
		expectedArgsContains string
	}{
		{
			name:                 "fallback to go run when binary does not exist",
			createBinary:         false,
			inputArgs:            []string{"env", "start"},
			expectedPathSuffix:   "go",
			expectedArgsContains: "run",
		},
		{
			name:                 "uses precompiled binary when it exists",
			createBinary:         true,
			inputArgs:            []string{"env", "start"},
			expectedPathSuffix:   "cre-env",
			expectedArgsContains: "env",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			environmentDir := filepath.Join(tmpDir, "core", "scripts", "cre", "environment")
			require.NoError(t, os.MkdirAll(environmentDir, 0755))

			if tt.createBinary {
				binDir := filepath.Join(tmpDir, "system-tests", "tests", "bin")
				require.NoError(t, os.MkdirAll(binDir, 0755))
				binPath := filepath.Join(binDir, "cre-env")
				require.NoError(t, os.WriteFile(binPath, []byte("#!/bin/sh\necho ok"), 0755))
			}

			cmd := resolveCreEnvCommand(ctx, tmpDir, environmentDir, tt.inputArgs...)
			require.NotNil(t, cmd)

			assert.Equal(t, environmentDir, cmd.Dir)
			if tt.createBinary {
				assert.True(t, filepath.IsAbs(cmd.Path) || filepath.Base(cmd.Path) == "cre-env", "path should be cre-env binary: %s", cmd.Path)
			} else {
				assert.Equal(t, "go", filepath.Base(cmd.Path))
			}
		})
	}
}
