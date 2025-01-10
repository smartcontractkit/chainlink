package test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
)

func TestSetupTestEnv(t *testing.T) {
	t.Parallel()
	ctx := tests.Context(t)
	for _, useMCMS := range []bool{true, false} {
		te := SetupTestEnv(t, TestConfig{
			WFDonConfig:     DonConfig{N: 4},
			AssetDonConfig:  DonConfig{N: 4},
			WriterDonConfig: DonConfig{N: 4},
			NumChains:       3,
			UseMCMS:         useMCMS,
		})
		t.Run(fmt.Sprintf("set up test env using MCMS: %t", useMCMS), func(t *testing.T) {
			require.NotNil(t, te.Env.ExistingAddresses)
			require.Len(t, te.Env.Chains, 3)
			require.NotEmpty(t, te.RegistrySelector)
			require.NotNil(t, te.Env.Offchain)
			r, err := te.Env.Offchain.ListNodes(ctx, &node.ListNodesRequest{})
			require.NoError(t, err)
			require.Len(t, r.Nodes, 12)
		})
	}
}
