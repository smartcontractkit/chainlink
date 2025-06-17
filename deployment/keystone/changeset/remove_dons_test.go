package changeset_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/test"
)

func Test_RemoveDONsRequest_validate(t *testing.T) {
	t.Parallel()

	env := test.SetupContractTestEnv(t, test.EnvWrapperConfig{
		WFDonConfig:     test.DonConfig{Name: "wfDon", N: 4},
		AssetDonConfig:  test.DonConfig{Name: "assetDon", N: 4},
		WriterDonConfig: test.DonConfig{Name: "writerDon", N: 4},
		NumChains:       1,
	})

	tests := []struct {
		name    string
		env     test.EnvWrapper
		req     *changeset.RemoveDONsRequest
		wantErr error
	}{
		{
			name: "empty dons => error",
			req: &changeset.RemoveDONsRequest{
				RegistryChainSel: 1,
				DONs:             []uint32{},
				RegistryRef:      env.CapabilityRegistryAddressRef(),
			},
			env:     env,
			wantErr: errors.New("dons is required"),
		},
		{
			name: "invalid registry chain selector => error",
			req: &changeset.RemoveDONsRequest{
				RegistryChainSel: 100,
				DONs:             []uint32{1, 2, 3},
				RegistryRef:      env.CapabilityRegistryAddressRef(),
			},
			env:     env,
			wantErr: errors.New("invalid registry chain selector 100: selector does not exist"),
		},
		{
			name: "valid request => no error",
			req: &changeset.RemoveDONsRequest{
				RegistryChainSel: env.RegistrySelector,
				DONs:             []uint32{1, 2, 3},
				RegistryRef:      env.CapabilityRegistryAddressRef(),
			},
			env:     env,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.req.Validate(tt.env.Env)
			if tt.wantErr != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestRemoveDONs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		useMCMS         bool // whether to generate a timelock proposal (MCMS) or remove immediately
		expectErr       bool // whether we expect an error from RemoveDONs
		finalizeOnChain bool // if we finalize the changeset on-chain (only relevant if useMCMS == true)
	}{
		{
			name:      "remove dons without MCMS",
			useMCMS:   false,
			expectErr: false,
		},
		{
			name:            "remove dons with MCMS",
			useMCMS:         true,
			expectErr:       false,
			finalizeOnChain: true, // we can show how to finalize and verify on-chain
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// 1) Setup a test environment that has multiple DONs. (Adjust DonConfigs as desired.)
			te := test.SetupContractTestEnv(t, test.EnvWrapperConfig{
				WFDonConfig:     test.DonConfig{Name: "wfDon", N: 4},
				AssetDonConfig:  test.DonConfig{Name: "assetDon", N: 4},
				WriterDonConfig: test.DonConfig{Name: "writerDon", N: 4},
				NumChains:       1,
				UseMCMS:         tt.useMCMS, // indicates if environment is MCMS-ready
			})

			// 2) If the test says “remove these DONs,” gather their actual IDs from the environment
			var donIDs []uint32

			donInfos, err := te.CapabilitiesRegistry().GetDONs(nil)
			require.NoError(t, err)
			for _, donInfo := range donInfos {
				donIDs = append(donIDs, donInfo.Id)
			}

			// 3) Build our request to RemoveDONs
			req := &changeset.RemoveDONsRequest{
				RegistryChainSel: te.RegistrySelector,
				DONs:             donIDs,
				RegistryRef:      te.CapabilityRegistryAddressRef(),
			}
			if tt.useMCMS {
				req.MCMSConfig = &changeset.MCMSConfig{MinDuration: 0}
			}

			// 4) Call the RemoveDONs changeset
			csOut, err := changeset.RemoveDONs(te.Env, req)

			if tt.expectErr {
				require.Error(t, err, "expected an error but got none")
				return
			}
			require.NoError(t, err, "RemoveDONs should not fail in this scenario")

			if !tt.useMCMS {
				// Non-MCMS => no proposals, no pending batch.
				assert.Empty(t, csOut.MCMSTimelockProposals, "expected no MCMS proposals when not using MCMS")

				// Also confirm on-chain that the DON is removed
				for _, donID := range donIDs {
					don, err := te.CapabilitiesRegistry().GetDON(nil, donID)
					require.NoError(t, err)
					require.NotEqual(t, don.Id, donID, "Expect donID to not be found")
					require.Equal(t, uint32(0), don.Id, "Expected returned donID to be 0")
				}
			} else {
				// MCMS => we expect a timelock proposal
				require.Len(t, csOut.MCMSTimelockProposals, 1, "should have exactly one proposal in MCMS mode")

				// If finalizeOnChain is true, we actually apply the proposals on-chain
				if tt.finalizeOnChain {
					err = applyProposal(t, te, commonchangeset.Configure(cldf.CreateLegacyChangeSet(changeset.RemoveDONs), req))
					require.NoError(t, err, "failed to finalize MCMS proposal on-chain")

					// Check on-chain that the DON is removed
					for _, donID := range donIDs {
						don, err := te.CapabilitiesRegistry().GetDON(nil, donID)
						require.NoError(t, err)
						require.NotEqual(t, don.Id, donID, "Expect donID to not be found")
						require.Equal(t, uint32(0), don.Id, "Expected returned donID to be 0")
					}
				} else {
					// If we do NOT finalize, the DON still exists on-chain because the proposal wasn't applied
					for _, donID := range donIDs {
						_, err := te.CapabilitiesRegistry().GetDON(nil, donID)
						require.NoError(t, err, "DON is not removed because MCMS proposal was not finalized")
					}
				}
			}
		})
	}
}
