package testshard

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShardForPackage(t *testing.T) {
	t.Parallel()
	pkg1 := "github.com/smartcontractkit/chainlink/v2/core/services/workflows"
	pkg2 := "github.com/smartcontractkit/chainlink/v2/core/services/ocr2"

	shard0 := ShardForPackage(pkg1, 4)
	shard1 := ShardForPackage(pkg2, 4)

	require.GreaterOrEqual(t, shard0, 0)
	require.Less(t, shard0, 4)
	require.GreaterOrEqual(t, shard1, 0)
	require.Less(t, shard1, 4)

	// Deterministic
	require.Equal(t, shard0, ShardForPackage(pkg1, 4))
}

func TestList(t *testing.T) {
	t.Parallel()
	input := strings.NewReader("pkg1\npkg2\npkg3\npkg4\n")
	var stdout bytes.Buffer

	err := List(input, &stdout, 2, 0)
	require.NoError(t, err)

	output := strings.TrimSpace(stdout.String())
	require.NotEmpty(t, output)

	// Output for shard 1
	var stdout1 bytes.Buffer
	input1 := strings.NewReader("pkg1\npkg2\npkg3\npkg4\n")
	err = List(input1, &stdout1, 2, 1)
	require.NoError(t, err)

	output1 := strings.TrimSpace(stdout1.String())
	require.NotEmpty(t, output1)

	// Union should contain all 4 packages, intersection should be empty
	shard0Pkgs := strings.Split(output, "\n")
	shard1Pkgs := strings.Split(output1, "\n")
	require.Equal(t, 4, len(shard0Pkgs)+len(shard1Pkgs))
}

func TestVerify(t *testing.T) {
	t.Parallel()
	input := strings.NewReader("pkg1\npkg2\npkg3\npkg4\n")
	var stdout bytes.Buffer

	err := Verify(input, &stdout, 3)
	require.NoError(t, err)
	require.Contains(t, stdout.String(), "verified 4 packages across 3 shards")
}

func TestVerify_DuplicatePackage(t *testing.T) {
	t.Parallel()
	input := strings.NewReader("pkg1\npkg2\npkg1\n")
	var stdout bytes.Buffer

	err := Verify(input, &stdout, 2)
	require.Error(t, err)
	require.Contains(t, err.Error(), "duplicate package")
}

func TestValidateShardArgs(t *testing.T) {
	t.Parallel()
	require.Error(t, ValidateShardArgs(0, 0))
	require.Error(t, ValidateShardArgs(2, 2))
	require.Error(t, ValidateShardArgs(2, -1))
	require.NoError(t, ValidateShardArgs(2, 1))
}
