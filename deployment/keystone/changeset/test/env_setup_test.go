package test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
)

func TestSetupEnv(t *testing.T) {
	t.Parallel()
	ctx := tests.Context(t)

	t.Run("test env with in memory nodes", func(t *testing.T) {
		for _, useMCMS := range []bool{true, false} {
			te := SetupContractTestEnv(t, EnvWrapperConfig{
				WFDonConfig:     DonConfig{Name: "wfDon", N: 4},
				AssetDonConfig:  DonConfig{Name: "assetDon", N: 4},
				WriterDonConfig: DonConfig{Name: "writerDon", N: 4},
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
				for _, donNames := range []string{"wfDon", "assetDon", "writerDon"} {
					require.Len(t, te.GetP2PIDs(donNames), 4, "don %s should have 4 p2p ids", donNames)
				}
			})
		}
	})

	t.Run("test env with view only, non functional node stubs", func(t *testing.T) {
		for _, useMCMS := range []bool{true, false} {
			te := SetupContractTestEnv(t, EnvWrapperConfig{
				WFDonConfig:     DonConfig{Name: "wfDon", N: 4},
				AssetDonConfig:  DonConfig{Name: "assetDon", N: 4},
				WriterDonConfig: DonConfig{Name: "writerDon", N: 4},
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
				for _, donNames := range []string{"wfDon", "assetDon", "writerDon"} {
					require.Len(t, te.GetP2PIDs(donNames), 4, "don %s should have 4 p2p ids", donNames)
				}
			})
		}
	})
}
