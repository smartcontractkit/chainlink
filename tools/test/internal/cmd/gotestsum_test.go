package cmd

import (
	"errors"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunGotestsumRunsCleanupWhenLookPathFails(t *testing.T) {
	var cleanups int
	err := runGotestsum(
		&cobra.Command{},
		nil,
		func(string) (string, error) { return "", errors.New("missing") },
		func() error { cleanups++; return nil },
	)
	require.Error(t, err)
	require.ErrorContains(t, err, "gotestsum not on PATH")
	assert.Equal(t, 1, cleanups, "cleanup must run when gotestsum is missing so ephemeral Postgres from PersistentPreRun is torn down")
}

func TestRunGotestsumRunsCleanupWhenConfigLoadFails(t *testing.T) {
	var cleanups int
	err := runGotestsum(
		nil,
		nil,
		func(string) (string, error) { return "/bin/gotestsum", nil },
		func() error { cleanups++; return nil },
	)
	require.Error(t, err)
	require.ErrorContains(t, err, "command is required")
	assert.Equal(t, 1, cleanups, "cleanup must run after LookPath succeeds but later steps fail")
}
