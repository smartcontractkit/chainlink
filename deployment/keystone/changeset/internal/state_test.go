package internal

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"

	"github.com/stretchr/testify/require"
)

func Test_GetContractSet(t *testing.T) {
	type testCase struct {
		name        string
		giveRequest *GetContractSetsRequest
		wantResp    *GetContractSetsResponse
		wantErr     string
	}
	lggr := logger.Test(t)
	chain := chainsel.ETHEREUM_TESTNET_SEPOLIA

	tt := []testCase{
		{
			name: "OK_picks unlabeled contracts by default",
			giveRequest: func() *GetContractSetsRequest {
				giveAB := map[uint64]map[string]cldf.TypeAndVersion{
					chain.Selector: {
						"0xabc": cldf.NewTypeAndVersion(
							cldf.ContractType("CapabilitiesRegistry"),
							deployment.Version1_0_0,
						),
						"0x123": func() cldf.TypeAndVersion {
							tv := cldf.NewTypeAndVersion(
								cldf.ContractType("CapabilitiesRegistry"),
								deployment.Version1_0_0,
							)
							tv.Labels.Add("SA")
							return tv
						}(),
					},
				}
				ab := cldf.NewMemoryAddressBookFromMap(
					giveAB,
				)
				req := &GetContractSetsRequest{
					Chains: map[uint64]cldf_evm.Chain{
						chain.Selector: {
							Selector: chain.Selector,
						},
					},
					AddressBook: ab,
				}
				return req
			}(),
			wantResp: &GetContractSetsResponse{
				ContractSets: map[uint64]ContractSet{
					chain.Selector: {
						MCMSWithTimelockState: commonchangeset.MCMSWithTimelockState{
							MCMSWithTimelockContracts: &proposalutils.MCMSWithTimelockContracts{},
						},
						CapabilitiesRegistry: func() *kcr.CapabilitiesRegistry {
							cr, err := kcr.NewCapabilitiesRegistry(common.HexToAddress("0xabc"), nil)
							require.NoError(t, err)
							return cr
						}(),
					},
				},
			},
		},
		{
			name: "OK_resolves labeled contracts",
			giveRequest: func() *GetContractSetsRequest {
				giveAB := map[uint64]map[string]cldf.TypeAndVersion{
					chain.Selector: {
						"0xabc": cldf.NewTypeAndVersion(
							cldf.ContractType("CapabilitiesRegistry"),
							deployment.Version1_0_0,
						),
						"0x123": func() cldf.TypeAndVersion {
							tv := cldf.NewTypeAndVersion(
								cldf.ContractType("CapabilitiesRegistry"),
								deployment.Version1_0_0,
							)
							tv.Labels.Add("SA")
							return tv
						}(),
					},
				}
				ab := cldf.NewMemoryAddressBookFromMap(
					giveAB,
				)
				req := &GetContractSetsRequest{
					Chains: map[uint64]cldf_evm.Chain{
						chain.Selector: {
							Selector: chain.Selector,
						},
					},
					Labels:      []string{"SA"},
					AddressBook: ab,
				}
				return req
			}(),
			wantResp: &GetContractSetsResponse{
				ContractSets: map[uint64]ContractSet{
					chain.Selector: {
						MCMSWithTimelockState: commonchangeset.MCMSWithTimelockState{
							MCMSWithTimelockContracts: &proposalutils.MCMSWithTimelockContracts{},
						},
						CapabilitiesRegistry: func() *kcr.CapabilitiesRegistry {
							cr, err := kcr.NewCapabilitiesRegistry(common.HexToAddress("0x123"), nil)
							require.NoError(t, err)
							return cr
						}(),
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gotResp, err := GetContractSets(lggr, tc.giveRequest)
			if tc.wantErr != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.wantErr)
				return
			}

			require.Equal(t, tc.wantResp, gotResp)
		})
	}
}
