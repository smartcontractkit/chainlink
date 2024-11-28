package changeset_test

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/stretchr/testify/assert"
)

func TestAcceptOwnershipConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  changeset.AcceptOwnershipConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: changeset.AcceptOwnershipConfig{
				TimelocksPerChain: map[uint64]common.Address{
					1: common.HexToAddress("0x1"),
				},
				ProposerMCMSes: map[uint64]*gethwrappers.ManyChainMultiSig{
					1: {},
				},
				Contracts: map[uint64][]changeset.OwnershipAcceptor{
					1: {},
				},
				MinDelay: 3 * time.Hour,
			},
			wantErr: false,
		},
		{
			name: "missing timelock",
			config: changeset.AcceptOwnershipConfig{
				TimelocksPerChain: map[uint64]common.Address{},
				ProposerMCMSes: map[uint64]*gethwrappers.ManyChainMultiSig{
					1: {},
				},
				Contracts: map[uint64][]changeset.OwnershipAcceptor{
					1: {},
				},
				MinDelay: 3 * time.Hour,
			},
			wantErr: true,
		},
		{
			name: "missing proposer MCMS",
			config: changeset.AcceptOwnershipConfig{
				TimelocksPerChain: map[uint64]common.Address{
					1: common.HexToAddress("0x1"),
				},
				ProposerMCMSes: map[uint64]*gethwrappers.ManyChainMultiSig{},
				Contracts: map[uint64][]changeset.OwnershipAcceptor{
					1: {},
				},
				MinDelay: 3 * time.Hour,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
