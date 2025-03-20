package aptos

import (
	"testing"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/stretchr/testify/assert"
)

func TestLoadOnchainStateAptos(t *testing.T) {
	tests := []struct {
		name       string
		env        deployment.Environment
		want       map[uint64]AptosCCIPChainState
		err        error
		wantErrStr string
	}{
		{
			name: "success - empty env.AptosChains",
			env:  deployment.Environment{},
			want: map[uint64]AptosCCIPChainState{},
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
							mockMCMSAddress: {Type: AptosMCMSType},
						},
					},
				),
			},
			want: map[uint64]AptosCCIPChainState{
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
							mockMCMSAddress: {Type: AptosMCMSType},
						},
						743186221051783445: {
							mockMCMSAddress: {Type: AptosMCMSType},
							mockCCIPAddress: {Type: AptosCCIPType},
						},
					},
				),
			},
			want: map[uint64]AptosCCIPChainState{
				4457093679053095497: {
					AptosMCMSObjAddr: mustParseAddress(t, mockMCMSAddress),
				},
				743186221051783445: {
					AptosMCMSObjAddr: mustParseAddress(t, mockMCMSAddress),
					AptosCCIPObjAddr: mustParseAddress(t, mockCCIPAddress),
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
							mockBadAddress: {Type: AptosMCMSType},
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
			got, err := LoadOnchainStateAptos(tt.env)
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
