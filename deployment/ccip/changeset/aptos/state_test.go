package aptos

import (
	"testing"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/stretchr/testify/assert"
)

func TestLoadOnchainStateAptos(t *testing.T) {
	tests := []struct {
		name       string
		env        deployment.Environment
		want       map[uint64]changeset.AptosCCIPChainState
		err        error
		wantErrStr string
	}{
		{
			name: "success - empty env.AptosChains",
			env:  deployment.Environment{},
			want: map[uint64]changeset.AptosCCIPChainState{},
			err:  nil,
		},
		{
			name: "success - chain not found in ab returns empty state",
			env: deployment.Environment{
				AptosChains: map[uint64]deployment.AptosChain{
					743186221051783445: {},
				},
				ExistingAddresses: getTestAddressBook(
					map[uint64]map[string]deployment.TypeAndVersion{
						4457093679053095497: {
							mockMCMSAddress: {Type: changeset.AptosMCMSType},
						},
					},
				),
			},
			want: map[uint64]changeset.AptosCCIPChainState{
				743186221051783445: {},
			},
			err: nil,
		},
		{
			name: "success - loads multiple aptos chains state",
			env: deployment.Environment{
				AptosChains: map[uint64]deployment.AptosChain{
					743186221051783445:  {},
					4457093679053095497: {},
				},
				ExistingAddresses: getTestAddressBook(
					map[uint64]map[string]deployment.TypeAndVersion{
						4457093679053095497: {
							mockMCMSAddress: {Type: changeset.AptosMCMSType},
						},
						743186221051783445: {
							mockMCMSAddress: {Type: changeset.AptosMCMSType},
							mockCCIPAddress: {Type: changeset.AptosCCIPType},
						},
					},
				),
			},
			want: map[uint64]changeset.AptosCCIPChainState{
				4457093679053095497: {
					MCMSAddress: mustParseAddress(t, mockMCMSAddress),
				},
				743186221051783445: {
					MCMSAddress: mustParseAddress(t, mockMCMSAddress),
					CCIPAddress: mustParseAddress(t, mockCCIPAddress),
				},
			},
			err: nil,
		},
		{
			name: "error - failed to parse address",
			env: deployment.Environment{
				AptosChains: map[uint64]deployment.AptosChain{
					743186221051783445: {},
				},
				ExistingAddresses: getTestAddressBook(
					map[uint64]map[string]deployment.TypeAndVersion{
						743186221051783445: {
							mockBadAddress: {Type: changeset.AptosMCMSType},
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
			got, err := changeset.LoadOnchainStateAptos(tt.env)
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
