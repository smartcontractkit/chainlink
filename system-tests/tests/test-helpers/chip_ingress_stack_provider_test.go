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
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name                 string
		createBinary         bool
		binaryMode           os.FileMode
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
			name:                 "uses precompiled binary when it exists and is executable",
			createBinary:         true,
			binaryMode:           0700,
			inputArgs:            []string{"env", "start"},
			expectedPathSuffix:   "cre-env",
			expectedArgsContains: "env",
		},
		{
			name:                 "fallback to go run when binary exists but is not executable",
			createBinary:         true,
			binaryMode:           0600,
			inputArgs:            []string{"env", "start"},
			expectedPathSuffix:   "go",
			expectedArgsContains: "run",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tmpDir := t.TempDir()
			environmentDir := filepath.Join(tmpDir, "core", "scripts", "cre", "environment")
			require.NoError(t, os.MkdirAll(environmentDir, 0700))

			if tt.createBinary {
				binDir := filepath.Join(tmpDir, "system-tests", "tests", "bin")
				require.NoError(t, os.MkdirAll(binDir, 0700))
				binPath := filepath.Join(binDir, "cre-env")
				require.NoError(t, os.WriteFile(binPath, []byte("#!/bin/sh\necho ok"), tt.binaryMode))
			}

			cmd := resolveCreEnvCommand(ctx, tmpDir, environmentDir, tt.inputArgs...)
			require.NotNil(t, cmd)

			assert.Equal(t, environmentDir, cmd.Dir)
			assert.Equal(t, tt.expectedPathSuffix, filepath.Base(cmd.Path))
			assert.Contains(t, cmd.Args, tt.expectedArgsContains)
		})
	}
}
