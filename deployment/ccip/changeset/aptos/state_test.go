package aptos

import (
	"testing"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/aptos"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	aptosstate "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/aptos"
)

func TestLoadOnchainStateAptos(t *testing.T) {
	tests := []struct {
		name       string
		env        cldf.Environment
		want       map[uint64]aptosstate.CCIPChainState
		err        error
		wantErrStr string
	}{
		{
			name: "success - empty env.AptosChains",
			env:  cldf.Environment{},
			want: map[uint64]aptosstate.CCIPChainState{},
			err:  nil,
		},
		{
			name: "success - chain not found in ab returns empty state",
			env: cldf.Environment{
				BlockChains: chain.NewBlockChains(
					map[uint64]chain.BlockChain{
						743186221051783445: aptos.Chain{},
					},
				),
				ExistingAddresses: getTestAddressBook(
					t,
					map[uint64]map[string]cldf.TypeAndVersion{
						4457093679053095497: {
							mockMCMSAddress: {Type: shared.AptosMCMSType},
						},
					},
				),
			},
			want: map[uint64]aptosstate.CCIPChainState{
				743186221051783445: {},
			},
			err: nil,
		},
		{
			name: "success - loads multiple aptos chains state",
			env: cldf.Environment{
				BlockChains: chain.NewBlockChains(
					map[uint64]chain.BlockChain{
						743186221051783445:  aptos.Chain{},
						4457093679053095497: aptos.Chain{},
					},
				),
				ExistingAddresses: getTestAddressBook(
					t,
					map[uint64]map[string]cldf.TypeAndVersion{
						4457093679053095497: {
							mockMCMSAddress: {Type: shared.AptosMCMSType},
						},
						743186221051783445: {
							mockMCMSAddress: {Type: shared.AptosMCMSType},
							mockAddress:     {Type: shared.AptosCCIPType},
						},
					},
				),
			},
			want: map[uint64]aptosstate.CCIPChainState{
				4457093679053095497: {
					MCMSAddress: mustParseAddress(t, mockMCMSAddress),
				},
				743186221051783445: {
					MCMSAddress: mustParseAddress(t, mockMCMSAddress),
					CCIPAddress: mustParseAddress(t, mockAddress),
				},
			},
			err: nil,
		},
		{
			name: "error - failed to parse address",
			env: cldf.Environment{
				BlockChains: chain.NewBlockChains(
					map[uint64]chain.BlockChain{
						743186221051783445: aptos.Chain{},
					},
				),
				ExistingAddresses: getTestAddressBook(
					t,
					map[uint64]map[string]cldf.TypeAndVersion{
						743186221051783445: {
							"invalidaddress": {Type: shared.AptosMCMSType},
						},
					},
				),
			},
			want:       nil,
			err:        assert.AnError,
			wantErrStr: "failed to parse address 0xinvalid for AptosManyChainMultisig:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := aptosstate.LoadOnchainStateAptos(tt.env)
			if tt.err != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErrStr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
