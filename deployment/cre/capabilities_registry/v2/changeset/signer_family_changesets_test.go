package changeset_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	cldf_offchain "github.com/smartcontractkit/chainlink-deployments-framework/offchain"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"

	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/operations/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/test"
)

// emptyJD knows no nodes, so no node can have an accepted job.
type emptyJD struct{ cldf_offchain.Client }

func (emptyJD) ListNodes(context.Context, *nodev1.ListNodesRequest, ...grpc.CallOption) (*nodev1.ListNodesResponse, error) {
	return &nodev1.ListNodesResponse{}, nil
}

func stellarOCR3Config() map[string]any {
	return map[string]any{"ocr3Configs": map[string]any{"__default__": map[string]any{
		"offchainConfig":      map[string]any{"maxFaultyOracles": 1},
		"extraSignerFamilies": []any{"stellar"},
	}}}
}

func TestAddCapabilities_VerifyPreconditions_SignerFamilyJobCheck(t *testing.T) {
	h := test.NewTestHarness(t)

	e := h.Runtime.Environment()
	e.Offchain = emptyJD{}

	in := changeset.AddCapabilitiesInput{
		RegistryChainSel:  h.RegistrySelector,
		RegistryQualifier: test.RegistryQualifier,
		DonCapabilityConfigs: map[string][]contracts.CapabilityConfig{
			test.DONName: {{Capability: contracts.Capability{CapabilityID: "consensus@1.0.0"}, Config: stellarOCR3Config()}},
		},
		SkipDuplicateCapabilityValidation: true,
	}

	require.ErrorContains(t, changeset.AddCapabilities{}.VerifyPreconditions(e, in), "not in JD")

	missing := in
	missing.DonCapabilityConfigs = map[string][]contracts.CapabilityConfig{"no-such-don": in.DonCapabilityConfigs[test.DONName]}
	require.ErrorContains(t, changeset.AddCapabilities{}.VerifyPreconditions(e, missing), "not found on-chain")

	in.SkipSignerFamilyJobCheck = true
	require.NoError(t, changeset.AddCapabilities{}.VerifyPreconditions(e, in))
}

func TestConfigureCapabilitiesRegistry_VerifyPreconditions_SignerFamilyJobCheck(t *testing.T) {
	h := test.NewTestHarness(t)

	e := h.Runtime.Environment()
	e.Offchain = emptyJD{}

	in := changeset.ConfigureCapabilitiesRegistryInput{
		ChainSelector: h.RegistrySelector,
		Qualifier:     test.RegistryQualifier,
		DONs: []changeset.CapabilitiesRegistryNewDONParams{{
			Name:  "new-don",
			Nodes: h.Don.GetP2PIDs().Strings(),
			CapabilityConfigurations: []changeset.CapabilitiesRegistryCapabilityConfiguration{
				{CapabilityID: "consensus@1.0.0", Config: stellarOCR3Config()},
			},
		}},
	}

	require.ErrorContains(t, changeset.ConfigureCapabilitiesRegistry{}.VerifyPreconditions(e, in), "not in JD")

	in.SkipSignerFamilyJobCheck = true
	require.NoError(t, changeset.ConfigureCapabilitiesRegistry{}.VerifyPreconditions(e, in))
}
