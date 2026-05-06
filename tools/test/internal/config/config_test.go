package config

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadBindsPersistentAndLocalFlags(t *testing.T) {
	t.Parallel()

	root := &cobra.Command{Use: "root"}
	root.PersistentFlags().String("database-url", "", "")
	sub := &cobra.Command{
		Use: "sub",
		Run: func(*cobra.Command, []string) {},
	}
	sub.Flags().Int("iterations", 1, "")
	root.AddCommand(sub)
	root.SetArgs([]string{"sub", "--database-url", "postgres://example", "--iterations", "7"})

	cmd, err := root.ExecuteC()
	require.NoError(t, err)

	conf, err := Load(cmd)
	require.NoError(t, err)
	assert.Equal(t, "postgres://example", conf.DatabaseURL)
	assert.Equal(t, 7, conf.Iterations)
}
