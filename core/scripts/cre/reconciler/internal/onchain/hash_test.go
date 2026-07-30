package onchain

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/domain"
)

func testDesiredStateForHashTests() *domain.DesiredState {
	return &domain.DesiredState{
		Chains: []domain.Chain{{ChainID: 1337, Registry: true}},
		DONs: []domain.DON{
			{
				Name:         "workflow",
				DONTypes:     []string{"workflow"},
				Capabilities: []string{"cron", "consensus"},
				CapabilityConfigs: map[string]domain.CapabilityConfig{
					"cron": {BinaryName: "cron", Values: map[string]any{"tick": "1s"}},
				},
			},
			{Name: "bootstrap", DONTypes: []string{"bootstrap"}},
		},
	}
}

func TestHashPreEnvStartupCapReg_DeterministicAndSensitiveToCapConfig(t *testing.T) {
	t.Parallel()

	topology := testWorkflowBootstrapTopology(t)
	desired := testDesiredStateForHashTests()

	hash1, err := hashPreEnvStartupCapReg(desired, topology, nil)
	require.NoError(t, err)
	hash2, err := hashPreEnvStartupCapReg(desired, topology, nil)
	require.NoError(t, err)
	require.Equal(t, hash1, hash2, "identical input must hash identically")

	desired.DONs[0].CapabilityConfigs["cron"] = domain.CapabilityConfig{
		BinaryName: "cron", Values: map[string]any{"tick": "5s"},
	}
	hash3, err := hashPreEnvStartupCapReg(desired, topology, nil)
	require.NoError(t, err)
	require.NotEqual(t, hash1, hash3, "changing a DON's capability config must change the hash")
}

func TestHashPreEnvStartupCapReg_IgnoresBootstrapOnlyDON(t *testing.T) {
	t.Parallel()

	topology := testWorkflowBootstrapTopology(t)
	desired := testDesiredStateForHashTests()

	input := preEnvStartupCapRegHashInput{
		RegistryChainID:         1337,
		GlobalCapabilityConfigs: desired.CapabilityConfigs,
		DONs:                    donHashInputs(desired, topology, nil, capRegExcludedDON),
	}
	require.Len(t, input.DONs, 1, "bootstrap-only DON must be excluded from CapReg's hash scope")
	require.Equal(t, "workflow", input.DONs[0].Name)
}

func TestHashJobs_IncludesAllDONTypes(t *testing.T) {
	t.Parallel()

	topology := testWorkflowBootstrapTopology(t)
	desired := testDesiredStateForHashTests()

	all := donHashInputs(desired, topology, nil, nil)
	require.Len(t, all, 2, "Jobs' DON scope includes bootstrap/gateway DONs, unlike CapReg's")

	hash1, err := hashJobs(desired, topology, nil, nil)
	require.NoError(t, err)
	hash2, err := hashJobs(desired, topology, nil, []domain.GatewayServiceConfigState{{ServiceName: "svc-a"}})
	require.NoError(t, err)
	require.NotEqual(t, hash1, hash2, "changing gateway service configs must change the Jobs hash")
}

func TestHashConfigureWorkflowReg_SensitiveToOwnerAndDONNames(t *testing.T) {
	t.Parallel()

	desired := testDesiredStateForHashTests()

	hash1, err := hashConfigureWorkflowReg(desired, "0xowner1")
	require.NoError(t, err)
	hash2, err := hashConfigureWorkflowReg(desired, "0xowner2")
	require.NoError(t, err)
	require.NotEqual(t, hash1, hash2, "changing the workflow owner must change the hash")

	desired.DONs = append(desired.DONs, domain.DON{Name: "workflow-2", DONTypes: []string{"workflow"}})
	hash3, err := hashConfigureWorkflowReg(desired, "0xowner1")
	require.NoError(t, err)
	require.NotEqual(t, hash1, hash3, "adding a workflow DON must change the hash")
}

func TestSortedCopy_DoesNotMutateInput(t *testing.T) {
	t.Parallel()

	original := []string{"c", "a", "b"}
	sorted := sortedCopy(original)
	require.Equal(t, []string{"a", "b", "c"}, sorted)
	require.Equal(t, []string{"c", "a", "b"}, original, "sortedCopy must not mutate its input")
}
